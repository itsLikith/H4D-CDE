// Copyright 2026 Likith Saragadam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package grpcserver provides the gRPC service interface for VoxelEngineService.
package grpcserver

import (
	"context"
	"fmt"
	"sync"
	"time"

	commonv1 "github.com/itsLikith/h4d-cde/gen/common"
	flightplanv1 "github.com/itsLikith/h4d-cde/gen/flightplan"
	riskscorerv1 "github.com/itsLikith/h4d-cde/gen/riskscorer"
	trajectorypredictorv1 "github.com/itsLikith/h4d-cde/gen/trajectorypredictor"
	voxelizerv1 "github.com/itsLikith/h4d-cde/gen/voxelizer"

	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/adaptive"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/advisory"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/conflict"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/occupancy"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/spatial"
	"github.com/itsLikith/h4d-cde/services/voxel-engine/internal/temporal"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	voxelizerv1.UnimplementedVoxelEngineServiceServer
	Occupancy        occupancy.Map
	RedisClient      *redis.Client
	TrajectoryClient trajectorypredictorv1.TrajectoryPredictorServiceClient
	RiskScorerClient riskscorerv1.RiskScorerServiceClient
	AdaptiveEngine   *adaptive.Engine
	AdvisoryWeights  advisory.Weights
	RiskThreshold    float64
	HorizontalSepNM  float64
	VerticalSepFt    float64
	Publisher        EventPublisher

	mu              sync.RWMutex
	storedConflicts []*commonv1.ConflictRecord
}

func NewServer(
	occ occupancy.Map,
	rClient *redis.Client,
	tpClient trajectorypredictorv1.TrajectoryPredictorServiceClient,
	rsClient riskscorerv1.RiskScorerServiceClient,
	adaptEngine *adaptive.Engine,
	weights advisory.Weights,
	riskThreshold float64,
	hSepNM, vSepFt float64,
	pub EventPublisher,
) *Server {
	if pub == nil {
		pub = &NoOpPublisher{}
	}
	if riskThreshold <= 0 {
		riskThreshold = 0.50
	}
	if hSepNM <= 0 {
		hSepNM = 5.0
	}
	if vSepFt <= 0 {
		vSepFt = 1000.0
	}

	return &Server{
		Occupancy:        occ,
		RedisClient:      rClient,
		TrajectoryClient: tpClient,
		RiskScorerClient: rsClient,
		AdaptiveEngine:   adaptEngine,
		AdvisoryWeights:  weights,
		RiskThreshold:    riskThreshold,
		HorizontalSepNM:  hSepNM,
		VerticalSepFt:    vSepFt,
		Publisher:        pub,
		storedConflicts:  make([]*commonv1.ConflictRecord, 0),
	}
}

// IngestFlightPlan coordinates the end-to-end conflict detection & resolution pipeline:
//  1. Trajectory Prediction / Physics refinement (Module 2)
//  2. 4D Hexagonal Voxelization (Module 1)
//  3. Same-Voxel & Neighbor-Voxel conflict scans (O(n) linear complexity)
//  4. Probabilistic Risk Scoring via XGBoost (Module 4)
//  5. Strategic Advisory Resolution Cascade (Module 5)
//  6. Asynchronous Audit Logging & Real-time WebSocket broadcasting
func (s *Server) IngestFlightPlan(
	ctx context.Context,
	req *voxelizerv1.IngestFlightPlanRequest,
) (*voxelizerv1.ProcessingResult, error) {
	fpl := req.FlightPlan
	if fpl == nil {
		return nil, fmt.Errorf("flight_plan cannot be nil")
	}

	// 1. Refine trajectory via Trajectory Predictor gRPC service
	var points []*commonv1.TrajectoryPoint
	if s.TrajectoryClient != nil {
		tpResp, err := s.TrajectoryClient.RefineTrajectory(ctx, &trajectorypredictorv1.RefineTrajectoryRequest{
			FlightPlan:       fpl,
			WindSpeedKt:      10.0,
			WindDirectionDeg: 45.0,
			MaxAccelMps2:     2.5,
			AirDensityKgm3:   1.225,
		})
		if err == nil && tpResp != nil && len(tpResp.Points) > 0 {
			points = tpResp.Points
		}
	}

	// Fallback local kinematic interpolation if Trajectory Predictor is unavailable
	if len(points) == 0 {
		points = fallbackInterpolate(fpl)
	}

	// 2. Voxelize all trajectory points
	var keys []temporal.VoxelKey
	positions := make(map[string]conflict.Position)

	for _, pt := range points {
		resolution := spatial.DefaultResolution
		if s.AdaptiveEngine != nil && s.RedisClient != nil {
			baseCell, err := spatial.PointToH3Cell(pt.Lat, pt.Lon, spatial.DefaultResolution)
			if err == nil {
				forecast, _ := adaptive.ForecastFor(ctx, s.RedisClient, baseCell.String())
				resolution = s.AdaptiveEngine.ResolutionFor(forecast)
			}
		}

		key, err := temporal.ToVoxelKey(pt.Lat, pt.Lon, pt.AltFt, pt.TS, resolution)
		if err != nil {
			continue
		}

		_ = s.Occupancy.Add(ctx, key, fpl.EntityId)
		keys = append(keys, key)
		positions[fpl.EntityId] = conflict.Position{Lat: pt.Lat, Lon: pt.Lon, AltFt: pt.AltFt}

		s.Publisher.PublishAudit(ctx, map[string]any{
			"type":      "voxel_write",
			"key":       key.String(),
			"entity_id": fpl.EntityId,
			"lat":       pt.Lat,
			"lon":       pt.Lon,
			"alt_ft":    pt.AltFt,
			"t_s":       pt.TS,
		})
	}

	// 3. Detect Same-Voxel (Stage A) & Neighbor-Voxel (Stage B) conflicts
	samePairs, _ := conflict.SameVoxelConflicts(ctx, s.Occupancy, keys)
	nbrPairs, _ := conflict.NeighborVoxelConflicts(ctx, s.Occupancy, keys, positions, s.HorizontalSepNM, s.VerticalSepFt)

	allPairs := append(samePairs, nbrPairs...)

	// 4. Score Candidate Conflicts with Risk Scorer (Module 4)
	var conflicts []*commonv1.ConflictRecord
	nowMs := time.Now().UnixMilli()

	for i, pair := range allPairs {
		riskScore := 0.69 // Default baseline matching paper Table II average
		if s.RiskScorerClient != nil {
			scoreResp, err := s.RiskScorerClient.ScoreConflict(ctx, &riskscorerv1.ScoreConflictRequest{
				NEntitiesInConflict:  2,
				ClosureRateMps:       30.0,
				HeadingDiffDeg:       90.0,
				LocalTrafficDensity:  float64(len(allPairs)),
				SectorLoadForecast:   5.0,
				WindShearKtPer_100Ft: 2.0,
				VisibilityKm:         10.0,
			})
			if err == nil && scoreResp != nil {
				riskScore = scoreResp.RiskScore
			}
		}

		cType := commonv1.ConflictType_CONFLICT_TYPE_SAME_VOXEL
		if pair.ConflictType == conflict.ConflictTypeNeighborVoxel {
			cType = commonv1.ConflictType_CONFLICT_TYPE_NEIGHBOR_VOXEL
		}

		cRecord := &commonv1.ConflictRecord{
			ConflictId: fmt.Sprintf("CONF-%s-%s-%d", pair.EntityA, pair.EntityB, i+1),
			VoxelKey: &commonv1.VoxelKey{
				H3Cell:   pair.Key.H3Cell.String(),
				AltBinFt: int32(pair.Key.AltBinFt),
				TimeBinS: int32(pair.Key.TimeBinS),
			},
			Entities:         []string{pair.EntityA, pair.EntityB},
			ConflictType:     cType,
			RiskScore:        riskScore,
			DetectedAtUnixMs: nowMs,
		}
		conflicts = append(conflicts, cRecord)
	}

	// Save conflicts in state
	s.mu.Lock()
	s.storedConflicts = append(s.storedConflicts, conflicts...)
	s.mu.Unlock()

	// 5. Synthesize advisories for severe conflicts (Risk >= threshold)
	var advisories []*commonv1.Advisory
	for _, c := range conflicts {
		if c.RiskScore >= s.RiskThreshold {
			advRes, err := advisory.SelectAdvisory(ctx, advisory.Conflict{
				ConflictID: c.ConflictId,
				EntityA:    c.Entities[0],
				EntityB:    c.Entities[1],
			}, nil, s.AdvisoryWeights, 3)

			if err == nil {
				advStrategy := commonv1.AdvisoryStrategy_ADVISORY_STRATEGY_DELAY
				if advRes.Strategy == advisory.StrategyReroute {
					advStrategy = commonv1.AdvisoryStrategy_ADVISORY_STRATEGY_REROUTE
				} else if advRes.Strategy == advisory.StrategyAltitudeChange {
					advStrategy = commonv1.AdvisoryStrategy_ADVISORY_STRATEGY_ALTITUDE_CHANGE
				}

				advisories = append(advisories, &commonv1.Advisory{
					ConflictId:               c.ConflictId,
					Strategy:                 advStrategy,
					Parameters:               advRes.Parameters,
					ExpectedRiskReductionPct: advRes.ExpectedRiskReductionPct,
				})
			}
		}
	}

	// 6. Asynchronous side-effects: Audit stream + live dashboard broadcast
	s.Publisher.PublishAudit(ctx, map[string]any{
		"type":             "conflicts_detected",
		"flight_plan_id":   fpl.EntityId,
		"conflicts_count":  len(conflicts),
		"advisories_count": len(advisories),
		"timestamp":        nowMs,
	})
	s.Publisher.PublishConflictDetected(ctx, conflicts)

	return &voxelizerv1.ProcessingResult{
		Conflicts:  conflicts,
		Advisories: advisories,
	}, nil
}

// GetConflicts lists recorded conflicts matching optional cell and minimum risk filters.
func (s *Server) GetConflicts(
	ctx context.Context,
	req *voxelizerv1.GetConflictsRequest,
) (*voxelizerv1.GetConflictsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*commonv1.ConflictRecord
	for _, c := range s.storedConflicts {
		if req.H3Cell != "" && c.VoxelKey.H3Cell != req.H3Cell {
			continue
		}
		if req.MinRisk > 0 && c.RiskScore < req.MinRisk {
			continue
		}
		result = append(result, c)
	}

	return &voxelizerv1.GetConflictsResponse{Conflicts: result}, nil
}

// fallbackInterpolate generates basic waypoints when ML trajectory predictor is offline.
func fallbackInterpolate(fpl *flightplanv1.FlightPlan) []*commonv1.TrajectoryPoint {
	var points []*commonv1.TrajectoryPoint
	if len(fpl.Waypoints) == 0 {
		return points
	}

	eobtS := float64(fpl.EobtUnixMs) / 1000.0
	altFt := fpl.CruiseAltitudeFt
	if altFt <= 0 {
		altFt = 1500.0
	}

	for i, wp := range fpl.Waypoints {
		points = append(points, &commonv1.TrajectoryPoint{
			EntityId: fpl.EntityId,
			TS:       eobtS + float64(i*60),
			Lat:      wp.Lat,
			Lon:      wp.Lon,
			AltFt:    altFt,
		})
	}
	return points
}
