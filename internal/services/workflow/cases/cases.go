package cases

import (
	"context"
	"errors"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/adapters/db/postgresql"
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
}

type CaseProvider interface {
	CaseByID(ctx context.Context, caseID int64) (models.Case, error)
	ListCases(ctx context.Context) ([]models.Case, error)
}

type ClusterProvider interface {
	ClusterByID(ctx context.Context, clusterID int64) (models.Cluster, error)
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

func (s *CaseService) GetCase(ctx context.Context, id int64) (models.Case, error) {
	const op = "CaseService.GetCase"
	log := s.log.WithField("op", op)

	caseItem, err := s.caseProvider.CaseByID(ctx, id)
	if err != nil {
		if errors.Is(err, postgresql.ErrCaseNotFound) {
			log.Warn("case not found", err)
			return models.Case{}, ErrInvalidCredentials
		}

		log.WithError(err).Error("failed to get case")
		return models.Case{}, err
	}

	return caseItem, nil
}

func (s *CaseService) ListCases(ctx context.Context) ([]models.Case, error) {
	const op = "CaseService.ListCases"
	log := s.log.WithField("op", op)

	cases, err := s.caseProvider.ListCases(ctx)
	if err != nil {
		log.WithError(err).Error("failed to list cases")
		return nil, err
	}

	return cases, nil
}

func (s *CaseService) GetCluster(ctx context.Context, clusterID int64) (models.Cluster, error) {
	const op = "CaseService.GetCluster"
	log := s.log.WithField("op", op)

	cluster, err := s.clusterProvider.ClusterByID(ctx, clusterID)
	if err != nil {
		if errors.Is(err, postgresql.ErrClusterNotFound) {
			log.Warn("cluster not found", err)
			return models.Cluster{}, ErrInvalidCredentials
		}

		log.WithError(err).Error("failed to get cluster")
		return models.Cluster{}, err
	}

	return cluster, nil
}
