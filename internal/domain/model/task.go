package model

import (
	"errors"
	"time"
	"todo_api/internal/lib/apperr"
	"unicode/utf8"
)

var (
	errInvalidTaskTitleLength  = errors.New("Task Title must be 1 to 50 characters")
	errInvalidTaskDetailLength = errors.New("Task Detail must be shorter or equal 400 characters")
	errInvalidTaskAssignment   = errors.New("Task cannot be assigned to an user of other companies")
)

const (
	minTaskTitleLength  = 1
	maxTaskTitleLength  = 50
	maxTaskDetailLength = 200
)

// 部門情報を持たせた方が楽かもしれない
type Task struct {
	ID             TaskIdentifier
	Title          string
	Detail         *string
	Status         TaskStatus
	Visibility     TaskVisibility
	PersonInCharge *User
	LimitDate      *time.Time

	CreateAt time.Time
	Creator  User
	UpdateAt time.Time
	Updator  User
}

type TaskDescription struct {
	Title          string
	Detail         *string
	Status         TaskStatus
	Visibility     TaskVisibility
	PersonInCharge *User
	LimitDate      *time.Time
	Creator        *User
	Updator        *User
}

type TaskIdentifier uint64

type TaskVisibility int

const (
	TaskVisibilityMe TaskVisibility = iota + 1
	TaskVisibilityCompany
)

type TaskStatus int

const (
	TaskStatusNew TaskStatus = iota + 1
	TaskStatusProcessing
	TaskStatusDone
)

func NewTask(desc TaskDescription) (*Task, apperr.AppErr) {
	task := new(Task)
	if err := task.Update(desc); err != nil {
		return nil, err
	}

	return task, nil
}

func (m *Task) Update(desc TaskDescription) apperr.AppErr {
	if err := desc.validate(); err != nil {
		return apperr.NewBadRequestError().Wrap(err)
	}

	m.Title = desc.Title
	m.Detail = desc.Detail
	m.Status = desc.Status
	m.Visibility = desc.Visibility
	m.PersonInCharge = desc.PersonInCharge
	m.LimitDate = desc.LimitDate
	if desc.Creator != nil {
		m.Creator = *desc.Creator
	}
	if desc.Updator != nil {
		m.Updator = *desc.Updator
	}
	return nil
}

func (d *TaskDescription) validate() error {
	titleLength := utf8.RuneCountInString(d.Title)
	if titleLength < minTaskTitleLength || titleLength > maxTaskTitleLength {
		return errInvalidTaskTitleLength
	}
	detailLength := utf8.RuneCountInString(*d.Detail)
	if detailLength > maxTaskDetailLength {
		return errInvalidTaskDetailLength
	}
	if d.PersonInCharge != nil {
		// CreatorかUpdatorは片方必ず存在する
		if d.Creator != nil {
			if d.PersonInCharge.Company.ID != d.Creator.Company.ID {
				return errInvalidTaskAssignment
			}
		}
		if d.Updator != nil {
			if d.PersonInCharge.Company.ID != d.Updator.Company.ID {
				return errInvalidTaskAssignment
			}
		}
	}
	return nil
}

func (e TaskStatus) String() string {
	switch e {
	case TaskStatusNew:
		return "NEW"
	case TaskStatusProcessing:
		return "PROCESSING"
	case TaskStatusDone:
		return "DONE"
	default:
		return ""
	}
}

func (e TaskVisibility) String() string {
	switch e {
	case TaskVisibilityMe:
		return "ME"
	case TaskVisibilityCompany:
		return "COMPANY"
	default:
		return ""
	}
}
