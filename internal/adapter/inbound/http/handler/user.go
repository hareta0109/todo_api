package handler

import (
	"net/http"
	"strconv"
	"todo_api/internal/adapter/inbound/http/model"
	domain "todo_api/internal/domain/model"
	"todo_api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	Get(c echo.Context) error
}

type userHandler struct {
	authUsecase usecase.AuthUsecase
	userUsecase usecase.UserUsecase
}

func NewUserHandler(
	authUsecase usecase.AuthUsecase,
	userUsecase usecase.UserUsecase,
) UserHandler {
	return &userHandler{
		authUsecase,
		userUsecase,
	}
}

// GetUser
//
//	@Summary		ユーザの取得
//	@Description	ユーザの情報をIDから取得する。
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			company_id		path		int		false	"企業ID"
//	@Param			user_id			path		int		false	"ユーザID"
//	@Success		200				{object}	model.User
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		500
//	@Router			/company/{company_id}/user/{user_id} [get]
func (h *userHandler) Get(c echo.Context) error {
	companyID, err := strconv.ParseUint(c.Param("company_id"), 10, 64)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	// 会社の閲覧権限を持つことの認証処理
	{
		authUserID, err := verifyJwtToken(c)
		if err != nil {
			return &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			}
		}
		authUser, aerr := h.authUsecase.Get(domain.UserIdentifier(authUserID))
		if aerr != nil {
			return aerr.HTTPError()
		}
		if !authUser.CanView(domain.CompanyIdentifier(companyID)) {
			return &echo.HTTPError{
				Code: http.StatusForbidden,
			}
		}
	}

	id, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	company, aerr := h.userUsecase.Get(domain.UserIdentifier(id))
	if aerr != nil {
		return aerr.HTTPError()
	}

	res := model.UnmarshalUser(company)

	return c.JSON(http.StatusOK, res)
}
