package model

import (
	domain "todo_api/internal/domain/model"
	"todo_api/internal/lib/apperr"
)

type User struct {
	ID        uint64
	Name      string `gorm:"column:user_name"`
	Role      string `gorm:"column:user_role"`
	UserType  string
	CompanyID uint64
	Company   Company
}

func (m *User) TableName() string {
	return "user"
}

func UnmarshalUser(d *domain.User) *User {
	if d == nil {
		return nil
	}
	return &User{
		ID:        uint64(d.ID),
		Name:      d.Name,
		Role:      d.Role.String(),
		UserType:  d.UserType.String(),
		CompanyID: uint64(d.Company.ID),
	}
}

func MarshalUser(m *User) (*domain.User, apperr.AppErr) {
	if m == nil {
		return nil, nil
	}
	role, err := marshalUserRole(m.Role)
	if err != nil {
		return nil, err
	}
	userType, err := marshalUserType(m.UserType)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:       domain.UserIdentifier(m.ID),
		Name:     m.Name,
		Role:     *role,
		UserType: *userType,
		Company:  *MarshalCompany(&m.Company),
	}, nil
}

func marshalUserRole(s string) (*domain.UserRole, apperr.AppErr) {
	var role domain.UserRole
	switch s {
	case "EDITOR":
		role = domain.UserRoleEditor
	case "VIEWER":
		role = domain.UserRoleViewer
	default:
		return nil, apperr.NewInternalServerError()
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
		return nil, apperr.NewInternalServerError()
	}
	return &userType, nil
}
