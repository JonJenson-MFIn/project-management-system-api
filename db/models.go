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
	ID                string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name              string         `gorm:"not null" json:"name"`
	Email             string         `gorm:"uniqueIndex;not null" json:"email"`
	Password          string         `json:"-"`
	Role              RoleDB         `gorm:"type:varchar(20);not null" json:"role"`
	Active            bool           `gorm:"not null;default:true" json:"active"`
	ProjectAssignedID *string        `gorm:"type:uuid;index" json:"project_assigned_id,omitempty"`
	ProjectAssigned   *Project       `gorm:"foreignKey:ProjectAssignedID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"project_assigned,omitempty"`
	Teams             []*Team        `gorm:"many2many:team_engineers;" json:"teams,omitempty"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

type Project struct {
	ID          string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ManagerID   *string        `gorm:"type:uuid;index" json:"manager_id,omitempty"`
	Manager     *Employee      `gorm:"foreignKey:ManagerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"manager,omitempty"`
	Name        string         `gorm:"not null" json:"name"`
	Status      StatusDB       `gorm:"type:varchar(20);not null" json:"status"`
	Description *string        `json:"description,omitempty"`
	StartDate   time.Time      `json:"start_date"`
	Teams       []*Team        `gorm:"many2many:project_teams;" json:"teams,omitempty"`
	Tickets     []*Ticket      `gorm:"foreignKey:ProjectID" json:"tickets,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Team struct {
	ID           string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TeamLeaderID *string        `gorm:"type:uuid;index" json:"team_leader_id,omitempty"`
	TeamLeader   *Employee      `gorm:"foreignKey:TeamLeaderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"team_leader,omitempty"`
	Engineers    []*Employee    `gorm:"many2many:team_engineers;" json:"engineers,omitempty"`
	Projects     []*Project     `gorm:"many2many:project_teams;" json:"projects,omitempty"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Ticket struct {
	ID          string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProjectID   string         `gorm:"type:uuid;not null;index" json:"project_id"`
	Project     *Project       `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"project,omitempty"`
	Status      StatusDB       `gorm:"type:varchar(20);not null" json:"status"`
	Title       string         `gorm:"not null" json:"title"`
	Description *string        `json:"description,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	CompletedAt *time.Time     `json:"completedAt,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Task struct {
	ID           string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title        string         `gorm:"not null" json:"title"`
	Description  *string        `json:"description,omitempty"`
	AssignedToID *string        `gorm:"type:uuid;index" json:"assigned_to_id,omitempty"`
	AssignedTo   *Employee      `gorm:"foreignKey:AssignedToID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"assigned_to,omitempty"`
	DueDate      *time.Time     `json:"due_date,omitempty"`
	Status       StatusDB       `gorm:"type:varchar(20);not null" json:"status"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	CompletedAt  *time.Time     `json:"completedAt,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Notification struct {
	ID         string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Message    string         `gorm:"not null" json:"message"`
	EmployeeID string         `gorm:"type:uuid;index" json:"employee_id"`
	Employee   *Employee      `gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"employee,omitempty"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	Read       bool           `gorm:"not null;default:false" json:"read"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func Migrate(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`).Error; err != nil {
		return err
	}
	return db.AutoMigrate(
		&Employee{},
		&Project{},
		&Team{},
		&Ticket{},
		&Task{},
		&Notification{},
	)
}
