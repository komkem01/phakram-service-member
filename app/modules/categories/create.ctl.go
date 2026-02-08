package categories

import (
	"log/slog"
	authmod "phakram/app/modules/auth"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateCategoryController struct {
	ParentID string `json:"parent_id"`
	NameTh   string `json:"name_th"`
	NameEn   string `json:"name_en"`
	IsActive bool   `json:"is_active"`
}

func (c *Controller) CreateCategoryController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	span.AddEvent(`categories.ctl.create.start`)

	var req CreateCategoryController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`categories.ctl.create.request`)

	var parentID *uuid.UUID
	if req.ParentID != "" {
		id, err := uuid.Parse(req.ParentID)
		if err != nil {
			log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		parentID = &id
	}
	memberID, ok := authmod.GetMemberID(ctx)
	if !ok {
		base.Unauthorized(ctx, i18n.Unauthorized, nil)
		return
	}

	if err := c.svc.CreateCategoryService(ctx.Request.Context(), &CreateCategoryService{
		ParentID: parentID,
		NameTh:   req.NameTh,
		NameEn:   req.NameEn,
		IsActive: req.IsActive,
		MemberID: memberID,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`categories.ctl.create.success`)
	base.Success(ctx, nil)
}
