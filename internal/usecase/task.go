package usecase

import (
	"time"
	"todo_api/internal/domain/model"
	"todo_api/internal/domain/repository"
	"todo_api/internal/lib/apperr"
)

type TaskUsecase interface {
	Find(userID model.UserIdentifier, id model.TaskIdentifier) (*model.Task, apperr.AppErr)
	ListByAssignedUserID(userID, assignedUserID model.UserIdentifier) ([]*model.Task, apperr.AppErr)
	ListByCompanyID(userID model.UserIdentifier, companyID model.CompanyIdentifier) ([]*model.Task, apperr.AppErr)

	Create(params TaskCreateParams) (*model.TaskIdentifier, apperr.AppErr)
	Update(id model.TaskIdentifier, params TaskUpdateParams) apperr.AppErr
	UpdateStatus(id model.TaskIdentifier, status model.TaskStatus) apperr.AppErr
}

type TaskCreateParams struct {
	Title            string
	Detail           *string
	Visibility       model.TaskVisibility
	PersonInChargeID *model.UserIdentifier
	LimitDate        *time.Time
	CreatorID        model.UserIdentifier
}

type TaskUpdateParams struct {
	Title            string
	Detail           *string
	Visibility       model.TaskVisibility
	Status           model.TaskStatus
	PersonInChargeID *model.UserIdentifier
	LimitDate        *time.Time
	UpdatorID        model.UserIdentifier
}

type taskUsecase struct {
	userRepository repository.UserRepository
	taskRepository repository.TaskRepository
}

func NewTaskUsecase(
	userRepository repository.UserRepository,
	taskRepository repository.TaskRepository,
) TaskUsecase {
	return &taskUsecase{
		userRepository,
		taskRepository,
	}
}

func (u *taskUsecase) Find(userID model.UserIdentifier, id model.TaskIdentifier) (*model.Task, apperr.AppErr) {
	task, err := u.taskRepository.Find(userID, id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (u *taskUsecase) ListByAssignedUserID(userID, assignedUserID model.UserIdentifier) ([]*model.Task, apperr.AppErr) {
	tasks, err := u.taskRepository.ListByAssignedUserID(userID, assignedUserID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (u *taskUsecase) ListByCompanyID(userID model.UserIdentifier, companyID model.CompanyIdentifier) ([]*model.Task, apperr.AppErr) {
	tasks, err := u.taskRepository.ListByCompanyID(userID, companyID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (u *taskUsecase) Create(params TaskCreateParams) (*model.TaskIdentifier, apperr.AppErr) {
	var personInCharge, creator *model.User
	var err apperr.AppErr

	if params.PersonInChargeID != nil {
		personInCharge, err = u.userRepository.Get(*params.PersonInChargeID)
		if err != nil {
			return nil, err
		}
	}

	creator, err = u.userRepository.Get(params.CreatorID)
	if err != nil {
		return nil, err
	}

	desc := model.TaskDescription{
		Title:          params.Title,
		Detail:         params.Detail,
		Status:         model.TaskStatusNew, // 新規タスクはNEW
		Visibility:     params.Visibility,
		PersonInCharge: personInCharge,
		LimitDate:      params.LimitDate,
		Creator:        creator,
		Updator:        creator,
	}
	task, err := model.NewTask(desc)
	if err != nil {
		return nil, err
	}

	taskID, err := u.taskRepository.Create(task)
	if err != nil {
		return nil, err
	}

	return taskID, err
}

func (u *taskUsecase) Update(id model.TaskIdentifier, params TaskUpdateParams) apperr.AppErr {
	task, err := u.taskRepository.Get(id)
	if err != nil {
		return err
	}

	var personInCharge, updator *model.User
	if params.PersonInChargeID != nil {
		personInCharge, err = u.userRepository.Get(*params.PersonInChargeID)
		if err != nil {
			return err
		}
	}
	updator, err = u.userRepository.Get(params.UpdatorID)
	if err != nil {
		return err
	}

	desc := model.TaskDescription{
		Title:          params.Title,
		Detail:         params.Detail,
		Status:         params.Status,
		Visibility:     params.Visibility,
		PersonInCharge: personInCharge,
		LimitDate:      params.LimitDate,
		Updator:        updator,
	}
	if err = task.Update(desc); err != nil {
		return err
	}

	if err := u.taskRepository.Update(task); err != nil {
		return err
	}

	return nil
}

func (u *taskUsecase) UpdateStatus(id model.TaskIdentifier, status model.TaskStatus) apperr.AppErr {
	_, err := u.taskRepository.Get(id)
	if err != nil {
		return err
	}

	if err = u.taskRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	return nil
}
