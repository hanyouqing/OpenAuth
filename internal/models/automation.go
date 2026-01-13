package models

import (
	"time"

	"gorm.io/gorm"
)

// AutomationWorkflow represents an automation workflow
type AutomationWorkflow struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description,omitempty"`
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	Trigger     string         `gorm:"type:jsonb;not null" json:"trigger"` // JSON: {type, conditions}
	Actions     string         `gorm:"type:jsonb;not null" json:"actions"` // JSON array of actions
	Priority    int            `gorm:"default:0" json:"priority"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Executions []AutomationExecution `gorm:"foreignKey:WorkflowID" json:"-"`
}

// AutomationExecution represents an execution of a workflow
type AutomationExecution struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	WorkflowID uint64    `gorm:"not null;index" json:"workflow_id"`
	Status     string    `gorm:"default:pending" json:"status"` // pending, running, completed, failed
	Input      string    `gorm:"type:jsonb" json:"input"`       // JSON input data
	Output     string    `gorm:"type:jsonb" json:"output"`       // JSON output data
	Error      string    `json:"error,omitempty"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`

	Workflow AutomationWorkflow `gorm:"foreignKey:WorkflowID" json:"-"`
}

// AutomationTrigger represents a workflow trigger
type AutomationTrigger struct {
	Type       string                 `json:"type"` // event, schedule, webhook, manual
	Event      string                 `json:"event,omitempty"` // user.created, user.updated, etc.
	Schedule   string                 `json:"schedule,omitempty"` // cron expression
	Conditions map[string]interface{} `json:"conditions,omitempty"` // trigger conditions
}

// AutomationAction represents a workflow action
type AutomationAction struct {
	Type    string                 `json:"type"` // send_email, create_user, assign_role, etc.
	Config  map[string]interface{} `json:"config"` // action configuration
	OnError string                 `json:"on_error,omitempty"` // continue, stop, retry
}
