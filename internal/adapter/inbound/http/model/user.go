package model

import (
	domain "todo_api/internal/domain/model"
)

type User struct {
	ID       uint64   `json:"id"`
	Name     string   `json:"name"`
	Role     string   `json:"role"`
	UserType string   `json:"user_type"`
	Company  *Company `json:"company"`
}

func unmarshalRole(d domain.UserRole) string {
	switch d {
	case domain.UserRoleEditor:
		return "EDITOR"
	case domain.UserRoleViewer:
		return "VIEWER"
	default:
		return ""
	}
}

func unmarshalUserType(d domain.UserType) string {
	switch d {
	case domain.UserTypeAdmin:
		return "ADMIN"
	case domain.UserTypeNormal:
		return "NORMAL"
	default:
		return ""
	}
}

func UnmarshalUser(d *domain.User) *User {
	if d == nil {
		return nil
	}
	return &User{
		ID:       uint64(d.ID),
		Name:     d.Name,
		Role:     unmarshalRole(d.Role),
		UserType: unmarshalUserType(d.UserType),
		Company:  UnmarshalCompany(&d.Company),
	}
}
