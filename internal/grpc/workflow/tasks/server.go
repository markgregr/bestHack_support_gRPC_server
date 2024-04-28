package tasks

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/domain/models"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/services/workflow/tasks"
	tasksv1 "github.com/markgregr/bestHack_support_protos/gen/go/workflow/tasks"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskService interface {
	CreateTask(ctx context.Context, title, description string, clusterIndex int64, clusterName string, frequency int64, avarage_duration float32) (models.Task, error)
	GetTask(ctx context.Context, taskID int64) (models.Task, error)
	ListTasks(ctx context.Context, status models.TaskStatus) ([]models.Task, error)
	ChangeTaskStatus(ctx context.Context, taskID int64) (models.Task, error)
	AddCaseToTask(ctx context.Context, taskID, caseID int64) (models.Task, error)
	AddSolutionToTask(ctx context.Context, taskID int64, solution string) (models.Task, error)
	RemoveCaseFromTask(ctx context.Context, taskID int64) (models.Task, error)
	RemoveSolutionFromTask(ctx context.Context, taskID int64) (models.Task, error)
	AppointUserToTask(ctx context.Context, taskID int64) (models.Task, error)
	FireTask(ctx context.Context, taskID int64) (models.Task, error)
	ListTasksByUserID(ctx context.Context, userID int64, status models.TaskStatus) ([]models.Task, error)
	ListUsers(ctx context.Context, empty *empty.Empty) ([]models.User, error)
}

type serverAPI struct {
	tasksv1.UnimplementedTaskServiceServer
	taskService TaskService
}

func Register(gRPC *grpc.Server, taskService TaskService) {
	tasksv1.RegisterTaskServiceServer(gRPC, &serverAPI{taskService: taskService})
}

func (s *serverAPI) CreateTask(ctx context.Context, req *tasksv1.CreateTaskRequest) (*tasksv1.Task, error) {
	task, err := s.taskService.CreateTask(ctx, req.GetTitle(), req.GetDescription(), req.GetClusterIndex(), req.GetClusterName(), req.GetFrequency(), req.GetAverageDuration())
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

func (s *serverAPI) AddCaseToTask(ctx context.Context, req *tasksv1.AddCaseToTaskRequest) (*tasksv1.Task, error) {
	task, err := s.taskService.AddCaseToTask(ctx, req.GetTaskId(), req.GetCaseId())
	if err != nil {
		if errors.Is(err, tasks.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return ConvertTaskToProto(task), nil
}

func (s *serverAPI) AddSolutionToTask(ctx context.Context, req *tasksv1.AddSolutionToTaskRequest) (*tasksv1.Task, error) {
	task, err := s.taskService.AddSolutionToTask(ctx, req.GetTaskId(), req.GetSolution())
	if err != nil {
		if errors.Is(err, tasks.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return ConvertTaskToProto(task), nil
}

func (s *serverAPI) RemoveCaseFromTask(ctx context.Context, req *tasksv1.RemoveCaseFromTaskRequest) (*tasksv1.Task, error) {
	task, err := s.taskService.RemoveCaseFromTask(ctx, req.GetTaskId())
	if err != nil {
		if errors.Is(err, tasks.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return ConvertTaskToProto(task), nil
}

func (s *serverAPI) RemoveSolutionFromTask(ctx context.Context, req *tasksv1.RemoveSolutionFromTaskRequest) (*tasksv1.Task, error) {
	task, err := s.taskService.RemoveSolutionFromTask(ctx, req.GetTaskId())
	if err != nil {
		if errors.Is(err, tasks.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return ConvertTaskToProto(task), nil
}

func (s *serverAPI) AppointUserToTask(ctx context.Context, req *tasksv1.AppointUserToTaskRequest) (*tasksv1.Task, error) {
	task, err := s.taskService.AppointUserToTask(ctx, req.GetTaskId())
	if err != nil {
		if errors.Is(err, tasks.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return ConvertTaskToProto(task), nil
}

func (s *serverAPI) FireTask(ctx context.Context, req *tasksv1.FireTaskRequest) (*tasksv1.Task, error) {
	task, err := s.taskService.FireTask(ctx, req.GetTaskId())
	if err != nil {
		if errors.Is(err, tasks.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return ConvertTaskToProto(task), nil
}

func (s *serverAPI) ListTasksByUserID(ctx context.Context, req *tasksv1.ListTasksByUserIDRequest) (*tasksv1.ListTasksResponse, error) {
	tasks, err := s.taskService.ListTasksByUserID(ctx, req.GetUserId(), models.TaskStatus(req.GetStatus()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &tasksv1.ListTasksResponse{Tasks: ConvertTaskListToProto(tasks)}, nil
}

func (s *serverAPI) ListUsers(ctx context.Context, req *empty.Empty) (*tasksv1.ListUsersResponse, error) {
	users, err := s.taskService.ListUsers(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &tasksv1.ListUsersResponse{Users: ConvertUserListToProto(users)}, nil
}
