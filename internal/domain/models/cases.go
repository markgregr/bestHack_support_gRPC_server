package models

type Case struct {
	ID       int64  `gorm:"primaryKey" json:"id"`
	Title    string `json:"title"`
	Solution string `json:"solution"`

	ClusterID *int64   `json:"cluster_id"`
	Cluster   *Cluster `gorm:"foreignKey:ClusterID" json:"cluster`
}
