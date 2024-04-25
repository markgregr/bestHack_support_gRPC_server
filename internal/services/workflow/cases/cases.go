package cases

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/domain/models"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/services/user"
	"github.com/sirupsen/logrus"
)

type CaseService struct {
	log             *logrus.Logger
	caseSaver       CaseSaver
	caseProvider    CaseProvider
	clusterProvider ClusterProvider

	userService user.UserService
}

type CaseSaver interface {
	SaveCase(ctx context.Context, caseItem models.Case) (createdCase models.Case, err error)
	UpdateCase(ctx context.Context, caseItem models.Case) (updatedCase models.Case, err error)
	DeleteCase(ctx context.Context, caseID int64) error
}

type CaseProvider interface {
	CaseByID(ctx context.Context, caseID int64) (models.Case, error)
	ListCasesByClusterID(ctx context.Context, clusterID int64) ([]models.Case, error)
}

type ClusterProvider interface {
	ClusterByID(ctx context.Context, clusterID int64) (models.Cluster, error)
	ListClusters(ctx context.Context) ([]models.Cluster, error)
	UpdateCluster(ctx context.Context, cluster models.Cluster) (models.Case, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func New(log *logrus.Logger, caseSaver CaseSaver, caseProvider CaseProvider, clusterProvider ClusterProvider, userService user.UserService) *CaseService {
	return &CaseService{
		log:             log,
		caseSaver:       caseSaver,
		caseProvider:    caseProvider,
		clusterProvider: clusterProvider,
		userService:     userService,
	}
}

func (s *CaseService) CreateCase(ctx context.Context, title string, solution string, clusterID int64) (models.Case, error) {
	const op = "CaseService.CreateCase"
	log := s.log.WithField("op", op)

	log.Info("get cluster by id")
	cluster, err := s.clusterProvider.ClusterByID(ctx, clusterID)
	if err != nil {
		log.WithError(err).Error("failed to get cluster")
		return models.Case{}, err
	}

	caseItem := models.Case{
		Title:    title,
		Solution: solution,
		Cluster:  &cluster,
	}

	createdCase, err := s.caseSaver.SaveCase(ctx, caseItem)
	if err != nil {
		log.WithError(err).Error("failed to save case")
		return models.Case{}, err
	}

	return createdCase, nil
}

func (s *CaseService) UpdateCase(ctx context.Context, id int64, title string, solution string) (models.Case, error) {
	const op = "CaseService.UpdateCase"
	log := s.log.WithField("op", op)

	caseItem, err := s.caseProvider.CaseByID(ctx, id)
	if err != nil {
		log.WithError(err).Error("failed to get case")
		return models.Case{}, err
	}

	caseItem.Title = title
	caseItem.Solution = solution

	updatedCase, err := s.caseSaver.UpdateCase(ctx, caseItem)
	if err != nil {
		log.WithError(err).Error("failed to update case")
		return models.Case{}, err
	}

	return updatedCase, nil
}

func (s *CaseService) DeleteCase(ctx context.Context, id int64) error {
	const op = "CaseService.DeleteCase"
	log := s.log.WithField("op", op)

	err := s.caseSaver.DeleteCase(ctx, id)
	if err != nil {
		log.WithError(err).Error("failed to delete case")
		return err
	}

	return nil
}

func (s *CaseService) ListClusters(ctx context.Context, empty *empty.Empty) ([]models.Cluster, error) {
	const op = "CaseService.ListClusters"
	log := s.log.WithField("op", op)

	clusters, err := s.clusterProvider.ListClusters(ctx)
	if err != nil {
		log.WithError(err).Error("failed to list clusters")
		return nil, err
	}

	return clusters, nil
}

func (s *CaseService) GetCasesFromCluster(ctx context.Context, clusterID int64) ([]models.Case, error) {
	const op = "CaseService.GetCasesFromCluster"
	log := s.log.WithField("op", op)

	cases, err := s.caseProvider.ListCasesByClusterID(ctx, clusterID)
	if err != nil {
		log.WithError(err).Error("failed to list cases by cluster id")
		return nil, err
	}

	return cases, nil
}

func (s *CaseService) UpdateClusterName(ctx context.Context, clusterID int64) (models.Case, error) {
	const op = "CaseService.UpdateClusterName"
	log := s.log.WithField("op", op)

	cluster, err := s.clusterProvider.ClusterByID(ctx, clusterID)
	if err != nil {
		log.WithError(err).Error("failed to get cluster")
		return models.Case{}, err
	}

	cluster.Name = "New name"

	updatedCluster, err := s.clusterProvider.UpdateCluster(ctx, cluster)
	if err != nil {
		log.WithError(err).Error("failed to update cluster")
		return models.Case{}, err
	}

	return updatedCluster, nil
}
