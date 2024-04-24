package cases

import (
	"github.com/markgregr/bestHack_support_gRPC_server/internal/domain/models"
	casesv1 "github.com/markgregr/bestHack_support_protos/gen/go/workflow/cases"
)

func ConvertCaseToProto(caseItem models.Case) *casesv1.Case {
	var clusterID int64
	var clusterName string
	var clusterFrequency int64
	if caseItem.Cluster != nil {
		clusterID = caseItem.Cluster.ID
		clusterName = caseItem.Cluster.Name
		clusterFrequency = caseItem.Cluster.Frequency
	}

	return &casesv1.Case{
		Id:       caseItem.ID,
		Title:    caseItem.Title,
		Solution: caseItem.Solution,
		Cluster: &casesv1.Cluster{Id: clusterID,
			Name:      clusterName,
			Frequency: clusterFrequency,
		},
	}
}

func ConvertCaseListToProto(cases []models.Case) []*casesv1.Case {
	protoCases := make([]*casesv1.Case, 0, len(cases))
	for _, caseItem := range cases {
		protoCases = append(protoCases, ConvertCaseToProto(caseItem))
	}
	return protoCases
}

//func ConvertTaskToProto(task models.Task) *casesv1.Task {
//	log.WithField("task", task).Info("convert task to proto")
//
//	var formedAt, completedAt *string
//	if task.FormedAt != nil {
//		formattedFormedAt := task.FormedAt.Format(time.RFC3339)
//		formedAt = &formattedFormedAt
//	}
//	if task.CompletedAt != nil {
//		formattedCompletedAt := task.CompletedAt.Format(time.RFC3339)
//		completedAt = &formattedCompletedAt
//	}
//
//	var caseID int64
//	var caseTitle, caseSolution string
//	if task.Case != nil {
//		caseID = task.Case.ID
//		caseTitle = task.Case.Title
//		caseSolution = task.Case.Solution
//	}
//
//	var userID int64
//	var userEmail string
//	if task.User != nil {
//		userID = task.User.ID
//		userEmail = task.User.Email
//	}
//
//	return &casesv1.Task{
//		Id:          task.ID,
//		Title:       task.Title,
//		Description: task.Description,
//		Status:      casesv1.TaskStatus(task.Status),
//		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
//		FormedAt:    formedAt,
//		CompletedAt: completedAt,
//		User: &casesv1.User{
//			Id:    userID,
//			Email: userEmail,
//		},
//		Case: &casesv1.Case{
//			Id:       caseID,
//			Title:    caseTitle,
//			Solution: caseSolution,
//		},
//	}
//}
//
//func ConvertTaskListToProto(tasks []models.Task) []*casesv1.Task {
//	protoTasks := make([]*casesv1.Task, 0, len(tasks))
//	for _, task := range tasks {
//		protoTasks = append(protoTasks, ConvertTaskToProto(task))
//	}
//	return protoTasks
//}

func ConvertClusterToProto(cluster models.Cluster) *casesv1.Cluster {
	return &casesv1.Cluster{
		Id:        cluster.ID,
		Name:      cluster.Name,
		Frequency: cluster.Frequency,
	}
}

func ConvertClusterListToProto(clusters []models.Cluster) []*casesv1.Cluster {
	protoClusters := make([]*casesv1.Cluster, 0, len(clusters))
	for _, cluster := range clusters {
		protoClusters = append(protoClusters, ConvertClusterToProto(cluster))
	}
	return protoClusters
}
