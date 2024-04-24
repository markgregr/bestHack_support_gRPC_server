package models

import "time"

type Task struct {
	ID          int64      `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"not null" json:"title`
	Description string     `gorm:"not null" json:"description`
	Status      TaskStatus `gorm:"not null" json:"status`
	CreatedAt   time.Time  `gorm:"autoCreateTime" gorm:"not null" json:"created_at`
	FormedAt    *time.Time `json:"formed_at`
	CompletedAt *time.Time `json:"completed_at`

	CaseID *int64 `json:"case_id`
	Case   *Case  `gorm:"foreignKey:CaseID" json:"case`

	ClusterID *int64   `json:"cluster_id`
	Cluster   *Cluster `gorm:"foreignKey:ClusterID" json:"cluster`

	UserID *int64 `json:"user_id`
	User   *User  `gorm:"foreignKey:UserID" json:"user`
}

type TaskStatus int32

const (
	TaskStatusOpen TaskStatus = iota
	TaskStatusInProgress
	TaskStatusClosed
)
