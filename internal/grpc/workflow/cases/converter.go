package cases

import (
	"github.com/markgregr/bestHack_support_gRPC_server/internal/domain/models"
	casesv1 "github.com/markgregr/bestHack_support_protos/gen/go/workflow/cases"
	"time"
)

func ConvertCaseToProto(caseItem models.Case) *casesv1.Case {
	return &casesv1.Case{
		Id:       caseItem.ID,
		Title:    caseItem.Title,
		Solution: caseItem.Solution,
		Cluster:  &casesv1.Cluster{Id: caseItem.Cluster.ID, Name: caseItem.Cluster.Name, Frequency: caseItem.Cluster.Frequency},
	}
}

func ConvertCaseListToProto(cases []models.Case) []*casesv1.Case {
	protoCases := make([]*casesv1.Case, 0, len(cases))
	for _, caseItem := range cases {
		protoCases = append(protoCases, ConvertCaseToProto(caseItem))
	}
	return protoCases
}

func ConvertTaskToProto(task models.Task) *casesv1.Task {
	var formedAt string
	if task.FormedAt != nil {
		formedAt = task.FormedAt.Format(time.RFC3339)
	}

	var completedAt string
	if task.CompletedAt != nil {
		completedAt = task.CompletedAt.Format(time.RFC3339)
	}

	var caseID int64
	var caseTitle, caseSolution string
	if task.Case != nil {
		caseID = task.Case.ID
		caseTitle = task.Case.Title
		caseSolution = task.Case.Solution
	}

	var userID int64
	var userEmail string
	if task.User != nil {
		userID = task.User.ID
		userEmail = task.User.Email
	}

	return &casesv1.Task{
		Id:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      casesv1.TaskStatus(task.Status),
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		FormedAt:    &formedAt,
		CompletedAt: &completedAt,
		User: &casesv1.User{
			Id:    userID,
			Email: userEmail,
		},
		Case: &casesv1.Case{
			Id:       caseID,
			Title:    caseTitle,
			Solution: caseSolution,
		},
	}
}

func ConvertTaskListToProto(tasks []models.Task) []*casesv1.Task {
	protoTasks := make([]*casesv1.Task, 0, len(tasks))
	for _, task := range tasks {
		protoTasks = append(protoTasks, ConvertTaskToProto(task))
	}
	return protoTasks
}

func ConvertClusterToProto(cluster models.Cluster) *casesv1.Cluster {
	return &casesv1.Cluster{
		Id:        cluster.ID,
		Name:      cluster.Name,
		Frequency: cluster.Frequency,
		Cases:     ConvertCaseListToProto(cluster.Cases),
		Tasks:     ConvertTaskListToProto(cluster.Tasks),
	}
}
