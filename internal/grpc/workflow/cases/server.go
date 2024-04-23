package cases

import (
	"context"
	"errors"
	casesv1 "github.com/markgregr/bestHack_support_protos/gen/go/workflow/cases"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CaseService interface {
	CreateCase(ctx context.Context, clusterID int64, title, description, solution string) (*casesv1.Case, error)
	GetCase(ctx context.Context, id int64) (*casesv1.Case, error)
	ListCases(ctx context.Context, clusterID int64) ([]*casesv1.Case, error)
}

type serverAPI struct {
	casesv1.UnimplementedCaseServiceServer
	caseService CaseService
}

func Register(gRPC *grpc.Server, caseService CaseService) {
	casesv1.RegisterCaseServiceServer(gRPC, &serverAPI{caseService: caseService})
}

func (s *serverAPI) CreateCase(ctx context.Context, req *casesv1.CreateCaseRequest) (*casesv1.Case, error) {
	caseItem, err := s.caseService.CreateCase(ctx, req.GetClusterId(), req.GetTitle(), req.GetDescription(), req.GetSolution())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return caseItem, nil
}

func (s *serverAPI) GetCase(ctx context.Context, req *casesv1.GetCaseRequest) (*casesv1.Case, error) {
	caseItem, err := s.caseService.GetCase(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, ErrCaseNotFound) {
			return nil, status.Error(codes.NotFound, "case not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return caseItem, nil
}

func (s *serverAPI) ListCases(ctx context.Context, req *casesv1.ListCasesRequest) (*casesv1.ListCasesResponse, error) {
	cases, err := s.caseService.ListCases(ctx, req.GetClusterId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &casesv1.ListCasesResponse{Cases: cases}, nil
}
