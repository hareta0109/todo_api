package repository

import (
	"errors"
	"todo_api/internal/adapter/outbound/mysql/model"
	domain "todo_api/internal/domain/model"
	"todo_api/internal/lib/apperr"

	"gorm.io/gorm"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepostiroy(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db}
}

func (r *CompanyRepository) Get(id domain.CompanyIdentifier) (*domain.Company, apperr.AppErr) {
	var row *model.Company
	if err := r.db.First(&row, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NewNotFoundError().Wrap(err)
		}
		return nil, apperr.NewInternalServerError().Wrap(err)
	}
	company := model.MarshalCompany(row)
	return company, nil
}

func (r *CompanyRepository) Create(company *domain.Company) (*domain.CompanyIdentifier, apperr.AppErr) {
	row := model.UnmarshalCompany(company)
	if err := r.db.Create(&row).Error; err != nil {
		return nil, apperr.NewInternalServerError().Wrap(err)
	}
	id := domain.CompanyIdentifier(row.ID)
	return &id, nil
}

func (r *CompanyRepository) Update(company *domain.Company) apperr.AppErr {
	row := model.UnmarshalCompany(company)
	if err := r.db.Save(&row).Error; err != nil {
		return apperr.NewInternalServerError().Wrap(err)
	}
	return nil
}
