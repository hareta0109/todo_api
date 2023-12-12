package usecase

import (
	"todo_api/internal/domain/model"
	"todo_api/internal/domain/repository"
	"todo_api/internal/lib/apperr"
)

type AuthUsecase interface {
	Create(params AuthCreateParams) (*model.UserIdentifier, apperr.AppErr)
	Update(id model.UserIdentifier, params AuthUpdateParams) apperr.AppErr
	Login(params AuthLoginParams) (*model.Auth, apperr.AppErr)
	Get(id model.UserIdentifier) (*model.Auth, apperr.AppErr)
}

type AuthCreateParams struct {
	Name      string
	Password  string
	Role      model.UserRole
	UserType  model.UserType
	CompanyID model.CompanyIdentifier
}

type AuthUpdateParams struct {
	Name     string
	Role     model.UserRole
	UserType model.UserType
}

type AuthLoginParams struct {
	ID       model.UserIdentifier
	Password string
}

type authUsecase struct {
	authRepository    repository.AuthRepository
	companyRepository repository.CompanyRepository
}

func NewAuthUsecase(authRepository repository.AuthRepository, companyRepository repository.CompanyRepository) AuthUsecase {
	return &authUsecase{
		authRepository, companyRepository,
	}
}

func (u *authUsecase) Create(
	params AuthCreateParams,
) (*model.UserIdentifier, apperr.AppErr) {
	company, err := u.companyRepository.Get(params.CompanyID)
	if err != nil {
		return nil, err
	}

	authDescription := model.AuthDescription{
		Name:     params.Name,
		Password: &params.Password,
		Role:     params.Role,
		UserType: params.UserType,
		Company:  company,
	}
	auth, err := model.NewAuth(authDescription)
	if err != nil {
		return nil, err
	}
	userID, err := u.authRepository.Create(auth)
	if err != nil {
		return nil, err
	}

	return userID, nil
}

func (u *authUsecase) Update(
	id model.UserIdentifier,
	params AuthUpdateParams,
) apperr.AppErr {
	auth, err := u.authRepository.Get(model.UserIdentifier(id))
	if err != nil {
		return err
	}

	desc := model.AuthDescription{
		Name:     params.Name,
		Role:     params.Role,
		UserType: params.UserType,
	}
	if err = auth.Update(desc); err != nil {
		return err
	}
	if err = u.authRepository.Update(auth); err != nil {
		return err
	}

	return nil
}

func (u *authUsecase) Login(
	params AuthLoginParams,
) (*model.Auth, apperr.AppErr) {
	auth, err := u.authRepository.Get(params.ID)
	if err != nil {
		return nil, err
	}

	hash := model.PasswordToHash(params.Password)
	if auth.Hash != hash {
		return nil, apperr.NewBadRequestError()
	}

	return auth, nil
}

func (u *authUsecase) Get(
	id model.UserIdentifier,
) (*model.Auth, apperr.AppErr) {
	auth, err := u.authRepository.Get(model.UserIdentifier(id))
	if err != nil {
		return nil, err
	}

	return auth, nil
}
