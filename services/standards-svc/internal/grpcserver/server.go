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

package grpcserver

import (
	"context"
	"fmt"

	standardsv1 "github.com/itsLikith/h4d-cde/gen/standards"
	voxelizerv1 "github.com/itsLikith/h4d-cde/gen/voxelizer"

	"github.com/itsLikith/h4d-cde/services/standards-svc/internal/icaofpl"
)

type Server struct {
	standardsv1.UnimplementedStandardsServiceServer
	VoxelEngineClient voxelizerv1.VoxelEngineServiceClient
}

func NewServer(veClient voxelizerv1.VoxelEngineServiceClient) *Server {
	return &Server{
		VoxelEngineClient: veClient,
	}
}

// SubmitFlightPlan validates the ICAO FPL and forwards it to the Voxel Engine.
func (s *Server) SubmitFlightPlan(
	ctx context.Context,
	req *standardsv1.SubmitFlightPlanRequest,
) (*voxelizerv1.ProcessingResult, error) {
	if err := icaofpl.Validate(req.FlightPlan); err != nil {
		return nil, fmt.Errorf("flight plan validation failed: %w", err)
	}

	if s.VoxelEngineClient == nil {
		return nil, fmt.Errorf("voxel-engine client is not connected")
	}

	return s.VoxelEngineClient.IngestFlightPlan(ctx, &voxelizerv1.IngestFlightPlanRequest{
		FlightPlan: req.FlightPlan,
	})
}

// ValidateFlightPlan provides non-blocking pre-submission validation.
func (s *Server) ValidateFlightPlan(
	ctx context.Context,
	req *standardsv1.ValidateFlightPlanRequest,
) (*standardsv1.ValidateFlightPlanResponse, error) {
	err := icaofpl.Validate(req.FlightPlan)
	if err != nil {
		if valErr, ok := err.(*icaofpl.ValidationError); ok {
			return &standardsv1.ValidateFlightPlanResponse{
				IsValid:          false,
				ValidationErrors: valErr.Issues,
			}, nil
		}
		return &standardsv1.ValidateFlightPlanResponse{
			IsValid:          false,
			ValidationErrors: []string{err.Error()},
		}, nil
	}

	return &standardsv1.ValidateFlightPlanResponse{
		IsValid:          true,
		ValidationErrors: []string{},
	}, nil
}
