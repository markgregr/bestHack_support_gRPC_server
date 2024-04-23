package postgresql

import (
	"context"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/domain/models"
)

func (p *Postgres) CreateCluster(ctx context.Context, cluster models.Cluster) error {
	return p.db.WithContext(ctx).Create(cluster).Error
}

func (p *Postgres) UpdateCluster(ctx context.Context, cluster models.Cluster) error {
	return p.db.WithContext(ctx).Save(&cluster).Error
}

func (p *Postgres) ClusterByIndex(ctx context.Context, index int64) (models.Cluster, error) {
	var cluster models.Cluster
	err := p.db.WithContext(ctx).Where("cluster_index = ?", index).First(&cluster).Error
	return cluster, err
}
