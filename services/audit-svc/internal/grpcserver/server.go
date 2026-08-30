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

	auditv1 "github.com/itsLikith/h4d-cde/gen/audit"

	"github.com/itsLikith/h4d-cde/services/audit-svc/internal/chain"
	"github.com/itsLikith/h4d-cde/services/audit-svc/internal/storage"
)

type Server struct {
	auditv1.UnimplementedAuditServiceServer
	Chain *chain.Chain
	Store storage.Store
}

func NewServer(c *chain.Chain, store storage.Store) *Server {
	return &Server{
		Chain: c,
		Store: store,
	}
}

func (s *Server) GetVoxelAuditTrail(
	ctx context.Context,
	req *auditv1.GetVoxelAuditTrailRequest,
) (*auditv1.GetVoxelAuditTrailResponse, error) {
	entries := s.Chain.Entries()
	valid := chain.Verify(entries)

	var protoEntries []*auditv1.AuditEntry
	for _, e := range entries {
		protoEntries = append(protoEntries, &auditv1.AuditEntry{
			TimestampUnixMs: e.TimestampUnixMs,
			EventType:       e.EventType,
			PayloadJson:     string(e.Event),
			PrevHash:        e.PrevHash,
			Hash:            e.Hash,
		})
	}

	return &auditv1.GetVoxelAuditTrailResponse{
		Entries:    protoEntries,
		ChainValid: valid,
	}, nil
}

func (s *Server) VerifyAuditChain(
	ctx context.Context,
	req *auditv1.VerifyAuditChainRequest,
) (*auditv1.VerifyAuditChainResponse, error) {
	entries := s.Chain.Entries()
	valid := chain.Verify(entries)

	lastHash := ""
	if len(entries) > 0 {
		lastHash = entries[len(entries)-1].Hash
	}

	return &auditv1.VerifyAuditChainResponse{
		IsValid:              valid,
		TotalEntriesVerified: int64(len(entries)),
		LastHash:             lastHash,
	}, nil
}
