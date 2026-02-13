package members

import (
	"phakram/app/modules/auth"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MemberFileURIRequest struct {
	MemberID string `uri:"id"`
	FileID   string `uri:"file_row_id"`
}

type CreateMemberFileControllerRequest struct {
	FileID string `json:"file_id"`
}

type ListMemberFilesControllerRequest struct {
	base.RequestPaginate
}

func (c *Controller) ListMemberFilesController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.file.list.start`)

	if !auth.GetIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	memberID, ok := c.parseMemberID(ctx)
	if !ok {
		return
	}

	var req ListMemberFilesControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, page, err := c.svc.ListMemberFilesService(ctx.Request.Context(), &entitiesdto.ListMemberFilesRequest{
		RequestPaginate: req.RequestPaginate,
		MemberID:        memberID,
	})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.file.list.success`)
	base.Paginate(ctx, data, page)
}

func (c *Controller) CreateMemberFileController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.file.create.start`)

	if !auth.GetIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	memberID, ok := c.parseMemberID(ctx)
	if !ok {
		return
	}

	var req CreateMemberFileControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	fileID, err := uuid.Parse(req.FileID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	actionBy := getActionBy(ctx)
	if err := c.svc.CreateMemberFileService(ctx.Request.Context(), memberID, fileID, actionBy); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.file.create.success`)
	base.Success(ctx, nil)
}

func (c *Controller) InfoMemberFileController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.file.info.start`)

	if !auth.GetIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	memberID, rowID, ok := c.parseMemberFileURI(ctx)
	if !ok {
		return
	}

	data, err := c.svc.InfoMemberFileService(ctx.Request.Context(), memberID, rowID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.file.info.success`)
	base.Success(ctx, data)
}

func (c *Controller) UpdateMemberFileController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.file.update.start`)

	if !auth.GetIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	memberID, rowID, ok := c.parseMemberFileURI(ctx)
	if !ok {
		return
	}

	var req CreateMemberFileControllerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	fileID, err := uuid.Parse(req.FileID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	actionBy := getActionBy(ctx)
	if err := c.svc.UpdateMemberFileService(ctx.Request.Context(), memberID, rowID, fileID, actionBy); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.file.update.success`)
	base.Success(ctx, nil)
}

func (c *Controller) DeleteMemberFileController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`members.ctl.file.delete.start`)

	if !auth.GetIsAdmin(ctx) {
		base.Forbidden(ctx, i18n.Forbidden, nil)
		return
	}

	memberID, rowID, ok := c.parseMemberFileURI(ctx)
	if !ok {
		return
	}

	actionBy := getActionBy(ctx)
	if err := c.svc.DeleteMemberFileService(ctx.Request.Context(), memberID, rowID, actionBy); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`members.ctl.file.delete.success`)
	base.Success(ctx, nil)
}

func (c *Controller) parseMemberFileURI(ctx *gin.Context) (uuid.UUID, uuid.UUID, bool) {
	var uri MemberFileURIRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}

	memberID, err := uuid.Parse(uri.MemberID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}
	rowID, err := uuid.Parse(uri.FileID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}

	return memberID, rowID, true
}
