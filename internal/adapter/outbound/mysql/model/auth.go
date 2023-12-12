package model

import (
	domain "todo_api/internal/domain/model"
	"todo_api/internal/lib/apperr"
)

type Auth struct {
	ID        uint64
	Name      string `gorm:"column:user_name"`
	Hash      string `gorm:"column:user_hash"`
	Role      string `gorm:"column:user_role"`
	UserType  string
	CompanyID uint64
	Company   Company
}

func (m *Auth) TableName() string {
	return "user"
}

func UnmarshalAuth(d *domain.Auth) *Auth {
	if d == nil {
		return nil
	}
	return &Auth{
		ID:       uint64(d.ID),
		Name:     d.Name,
		Hash:     d.Hash,
		Role:     d.Role.String(),
		UserType: d.UserType.String(),
		Company:  *UnmarshalCompany(&d.Company),
	}
}

func MarshalAuth(m *Auth) (*domain.Auth, apperr.AppErr) {
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

	return &domain.Auth{
		ID:       domain.UserIdentifier(m.ID),
		Name:     m.Name,
		Hash:     m.Hash,
		Role:     *role,
		UserType: *userType,
		Company:  *MarshalCompany(&m.Company),
	}, nil
}
