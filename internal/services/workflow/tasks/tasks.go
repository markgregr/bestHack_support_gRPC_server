package tasks

import (
	"context"
	"errors"
	"fmt"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/adapters/db/postgresql"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/domain/models"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/services/user"
	"github.com/sirupsen/logrus"
	"time"
)

type TaskService struct {
	log             *logrus.Logger
	taskSaver       TaskSaver
	taskProvider    TaskProvider
	clusterSaver    ClusterSaver
	clusterProvider ClusterProvider

	userService user.UserService
}

type TaskSaver interface {
	SaveTask(ctx context.Context, task models.Task) (createdTask models.Task, err error)
	UpdateTask(ctx context.Context, id int64, task models.Task) error
}

type TaskProvider interface {
	TaskByID(ctx context.Context, taskID int64) (models.Task, error)
	ListTasks(ctx context.Context, status models.TaskStatus) ([]models.Task, error)
}

type ClusterSaver interface {
	SaveCluster(ctx context.Context, cluster models.Cluster) error
	UpdateCluster(ctx context.Context, cluster models.Cluster) error
}

type ClusterProvider interface {
	ClusterByIndex(ctx context.Context, index int64) (models.Cluster, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func New(log *logrus.Logger, taskSaver TaskSaver, taskProvider TaskProvider, clusterProvider ClusterProvider, clusterSaver ClusterSaver, userService user.UserService) *TaskService {
	return &TaskService{
		log:             log,
		taskSaver:       taskSaver,
		taskProvider:    taskProvider,
		clusterSaver:    clusterSaver,
		clusterProvider: clusterProvider,
		userService:     userService,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, title string, description string, clusterIndex int64) (models.Task, error) {
	const op = "TaskService.CreateTask"
	log := s.log.WithField("op", op)

	log.Info("create by index")
	cluster, err := s.clusterProvider.ClusterByIndex(ctx, clusterIndex)
	if err != nil {
		log.Warn("cluster not found", err)
		cluster = models.Cluster{
			ClusterIndex: clusterIndex,
			Name:         "Cluster " + fmt.Sprint(clusterIndex),
			Frequency:    0,
		}

		log.Info("create cluster")
		if err := s.clusterSaver.SaveCluster(ctx, cluster); err != nil {
			log.WithError(err).Error("failed to create cluster")
			return models.Task{}, err
		}
	}

	task := models.Task{
		Title:       title,
		Description: description,
		Status:      models.TaskStatusOpen,
		ClusterID:   &cluster.ID,
		Cluster:     &cluster,
	}

	log.Info("create tasks")
	task, err = s.taskSaver.SaveTask(ctx, task)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, id int64) (models.Task, error) {
	const op = "TaskService.GetTask"
	log := s.log.WithField("op", op)

	log.Info("get tasks")
	task, err := s.taskProvider.TaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, postgresql.ErrTaskNotFound) {
			log.Warn("tasks not found", err)
			return models.Task{}, ErrInvalidCredentials
		}

		log.WithError(err).Error("failed to get tasks")
		return models.Task{}, err
	}

	return task, nil
}

func (s *TaskService) ListTasks(ctx context.Context, status models.TaskStatus) ([]models.Task, error) {
	const op = "TaskService.ListTasks"
	log := s.log.WithField("op", op)

	log.Info("list tasks")
	tasks, err := s.taskProvider.ListTasks(ctx, status)
	if err != nil {
		log.WithError(err).Error("failed to list tasks")
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) ChangeTaskStatus(ctx context.Context, taskID int64) (models.Task, error) {
	const op = "TaskService.ChangeTaskStatus"
	log := s.log.WithField("op", op)

	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		log.Error("failed to get userID from context")
		return models.Task{}, errors.New("failed to get userID from context")
	}

	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, postgresql.ErrUsersNotFound) {
			log.Warn("user not found", err)
			return models.Task{}, ErrInvalidCredentials
		}
		log.WithError(err).Error("failed to get user")
		return models.Task{}, err
	}

	task, err := s.taskProvider.TaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, postgresql.ErrTaskNotFound) {
			log.Warn("tasks not found", err)
			return models.Task{}, ErrInvalidCredentials
		}

		log.WithError(err).Error("failed to get tasks")
		return models.Task{}, err
	}

	switch task.Status {
	case models.TaskStatusOpen:
		task.Status = models.TaskStatusInProgress

		task.UserID = &user.ID
		task.User = &user

		currTime := time.Now()
		task.FormedAt = &currTime
	case models.TaskStatusInProgress:
		task.Status = models.TaskStatusClosed
		currTime := time.Now()
		task.CompletedAt = &currTime
	default:
		return models.Task{}, errors.New("invalid task status")
	}

	log.Info("change tasks status")
	if err := s.taskSaver.UpdateTask(ctx, taskID, task); err != nil {
		log.WithError(err).Error("failed to update tasks")
		return models.Task{}, err
	}

	task, err = s.taskProvider.TaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, postgresql.ErrTaskNotFound) {
			log.Warn("tasks not found", err)
			return models.Task{}, ErrInvalidCredentials
		}

		log.WithError(err).Error("failed to get tasks")
		return models.Task{}, err
	}

	return task, nil
}
