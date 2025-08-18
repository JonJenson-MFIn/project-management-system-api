package db

import (
	"time"
	"gorm.io/gorm"
)

type RoleDB string
type StatusDB string

const (
	RoleAdminDB    RoleDB = "ADMIN"
	RoleEmployeeDB RoleDB = "EMPLOYEE"
	RoleTlDB       RoleDB = "TL"
	RoleManagerDB  RoleDB = "MANAGER"

	StatusNotStartedDB StatusDB = "NOT_STARTED"
	StatusInProgressDB StatusDB = "IN_PROGRESS"
	StatusCompletedDB  StatusDB = "COMPLETED"
	StatusOnHoldDB     StatusDB = "ON_HOLD"
	StatusCancelledDB  StatusDB = "CANCELLED"
)

type Employee struct {
	ID                int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name              string         `gorm:"not null" json:"name"`
	Email             string         `gorm:"uniqueIndex;not null" json:"email"`
	Password          string         `json:"-"`
	Role              RoleDB         `gorm:"type:varchar(20);not null" json:"role"`
	Active            bool           `gorm:"not null;default:true" json:"active"`
	ProjectAssignedID *int           `gorm:"index" json:"project_assigned_id,omitempty"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

type Team struct {
	ID           int            `gorm:"primaryKey;autoIncrement" json:"id"`
	TeamLeaderID *int           `gorm:"index" json:"team_leader_id,omitempty"`
	Name         string         `gorm:"not null" json:"name"`
	Description  *string        `json:"description,omitempty"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Project struct {
	ID          int            `gorm:"primaryKey;autoIncrement" json:"id"`
	ManagerID   *int           `gorm:"index" json:"manager_id,omitempty"`
	Name        string         `gorm:"not null" json:"name"`
	Status      StatusDB       `gorm:"type:varchar(20);not null" json:"status"`
	Description *string        `json:"description,omitempty"`
	StartDate   time.Time      `json:"start_date"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Ticket struct {
	ID           int            `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID    int            `gorm:"not null;index" json:"project_id"`
	AssignedToID *int           `gorm:"index" json:"assigned_to_id,omitempty"`
	Status       StatusDB       `gorm:"type:varchar(20);not null" json:"status"`
	Title        string         `gorm:"not null" json:"title"`
	Description  *string        `json:"description,omitempty"`
	Priority     string         `gorm:"type:varchar(20);default:'MEDIUM'" json:"priority"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	CompletedAt  *time.Time     `json:"completedAt,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Task struct {
	ID           int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Title        string         `gorm:"not null" json:"title"`
	Description  *string        `json:"description,omitempty"`
	AssignedToID *int           `gorm:"index" json:"assigned_to_id,omitempty"`
	ProjectID    *int           `gorm:"index" json:"project_id,omitempty"`
	DueDate      *time.Time     `json:"due_date,omitempty"`
	Status       StatusDB       `gorm:"type:varchar(20);not null" json:"status"`
	Priority     string         `gorm:"type:varchar(20);default:'MEDIUM'" json:"priority"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	CompletedAt  *time.Time     `json:"completedAt,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Notification struct {
	ID         int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Message    string         `gorm:"not null" json:"message"`
	EmployeeID int            `gorm:"index" json:"employee_id"`
	Type       string         `gorm:"type:varchar(50);default:'INFO'" json:"type"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	Read       bool           `gorm:"not null;default:false" json:"read"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// Junction tables for many-to-many relationships
type TeamEngineer struct {
	TeamID     int       `gorm:"primaryKey" json:"team_id"`
	EngineerID int       `gorm:"primaryKey" json:"engineer_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

type ProjectTeam struct {
	ProjectID int       `gorm:"primaryKey" json:"project_id"`
	TeamID    int       `gorm:"primaryKey" json:"team_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

type ProjectEmployee struct {
	ProjectID  int       `gorm:"primaryKey" json:"project_id"`
	EmployeeID int       `gorm:"primaryKey" json:"employee_id"`
	Role       string    `gorm:"type:varchar(50);default:'MEMBER'" json:"role"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
}
