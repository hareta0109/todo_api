package request

import (
	"time"
	domain "todo_api/internal/domain/model"
	"todo_api/internal/lib/apperr"
	"todo_api/internal/usecase"
)

type TaskCreate struct {
	Title            string     `json:"title"`
	Detail           *string    `json:"detail"`
	Visibility       string     `json:"visibility"`
	PersonInChargeID *uint64    `json:"person_in_charge_id"`
	LimitDate        *time.Time `json:"limit_date"`
}

type TaskUpdate struct {
	Title            string     `json:"title"`
	Detail           *string    `json:"detail"`
	Visibility       string     `json:"visibility"`
	Status           string     `json:"status"`
	PersonInChargeID *uint64    `json:"person_in_charge_id"`
	LimitDate        *time.Time `json:"limit_date"`
}

func MarshalTaskCreateParams(userID uint64, req *TaskCreate) (*usecase.TaskCreateParams, apperr.AppErr) {
	if req == nil {
		return nil, nil
	}
	visibility, err := marshalTaskVisibility(req.Visibility)
	if err != nil {
		return nil, err
	}
	return &usecase.TaskCreateParams{
		Title:            req.Title,
		Detail:           req.Detail,
		Visibility:       *visibility,
		PersonInChargeID: (*domain.UserIdentifier)(req.PersonInChargeID),
		LimitDate:        req.LimitDate,
		CreatorID:        domain.UserIdentifier(userID),
	}, nil
}

func MarshalTaskUpdateParams(userID uint64, req *TaskUpdate) (*usecase.TaskUpdateParams, apperr.AppErr) {
	if req == nil {
		return nil, nil
	}
	visibility, err := marshalTaskVisibility(req.Visibility)
	if err != nil {
		return nil, err
	}
	status, err := MarshalTaskStatus(req.Status)
	if err != nil {
		return nil, err
	}
	return &usecase.TaskUpdateParams{
		Title:            req.Title,
		Detail:           req.Detail,
		Visibility:       *visibility,
		Status:           *status,
		PersonInChargeID: (*domain.UserIdentifier)(req.PersonInChargeID),
		LimitDate:        req.LimitDate,
		UpdatorID:        domain.UserIdentifier(userID),
	}, nil
}

func MarshalTaskStatus(s string) (*domain.TaskStatus, apperr.AppErr) {
	var status domain.TaskStatus
	switch s {
	case "NEW":
		status = domain.TaskStatusNew
	case "PROCESSING":
		status = domain.TaskStatusProcessing
	case "DONE":
		status = domain.TaskStatusDone
	default:
		return nil, apperr.NewBadRequestError()
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
		return nil, apperr.NewBadRequestError()
	}
	return &visibility, nil
}
