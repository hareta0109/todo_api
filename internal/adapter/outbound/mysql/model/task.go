package model

import (
	"time"
	domain "todo_api/internal/domain/model"
	"todo_api/internal/lib/apperr"
)

type Task struct {
	ID               uint64
	Title            string
	Detail           *string
	Status           string `gorm:"column:task_status"`
	Visibility       string
	PersonInChargeID *uint64
	PersonInCharge   *User `gorm:"foreignKey:PersonInChargeID"`
	LimitDate        *time.Time

	CreateAt  time.Time `gorm:"autoCreateTime"`
	CreatorID uint64
	Creator   User      `gorm:"foreignKey:CreatorID"`
	UpdateAt  time.Time `gorm:"autoUpdateTime"`
	UpdatorID uint64
	Updator   User `gorm:"foreignKey:UpdatorID"`
}

func (m *Task) TableName() string {
	return "task"
}

func UnmarshalTask(d *domain.Task) *Task {
	if d == nil {
		return nil
	}
	return &Task{
		ID:               uint64(d.ID),
		Title:            d.Title,
		Detail:           d.Detail,
		Status:           d.Status.String(),
		Visibility:       d.Visibility.String(),
		PersonInChargeID: (*uint64)(&d.PersonInCharge.ID),
		LimitDate:        d.LimitDate,
		CreateAt:         d.CreateAt,
		CreatorID:        uint64(d.Creator.ID),
		UpdateAt:         d.UpdateAt,
		UpdatorID:        uint64(d.Updator.ID),
	}
}

func MarshalTask(m *Task) (*domain.Task, apperr.AppErr) {
	if m == nil {
		return nil, nil
	}
	status, err := marshalTaskStatus(m.Status)
	if err != nil {
		return nil, err
	}
	visibility, err := marshalTaskVisibility(m.Visibility)
	if err != nil {
		return nil, err
	}

	personInCharge, err := MarshalUser(m.PersonInCharge)
	if err != nil {
		return nil, err
	}
	creator, err := MarshalUser(&m.Creator)
	if err != nil {
		return nil, err
	}
	updator, err := MarshalUser(&m.Updator)
	if err != nil {
		return nil, err
	}
	return &domain.Task{
		ID:             domain.TaskIdentifier(m.ID),
		Title:          m.Title,
		Detail:         m.Detail,
		Status:         *status,
		Visibility:     *visibility,
		PersonInCharge: personInCharge,
		LimitDate:      m.LimitDate,
		CreateAt:       m.CreateAt,
		Creator:        *creator,
		UpdateAt:       m.UpdateAt,
		Updator:        *updator,
	}, nil
}

func marshalTaskStatus(s string) (*domain.TaskStatus, apperr.AppErr) {
	var status domain.TaskStatus
	switch s {
	case "NEW":
		status = domain.TaskStatusNew
	case "PROCESSING":
		status = domain.TaskStatusProcessing
	case "DONE":
		status = domain.TaskStatusDone
	default:
		return nil, apperr.NewInternalServerError()
	}
	return &status, nil
}

func marshalTaskVisibility(s string) (*domain.TaskVisibility, apperr.AppErr) {
	var visibility domain.TaskVisibility
	switch s {
	case "ME":
		visibility = domain.TaskVisibilityMe
	case "COMPANY":
		visibility = domain.TaskVisibilityCompany
	default:
		return nil, apperr.NewInternalServerError()
	}
	return &visibility, nil
}
