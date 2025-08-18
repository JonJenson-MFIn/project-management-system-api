package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ResetDatabase manually resets the database - use this only when you need to clear all data
func ResetDatabase() error {
	if DB == nil {
		return nil
	}

	log.Println("Manually resetting database...")

	// Drop all tables in the correct order to avoid foreign key constraint issues
	tables := []string{
		"project_employees",
		"project_teams",
		"team_engineers",
		"notifications",
		"tasks",
		"tickets",
		"projects",
		"teams",
		"employees",
	}

	for _, table := range tables {
		if err := DB.Exec("DROP TABLE IF EXISTS " + table + " CASCADE").Error; err != nil {
			log.Printf("Warning: Could not drop table %s: %v", table, err)
		}
	}

	log.Println("Database reset completed")
	return nil
}

func Migrate(db *gorm.DB) error {
	// Create extension if it doesn't exist
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`).Error; err != nil {
		log.Printf("Warning: Could not create pgcrypto extension: %v", err)
	}

	// Check if tables already exist to avoid unnecessary migration attempts
	var tableCount int64
	db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = 'employees'").Scan(&tableCount)

	if tableCount == 0 {
		// Only run AutoMigrate if tables don't exist
		models := []interface{}{
			&Employee{},
			&Team{},
			&Project{},
			&Ticket{},
			&Task{},
			&Notification{},
			&TeamEngineer{},
			&ProjectTeam{},
			&ProjectEmployee{},
		}

		if err := db.AutoMigrate(models...); err != nil {
			log.Printf("Warning: AutoMigrate had some issues: %v", err)
		}
		log.Println("Database tables created successfully")
	} else {
		log.Println("Database tables already exist, skipping table creation")
	}

	// Add foreign key constraints after tables are created
	if err := addForeignKeys(db); err != nil {
		log.Printf("Warning: Could not add all foreign key constraints: %v", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

func addForeignKeys(db *gorm.DB) error {
	// Add foreign key constraints for Employee -> Project
	if err := db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_name = 'fk_employee_project' 
				AND table_name = 'employees'
			) THEN
				ALTER TABLE employees 
				ADD CONSTRAINT fk_employee_project 
				FOREIGN KEY (project_assigned_id) 
				REFERENCES projects(id) ON DELETE SET NULL ON UPDATE CASCADE;
			END IF;
		END $$;
	`).Error; err != nil {
		log.Printf("Warning: Could not add employee-project foreign key: %v", err)
	}

	// Add foreign key constraints for Team -> Employee (TeamLeader)
	if err := db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_name = 'fk_team_leader' 
				AND table_name = 'teams'
			) THEN
				ALTER TABLE teams 
				ADD CONSTRAINT fk_team_leader 
				FOREIGN KEY (team_leader_id) 
				REFERENCES employees(id) ON DELETE SET NULL ON UPDATE CASCADE;
			END IF;
		END $$;
	`).Error; err != nil {
		log.Printf("Warning: Could not add team-leader foreign key: %v", err)
	}

	// Add foreign key constraints for Project -> Employee (Manager)
	if err := db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_name = 'fk_project_manager' 
				AND table_name = 'projects'
			) THEN
				ALTER TABLE projects 
				ADD CONSTRAINT fk_project_manager 
				FOREIGN KEY (manager_id) 
				REFERENCES employees(id) ON DELETE SET NULL ON UPDATE CASCADE;
			END IF;
		END $$;
	`).Error; err != nil {
		log.Printf("Warning: Could not add project-manager foreign key: %v", err)
	}

	// Add foreign key constraints for Ticket -> Project
	if err := db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_name = 'fk_ticket_project' 
				AND table_name = 'tickets'
			) THEN
				ALTER TABLE tickets 
				ADD CONSTRAINT fk_ticket_project 
				FOREIGN KEY (project_id) 
				REFERENCES projects(id) ON DELETE CASCADE ON UPDATE CASCADE;
			END IF;
		END $$;
	`).Error; err != nil {
		log.Printf("Warning: Could not add ticket-project foreign key: %v", err)
	}

	// Add foreign key constraints for Ticket -> Employee (AssignedTo)
	if err := db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_name = 'fk_ticket_assigned_to' 
				AND table_name = 'tickets'
			) THEN
				ALTER TABLE tickets 
				ADD CONSTRAINT fk_ticket_assigned_to 
				FOREIGN KEY (assigned_to_id) 
				REFERENCES employees(id) ON DELETE SET NULL ON UPDATE CASCADE;
			END IF;
		END $$;
	`).Error; err != nil {
		log.Printf("Warning: Could not add ticket-assigned_to foreign key: %v", err)
	}

	// Add foreign key constraints for Task -> Employee (AssignedTo)
	if err := db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_name = 'fk_task_assigned_to' 
				AND table_name = 'tasks'
			) THEN
				ALTER TABLE tasks 
				ADD CONSTRAINT fk_task_assigned_to 
				FOREIGN KEY (assigned_to_id) 
				REFERENCES employees(id) ON DELETE SET NULL ON UPDATE CASCADE;
			END IF;
		END $$;
	`).Error; err != nil {
		log.Printf("Warning: Could not add task-assigned_to foreign key: %v", err)
	}

	// Add foreign key constraints for Task -> Project
	if err := db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_name = 'fk_task_project' 
				AND table_name = 'tasks'
			) THEN
				ALTER TABLE tasks 
				ADD CONSTRAINT fk_task_project 
				FOREIGN KEY (project_id) 
				REFERENCES projects(id) ON DELETE SET NULL ON UPDATE CASCADE;
			END IF;
		END $$;
	`).Error; err != nil {
		log.Printf("Warning: Could not add task-project foreign key: %v", err)
	}

	// Add foreign key constraints for Notification -> Employee
	if err := db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_name = 'fk_notification_employee' 
				AND table_name = 'notifications'
			) THEN
				ALTER TABLE notifications 
				ADD CONSTRAINT fk_notification_employee 
				FOREIGN KEY (employee_id) 
				REFERENCES employees(id) ON DELETE CASCADE ON UPDATE CASCADE;
			END IF;
		END $$;
	`).Error; err != nil {
		log.Printf("Warning: Could not add notification-employee foreign key: %v", err)
	}

	// Add foreign key constraints for junction tables
	if err := db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_name = 'fk_team_engineer_team' 
				AND table_name = 'team_engineers'
			) THEN
				ALTER TABLE team_engineers 
				ADD CONSTRAINT fk_team_engineer_team 
				FOREIGN KEY (team_id) 
				REFERENCES teams(id) ON DELETE CASCADE ON UPDATE CASCADE;
			END IF;
		END $$;
	`).Error; err != nil {
		log.Printf("Warning: Could not add team_engineer-team foreign key: %v", err)
	}

	if err := db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_name = 'fk_team_engineer_engineer' 
				AND table_name = 'team_engineers'
			) THEN
				ALTER TABLE team_engineers 
				ADD CONSTRAINT fk_team_engineer_engineer 
				FOREIGN KEY (engineer_id) 
				REFERENCES employees(id) ON DELETE CASCADE ON UPDATE CASCADE;
			END IF;
		END $$;
	`).Error; err != nil {
		log.Printf("Warning: Could not add team_engineer-engineer foreign key: %v", err)
	}

	if err := db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_name = 'fk_project_team_project' 
				AND table_name = 'project_teams'
			) THEN
				ALTER TABLE project_teams 
				ADD CONSTRAINT fk_project_team_project 
				FOREIGN KEY (project_id) 
				REFERENCES projects(id) ON DELETE CASCADE ON UPDATE CASCADE;
			END IF;
		END $$;
	`).Error; err != nil {
		log.Printf("Warning: Could not add project_team-project foreign key: %v", err)
	}

	if err := db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_name = 'fk_project_team_team' 
				AND table_name = 'project_teams'
			) THEN
				ALTER TABLE project_teams 
				ADD CONSTRAINT fk_project_team_team 
				FOREIGN KEY (team_id) 
				REFERENCES teams(id) ON DELETE CASCADE ON UPDATE CASCADE;
			END IF;
		END $$;
	`).Error; err != nil {
		log.Printf("Warning: Could not add project_team-team foreign key: %v", err)
	}

	if err := db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_name = 'fk_project_employee_project' 
				AND table_name = 'project_employees'
			) THEN
				ALTER TABLE project_employees 
				ADD CONSTRAINT fk_project_employee_project 
				FOREIGN KEY (project_id) 
				REFERENCES projects(id) ON DELETE CASCADE ON UPDATE CASCADE;
			END IF;
		END $$;
	`).Error; err != nil {
		log.Printf("Warning: Could not add project_employee-project foreign key: %v", err)
	}

	if err := db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.table_constraints 
				WHERE constraint_name = 'fk_project_employee_employee' 
				AND table_name = 'project_employees'
			) THEN
				ALTER TABLE project_employees 
				ADD CONSTRAINT fk_project_employee_employee 
				FOREIGN KEY (employee_id) 
				REFERENCES employees(id) ON DELETE CASCADE ON UPDATE CASCADE;
			END IF;
		END $$;
	`).Error; err != nil {
		log.Printf("Warning: Could not add project_employee-employee foreign key: %v", err)
	}

	return nil
}

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=aws-0-ap-south-1.pooler.supabase.com user=postgres.lwybefbgqqmvzdzkvqnt password=Joans88@joejon dbname=postgres port=6543 sslmode=disable TimeZone=Asia/Kolkata"
	}

	// Configure GORM with better connection settings
	config := &gorm.Config{
		PrepareStmt: false, // Disable prepared statements to avoid conflicts
	}

	database, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	// Get the underlying SQL DB to configure connection pool
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("failed to get underlying SQL DB:", err)
	}

	// Configure connection pool to avoid prepared statement conflicts
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(0)

	if err := Migrate(database); err != nil {
		log.Printf("Warning: Migration had issues but continuing: %v", err)
	}

	DB = database
	log.Println("Database connected successfully")
}
