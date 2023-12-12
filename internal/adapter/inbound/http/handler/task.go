package handler

import (
	"net/http"
	"strconv"
	"strings"
	"todo_api/internal/adapter/inbound/http/model"
	"todo_api/internal/adapter/inbound/http/request"
	domain "todo_api/internal/domain/model"
	"todo_api/internal/usecase"

	"github.com/labstack/echo/v4"
)

type TaskHandler interface {
	Find(c echo.Context) error
	ListByAssignedUserID(c echo.Context) error
	ListByCompanyID(c echo.Context) error
	Create(c echo.Context) error
	Update(c echo.Context) error
	UpdateStatus(c echo.Context) error
}

type taskHandler struct {
	authUsecase usecase.AuthUsecase
	taskUsecase usecase.TaskUsecase
}

func NewTaskHandler(
	authUsecase usecase.AuthUsecase,
	taskUsecase usecase.TaskUsecase,
) TaskHandler {
	return &taskHandler{
		authUsecase,
		taskUsecase,
	}
}

// FindTask
//
//	@Summary		タスクの取得
//	@Description	閲覧可能なタスクの情報をIDから取得する。
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			company_id		path		int		false	"企業ID"
//	@Param			task_id			path		int		false	"タスクID"
//	@Success		200				{object}	model.Task
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		500
//	@Router			/company/{company_id}/task/{task_id} [get]
func (h *taskHandler) Find(c echo.Context) error {
	companyID, err := strconv.ParseUint(c.Param("company_id"), 10, 64)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	// 会社の閲覧権限を持つことの認証処理
	var authUserID uint64
	{
		authUserID, err = verifyJwtToken(c)
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

	id, err := strconv.ParseUint(c.Param("task_id"), 10, 64)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	task, aerr := h.taskUsecase.Find(domain.UserIdentifier(authUserID), domain.TaskIdentifier(id))
	if aerr != nil {
		return aerr.HTTPError()
	}

	res := model.UnmarshalTask(task)

	return c.JSON(http.StatusOK, res)
}

// ListTaskByAssignedUserID
//
//	@Summary		ユーザに割り当てられたタスク一覧の取得
//	@Description	閲覧可能なタスクの一覧を割り当てユーザIDから取得する。
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Param			Authorization		header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			company_id			path	int		false	"企業ID"
//	@Param			assigned_user_id	path	int		false	"ユーザID"
//	@Success		200					{array}	model.Task
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		500
//	@Router			/company/{company_id}/task/list_by_assigned_user_id/{assigned_user_id} [get]
func (h *taskHandler) ListByAssignedUserID(c echo.Context) error {
	companyID, err := strconv.ParseUint(c.Param("company_id"), 10, 64)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	// 会社の閲覧権限を持つことの認証処理
	var authUserID uint64
	{
		authUserID, err = verifyJwtToken(c)
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

	assignedUserID, err := strconv.ParseUint(c.Param("assigned_user_id"), 10, 64)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	tasks, aerr := h.taskUsecase.ListByAssignedUserID(domain.UserIdentifier(authUserID), domain.UserIdentifier(assignedUserID))
	if aerr != nil {
		return aerr.HTTPError()
	}

	var res []*model.Task
	for _, task := range tasks {
		res = append(res, model.UnmarshalTask(task))
	}

	return c.JSON(http.StatusOK, res)
}

// ListTaskByCompanyID
//
//	@Summary		企業のタスク一覧の取得
//	@Description	閲覧可能なタスクの一覧を企業IDから取得する。
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			company_id		path	int		false	"企業ID"
//	@Success		200				{array}	model.Task
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		500
//	@Router			/company/{company_id}/task/list [get]
func (h *taskHandler) ListByCompanyID(c echo.Context) error {
	companyID, err := strconv.ParseUint(c.Param("company_id"), 10, 64)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	// 会社の閲覧権限を持つことの認証処理
	var authUserID uint64
	{
		authUserID, err = verifyJwtToken(c)
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

	tasks, aerr := h.taskUsecase.ListByCompanyID(domain.UserIdentifier(authUserID), domain.CompanyIdentifier(companyID))
	if aerr != nil {
		return aerr.HTTPError()
	}

	var res []*model.Task
	for _, task := range tasks {
		res = append(res, model.UnmarshalTask(task))
	}

	return c.JSON(http.StatusOK, res)
}

// CreateTask
//
//	@Summary		タスクの作成
//	@Description	タスクを作成する。編集者のみ可能。
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			company_id		path		int					false	"企業ID"
//	@Param			body			body		request.TaskCreate	false	"タスク作成用リクエスト"
//	@Success		200				{object}	integer
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		500
//	@Router			/company/{company_id}/task/create [post]
func (h *taskHandler) Create(c echo.Context) error {
	companyID, err := strconv.ParseUint(c.Param("company_id"), 10, 64)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	// 会社の編集権限を持つことの認証処理
	var authUserID uint64
	{
		authUserID, err = verifyJwtToken(c)
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
		if !authUser.CanEdit(domain.CompanyIdentifier(companyID)) {
			return &echo.HTTPError{
				Code: http.StatusForbidden,
			}
		}
	}

	var req *request.TaskCreate
	if err := c.Bind(&req); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	params, aerr := request.MarshalTaskCreateParams(authUserID, req)
	if aerr != nil {
		return aerr.HTTPError()
	}

	id, aerr := h.taskUsecase.Create(*params)
	if aerr != nil {
		return aerr.HTTPError()
	}

	return c.JSON(http.StatusCreated, uint64(*id))
}

// UpdateTask
//
//	@Summary		タスクの更新
//	@Description	タスクの情報を更新する。編集者のみ可能。
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			company_id		path		int					false	"企業ID"
//	@Param			task_id			path		int					false	"タスクID"
//	@Param			body			body		request.TaskUpdate	false	"タスク更新用リクエスト"
//	@Success		200				{object}	integer
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		500
//	@Router			/company/{company_id}/task/{task_id}/update [put]
func (h *taskHandler) Update(c echo.Context) error {
	companyID, err := strconv.ParseUint(c.Param("company_id"), 10, 64)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	// 会社の編集権限を持つことの認証処理
	var authUserID uint64
	{
		authUserID, err = verifyJwtToken(c)
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
		if !authUser.CanEdit(domain.CompanyIdentifier(companyID)) {
			return &echo.HTTPError{
				Code: http.StatusForbidden,
			}
		}
	}

	id, err := strconv.ParseUint(c.Param("task_id"), 10, 64)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	var req *request.TaskUpdate
	if err := c.Bind(&req); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	params, aerr := request.MarshalTaskUpdateParams(authUserID, req)
	if aerr != nil {
		return aerr.HTTPError()
	}

	aerr = h.taskUsecase.Update(domain.TaskIdentifier(id), *params)
	if aerr != nil {
		return aerr.HTTPError()
	}

	return c.NoContent(http.StatusOK)
}

// UpdateTaskStatus
//
//	@Summary		タスクステータスの更新
//	@Description	タスクのステータスを指定した状態へと更新する。編集者のみ可能。
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			company_id		path	int		false	"企業ID"
//	@Param			task_id			path	int		false	"タスクID"
//	@Param			task_status		path	string	false	"タスクステータス"
//	@Success		200
//	@Failure		400
//	@Failure		401
//	@Failure		403
//	@Failure		500
//	@Router			/company/{company_id}/task/{task_id}/status/{task_status} [put]
func (h *taskHandler) UpdateStatus(c echo.Context) error {
	companyID, err := strconv.ParseUint(c.Param("company_id"), 10, 64)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	// 会社の編集権限を持つことの認証処理
	var authUserID uint64
	{
		authUserID, err = verifyJwtToken(c)
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
		if !authUser.CanEdit(domain.CompanyIdentifier(companyID)) {
			return &echo.HTTPError{
				Code: http.StatusForbidden,
			}
		}
	}

	id, err := strconv.ParseUint(c.Param("task_id"), 10, 64)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	status, aerr := request.MarshalTaskStatus(strings.ToUpper(c.Param("status")))
	if aerr != nil {
		return aerr.HTTPError()
	}

	if aerr = h.taskUsecase.UpdateStatus(domain.TaskIdentifier(id), *status); aerr != nil {
		return aerr.HTTPError()
	}

	return c.NoContent(http.StatusOK)
}
