package tasks

import (
	"github.com/markgregr/bestHack_support_gRPC_server/internal/domain/models"
	tasksv1 "github.com/markgregr/bestHack_support_protos/gen/go/workflow/tasks"
	"time"
)

func ConvertTaskToProto(task models.Task) *tasksv1.Task {
	createdAt := task.CreatedAt.Format(time.RFC3339)
	formedAt := ""
	if task.FormedAt != nil {
		formedAt = task.FormedAt.Format(time.RFC3339)
	}
	completedAt := ""
	if task.CompletedAt != nil {
		completedAt = task.CompletedAt.Format(time.RFC3339)
	}
	var caseID, clusterID int64
	if task.CaseID != nil {
		caseID = *task.CaseID
	}
	if task.ClusterID != nil {
		clusterID = *task.ClusterID
	}
	return &tasksv1.Task{
		Id:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      tasksv1.TaskStatus(task.Status),
		CaseId:      caseID,
		ClusterId:   clusterID,
		CreatedAt:   createdAt,
		FormedAt:    formedAt,
		CompletedAt: completedAt,
	}
}

func ConvertTaskListToProto(tasks []models.Task) []*tasksv1.Task {
	protoTasks := make([]*tasksv1.Task, 0, len(tasks))
	for _, task := range tasks {
		protoTasks = append(protoTasks, ConvertTaskToProto(task))
	}
	return protoTasks
}
