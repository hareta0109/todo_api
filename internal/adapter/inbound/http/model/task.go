package model

import (
	"time"
	domain "todo_api/internal/domain/model"
)

type Task struct {
	ID             uint64     `json:"id"`
	Title          string     `json:"title"`
	Detail         *string    `json:"detail,omitempty"`
	Status         string     `json:"status"`
	Visibility     string     `json:"visibility"`
	PersonInCharge *User      `json:"person_in_charge,omitempty"`
	LimitDate      *time.Time `json:"limit_date,omitempty"`

	CreateAt time.Time `json:"create_at"`
	Creator  User      `json:"creator"`
	UpdateAt time.Time `json:"update_at"`
	Updator  User     `json:"updator"`
}

func unmarshalStatus(d domain.TaskStatus) string {
	switch d {
	case domain.TaskStatusNew:
		return "NEW"
	case domain.TaskStatusProcessing:
		return "PROCESSING"
	case domain.TaskStatusDone:
		return "DONE"
	default:
		return ""
	}
}

func unmarshalVisibility(d domain.TaskVisibility) string {
	switch d {
	case domain.TaskVisibilityMe:
		return "ME"
	case domain.TaskVisibilityCompany:
		return "COMPANY"
	default:
		return ""
	}
}

func UnmarshalTask(d *domain.Task) *Task {
	if d == nil {
		return nil
	}
	return &Task{
		ID:             uint64(d.ID),
		Title:          d.Title,
		Detail:         d.Detail,
		Status:         unmarshalStatus(d.Status),
		Visibility:     unmarshalVisibility(d.Visibility),
		PersonInCharge: UnmarshalUser(d.PersonInCharge),
		LimitDate:      d.LimitDate,
		CreateAt:       d.CreateAt,
		Creator:        *UnmarshalUser(&d.Creator),
		UpdateAt:       d.UpdateAt,
		Updator:        *UnmarshalUser(&d.Updator),
	}
}
