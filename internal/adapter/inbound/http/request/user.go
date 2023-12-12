package request

import (
	domain "todo_api/internal/domain/model"
	"todo_api/internal/lib/apperr"
)

type UserCreate struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
	UserType string `json:"user_type"`
}

type UserUpdate struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	UserType string `json:"user_type"`
}

func marshalUserRole(s string) (*domain.UserRole, apperr.AppErr) {
	var role domain.UserRole
	switch s {
	case "EDITOR":
		role = domain.UserRoleEditor
	case "VIEWER":
		role = domain.UserRoleViewer
	default:
		return nil, apperr.NewBadRequestError()
	}
	return &role, nil
}

func marshalUserType(s string) (*domain.UserType, apperr.AppErr) {
	var userType domain.UserType
	switch s {
	case "ADMIN":
		userType = domain.UserTypeAdmin
	case "NORMAL":
		userType = domain.UserTypeNormal
	default:
		return nil, apperr.NewBadRequestError()
	}
	return &userType, nil
}
