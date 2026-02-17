package productstocks

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"phakram/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductStockRequestUri struct {
	ID string `uri:"id"`
}

type ProductStockRequest struct {
	StockAmount   int    `json:"stock_amount"`
	Remaining     int    `json:"remaining"`
	Action        string `json:"action"`
	AdjustmentQty int    `json:"adjustment_qty"`
}

type ListProductStockControllerRequest struct {
	base.RequestPaginate
}

func (c *Controller) ListController(ctx *gin.Context) {
	span, _ := utils.LogSpanFromGin(ctx)
	span.AddEvent(`product_stocks.ctl.list.start`)

	var req ListProductStockControllerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	data, page, err := c.svc.ListService(ctx.Request.Context(), &ListProductStocksServiceRequest{RequestPaginate: req.RequestPaginate})
	if err != nil {
		base.HandleError(ctx, err)
		return
	}

	span.AddEvent(`product_stocks.ctl.list.success`)
	base.Paginate(ctx, data, page)
}

func applyStockAdjustment(currentStock, currentRemaining int, action string, adjustmentQty int) (int, int, error) {
	if adjustmentQty <= 0 {
		return 0, 0, fmt.Errorf("adjustment_qty must be greater than 0")
	}

	switch action {
	case "increase":
		return currentStock + adjustmentQty, currentRemaining + adjustmentQty, nil
	case "decrease":
		if adjustmentQty > currentRemaining {
			return 0, 0, fmt.Errorf("adjustment_qty cannot exceed current remaining")
		}
		if adjustmentQty > currentStock {
			return 0, 0, fmt.Errorf("adjustment_qty cannot exceed current stock_amount")
		}
		return currentStock - adjustmentQty, currentRemaining - adjustmentQty, nil
	default:
		return 0, 0, fmt.Errorf("invalid action")
	}
}

func (c *Controller) InfoController(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	var req ProductStockRequestUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.With(slog.Any(`body`, req)).Errf(`internal: %s`, err)
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	productID, err := uuid.Parse(req.ID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	data, err := c.svc.GetByProductID(ctx, productID)
	if err != nil {
		base.HandleError(ctx, err)
		return
	}
	span.AddEvent(`product_stocks.ctl.info.success`)
	base.Success(ctx, data)
}

func (c *Controller) CreateController(ctx *gin.Context) {
	var reqUri ProductStockRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	productID, err := uuid.Parse(reqUri.ID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req ProductStockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	stockAmount := req.StockAmount
	remaining := req.Remaining
	if req.Action != "" {
		if req.Action == "decrease" {
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		calculatedStockAmount, calculatedRemaining, err := applyStockAdjustment(0, 0, req.Action, req.AdjustmentQty)
		if err != nil {
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		stockAmount = calculatedStockAmount
		remaining = calculatedRemaining
	}

	if stockAmount < 0 || remaining < 0 || remaining > stockAmount {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	payload := &ent.ProductStockEntity{
		StockAmount: stockAmount,
		Remaining:   remaining,
	}
	if err := c.svc.CreateByProductID(ctx, productID, payload); err != nil {
		base.HandleError(ctx, err)
		return
	}
	base.Success(ctx, nil)
}

func (c *Controller) UpdateController(ctx *gin.Context) {
	var reqUri ProductStockRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	productID, err := uuid.Parse(reqUri.ID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	var req ProductStockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	stockAmount := req.StockAmount
	remaining := req.Remaining
	if req.Action != "" {
		currentStock, err := c.svc.GetByProductID(ctx, productID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				base.BadRequest(ctx, i18n.BadRequest, nil)
				return
			}
			base.HandleError(ctx, err)
			return
		}
		calculatedStockAmount, calculatedRemaining, err := applyStockAdjustment(currentStock.StockAmount, currentStock.Remaining, req.Action, req.AdjustmentQty)
		if err != nil {
			base.BadRequest(ctx, i18n.BadRequest, nil)
			return
		}
		stockAmount = calculatedStockAmount
		remaining = calculatedRemaining
	}

	if stockAmount < 0 || remaining < 0 || remaining > stockAmount {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}

	payload := &ent.ProductStockEntity{
		StockAmount: stockAmount,
		Remaining:   remaining,
	}
	if err := c.svc.UpdateByProductID(ctx, productID, payload); err != nil {
		base.HandleError(ctx, err)
		return
	}
	base.Success(ctx, nil)
}

func (c *Controller) DeleteController(ctx *gin.Context) {
	var reqUri ProductStockRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	productID, err := uuid.Parse(reqUri.ID)
	if err != nil {
		base.BadRequest(ctx, i18n.BadRequest, nil)
		return
	}
	if err := c.svc.DeleteByProductID(ctx, productID); err != nil {
		base.HandleError(ctx, err)
		return
	}
	base.Success(ctx, nil)
}
