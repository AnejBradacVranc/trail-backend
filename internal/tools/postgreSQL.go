package tools

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type postgreSQL struct {
	db *pgxpool.Pool
}

func (p *postgreSQL) GetStatisticsSummary(email string) ([]*StatisticsSummary, error) {

	var period = "month"

	//Skupni imenovalec mogoce da gre dolocit za kateri time span in za kateri status gledamo
	//Total applications, active interviews, Response rate, Next coming up interview
	applicationsStatsQuery := `
	SELECT
		this_month_count,
		CASE
			WHEN previous_month_count = 0 THEN NULL
			ELSE ROUND(
				((this_month_count - previous_month_count)::numeric
				 / previous_month_count) * 100,
				2
			)
		END AS delta,
		
		CASE
			WHEN previous_month_count = 0 AND this_month_count > 0 THEN 'new'
			WHEN this_month_count > previous_month_count THEN 'up'
			WHEN this_month_count < previous_month_count THEN 'down'
			ELSE 'flat'
		END AS trend
	FROM (
		SELECT
			COUNT(*) FILTER (
				WHERE a.created_at >= date_trunc('month', CURRENT_DATE)
			) AS this_month_count,
			COUNT(*) FILTER (
				WHERE a.created_at >= date_trunc('month', CURRENT_DATE) - INTERVAL '1 month'
			  AND a.created_at <  date_trunc('month', CURRENT_DATE)
			) AS previous_month_count
		FROM applications a
		JOIN users u ON u.user_id = a.user_id
		WHERE u.email = $1
	) as s`

	row := p.db.QueryRow(context.Background(), applicationsStatsQuery, email)

	applicationsStat := &StatisticsSummary{}
	applicationsStat.Period = &period
	applicationsStat.Name = "Total applications"
	appStatText := "Total applications you applied for this month."
	applicationsStat.SummaryText = &appStatText

	err := row.Scan(&applicationsStat.Value, &applicationsStat.Delta, &applicationsStat.Trend)

	if err != nil {
		return nil, err
	}

	//Applications that are over the applied state
	interviewsStatsQuery := `
	SELECT
		this_month_count,
		CASE
			WHEN previous_month_count = 0 THEN NULL
			ELSE ROUND(
				((this_month_count - previous_month_count)::numeric
				 / previous_month_count) * 100,
				2
			)
		END AS delta,
		
		CASE
			WHEN previous_month_count = 0 AND this_month_count > 0 THEN 'new'
			WHEN this_month_count > previous_month_count THEN 'up'
			WHEN this_month_count < previous_month_count THEN 'down'
			ELSE 'flat'
		END AS trend
	FROM (
		SELECT
			COUNT(*) FILTER (
				WHERE a.created_at >= date_trunc('month', CURRENT_DATE)
			) AS this_month_count,
			COUNT(*) FILTER (
				WHERE a.created_at >= date_trunc('month', CURRENT_DATE) - INTERVAL '1 month'
			  AND a.created_at <  date_trunc('month', CURRENT_DATE)
			) AS previous_month_count
		FROM applications a
		JOIN users u ON u.user_id = a.user_id
		WHERE u.email = $1 AND a.status_id IN (2, 3, 4)
	) as s`

	row = p.db.QueryRow(context.Background(), interviewsStatsQuery, email)

	interviewsStat := &StatisticsSummary{}
	interviewsStat.Period = &period
	interviewsStat.Name = "Active applications"
	interviewsStatText := "Number of currently open job applications."
	interviewsStat.SummaryText = &interviewsStatText

	err = row.Scan(&interviewsStat.Value, &interviewsStat.Delta, &interviewsStat.Trend)

	if err != nil {
		return nil, err
	}

	//Share of applications that are over the applied state in relation to total applications

	processedApplicationsQuery := `SELECT this_month_count,
	CASE
        WHEN this_month_count = 0 THEN NULL
        ELSE ROUND(
            ((this_month_processed_count)::numeric
             / this_month_count) * 100,
            2
        )
    END AS value,	
	CASE
        WHEN previous_month_count = 0 THEN NULL
        ELSE ROUND(
            ((previous_month_processed_count)::numeric
             / previous_month_count) * 100,
            2
        )
    END AS delta,
	
    CASE
		WHEN previous_month_processed_count = 0 AND this_month_processed_count > 0 THEN 'new'
        WHEN this_month_processed_count > previous_month_processed_count THEN 'up'
        WHEN this_month_processed_count < previous_month_processed_count THEN 'down'
        ELSE 'flat'
    END AS trend
FROM (
    SELECT
        COUNT(*) FILTER (
            WHERE a.created_at >= date_trunc('month', CURRENT_DATE)
        ) AS this_month_count,
		COUNT (*) FILTER (  
			WHERE a.created_at >= date_trunc('month', CURRENT_DATE) AND a.status_id != 1 
		) as this_month_processed_count,	
        COUNT(*) FILTER (
            WHERE a.created_at >= date_trunc('month', CURRENT_DATE) - INTERVAL '1 month'
          AND a.created_at <  date_trunc('month', CURRENT_DATE)
        ) AS previous_month_count,
			COUNT (*) FILTER (  
			WHERE a.created_at >= date_trunc('month', CURRENT_DATE) - INTERVAL '1 month'
          AND a.created_at <  date_trunc('month', CURRENT_DATE) AND a.status_id != 1 
		) as previous_month_processed_count
    FROM applications a
	JOIN users u ON u.user_id = a.user_id
	WHERE u.email = $1
	) as p`

	row = p.db.QueryRow(context.Background(), processedApplicationsQuery, email)

	var percentUnit = "%"
	var thisMonthApplicationsCount string
	processedStat := &StatisticsSummary{}
	processedStat.Period = &period
	processedStat.Unit = &percentUnit

	err = row.Scan(&thisMonthApplicationsCount, &processedStat.Value, &processedStat.Delta, &processedStat.Trend)

	if err != nil {
		return nil, err
	}

	processedStat.Name = "Processed applications"
	processedStatText := fmt.Sprintf("Percentage of processed applications out of %s applied this month.", thisMonthApplicationsCount)
	processedStat.SummaryText = &processedStatText

	//Time to next interview

	nextInterviewQuery := `SELECT
    	a.interview_at,
		EXTRACT(DAY FROM (a.interview_at - CURRENT_DATE)) AS days_until_interview
		FROM applications a
		JOIN users u ON u.user_id = a.user_id
		WHERE a.status_id = 3
		  AND a.interview_at >= CURRENT_DATE AND u.email = $1
		ORDER BY a.interview_at
		LIMIT 1;`

	nextInterviewStatUnit := "d"
	nextInterviewStat := &StatisticsSummary{}
	nextInterviewStat.Unit = &nextInterviewStatUnit
	nextInterviewStat.Name = "Time to next interview"
	nextInterviewDefaultText := "No interviews scheduled yet."

	var nextInterviewDate *time.Time

	row = p.db.QueryRow(context.Background(), nextInterviewQuery, email)

	err = row.Scan(&nextInterviewDate, &nextInterviewStat.Value)

	if nextInterviewDate == nil || nextInterviewStat.Value == 0 {
		nextInterviewStat.SummaryText = &nextInterviewDefaultText
	} else {
		text := fmt.Sprintf("Interview scheduled for %s.", nextInterviewDate.Format("02.01.2006"))
		nextInterviewStat.SummaryText = &text
	}

	if err != nil {
		return nil, err
	}

	var stats []*StatisticsSummary

	stats = append(stats, applicationsStat)
	stats = append(stats, interviewsStat)
	stats = append(stats, processedStat)
	stats = append(stats, nextInterviewStat)

	return stats, nil
}

func (p *postgreSQL) GetPlatforms() ([]*Platform, error) {

	query := `SELECT platform_id, name, website, is_active FROM platforms`
	rows, err := p.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var platforms []*Platform

	for rows.Next() {
		platform := &Platform{}

		err := rows.Scan(&platform.PlatformID, &platform.Name, &platform.Website, &platform.IsActive)
		if err != nil {
			return nil, err
		}
		platforms = append(platforms, platform)

	}

	return platforms, nil
}

func (p *postgreSQL) GetStatuses() ([]*ApplicationStatus, error) {
	query := `SELECT status_id, status_name FROM application_statuses`
	rows, err := p.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []*ApplicationStatus

	for rows.Next() {
		applicationStatus := &ApplicationStatus{}
		err := rows.Scan(&applicationStatus.StatusID, &applicationStatus.StatusName)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, applicationStatus)
	}

	return statuses, nil

}

func (p *postgreSQL) BeginTx(ctx context.Context) (*pgx.Tx, error) {
	tx, err := p.db.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return nil, err
	}

	return &tx, nil
}

func (p *postgreSQL) CreateFile(ctx context.Context, tx *pgx.Tx, applicationId int64, filename string) (int64, error) {
	query := `INSERT INTO files(application_id, filename) VALUES ($1, $2) RETURNING file_id`

	var id int64

	if tx != nil {
		err := (*tx).QueryRow(ctx, query, applicationId, filename).Scan(&id)
		if err != nil {
			return -1, err
		}
	} else {
		err := p.db.QueryRow(ctx, query, applicationId, filename).Scan(&id)
		if err != nil {
			return -1, err
		}
	}

	return id, nil
}

func (p *postgreSQL) CreateNote(ctx context.Context, tx *pgx.Tx, applicationId int64, noteContent string) (int64, error) {
	query := `INSERT INTO notes(application_id, note_content) VALUES ($1, $2) RETURNING note_id`

	var id int64

	if tx != nil {
		err := (*tx).QueryRow(ctx, query, applicationId, noteContent).Scan(&id)
		if err != nil {
			return -1, err
		}
	} else {
		err := p.db.QueryRow(ctx, query, applicationId, noteContent).Scan(&id)
		if err != nil {
			return -1, err
		}
	}

	return id, nil
}

func (p *postgreSQL) CreateCompany(ctx context.Context, tx *pgx.Tx, name string, location string, employeesCount *int) (int64, error) {
	query := `INSERT INTO companies(name, location, employee_count) 
				VALUES ($1, $2, $3) 
				ON CONFLICT (name) DO UPDATE
				SET
    				location = COALESCE(EXCLUDED.location, companies.location),
    				employee_count = COALESCE(EXCLUDED.employee_count, companies.employee_count)
				RETURNING company_id`

	var id int64

	if tx != nil {
		err := (*tx).QueryRow(ctx, query, name, location, employeesCount).Scan(&id)
		if err != nil {
			return -1, err
		}
	} else {
		err := p.db.QueryRow(ctx, query, name, location, employeesCount).Scan(&id)
		if err != nil {
			return -1, err
		}
	}

	return id, nil
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

func (p *postgreSQL) CreateApplication(ctx context.Context, tx *pgx.Tx, userId int64, statusId int64, companyId int64, jobTitle string, platformId int64, jobUrl *string, salaryMin *int, salaryMax *int, appliedAt time.Time, interviewAt *time.Time) (int64, error) {

	query := `INSERT INTO applications(USER_ID, STATUS_ID, COMPANY_ID, JOB_TITLE, PLATFORM_ID, JOB_URL, SALARY_MAX, SALARY_MIN, APPLIED_AT, INTERVIEW_AT) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING application_id`

	var id int64

	if tx != nil {
		err := (*tx).QueryRow(ctx, query, userId, statusId, companyId, jobTitle, platformId, jobUrl, salaryMax, salaryMin, appliedAt, interviewAt).Scan(&id)
		if err != nil {
			return -1, err
		}
	} else {
		err := p.db.QueryRow(ctx, query, userId, statusId, companyId, jobTitle, platformId, jobUrl, salaryMax, salaryMin, appliedAt, interviewAt).Scan(&id)
		if err != nil {
			return -1, err
		}
	}

	return id, nil
}

func (p *postgreSQL) GetApplicationByID(applicationId int64) (*ApplicationDetail, error) {
	query := `
		SELECT
			a.application_id,
			a.job_title,
			a.job_url,
			a.salary_min,
			a.salary_max,
			a.created_at,
			a.modified_at,
			s.status_name,
			c.company_id,
			c.name,
			c.location,
			c.employee_count,
			p.platform_id,
			p.name,
			p.website,
			p.is_active
		FROM applications a
		JOIN application_statuses s ON s.status_id = a.status_id
		JOIN companies c ON c.company_id = a.company_id
		JOIN platforms p ON p.platform_id = a.platform_id
		WHERE a.application_id = $1
	`

	row := p.db.QueryRow(context.Background(), query, applicationId)

	detail := &ApplicationDetail{}

	err := row.Scan(
		&detail.ApplicationID,
		&detail.JobTitle,
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
		&detail.Platform.PlatformID,
		&detail.Platform.Name,
		&detail.Platform.Website,
		&detail.Platform.IsActive,
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

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return err
	}

	p.db = pool
	return nil
}

func (p *postgreSQL) CloseDatabase() {
	if p.db != nil {
		p.db.Close()
	}
}
