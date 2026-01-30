package services

import (
	"backend/internal/tools"
	"context"
	"time"
)

type TransactionService struct {
	db  tools.DbInterface
	ctx context.Context
}

func (a TransactionService) CreateApplication(
	userId int64,
	statusId int64,
	companyName string,
	location string,
	employeeCount *int,
	jobTitle string,
	platformId int64,
	jobUrl *string,
	salaryMin *int,
	salaryMax *int,
	appliedAt time.Time,
	noteContent *string,
) (*int64, error) {

	tx, err := a.db.BeginTx(a.ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err := (*tx).Rollback(a.ctx)
			if err != nil {
				return
			}
		} else {
			err := (*tx).Commit(a.ctx)
			if err != nil {
				return
			}
		}
	}()

	companyID, err := a.db.CreateCompany(a.ctx, tx, companyName, location, employeeCount)
	if err != nil {
		return nil, err
	}

	appID, err := a.db.CreateApplication(a.ctx, tx, userId, statusId, companyID, jobTitle, platformId, jobUrl, salaryMin, salaryMax, appliedAt)
	if err != nil {
		return nil, err
	}

	if noteContent != nil {
		_, err := a.db.CreateNote(a.ctx, tx, appID, *noteContent)

		if err != nil {
			return nil, err
		}
	}

	return &appID, nil
}

type AppServiceInterface interface {
	CreateApplication(
		userId int64,
		statusId int64,
		companyName string,
		location string,
		employeesCount *int,
		jobTitle string,
		platformId int64,
		jobUrl *string,
		salaryMin *int,
		salaryMax *int,
		appliedAt time.Time,
		noteContent *string,
	) (*int64, error)
}

func NewTransactionService(ctx context.Context, db tools.DbInterface) *AppServiceInterface {

	var t AppServiceInterface = &TransactionService{
		db:  db,
		ctx: ctx,
	}

	return &t
}
