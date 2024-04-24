package postgresql

import (
	"context"
	"errors"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/domain/models"
)

var (
	ErrCaseNotFound = errors.New("case not found")
)

func (p *Postgres) SaveCase(ctx context.Context, caseItem models.Case) (models.Case, error) {
	err := p.db.WithContext(ctx).Create(&caseItem).Error
	return caseItem, err
}

func (p *Postgres) DeleteCase(ctx context.Context, id int64) error {
	err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Case{}).Error
	return err
}

func (p *Postgres) UpdateCase(ctx context.Context, caseItem models.Case) (models.Case, error) {
	err := p.db.WithContext(ctx).Save(&caseItem).Error
	return caseItem, err

}

func (p *Postgres) CaseByID(ctx context.Context, id int64) (models.Case, error) {
	var caseItem models.Case
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&caseItem).Error
	return caseItem, err
}

func (p *Postgres) ListCasesByClusterID(ctx context.Context, clusterID int64) ([]models.Case, error) {
	var cases []models.Case
	err := p.db.WithContext(ctx).Joins("Cluster").Where("cluster_id = ?", clusterID).Find(&cases).Error
	return cases, err
}
