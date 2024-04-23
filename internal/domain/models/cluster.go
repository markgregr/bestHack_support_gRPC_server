package models

type Cluster struct {
	ID             int64  `gorm:"primaryKey" json:"id"`
	Name           string `json:"name"`
	TreatmentCount int64  `json:"treatment_count"`
	Tasks          []Task `json:"tasks"`
	Cases          []Case `json:"cases"`
}
