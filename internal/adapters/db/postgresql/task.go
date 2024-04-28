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
	if err := p.db.WithContext(ctx).Joins("User").Joins("Case").Joins("Cluster").First(&task, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Task{}, fmt.Errorf("%s: %w", op, ErrTaskNotFound)
		}
	}

	// Дополнительная проверка, чтобы убедиться, что поле Cluster заполнено
	if task.Cluster == nil {
		return models.Task{}, fmt.Errorf("%s: cluster is nil for task with ID %d", op, id)
	}

	return task, nil
}

func (p *Postgres) ListTasks(ctx context.Context, status models.TaskStatus) ([]models.Task, error) {
	const op = "postgresql.Postgres.ListTasks"

	var tasks []models.Task
	if err := p.db.WithContext(ctx).Joins("User").Joins("Case").Joins("Cluster").Where("tasks.status = ?", status).Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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
	err := p.db.WithContext(ctx).Save(&task).Error
	return err
}

func (p *Postgres) ListTaskByUserID(ctx context.Context, userID int64) ([]models.Task, error) {
	const op = "postgresql.Postgres.ListTaskByUserID"

	var tasks []models.Task
	if err := p.db.WithContext(ctx).Joins("User").Joins("Case").Joins("Cluster").Where("tasks.user_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}
