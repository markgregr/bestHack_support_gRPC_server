package models

type Cluster struct {
	ID           int64  `gorm:"primaryKey" json:"id"`
	ClusterIndex int64  `gorm:"unique" json:"cluster_index"`
	Name         string `json:"name"`
	Frequency    int64  `json:"frequency"`
	Tasks        []Task `json:"tasks"`
	Cases        []Case `json:"cases"`
}
