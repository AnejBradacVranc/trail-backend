package tools

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

type ApplicationSummary struct {
	ApplicationID int64     `json:"application_id"`
	JobTitle      string    `json:"job_title"`
	AppliedAt     time.Time `json:"applied_at"`
	StatusName    string    `json:"status_name"`
	SalaryMax     int       `json:"salary_max"`
	SalaryMin     int       `json:"salary_min"`
	Location      string    `json:"location"`
	CompanyName   string    `json:"company_name"`
}

type ApplicationStatus struct {
	StatusID   int64  `json:"status_id"`
	StatusName string `json:"status_name"`
}

type ApplicationDetail struct {
	ApplicationID int64     `json:"application_id"`
	JobTitle      string    `json:"job_title"`
	Platform      Platform  `json:"platform"`
	JobURL        *string   `json:"job_url"`
	SalaryMin     *int      `json:"salary_min"`
	SalaryMax     *int      `json:"salary_max"`
	CreatedAt     time.Time `json:"created_at"`
	ModifiedAt    time.Time `json:"modified_at"`

	StatusName string `json:"status_name"`

	Company CompanyDetail `json:"company"`

	Events []*ApplicationEvent `json:"events"`
	Notes  []*Note             `json:"notes"`
	Files  []*File             `json:"files"`
}

type Platform struct {
	PlatformID int64   `json:"platform_id"`
	Name       string  `json:"name"`
	Website    *string `json:"website"`
	IsActive   bool    `json:"is_active"`
}

type ApplicationEvent struct {
	EventID         int64      `json:"event_id"`
	ApplicationID   int64      `json:"application_id"`
	EventType       string     `json:"event_type"`
	Note            *string    `json:"note"`
	EventStartTime  time.Time  `json:"event_start_time"`
	EventEstEndTime *time.Time `json:"event_est_end_time"`
	CreatedAt       time.Time  `json:"created_at"`
}

type CompanyDetail struct {
	CompanyID            int64   `json:"company_id"`
	Name                 string  `json:"name"`
	HeadquartersLocation *string `json:"headquarters_location"`
	EmployeesCount       *int    `json:"employees_count"`
}

/*type CompanyContact struct {
	CompanyContactID int64     `json:"company_contact_id"`
	CompanyID        int64     `json:"company_id"`
	Name             string    `json:"name"`
	Surname          string    `json:"surname"`
	Email            *string   `json:"email"`
	Phone            *string   `json:"phone"`
	Role             *string   `json:"role"`
	CreatedAt        time.Time `json:"created_at"`
	ModifiedAt       time.Time `json:"modified_at"`
}*/

type Note struct {
	NoteID        int64     `json:"note_id"`
	ApplicationID int64     `json:"application_id"`
	NoteContent   *string   `json:"note_content"`
	CreatedAt     time.Time `json:"created_at"`
	ModifiedAt    time.Time `json:"modified_at"`
}

type File struct {
	FileID        int64     `json:"file_id"`
	ApplicationID int64     `json:"application_id"`
	Filename      string    `json:"filename"`
	CreatedAt     time.Time `json:"created_at"`
	ModifiedAt    time.Time `json:"modified_at"`
}

type User struct {
	UserID     int64     `json:"user_id"`
	Name       string    `json:"name"`
	Surname    string    `json:"surname"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

type DbInterface interface {
	SetupDatabase() error
	CloseDatabase()
	CreateApplication(ctx context.Context, tx *pgx.Tx, userId int64, statusId int64, companyId int64, jobTitle string, platformId int64, jobUrl *string, salaryMin *int, salaryMax *int, appliedAt time.Time) (int64, error)
	CreateCompany(ctx context.Context, tx *pgx.Tx, name string, location string, employeesCount *int) (int64, error)
	CreateNote(ctx context.Context, tx *pgx.Tx, applicationId int64, noteContent string) (int64, error)
	CreateFile(ctx context.Context, tx *pgx.Tx, applicationId int64, filename string) (int64, error)
	GetApplicationByID(applicationId int64) (*ApplicationDetail, error)
	GetApplicationsFromUserByEmail(email string) ([]*ApplicationSummary, error)
	GetStatuses() ([]*ApplicationStatus, error)
	GetPlatforms() ([]*Platform, error)
	LoginUser(email string, password string) (*User, error)
	CreateUser(name string, surname string, email string, password string) (int64, error)
	BeginTx(ctx context.Context) (*pgx.Tx, error)
}

func NewDatabase() (*DbInterface, error) {

	var db DbInterface = &postgreSQL{}

	err := db.SetupDatabase()

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &db, nil
}
