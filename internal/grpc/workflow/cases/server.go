package cases

import (
	"context"
	"errors"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/domain/models"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/services/workflow/cases"
	casesv1 "github.com/markgregr/bestHack_support_protos/gen/go/workflow/cases"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CaseService interface {
	CreateCase(ctx context.Context, title, solution string, clusterID int64) (models.Case, error)
	GetCase(ctx context.Context, id int64) (models.Case, error)
	ListCases(ctx context.Context) ([]models.Case, error)
	GetCluster(ctx context.Context, id int64) (models.Cluster, error)
}

type serverAPI struct {
	casesv1.UnimplementedCaseServiceServer
	caseService CaseService
}

func Register(gRPC *grpc.Server, caseService CaseService) {
	casesv1.RegisterCaseServiceServer(gRPC, &serverAPI{caseService: caseService})
}

func (s *serverAPI) CreateCase(ctx context.Context, req *casesv1.CreateCaseRequest) (*casesv1.Case, error) {
	caseItem, err := s.caseService.CreateCase(ctx, req.GetTitle(), req.GetSolution(), req.GetClusterId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return ConvertCaseToProto(caseItem), nil
}

func (s *serverAPI) GetCase(ctx context.Context, req *casesv1.GetCaseRequest) (*casesv1.Case, error) {
	caseItem, err := s.caseService.GetCase(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, cases.ErrInvalidCredentials) {
			return nil, status.Error(codes.NotFound, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return ConvertCaseToProto(caseItem), nil
}

func (s *serverAPI) ListCases(ctx context.Context, empty *emptypb.Empty) (*casesv1.ListCasesResponse, error) {
	cases, err := s.caseService.ListCases(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &casesv1.ListCasesResponse{Cases: ConvertCaseListToProto(cases)}, nil
}

func (s *serverAPI) GetCluster(ctx context.Context, req *casesv1.GetClusterRequest) (*casesv1.Cluster, error) {
	cluster, err := s.caseService.GetCluster(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, cases.ErrInvalidCredentials) {
			return nil, status.Error(codes.NotFound, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return ConvertClusterToProto(cluster), nil
}
