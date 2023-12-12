package model

import domain "todo_api/internal/domain/model"

type Company struct {
	ID   uint64
	Name string `gorm:"column:company_name"`
}

func (m *Company) TableName() string {
	return "company"
}

func UnmarshalCompany(d *domain.Company) *Company {
	if d == nil {
		return nil
	}
	return &Company{
		ID:   uint64(d.ID),
		Name: d.Name,
	}
}

func MarshalCompany(m *Company) *domain.Company {
	if m == nil {
		return nil
	}
	return &domain.Company{
		ID:   domain.CompanyIdentifier(m.ID),
		Name: m.Name,
	}
}
