package tasks

import (
	"github.com/markgregr/bestHack_support_gRPC_server/internal/domain/models"
	tasksv1 "github.com/markgregr/bestHack_support_protos/gen/go/workflow/tasks"
	log "github.com/sirupsen/logrus"
	"time"
)

func ConvertTaskToProto(task models.Task) *tasksv1.Task {
	log.WithField("task", task).Info("convert task to proto")

	var formedAt, completedAt *string
	if task.FormedAt != nil {
		formattedFormedAt := task.FormedAt.Format(time.RFC3339)
		formedAt = &formattedFormedAt
	}
	if task.CompletedAt != nil {
		formattedCompletedAt := task.CompletedAt.Format(time.RFC3339)
		completedAt = &formattedCompletedAt
	}

	var solution *string
	if task.Solution != nil {
		sol := *task.Solution
		solution = &sol
	}

	var caseID, caseClusterID int64
	var caseTitle, caseSolution string
	if task.Case != nil {
		caseID = task.Case.ID
		if task.Case.Cluster != nil {
			caseClusterID = task.Case.Cluster.ID
		}
		caseTitle = task.Case.Title
		caseSolution = task.Case.Solution
	}

	var clusterID int64
	var clusterName string
	var clusterFrequency int64
	if task.Cluster != nil {
		clusterID = task.Cluster.ID
		clusterName = task.Cluster.Name
		clusterFrequency = task.Cluster.Frequency
	}

	var userID int64
	var userEmail string
	var avarageDuration float32
	if task.User != nil {
		userID = task.User.ID
		userEmail = task.User.Email
		avarageDuration = task.User.AvarageDuration
	}

	return &tasksv1.Task{
		Id:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Solution:    solution,
		Fire:        task.Fire,
		Status:      tasksv1.TaskStatus(task.Status),
		Case: &tasksv1.Case{Id: caseID,
			ClusterId: caseClusterID,
			Title:     caseTitle,
			Solution:  caseSolution,
		},
		Cluster: &tasksv1.Cluster{Id: clusterID,
			Name:      clusterName,
			Frequency: clusterFrequency,
		},
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		FormedAt:    formedAt,
		CompletedAt: completedAt,
		User: &tasksv1.User{
			Id:              userID,
			Email:           userEmail,
			AvarageDuration: avarageDuration,
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

func ConvertUserToProto(user models.User) *tasksv1.User {
	return &tasksv1.User{
		Id:              user.ID,
		Email:           user.Email,
		AvarageDuration: user.AvarageDuration,
	}
}

func ConvertUserListToProto(users []models.User) []*tasksv1.User {
	protoUsers := make([]*tasksv1.User, 0, len(users))
	for _, user := range users {
		protoUsers = append(protoUsers, ConvertUserToProto(user))
	}
	return protoUsers
}
