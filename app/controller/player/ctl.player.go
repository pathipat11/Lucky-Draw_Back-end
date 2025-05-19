package player

import (
	"app/app/request"
	"app/app/response"
	"app/internal/logger"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Create(ctx *gin.Context) {
	body := request.CreatePlayer{}

	if err := ctx.Bind(&body); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	_, mserr, err := ctl.Service.Create(ctx, body)
	if err != nil {
		ms := "internal server error"
		if mserr {
			ms = err.Error()
		}
		logger.Errf(err.Error())
		response.InternalError(ctx, ms)
		return
	}

	response.Success(ctx, nil)
}

func (ctl *Controller) Update(ctx *gin.Context) {
	ID := request.GetByIDPlayer{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	body := request.UpdatePlayer{}
	if err := ctx.Bind(&body); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	_, mserr, err := ctl.Service.Update(ctx, body, ID)
	if err != nil {
		ms := "internal server error"
		if mserr {
			ms = err.Error()
		}
		logger.Errf(err.Error())
		response.InternalError(ctx, ms)
		return
	}

	response.Success(ctx, nil)
}

func (ctl *Controller) List(ctx *gin.Context) {
	req := request.ListPlayer{}
	if err := ctx.Bind(&req); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.Page == 0 {
		req.Page = 10
	}

	if req.OrderBy == "" {
		req.OrderBy = "asc"
	}

	if req.SortBy == "" {
		req.SortBy = "created_at"
	}

	data, total, err := ctl.Service.List(ctx, req)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalError(ctx, err.Error())
		return
	}
	response.SuccessWithPaginate(ctx, data, req.Size, req.Page, total)
}

func (ctl *Controller) Get(ctx *gin.Context) {
	ID := request.GetByIDPlayer{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	data, err := ctl.Service.Get(ctx, ID)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalError(ctx, err.Error())
		return
	}
	response.Success(ctx, data)
}

func (ctl *Controller) Delete(ctx *gin.Context) {
	ID := request.GetByIDPlayer{}
	if err := ctx.BindUri(&ID); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	err := ctl.Service.Delete(ctx, ID)
	if err != nil {
		logger.Errf(err.Error())
		response.InternalError(ctx, err.Error())
		return
	}
	response.Success(ctx, nil)
}

// new fucntion

func (ctl *Controller) ImportPlayerFile(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		logger.Errf("failed to get form file: %v", err)
		response.BadRequest(ctx, "invalid file upload")
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		logger.Errf("failed to open uploaded file: %v", err)
		response.InternalError(ctx, "unable to read file")
		return
	}
	defer src.Close()

	roomID := ctx.PostForm("room_id")
	if roomID == "" {
		response.BadRequest(ctx, "room_id is required")
		return
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	switch ext {
	case ".csv":
		err = ctl.Service.ImportPlayersFromCSV(ctx, src, roomID)
	case ".xlsx":
		err = ctl.Service.ImportPlayersFromXLSX(ctx, src, roomID)
	default:
		response.BadRequest(ctx, "unsupported file type: only .csv and .xlsx allowed")
		return
	}

	if err != nil {
		logger.Errf("failed to import player data: %v", err)
		response.InternalError(ctx, "import failed")
		return
	}

	response.Success(ctx, "import successful")
}
