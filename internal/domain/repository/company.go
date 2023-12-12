package repository

import (
	"todo_api/internal/domain/model"
	"todo_api/internal/lib/apperr"
)

type CompanyRepository interface {
	Get(id model.CompanyIdentifier) (*model.Company, apperr.AppErr)

	Create(company *model.Company) (*model.CompanyIdentifier, apperr.AppErr)
	Update(company *model.Company) apperr.AppErr
}
