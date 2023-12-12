package model

import (
	domain "todo_api/internal/domain/model"
)

type Company struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
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
