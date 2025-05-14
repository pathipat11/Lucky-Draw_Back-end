package prize

import (
	"app/app/helper"
	"app/app/request"
	"app/app/response"
	"app/internal/logger"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (ctl *Controller) Create(ctx *gin.Context) {
	var body request.CreatePrize

	// Bind form fields
	body.Name = ctx.PostForm("name")
	body.Quantity, _ = strconv.ParseInt(ctx.PostForm("quantity"), 10, 64)
	body.RoomID = ctx.PostForm("room_id")

	// ตรวจสอบว่ามีไฟล์ image แนบมาหรือไม่
	file, err := ctx.FormFile("image")
	if err == nil {
		// มีไฟล์ -> ทำการเปิดและอัปโหลด
		src, err := file.Open()
		if err != nil {
			response.InternalError(ctx, "failed to open image")
			return
		}
		defer src.Close()

		imageURL, err := helper.UploadToCloudinary(ctx.Request.Context(), src, "prize_images")
		if err != nil {
			response.InternalError(ctx, "failed to upload image")
			return
		}
		body.ImageURL = imageURL
	} else {
		// ไม่มีไฟล์ -> ไม่ต้องแนบรูป
		body.ImageURL = ""
	}

	data, mserr, err := ctl.Service.Create(ctx, body)
	if err != nil {
		ms := "internal server error"
		if mserr {
			ms = err.Error()
		}
		response.InternalError(ctx, ms)
		return
	}

	response.Success(ctx, data)
}

func (ctl *Controller) Update(ctx *gin.Context) {
	ID := request.GetByIDPrize{}
	if err := ctx.BindUri(&ID); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	body := request.UpdatePrize{}
	body.Name = ctx.PostForm("name")
	body.Quantity, _ = strconv.ParseInt(ctx.PostForm("quantity"), 10, 64)
	body.RoomID = ctx.PostForm("room_id")

	// ถ้ามีไฟล์มาให้ อัปโหลดใหม่
	file, err := ctx.FormFile("image")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			response.InternalError(ctx, "failed to open image")
			return
		}
		defer src.Close()

		imageURL, err := helper.UploadToCloudinary(ctx.Request.Context(), src, "prize_images")
		if err != nil {
			response.InternalError(ctx, "failed to upload image")
			return
		}
		body.ImageURL = imageURL
	} else {
		// ไม่มีไฟล์ใหม่ -> ดึงค่ารูปเดิมจาก database มาก่อน
		oldData, err := ctl.Service.Get(ctx, ID)
		if err != nil {
			response.InternalError(ctx, "failed to get existing prize data")
			return
		}
		body.ImageURL = oldData.ImageURL
	}

	_, mserr, err := ctl.Service.Update(ctx, body, ID)
	if err != nil {
		ms := "internal server error"
		if mserr {
			ms = err.Error()
		}
		response.InternalError(ctx, ms)
		return
	}

	response.Success(ctx, nil)
}

func (ctl *Controller) List(ctx *gin.Context) {
	req := request.ListPrize{}
	if err := ctx.Bind(&req); err != nil {
		logger.Errf(err.Error())
		response.BadRequest(ctx, err.Error())
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.Size == 0 {
		req.Size = 10
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
	ID := request.GetByIDPrize{}
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
	ID := request.GetByIDPrize{}
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
