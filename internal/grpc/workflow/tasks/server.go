package tasks

import (
	"context"
	"errors"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/domain/models"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/services/workflow/tasks"
	tasksv1 "github.com/markgregr/bestHack_support_protos/gen/go/workflow/tasks"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskService interface {
	CreateTask(ctx context.Context, title, description string, clusterID int64) (models.Task, error)
	GetTask(ctx context.Context, taskID int64) (models.Task, error)
	ListTasks(ctx context.Context, status models.TaskStatus) ([]models.Task, error)
	ChangeTaskStatus(ctx context.Context, taskID int64) (models.Task, error)
}

type serverAPI struct {
	tasksv1.UnimplementedTaskServiceServer
	taskService TaskService
}

func Register(gRPC *grpc.Server, taskService TaskService) {
	tasksv1.RegisterTaskServiceServer(gRPC, &serverAPI{taskService: taskService})
}

func (s *serverAPI) CreateTask(ctx context.Context, req *tasksv1.CreateTaskRequest) (*tasksv1.Task, error) {
	task, err := s.taskService.CreateTask(ctx, req.GetTitle(), req.GetDescription(), req.GetClusterId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return ConvertTaskToProto(task), nil
}

func (s *serverAPI) GetTask(ctx context.Context, req *tasksv1.GetTaskRequest) (*tasksv1.Task, error) {
	task, err := s.taskService.GetTask(ctx, req.GetTaskId())
	if err != nil {
		if errors.Is(err, tasks.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return ConvertTaskToProto(task), nil
}

func (s *serverAPI) ListTasks(ctx context.Context, req *tasksv1.ListTasksRequest) (*tasksv1.ListTasksResponse, error) {
	tasks, err := s.taskService.ListTasks(ctx, models.TaskStatus(req.GetStatus()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &tasksv1.ListTasksResponse{Tasks: ConvertTaskListToProto(tasks)}, nil
}

func (s *serverAPI) ChangeTaskStatus(ctx context.Context, req *tasksv1.ChangeTaskStatusRequest) (*tasksv1.Task, error) {
	task, err := s.taskService.ChangeTaskStatus(ctx, req.GetTaskId())
	if err != nil {
		if errors.Is(err, tasks.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return ConvertTaskToProto(task), nil
}
