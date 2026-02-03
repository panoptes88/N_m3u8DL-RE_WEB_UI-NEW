package model

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	URL               string         `gorm:"size:1024" json:"url"`
	Status            string         `gorm:"size:32;default:'pending'" json:"status"`
	Progress          int            `gorm:"default:0" json:"progress"`
	Speed             string         `gorm:"size:32" json:"speed"`
	OutputName        string         `gorm:"size:512" json:"output_name"`
	ThreadCount       int            `gorm:"default:32" json:"thread_count"`
	RetryCount        int            `gorm:"default:15" json:"retry_count"`
	Headers           string         `gorm:"size:2048" json:"headers,omitempty"`
	BaseURL           string         `gorm:"size:1024" json:"base_url,omitempty"`
	DelAfterDone      bool           `gorm:"default:true" json:"del_after_done"`
	BinaryMerge       bool           `gorm:"default:false" json:"binary_merge"`
	Key               string         `gorm:"size:512" json:"key,omitempty"`
	DecryptionEngine  string         `gorm:"size:32;default:'MP4DECRYPT'" json:"decryption_engine"`
	CustomArgs        string         `gorm:"size:2048" json:"custom_args,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	FinishedAt        *time.Time     `json:"finished_at,omitempty"`
	ErrorMsg          string         `gorm:"size:2048" json:"error_msg,omitempty"`
	PID               int            `gorm:"default:0" json:"pid,omitempty"`
	LogFile           string         `gorm:"size:512" json:"log_file,omitempty"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

const (
	TaskStatusPending     = "pending"
	TaskStatusDownloading = "downloading"
	TaskStatusCompleted   = "completed"
	TaskStatusFailed      = "failed"
)

func (Task) TableName() string {
	return "tasks"
}
