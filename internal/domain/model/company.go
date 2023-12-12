package model

import (
	"errors"
	"todo_api/internal/lib/apperr"
	"unicode/utf8"
)

var (
	errInvalidCompanyNameLength = errors.New("Company Name must be 1 to 20 characters")
)

const (
	AdminCompanyID = 1

	minCompanyNameLength = 1
	maxCompanyNameLength = 20
)

type Company struct {
	ID   CompanyIdentifier
	Name string
}

type CompanyDescipriton struct {
	Name string
}

type CompanyIdentifier uint64

func NewCompany(desc CompanyDescipriton) (*Company, apperr.AppErr) {
	company := new(Company)
	if err := company.Update(desc); err != nil {
		return nil, err
	}

	return company, nil
}

func (m *Company) Update(desc CompanyDescipriton) apperr.AppErr {
	if err := desc.validate(); err != nil {
		return apperr.NewBadRequestError().Wrap(err)
	}

	m.Name = desc.Name
	return nil
}

func (d *CompanyDescipriton) validate() error {
	nameLength := utf8.RuneCountInString(d.Name)
	if nameLength < minCompanyNameLength || nameLength > maxCompanyNameLength {
		return errInvalidCompanyNameLength
	}
	return nil
}
