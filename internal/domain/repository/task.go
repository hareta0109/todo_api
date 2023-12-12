package repository

import (
	"todo_api/internal/domain/model"
	"todo_api/internal/lib/apperr"
)

type TaskRepository interface {
	// Get / タスクIDを元にタスクを取得する。
	Get(id model.TaskIdentifier) (*model.Task, apperr.AppErr)
	// Find / 取得者のIDで表示可能なタスクを取得する
	Find(userID model.UserIdentifier, id model.TaskIdentifier) (*model.Task, apperr.AppErr)
	// ListByUserID / 担当者IDを元にタスクを取得する。
	ListByAssignedUserID(userID, assignedUserID model.UserIdentifier) ([]*model.Task, apperr.AppErr)
	// ListByCompanyID / 組織IDを元にタスクを取得する。
	ListByCompanyID(userID model.UserIdentifier, companyID model.CompanyIdentifier) ([]*model.Task, apperr.AppErr)

	Create(task *model.Task) (*model.TaskIdentifier, apperr.AppErr)
	Update(task *model.Task) apperr.AppErr
	UpdateStatus(id model.TaskIdentifier, status model.TaskStatus) apperr.AppErr
}
