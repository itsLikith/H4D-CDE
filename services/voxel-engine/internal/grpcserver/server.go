// services/voxel-engine/internal/grpcserver/server.go
package grpcserver

import (
	"context"

	commonv1 "hive/gen/common"
	riskscorerv1 "hive/gen/riskscorer"
	trajectorypredictorv1 "hive/gen/trajectorypredictor"
	voxelizerv1 "hive/gen/voxelizer"

	"hive/voxel-engine/internal/adaptive"
	"hive/voxel-engine/internal/advisory"
	"hive/voxel-engine/internal/conflict"
	"hive/voxel-engine/internal/occupancy"
	"hive/voxel-engine/internal/spatial"
	"hive/voxel-engine/internal/temporal"

)

type Server struct {
	voxelizerv1.UnimplementedVoxelEngineServiceServer
	Occupancy        occupancy.Map
	TrajectoryClient trajectorypredictorv1.TrajectoryPredictorServiceClient // gRPC, Part 10
	RiskScorerClient riskscorerv1.RiskScorerServiceClient                  // gRPC, Part 12
	AdaptiveEngine   *adaptive.Engine                                       // in-process, Part 14
	AdvisoryWeights  advisory.Weights                                       // in-process, Part 13
	RiskThreshold    float64
	Publisher        EventPublisher // Kafka, Part 17.3
}

func (s *Server) IngestFlightPlan(ctx context.Context, req *voxelizerv1.IngestFlightPlanRequest) (*voxelizerv1.ProcessingResult, error) {
	fpl := req.FlightPlan

	// 1. Dense trajectory from trajectory-predictor-svc (Python, gRPC call #1)
	traj, err := s.TrajectoryClient.RefineTrajectory(ctx, &trajectorypredictorv1.RefineTrajectoryRequest{FlightPlan: fpl})
	if err != nil {
		return nil, err
	}

	// 2. Voxelize every point, consulting Adaptive Discretization (Part 14)
	var keys []temporal.VoxelKey
	positions := map[string]conflict.Position{}
	for _, pt := range traj.Points {
		baseCell, _ := spatial.PointToH3Cell(pt.Lat, pt.Lon, spatial.DefaultResolution)
		forecast, _ := adaptive.ForecastFor(ctx, s.redisClient(), baseCell.String())
		resolution := s.AdaptiveEngine.ResolutionFor(forecast)

		key, err := temporal.ToVoxelKey(pt.Lat, pt.Lon, pt.AltFt, pt.TS, resolution)
		if err != nil {
			continue
		}
		_ = s.Occupancy.Add(ctx, key, fpl.EntityId)
		s.Publisher.PublishAudit(ctx, map[string]any{"type": "voxel_write", "key": key.String(), "entity": fpl.EntityId})
		keys = append(keys, key)
		positions[fpl.EntityId] = conflict.Position{Lat: pt.Lat, Lon: pt.Lon, AltFt: pt.AltFt}
	}

	// 3. Same-voxel + neighbour-voxel checks (Part 9.5-9.6), no gRPC involved
	pairs, _ := conflict.SameVoxelConflicts(ctx, s.Occupancy, keys)
	nbrPairs, _ := conflict.NeighborVoxelConflicts(ctx, s.Occupancy, keys, positions, 5.0, 1000.0)
	pairs = append(pairs, nbrPairs...)

	// 4. Score every candidate via risk-scorer-svc (Python, gRPC call #2 per candidate)
	var conflicts []*commonv1.ConflictRecord
	for _, p := range pairs {
		scoreResp, err := s.RiskScorerClient.ScoreConflict(ctx, buildRiskRequest(p))
		if err != nil {
			continue
		}
		conflicts = append(conflicts, toConflictRecord(p, scoreResp.RiskScore))
	}

	// 5. Advisory Selector (Part 13), in-process, for conflicts above threshold
	var advisories []*commonv1.Advisory
	for _, c := range conflicts {
		if c.RiskScore < s.RiskThreshold {
			continue
		}
		result, err := advisory.SelectAdvisory(ctx, toAdvisoryConflict(c), s, s.AdvisoryWeights, 3)
		if err == nil {
			advisories = append(advisories, toProtoAdvisory(result))
		}
	}

	// 6. Fire-and-forget: audit trail + live dashboard updates (Part 4.4) --
	// neither blocks the response below.
	s.Publisher.PublishAudit(ctx, map[string]any{"type": "conflicts_detected", "count": len(conflicts)})
	s.Publisher.PublishConflictDetected(ctx, conflicts)

	return &voxelizerv1.ProcessingResult{Conflicts: conflicts, Advisories: advisories}, nil
}