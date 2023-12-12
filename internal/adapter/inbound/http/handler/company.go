package handler

import (
	"net/http"
	"strconv"
	"todo_api/internal/adapter/inbound/http/model"
	"todo_api/internal/adapter/inbound/http/request"
	domain "todo_api/internal/domain/model"
	"todo_api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type CompanyHandler interface {
	Get(c echo.Context) error
	Create(c echo.Context) error
	Update(c echo.Context) error
}

type companyHandler struct {
	authUsecase    usecase.AuthUsecase
	companyUsecase usecase.CompanyUsecase
}

func NewCompanyHandler(
	authUsecase usecase.AuthUsecase,
	companyUsecase usecase.CompanyUsecase,
) CompanyHandler {
	return &companyHandler{
		authUsecase,
		companyUsecase,
	}
}

// GetCompany
//
//	@Summary		企業の取得
//	@Description	企業の情報をIDから取得する。管理会社のユーザのみ実行可能。
//	@Tags			company
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			company_id		path		int		false	"企業ID"
//	@Success		200				{object}	model.Company
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		500
//	@Router			/company/{company_id} [get]
func (h *companyHandler) Get(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("company_id"), 10, 64)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	// 管理会社のユーザであることの認証処理
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
			return &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			}
		}
		if !authUser.IsSuperUser() {
			return &echo.HTTPError{
				Code: http.StatusForbidden,
			}
		}
	}

	company, aerr := h.companyUsecase.Get(domain.CompanyIdentifier(id))
	if aerr != nil {
		return aerr.HTTPError()
	}

	res := model.UnmarshalCompany(company)

	return c.JSON(http.StatusOK, res)
}

// CreateCompany
//
//	@Summary		企業の作成
//	@Description	企業を作成する。管理会社の管理者のみ実行可能。
//	@Tags			company
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		request.CompanyCreate	false	"企業作成用リクエスト"
//	@Success		200				{object}	integer
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		500
//	@Router			/company/create [post]
func (h *companyHandler) Create(c echo.Context) error {
	// 管理会社の管理者であることの認証処理
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
		if !authUser.IsSuperAdmin() {
			return &echo.HTTPError{
				Code: http.StatusForbidden,
			}
		}
	}

	var req *request.CompanyCreate
	if err := c.Bind(&req); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	id, aerr := h.companyUsecase.Create(req.Name)
	if aerr != nil {
		return aerr.HTTPError()
	}

	return c.JSON(http.StatusCreated, uint64(*id))
}

// UpdateCompany
//
//	@Summary		企業の更新
//	@Description	企業を作成する。管理会社の管理者のみ実行可能。
//	@Tags			company
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string					true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			company_id		path	int						false	"企業ID"
//	@Param			body			body	request.CompanyUpdate	false	"企業更新用リクエスト"
//	@Success		200				
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		500
//	@Router			/company/{company_id}/update [put]
func (h *companyHandler) Update(c echo.Context) error {
	// 管理会社の管理者であることの認証処理
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
		if !authUser.IsSuperAdmin() {
			return &echo.HTTPError{
				Code: http.StatusForbidden,
			}
		}
	}

	id, err := strconv.ParseUint(c.Param("company_id"), 10, 64)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	var req *request.CompanyUpdate
	if err := c.Bind(&req); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	aerr := h.companyUsecase.Update(domain.CompanyIdentifier(id), req.Name)
	if aerr != nil {
		return aerr.HTTPError()
	}

	return c.NoContent(http.StatusOK)
}
