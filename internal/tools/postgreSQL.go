package tools

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type postgreSQL struct {
	db *pgx.Conn
}

func (p *postgreSQL) GetApplicationsFromUserByEmail(email string) ([]*ApplicationSummary, error) {
	query := `SELECT 	
		a.application_id,
		a.job_title,    
		a.applied_at,
		a.salary_max,
		a.salary_min,
		c.name AS company_name,
		c.location,
		s.status_name
	FROM applications a
	JOIN users u ON u.user_id = a.user_id
	JOIN application_statuses s ON s.status_id = a.status_id
	JOIN companies c ON c.company_id = a.company_id
	WHERE u.email = $1`

	rows, err := p.db.Query(context.Background(), query, email)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var applications []*ApplicationSummary

	for rows.Next() {
		application := &ApplicationSummary{}

		err := rows.Scan(
			&application.ApplicationID,
			&application.JobTitle,
			&application.AppliedAt,
			&application.SalaryMax,
			&application.SalaryMin,
			&application.CompanyName,
			&application.Location,
			&application.StatusName)

		if err != nil {
			return nil, err
		}

		applications = append(applications, application)
	}

	return applications, nil
}

func (p *postgreSQL) CreateApplication(userId int64, statusId int64, companyId int64, jobTitle string, platform string, jobUrl *string, salaryMin *int, salaryMax *int, appliedAt time.Time) (int64, error) {

	query := `INSERT INTO applications(USER_ID, STATUS_ID, COMPANY_ID, JOB_TITLE, PLATFORM, JOB_URL, SALARY_MAX, SALARY_MIN, APPLIED_AT) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING application_id`

	var id int64

	err := p.db.QueryRow(context.Background(), query, userId, statusId, companyId, jobTitle, platform, jobUrl, salaryMax, salaryMin, appliedAt).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
}

func (p *postgreSQL) GetApplicationByID(applicationId int64) (*ApplicationDetail, error) {
	query := `
		SELECT
			a.application_id,
			a.job_title,
			a.platform,
			a.job_url,
			a.salary_min,
			a.salary_max,
			a.created_at,
			a.modified_at,
			s.status_name,
			c.company_id,
			c.name,
			c.location,
			c.employee_count
		FROM applications a
		JOIN application_statuses s ON s.status_id = a.status_id
		JOIN companies c ON c.company_id = a.company_id
		WHERE a.application_id = $1
	`

	row := p.db.QueryRow(context.Background(), query, applicationId)

	detail := &ApplicationDetail{}

	err := row.Scan(
		&detail.ApplicationID,
		&detail.JobTitle,
		&detail.Platform,
		&detail.JobURL,
		&detail.SalaryMin,
		&detail.SalaryMax,
		&detail.CreatedAt,
		&detail.ModifiedAt,
		&detail.StatusName,
		&detail.Company.CompanyID,
		&detail.Company.Name,
		&detail.Company.HeadquartersLocation,
		&detail.Company.EmployeesCount,
	)

	if err != nil {
		return nil, err
	}

	eventsQuery := `
		SELECT
			event_id,
			application_id,
			event_type,
			note,
			event_start_time,
			event_est_end_time,
			created_at
		FROM application_events
		WHERE application_id = $1
		ORDER BY event_start_time
	`

	rows, err := p.db.Query(context.Background(), eventsQuery, applicationId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		event := &ApplicationEvent{}
		err := rows.Scan(
			&event.EventID,
			&event.ApplicationID,
			&event.EventType,
			&event.Note,
			&event.EventStartTime,
			&event.EventEstEndTime,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		detail.Events = append(detail.Events, event)
	}

	notesQuery := `
		SELECT
			note_id,
			application_id,
			note_content,
			created_at,
			modified_at
		FROM notes
		WHERE application_id = $1
		ORDER BY created_at DESC
	`

	rows, err = p.db.Query(context.Background(), notesQuery, applicationId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		note := &Note{}
		err := rows.Scan(
			&note.NoteID,
			&note.ApplicationID,
			&note.NoteContent,
			&note.CreatedAt,
			&note.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}
		detail.Notes = append(detail.Notes, note)
	}

	filesQuery := `
		SELECT
			file_id,
			application_id,
			filename,
			created_at,
			modified_at
		FROM files
		WHERE application_id = $1
		ORDER BY created_at
	`

	rows, err = p.db.Query(context.Background(), filesQuery, applicationId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		file := &File{}
		err := rows.Scan(
			&file.FileID,
			&file.ApplicationID,
			&file.Filename,
			&file.CreatedAt,
			&file.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}
		detail.Files = append(detail.Files, file)
	}

	return detail, nil
}

func (p *postgreSQL) LoginUser(email string, password string) (*User, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postgreSQL) CreateUser(name string, surname string, email string, password string) (int64, error) {
	query := `INSERT INTO users(NAME, SURNAME, EMAIL, PASSWORD) VALUES ($1, $2, $3, $4) RETURNING user_id`

	var id int64

	err := p.db.QueryRow(context.Background(), query, name, surname, email, password).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
}

func (p *postgreSQL) SetupDatabase() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Error(err)
	}
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	port := os.Getenv("PORT")

	fmt.Println("DB User:", dbUser)
	fmt.Println("Port:", port)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, port, dbName)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return err
	}

	p.db = conn
	return nil
}

func (p *postgreSQL) CloseDatabase() error {
	return p.db.Close(context.Background())
}
