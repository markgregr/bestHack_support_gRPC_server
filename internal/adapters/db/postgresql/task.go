package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/domain/models"
	"gorm.io/gorm"
)

var (
	ErrTaskNotFound = errors.New("tasks not found")
)

func (p *Postgres) TaskByID(ctx context.Context, id int64) (models.Task, error) {
	const op = "postgresql.Postgres.TaskByID"

	var task models.Task

	if err := p.db.WithContext(ctx).Where("id = ?", id).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return task, fmt.Errorf("%s: %w", op, ErrTaskNotFound)
		}
		return task, fmt.Errorf("%s: %w", op, err)
	}

	return task, nil
}

func (p *Postgres) ListTasks(ctx context.Context, status models.TaskStatus) ([]models.Task, error) {
	const op = "postgresql.Postgres.ListTasks"

	var tasks []models.Task

	if err := p.db.WithContext(ctx).Where("status = ?", status).Find(&tasks).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tasks, fmt.Errorf("%s: %w", op, ErrTaskNotFound)
		}
		return tasks, fmt.Errorf("%s: %w", op, err)
	}
	return tasks, nil
}

func (p *Postgres) SaveTask(ctx context.Context, task models.Task) (models.Task, error) {
	const op = "postgresql.Postgres.SaveTask"

	if err := p.db.WithContext(ctx).Create(&task).Error; err != nil {
		return models.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	return task, nil
}

func (p *Postgres) UpdateTask(ctx context.Context, id int64, task models.Task) error {
	const op = "postgresql.Postgres.UpdateTask"

	if err := p.db.WithContext(ctx).Model(&models.Task{}).Where("id = ?", id).Updates(&task).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
