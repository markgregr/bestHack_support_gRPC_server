package cases

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/domain/models"
	casesv1 "github.com/markgregr/bestHack_support_protos/gen/go/workflow/cases"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CaseService interface {
	CreateCase(ctx context.Context, title, solution string, clusterID int64) (models.Case, error)
	UpdateCase(ctx context.Context, id int64, title, solution string) (models.Case, error)
	DeleteCase(ctx context.Context, id int64) error
	ListClusters(ctx context.Context, empty *empty.Empty) ([]models.Cluster, error)
	GetCasesFromCluster(ctx context.Context, clusterID int64) ([]models.Case, error)
	UpdateClusterName(ctx context.Context, clusterID int64) (models.Cluster, error)
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

func (s *serverAPI) UpdateCase(ctx context.Context, req *casesv1.UpdateCaseRequest) (*casesv1.Case, error) {
	caseItem, err := s.caseService.UpdateCase(ctx, req.GetId(), req.GetTitle(), req.GetSolution())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return ConvertCaseToProto(caseItem), nil
}

func (s *serverAPI) DeleteCase(ctx context.Context, req *casesv1.DeleteCaseRequest) (*empty.Empty, error) {
	err := s.caseService.DeleteCase(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &empty.Empty{}, nil
}

func (s *serverAPI) ListClusters(ctx context.Context, req *empty.Empty) (*casesv1.ListClustersResponse, error) {
	clusters, err := s.caseService.ListClusters(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &casesv1.ListClustersResponse{Clusters: ConvertClusterListToProto(clusters)}, nil
}

func (s *serverAPI) GetCasesFromCluster(ctx context.Context, req *casesv1.GetCasesFromClusterRequest) (*casesv1.GetCasesFromClusterResponse, error) {
	cases, err := s.caseService.GetCasesFromCluster(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &casesv1.GetCasesFromClusterResponse{Cases: ConvertCaseListToProto(cases)}, nil
}

func (s *serverAPI) UpdateClusterName(ctx context.Context, req *casesv1.UpdateClusterNameRequest) (*casesv1.Cluster, error) {
	cluster, err := s.caseService.UpdateClusterName(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return ConvertClusterToProto(cluster), nil
}
