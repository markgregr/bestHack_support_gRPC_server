package models

type Case struct {
	ID          int64  `gorm:"primaryKey" json:"id"`
	ClusterID   int64  `json:"cluster_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Solution    string `json:"solution"`
	CreatedAt   int64  `json:"created_at"`
}
