package model

import (
	"errors"
	"todo_api/internal/lib/apperr"
	"unicode/utf8"
)

var (
	errInvalidUserNameLength = errors.New("User Name must be 1 to 20 characters")
)

const (
	minUserNameLength = 1
	maxUserNameLength = 20
)

type User struct {
	ID       UserIdentifier
	Name     string
	Role     UserRole
	UserType UserType
	Company  Company
}

type UserDescription struct {
	Name     string
	Role     UserRole
	UserType UserType
	Company  Company
}

type UserIdentifier uint64

type UserType int

const (
	UserTypeAdmin UserType = iota + 1
	UserTypeNormal
)

type UserRole int

const (
	UserRoleEditor UserRole = iota + 1
	UserRoleViewer
)

func NewUser(desc UserDescription) (*User, apperr.AppErr) {
	user := new(User)
	if err := user.Update(desc); err != nil {
		return nil, err
	}

	return user, nil
}

func (m *User) Update(desc UserDescription) apperr.AppErr {
	if err := desc.validate(); err != nil {
		return apperr.NewBadRequestError().Wrap(err)
	}

	m.Name = desc.Name
	m.Role = desc.Role
	m.UserType = desc.UserType
	m.Company = desc.Company
	return nil
}

func (d *UserDescription) validate() error {
	nameLength := utf8.RuneCountInString(d.Name)
	if nameLength < minUserNameLength || nameLength > maxUserNameLength {
		return errInvalidUserNameLength
	}

	return nil
}

func (e UserType) String() string {
	switch e {
	case UserTypeAdmin:
		return "ADMIN"
	case UserTypeNormal:
		return "NORMAL"
	default:
		return ""
	}
}

func (e UserRole) String() string {
	switch e {
	case UserRoleEditor:
		return "EDITOR"
	case UserRoleViewer:
		return "VIEWER"
	default:
		return ""
	}
}
