package banks

import (
	"log/slog"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateBankControllerRequestUri struct {
	ID string `uri:"id"`
}

type UpdateBankController struct {
	NameTh    string `json:"name_th"`
	NameAbbTh string `json:"name_abb_th"`
	NameEn    string `json:"name_en"`
	NameAbbEn string `json:"name_abb_en"`
	IsActive  bool   `json:"is_active"`
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)

	var reqUri UpdateBankControllerRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`banks.ctl.update.request_uri`)

	id, err := uuid.Parse(reqUri.ID)
	if err != nil {
		log.With(slog.Any(`body`, reqUri)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req UpdateBankController
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	span.AddEvent(`banks.ctl.update.request_body`)

	if err := c.svc.UpdateService(ctx, id, &UpdateBankService{
		NameTh:    req.NameTh,
		NameAbbTh: req.NameAbbTh,
		NameEn:    req.NameEn,
		NameAbbEn: req.NameAbbEn,
	}); err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`banks.ctl.update.callsvc`)

	base.Success(ctx, nil)
}

func (c *Controller) BanksUpdate(ctx *gin.Context) {
	c.UpdateController(ctx)
}
