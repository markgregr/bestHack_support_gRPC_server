package tasks

import (
	"github.com/markgregr/bestHack_support_gRPC_server/internal/domain/models"
	tasksv1 "github.com/markgregr/bestHack_support_protos/gen/go/workflow/tasks"
	"time"
)

func ConvertTaskToProto(task models.Task) *tasksv1.Task {
	//createdAt := task.CreatedAt.Format(time.RFC3339)

	var formedAt string
	if task.FormedAt != nil {
		formedAt = task.FormedAt.Format(time.RFC3339)
	}

	var completedAt string
	if task.CompletedAt != nil {
		completedAt = task.CompletedAt.Format(time.RFC3339)
	}

	//var caseID, clusterID int64
	//var caseTitle, caseSolution string
	//if task.Case != nil {
	//	caseID = task.Case.ID
	//	if task.Case.Cluster != nil {
	//		clusterID = task.Case.Cluster.ID
	//	}
	//	caseTitle = task.Case.Title
	//	caseSolution = task.Case.Solution
	//}
	//
	//var userID int64
	//var userEmail string
	//if task.User != nil {
	//	userID = task.User.ID
	//	userEmail = task.User.Email
	//}

	return &tasksv1.Task{
		Id:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      tasksv1.TaskStatus(task.Status),
		Case: &tasksv1.Case{Id: task.Case.ID,
			ClusterId: task.Case.Cluster.ID,
			Title:     task.Case.Title,
			Solution:  task.Case.Solution,
		},
		Cluster: &tasksv1.Cluster{Id: task.Cluster.ID,
			Name:      task.Cluster.Name,
			Frequency: task.Cluster.Frequency,
		},
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		FormedAt:    &formedAt,
		CompletedAt: &completedAt,
		User: &tasksv1.User{
			Id:    task.User.ID,
			Email: task.User.Email,
		},
	}
}

func ConvertTaskListToProto(tasks []models.Task) []*tasksv1.Task {
	protoTasks := make([]*tasksv1.Task, 0, len(tasks))
	for _, task := range tasks {
		protoTasks = append(protoTasks, ConvertTaskToProto(task))
	}
	return protoTasks
}
