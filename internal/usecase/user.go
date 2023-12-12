package usecase

import (
	"todo_api/internal/domain/model"
	"todo_api/internal/domain/repository"
	"todo_api/internal/lib/apperr"
)

type UserUsecase interface {
	Get(id model.UserIdentifier) (*model.User, apperr.AppErr)
}

type userUsecase struct {
	userRepository    repository.UserRepository
	companyRepository repository.CompanyRepository
}

func NewUserUsecase(userRepository repository.UserRepository, companyRepository repository.CompanyRepository) UserUsecase {
	return &userUsecase{
		userRepository, companyRepository,
	}
}

func (u *userUsecase) Get(
	id model.UserIdentifier,
) (*model.User, apperr.AppErr) {
	user, err := u.userRepository.Get(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
