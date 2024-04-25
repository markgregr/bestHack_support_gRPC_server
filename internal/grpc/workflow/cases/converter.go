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
