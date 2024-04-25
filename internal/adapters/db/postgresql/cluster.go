package postgresql

import (
	"context"
	"errors"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/domain/models"
)

var (
	ErrClusterNotFound = errors.New("cluster not found")
)

func (p *Postgres) SaveCluster(ctx context.Context, cluster models.Cluster) error {
	const op = "postgresql.Postgres.SaveCluster"

	return p.db.WithContext(ctx).Create(&cluster).Error
}

func (p *Postgres) UpdateCluster(ctx context.Context, cluster models.Cluster) (models.Cluster, error) {
	const op = "postgresql.Postgres.UpdateCluster"

	if err := p.db.WithContext(ctx).Save(&cluster).Error; err != nil {
		return models.Cluster{}, err
	}

	return cluster, nil
}

func (p *Postgres) ClusterByID(ctx context.Context, id int64) (models.Cluster, error) {
	const op = "postgresql.Postgres.ClusterByID"

	var cluster models.Cluster
	if err := p.db.WithContext(ctx).Preload("Tasks").Preload("Cases").First(&cluster, id).Error; err != nil {
		return models.Cluster{}, err
	}

	return cluster, nil
}

func (p *Postgres) ListClusters(ctx context.Context) ([]models.Cluster, error) {
	const op = "postgresql.Postgres.ListClusters"

	var clusters []models.Cluster
	if err := p.db.WithContext(ctx).Preload("Tasks").Preload("Cases").Find(&clusters).Error; err != nil {
		return nil, err
	}

	return clusters, nil
}

func (p *Postgres) ClusterByIndex(ctx context.Context, index int64) (models.Cluster, error) {
	const op = "postgresql.Postgres.ClusterByIndex"

	var cluster models.Cluster
	if err := p.db.WithContext(ctx).Where("cluster_index = ?", index).Preload("Tasks").Preload("Cases").First(&cluster).Error; err != nil {
		return models.Cluster{}, err
	}

	return cluster, nil

}
