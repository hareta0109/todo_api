package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"todo_api/internal/adapter/inbound/http/request"
	"todo_api/internal/adapter/inbound/http/response"
	domain "todo_api/internal/domain/model"
	"todo_api/internal/usecase"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

const jwtDuration = time.Hour * 24

type AuthHandler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Login(c echo.Context) error
}

type authHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(
	authUsecase usecase.AuthUsecase,
) AuthHandler {
	return &authHandler{
		authUsecase,
	}
}

// CreateUser
//
//	@Summary		ユーザの登録
//	@Description	ユーザを登録する。管理会社の管理者と企業の管理者に実行可能。
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			company_id		path		int					false	"企業ID"
//	@Param			body			body		request.AuthCreate	false	"ユーザ作成用リクエスト"
//	@Success		200				{object}	integer				"登録されたユーザID"
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		500
//	@Router			/company/{company_id}/user/create [post]
func (h *authHandler) Create(c echo.Context) error {
	companyID, err := strconv.ParseUint(c.Param("company_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// 会社の管理者権限をもつことの認証処理
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
		if !authUser.CanAdminAction(domain.CompanyIdentifier(companyID)) {
			return &echo.HTTPError{
				Code: http.StatusForbidden,
			}
		}
	}

	var req *request.AuthCreate
	if err := c.Bind(&req); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	params, aerr := request.MarshalAuthCreateParams(companyID, req)
	if aerr != nil {
		return aerr.HTTPError()
	}

	id, aerr := h.authUsecase.Create(*params)
	if aerr != nil {
		return aerr.HTTPError()
	}

	return c.JSON(http.StatusCreated, uint64(*id))
}

// UpdateUser
//
//	@Summary		ユーザの更新
//	@Description	ユーザの情報を更新する。管理会社の管理者と企業の管理者に実行可能。
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string				true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			company_id		path	int					false	"企業ID"
//	@Param			user_id			path	int					false	"ユーザID"
//	@Param			body			body	request.AuthUpdate	false	"ユーザ更新用リクエスト"
//	@Success		200				
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		500
//	@Router			/company/{company_id}/user/{user_id}/update [post]
func (h *authHandler) Update(c echo.Context) error {
	companyID, err := strconv.ParseUint(c.Param("company_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// 会社の管理者権限をもつことの認証処理
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
		if !authUser.CanAdminAction(domain.CompanyIdentifier(companyID)) {
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

	var req *request.AuthUpdate
	if err := c.Bind(&req); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	params, aerr := request.MarshalAuthUpdateParams(companyID, req)
	if aerr != nil {
		return aerr.HTTPError()
	}

	if aerr = h.authUsecase.Update(domain.UserIdentifier(id), *params); aerr != nil {
		return aerr.HTTPError()
	}

	return c.NoContent(http.StatusOK)
}

// LoginUser
//
//	@Summary		ログイン
//	@Description	ログインを行い、トークンを発行する。
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.AuthLogin	false	"ログイン用リクエスト"
//	@Success		200		{object}	response.AuthLogin
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		500
//	@Router			/auth/login [post]
func (h *authHandler) Login(c echo.Context) error {
	var req *request.AuthLogin
	if err := c.Bind(&req); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	params := request.MarshalAuthLoginParams(req)
	if params == nil {
		return &echo.HTTPError{
			Code: http.StatusBadRequest,
		}
	}

	auth, aerr := h.authUsecase.Login(*params)
	if aerr != nil {
		return aerr.HTTPError()
	}

	t, err := generateToken(uint64(auth.ID))
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	res := &response.AuthLogin{
		Token:     t,
		CompanyID: uint64(auth.Company.ID),
		UserID:    uint64(auth.ID),
	}

	return c.JSON(http.StatusOK, res)
}

var signingKey = []byte("secret") // TODO: 生成方法を検討する

var Config = echojwt.Config{
	SigningKey: signingKey,
}

func generateToken(userID uint64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": fmt.Sprintf("%d", userID),
		"exp":     jwt.NewNumericDate(time.Now().Add(jwtDuration)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return t, nil
}

func verifyJwtToken(c echo.Context) (uint64, error) {
	user := c.Get("user").(*jwt.Token)
	claims, _ := user.Claims.(jwt.MapClaims)
	strID := claims["user_id"].(string)

	id, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
