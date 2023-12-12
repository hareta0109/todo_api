package repository

import (
	"errors"
	"todo_api/internal/adapter/outbound/mysql/model"
	domain "todo_api/internal/domain/model"
	"todo_api/internal/lib/apperr"

	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db}
}

func (r *TaskRepository) Get(id domain.TaskIdentifier) (*domain.Task, apperr.AppErr) {
	var row *model.Task
	if err := r.db.
		Preload("PersonInCharge.Company").
		Preload("Creator.Company").
		Preload("Updator.Company").
		First(&row, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NewNotFoundError().Wrap(err)
		}
		return nil, apperr.NewInternalServerError().Wrap(err)
	}
	task, aerr := model.MarshalTask(row)
	if aerr != nil {
		return nil, aerr
	}
	return task, nil
}

func (r *TaskRepository) Find(userID domain.UserIdentifier, id domain.TaskIdentifier) (*domain.Task, apperr.AppErr) {
	var rowUser *model.User
	if err := r.db.Preload("Company").First(&rowUser, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NewNotFoundError().Wrap(err)
		}
		return nil, apperr.NewInternalServerError().Wrap(err)
	}
	user, aerr := model.MarshalUser(rowUser)
	if aerr != nil {
		return nil, aerr
	}

	var rowTask *model.Task
	if err := r.db.
		Preload("PersonInCharge.Company").
		Preload("Creator.Company").
		Preload("Updator.Company").
		First(&rowTask, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NewNotFoundError().Wrap(err)
		}
		return nil, apperr.NewInternalServerError().Wrap(err)
	}
	task, aerr := model.MarshalTask(rowTask)
	if aerr != nil {
		return nil, aerr
	}

	// 閲覧可能かどうかのチェック
	if user.Company.ID != domain.AdminCompanyID {
		// 自社のタスクでない場合は閲覧不可
		if task.Creator.Company.ID != user.Company.ID {
			return nil, apperr.NewForbiddenError()
		}
		// 作成者が自身でなく、会社に公開でない場合は閲覧不可
		if task.Creator.ID != user.ID && task.Visibility == domain.TaskVisibilityMe {
			return nil, apperr.NewForbiddenError()
		}
	}

	return task, nil
}

func (r *TaskRepository) ListByAssignedUserID(userID, assignedUserID domain.UserIdentifier) ([]*domain.Task, apperr.AppErr) {
	var rows []*model.Task
	query := r.db.
		Preload("PersonInCharge.Company").
		Preload("Creator.Company").
		Preload("Updator.Company")
	if userID != assignedUserID {
		query = query.Where("visibility", "COMPANY")
	}
	if err := query.Where("person_in_charge_id", assignedUserID).
		Find(&rows).Error; err != nil {
		return nil, apperr.NewInternalServerError().Wrap(err)
	}
	if len(rows) == 0 {
		return nil, apperr.NewNotFoundError()
	}
	var tasks []*domain.Task
	for _, row := range rows {
		task, aerr := model.MarshalTask(row)
		if aerr != nil {
			return nil, aerr
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) ListByCompanyID(userID domain.UserIdentifier, companyID domain.CompanyIdentifier) ([]*domain.Task, apperr.AppErr) {
	var companyUserIDs []uint64
	if err := r.db.Table("user").
		Where("company_id", companyID).
		Select("id").Take(&companyUserIDs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NewNotFoundError().Wrap(err)
		}
		return nil, apperr.NewInternalServerError().Wrap(err)
	}
	var rows []*model.Task
	if err := r.db.
		Preload("PersonInCharge.Company").
		Preload("Creator.Company").
		Preload("Updator.Company").
		Where("person_in_charge_id = ? or visibility = ?", userID, "COMPANY").
		Where("creator_id", companyUserIDs).Find(&rows).Error; err != nil {
		return nil, apperr.NewInternalServerError().Wrap(err)
	}
	if len(rows) == 0 {
		return nil, apperr.NewNotFoundError()
	}
	var tasks []*domain.Task
	for _, row := range rows {
		task, aerr := model.MarshalTask(row)
		if aerr != nil {
			return nil, aerr
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepository) Create(task *domain.Task) (*domain.TaskIdentifier, apperr.AppErr) {
	row := model.UnmarshalTask(task)
	if err := r.db.Create(&row).Error; err != nil {
		return nil, apperr.NewInternalServerError().Wrap(err)
	}
	id := domain.TaskIdentifier(row.ID)
	return &id, nil
}

func (r *TaskRepository) Update(task *domain.Task) apperr.AppErr {
	row := model.UnmarshalTask(task)
	if err := r.db.Save(&row).Error; err != nil {
		return apperr.NewInternalServerError().Wrap(err)
	}
	return nil
}

func (r *TaskRepository) UpdateStatus(id domain.TaskIdentifier, status domain.TaskStatus) apperr.AppErr {
	if err := r.db.Model(&model.Task{}).
		Where("id", id).Update("task_status", status.String()).
		Error; err != nil {
		return apperr.NewInternalServerError().Wrap(err)
	}
	return nil
}
