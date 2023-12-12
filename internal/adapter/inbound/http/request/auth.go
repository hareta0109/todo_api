package request

import (
	domain "todo_api/internal/domain/model"
	"todo_api/internal/lib/apperr"
	"todo_api/internal/usecase"
)

type AuthCreate struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
	UserType string `json:"user_type"`
}

func MarshalAuthCreateParams(companyID uint64, req *AuthCreate) (*usecase.AuthCreateParams, apperr.AppErr) {
	if req == nil {
		return nil, nil
	}
	role, err := marshalUserRole(req.Role)
	if err != nil {
		return nil, err
	}
	userType, err := marshalUserType(req.UserType)
	if err != nil {
		return nil, err
	}
	return &usecase.AuthCreateParams{
		Name:      req.Name,
		Password:  req.Password,
		Role:      *role,
		UserType:  *userType,
		CompanyID: domain.CompanyIdentifier(companyID),
	}, nil
}

type AuthUpdate struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	UserType string `json:"user_type"`
}

func MarshalAuthUpdateParams(companyID uint64, req *AuthUpdate) (*usecase.AuthUpdateParams, apperr.AppErr) {
	if req == nil {
		return nil, nil
	}
	role, err := marshalUserRole(req.Role)
	if err != nil {
		return nil, err
	}
	userType, err := marshalUserType(req.UserType)
	if err != nil {
		return nil, err
	}
	return &usecase.AuthUpdateParams{
		Name:     req.Name,
		Role:     *role,
		UserType: *userType,
	}, nil
}

type AuthLogin struct {
	ID       uint64 `json:"id"`
	Password string `json:"password"`
}

func MarshalAuthLoginParams(req *AuthLogin) *usecase.AuthLoginParams {
	if req == nil {
		return nil
	}
	return &usecase.AuthLoginParams{
		ID:       domain.UserIdentifier(req.ID),
		Password: req.Password,
	}
}
