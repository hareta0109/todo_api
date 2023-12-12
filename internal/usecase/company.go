package usecase

import (
	"todo_api/internal/domain/model"
	"todo_api/internal/domain/repository"
	"todo_api/internal/lib/apperr"
)

type CompanyUsecase interface {
	Get(id model.CompanyIdentifier) (*model.Company, apperr.AppErr)
	Create(name string) (*model.CompanyIdentifier, apperr.AppErr)
	Update(id model.CompanyIdentifier, name string) apperr.AppErr
}

type companyUsecase struct {
	companyRepository repository.CompanyRepository
}

func NewCompanyUsecase(companyRepository repository.CompanyRepository) CompanyUsecase {
	return &companyUsecase{
		companyRepository,
	}
}

func (u *companyUsecase) Get(
	id model.CompanyIdentifier,
) (*model.Company, apperr.AppErr) {
	company, err := u.companyRepository.Get(id)
	if err != nil {
		return nil, err
	}

	return company, nil
}

func (u *companyUsecase) Create(
	name string,
) (*model.CompanyIdentifier, apperr.AppErr) {
	desc := model.CompanyDescipriton{
		Name: name,
	}
	company, err := model.NewCompany(desc)
	if err != nil {
		return nil, err
	}
	companyID, err := u.companyRepository.Create(company)
	if err != nil {
		return nil, err
	}

	return companyID, nil
}

func (u *companyUsecase) Update(
	id model.CompanyIdentifier,
	name string,
) apperr.AppErr {
	company, err := u.companyRepository.Get(id)
	if err != nil {
		return err
	}

	desc := model.CompanyDescipriton{
		Name: name,
	}
	if err = company.Update(desc); err != nil {
		return err
	}

	if err := u.companyRepository.Update(company); err != nil {
		return err
	}

	return nil
}
