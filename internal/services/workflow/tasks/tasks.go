package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/adapters/db/postgresql"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/domain/models"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/services/user"
	"github.com/markgregr/bestHack_support_gRPC_server/pkg/csvsaver"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"time"
)

type TaskService struct {
	log             *logrus.Logger
	outputFileData  string
	inputFileData   string
	AnalURL         string
	taskSaver       TaskSaver
	taskProvider    TaskProvider
	clusterSaver    ClusterSaver
	clusterProvider ClusterProvider
	caseProvider    CaseProvider

	userService user.UserService
}

type TaskSaver interface {
	SaveTask(ctx context.Context, task models.Task) (createdTask models.Task, err error)
	UpdateTask(ctx context.Context, id int64, task models.Task) error
	UserList(ctx context.Context) ([]models.User, error)
}

type TaskProvider interface {
	TaskByID(ctx context.Context, taskID int64) (models.Task, error)
	ListTasks(ctx context.Context, status models.TaskStatus) ([]models.Task, error)
	UserWithMinAverageDuration(ctx context.Context) (models.User, error)
	ListTasksByUserID(ctx context.Context, userID int64, status models.TaskStatus) ([]models.Task, error)
}

type ClusterSaver interface {
	SaveCluster(ctx context.Context, cluster models.Cluster) error
}

type ClusterProvider interface {
	ClusterByIndex(ctx context.Context, index int64) (models.Cluster, error)
	UpdateCluster(ctx context.Context, cluster models.Cluster) (models.Cluster, error)
}

type CaseProvider interface {
	CaseByID(ctx context.Context, caseID int64) (models.Case, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type NotficationRequest struct {
	Username string `json:"username"`
}

func New(log *logrus.Logger, inputFileData, outputFileData, AnalURL string, taskSaver TaskSaver, taskProvider TaskProvider, clusterProvider ClusterProvider, clusterSaver ClusterSaver, caseProvider CaseProvider, userService user.UserService) *TaskService {
	return &TaskService{
		log:             log,
		outputFileData:  outputFileData,
		inputFileData:   inputFileData,
		AnalURL:         AnalURL,
		taskSaver:       taskSaver,
		taskProvider:    taskProvider,
		clusterSaver:    clusterSaver,
		clusterProvider: clusterProvider,
		caseProvider:    caseProvider,
		userService:     userService,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, title string, description string, clusterIndex int64, clusterName string, frequency int64, avarage_duration float32) (models.Task, error) {
	const op = "TaskService.CreateTask"
	log := s.log.WithField("op", op)

	log.Info("create by index")
	cluster, err := s.clusterProvider.ClusterByIndex(ctx, clusterIndex)
	if err != nil {
		log.Warn("cluster not found", err)
		cluster = models.Cluster{
			ClusterIndex: clusterIndex,
			Name:         clusterName,
			Frequency:    frequency,
		}

		log.Info("create cluster")
		if err := s.clusterSaver.SaveCluster(ctx, cluster); err != nil {
			log.WithError(err).Error("failed to create cluster")
			return models.Task{}, err
		}
	}

	task := models.Task{
		Title:           title,
		Description:     description,
		Status:          models.TaskStatusOpen,
		AvarageDuration: avarage_duration,
		ClusterID:       &cluster.ID,
		Cluster:         &cluster,
		Fire:            false,
	}

	log.WithField("task", task).Info("create tasks")
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

		err = s.userService.UpdateUserAvarageDuration(ctx, user.ID, user.AvarageDuration+task.AvarageDuration)
		if err != nil {
			log.WithError(err).Error("failed to update user avarage duration")
			return models.Task{}, err
		}
	case models.TaskStatusInProgress:
		task.Status = models.TaskStatusClosed
		currTime := time.Now()
		task.CompletedAt = &currTime

		// Вычисление времени формирования и времени начала выполнения в секундах
		formedAtUnix := task.FormedAt.Unix()
		startedAtUnix := task.CreatedAt.Unix()

		// Вычисление разницы в секундах
		reactionTimeInSeconds := int(formedAtUnix - startedAtUnix)
		durationInSeconds := int(currTime.Unix() - formedAtUnix)

		clusterData := csvsaver.ClusterData{
			ClusterIndex: int(task.Cluster.ClusterIndex),
			ReactionTime: reactionTimeInSeconds,
			DurationTime: durationInSeconds,
		}
		err = csvsaver.AddDataToJSON(s.outputFileData, clusterData, log.Logger)
		if err != nil {
			log.WithError(err).Error("failed to add data to JSON")
			return models.Task{}, err
		}
		err = csvsaver.AvgCsv(s.inputFileData, s.outputFileData, log.Logger)
		if err != nil {
			log.WithError(err).Error("failed to calculate average")
			return models.Task{}, err
		}

		err = s.userService.UpdateUserAvarageDuration(ctx, user.ID, user.AvarageDuration-task.AvarageDuration)
		if err != nil {
			log.WithError(err).Error("failed to update user avarage duration")
			return models.Task{}, err
		}

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

func (s *TaskService) AddCaseToTask(ctx context.Context, taskID, caseID int64) (models.Task, error) {
	const op = "TaskService.AddCaseToTask"
	log := s.log.WithField("op", op)

	log.Info("add case to task")
	task, err := s.taskProvider.TaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, postgresql.ErrTaskNotFound) {
			log.Warn("tasks not found", err)
			return models.Task{}, ErrInvalidCredentials
		}

		log.WithError(err).Error("failed to get tasks")
		return models.Task{}, err
	}

	caseItem, err := s.caseProvider.CaseByID(ctx, caseID)
	if err != nil {
		if errors.Is(err, postgresql.ErrCaseNotFound) {
			log.Warn("case not found", err)
			return models.Task{}, ErrInvalidCredentials
		}

		log.WithError(err).Error("failed to get case")
		return models.Task{}, err

	}
	task.CaseID = &caseID
	task.Case = &caseItem

	log.Info("change tasks status")
	if err := s.taskSaver.UpdateTask(ctx, taskID, task); err != nil {
		log.WithError(err).Error("failed to update tasks")
		return models.Task{}, err
	}

	return task, nil
}

func (s *TaskService) AddSolutionToTask(ctx context.Context, taskID int64, solution string) (models.Task, error) {
	const op = "TaskService.AddCaseToTask"
	log := s.log.WithField("op", op)

	log.Info("add case to task")
	task, err := s.taskProvider.TaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, postgresql.ErrTaskNotFound) {
			log.Warn("tasks not found", err)
			return models.Task{}, ErrInvalidCredentials
		}

		log.WithError(err).Error("failed to get tasks")
		return models.Task{}, err
	}

	task.Solution = &solution

	log.Info("change tasks status")
	if err := s.taskSaver.UpdateTask(ctx, taskID, task); err != nil {
		log.WithError(err).Error("failed to update tasks")
		return models.Task{}, err
	}

	return task, nil
}

func (s *TaskService) RemoveCaseFromTask(ctx context.Context, taskID int64) (models.Task, error) {
	const op = "TaskService.RemoveCaseFromTask"
	log := s.log.WithField("op", op)

	log.Info("remove case from task")
	task, err := s.taskProvider.TaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, postgresql.ErrTaskNotFound) {
			log.Warn("tasks not found", err)
			return models.Task{}, ErrInvalidCredentials
		}

		log.WithError(err).Error("failed to get tasks")
		return models.Task{}, err
	}

	task.CaseID = nil
	task.Case = nil

	log.Info("change tasks status")
	if err := s.taskSaver.UpdateTask(ctx, taskID, task); err != nil {
		log.WithError(err).Error("failed to update tasks")
		return models.Task{}, err
	}

	return task, nil
}

func (s *TaskService) RemoveSolutionFromTask(ctx context.Context, taskID int64) (models.Task, error) {
	const op = "TaskService.RemoveSolutionFromTask"
	log := s.log.WithField("op", op)

	log.Info("remove solution from task")
	task, err := s.taskProvider.TaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, postgresql.ErrTaskNotFound) {
			log.Warn("tasks not found", err)
			return models.Task{}, ErrInvalidCredentials
		}

		log.WithError(err).Error("failed to get tasks")
		return models.Task{}, err
	}

	task.Solution = nil

	log.Info("change tasks status")
	if err := s.taskSaver.UpdateTask(ctx, taskID, task); err != nil {
		log.WithError(err).Error("failed to update tasks")
		return models.Task{}, err
	}

	return task, nil
}

func (s *TaskService) AppointUserToTask(ctx context.Context, taskID int64) (models.Task, error) {
	const op = "TaskService.AppointUserToTask"
	log := s.log.WithField("op", op)

	user, err := s.taskProvider.UserWithMinAverageDuration(ctx)
	if err != nil {
		log.WithError(err).Error("failed to get user with min avarage duration")
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

	if task.User != nil {
		return models.Task{}, errors.New("user already appointed")
	}

	task.UserID = &user.ID
	task.User = &user

	task.Status = models.TaskStatusInProgress
	currTime := time.Now()
	task.FormedAt = &currTime

	log.Info("change tasks status")
	if err := s.taskSaver.UpdateTask(ctx, taskID, task); err != nil {
		log.WithError(err).Error("failed to update tasks")
		return models.Task{}, err
	}
	err = s.userService.UpdateUserAvarageDuration(ctx, user.ID, user.AvarageDuration+task.AvarageDuration)
	if err != nil {
		log.WithError(err).Error("failed to update user avarage duration")
		return models.Task{}, err
	}

	requestBody, err := json.Marshal(NotficationRequest{
		Username: user.TelegramUsername,
	})
	if err != nil {
		log.WithError(err).Error("failed to marshal request body")
		return models.Task{}, err
	}

	resp, err := http.Post(s.AnalURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.WithError(err).Error("failed to send request")
		return models.Task{}, err
	}
	defer resp.Body.Close()

	return task, nil
}

func (s *TaskService) FireTask(ctx context.Context, taskID int64) (models.Task, error) {
	const op = "TaskService.FireTask"
	log := s.log.WithField("op", op)

	task, err := s.taskProvider.TaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, postgresql.ErrTaskNotFound) {
			log.Warn("tasks not found", err)
			return models.Task{}, ErrInvalidCredentials
		}

		log.WithError(err).Error("failed to get tasks")
		return models.Task{}, err
	}

	task.Fire = true

	log.Info("change tasks status")
	if err := s.taskSaver.UpdateTask(ctx, taskID, task); err != nil {
		log.WithError(err).Error("failed to update tasks")
		return models.Task{}, err
	}

	return task, nil
}

func (s *TaskService) ListTasksByUserID(ctx context.Context, userID int64, status models.TaskStatus) ([]models.Task, error) {
	const op = "TaskService.ListTasksByUserID"
	log := s.log.WithField("op", op)

	log.Info("list tasks by user id")
	tasks, err := s.taskProvider.ListTasksByUserID(ctx, userID, status)
	if err != nil {
		log.WithError(err).Error("failed to list tasks")
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) ListUsers(ctx context.Context, empty *emptypb.Empty) ([]models.User, error) {
	const op = "TaskService.ListUsers"
	log := s.log.WithField("op", op)

	log.Info("list users")
	users, err := s.taskSaver.UserList(ctx)
	if err != nil {
		log.WithError(err).Error("failed to list users")
		return nil, err
	}

	return users, nil
}
