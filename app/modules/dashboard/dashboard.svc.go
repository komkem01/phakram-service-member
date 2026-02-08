package dashboard

import (
	"context"
	"time"

	"phakram/app/modules/entities/ent"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type SummaryResponse struct {
	Cards        []StatCard        `json:"cards"`
	SalesChart   []SalesChartItem  `json:"sales_chart"`
	LowStock     []LowStockItem    `json:"low_stock"`
	RecentOrders []RecentOrderItem `json:"recent_orders"`
	TopProducts  []TopProductItem  `json:"top_products"`
}

type StatCard struct {
	Key           string          `json:"key"`
	Value         decimal.Decimal `json:"value"`
	ChangePercent decimal.Decimal `json:"change_percent"`
}

type SalesChartItem struct {
	Month string          `json:"month"`
	Value decimal.Decimal `json:"value"`
}

type LowStockItem struct {
	ProductID   uuid.UUID `json:"product_id"`
	Name        string    `json:"name"`
	Remaining   int       `json:"remaining"`
	StockAmount int       `json:"stock_amount"`
	Status      string    `json:"status"`
	Progress    int       `json:"progress"`
}

type RecentOrderItem struct {
	ID        uuid.UUID       `json:"id"`
	OrderNo   string          `json:"order_no"`
	MemberID  uuid.UUID       `json:"member_id"`
	Status    string          `json:"status"`
	NetAmount decimal.Decimal `json:"net_amount"`
	CreatedAt string          `json:"created_at"`
}

type TopProductItem struct {
	ProductID uuid.UUID       `json:"product_id"`
	Name      string          `json:"name"`
	Quantity  int64           `json:"quantity"`
	Revenue   decimal.Decimal `json:"revenue"`
	Remaining int             `json:"remaining"`
}

func (s *Service) SummaryService(ctx context.Context, rangeKey string) (*SummaryResponse, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`dashboard.svc.summary.start`)

	start, end, prevStart, prevEnd := rangeWindow(rangeKey)

	salesUnits, err := s.sumOrderItemQuantity(ctx, start, end)
	if err != nil {
		return nil, err
	}
	salesUnitsPrev, err := s.sumOrderItemQuantity(ctx, prevStart, prevEnd)
	if err != nil {
		return nil, err
	}

	ordersCount, err := s.countOrders(ctx, start, end)
	if err != nil {
		return nil, err
	}
	ordersCountPrev, err := s.countOrders(ctx, prevStart, prevEnd)
	if err != nil {
		return nil, err
	}

	newMembers, err := s.countNewMembers(ctx, start, end)
	if err != nil {
		return nil, err
	}
	newMembersPrev, err := s.countNewMembers(ctx, prevStart, prevEnd)
	if err != nil {
		return nil, err
	}

	revenue, err := s.sumOrderNetAmount(ctx, start, end)
	if err != nil {
		return nil, err
	}
	revenuePrev, err := s.sumOrderNetAmount(ctx, prevStart, prevEnd)
	if err != nil {
		return nil, err
	}

	cards := []StatCard{
		{
			Key:           "sales",
			Value:         decimal.NewFromInt(salesUnits),
			ChangePercent: percentChangeDecimal(decimal.NewFromInt(salesUnitsPrev), decimal.NewFromInt(salesUnits)),
		},
		{
			Key:           "orders",
			Value:         decimal.NewFromInt(ordersCount),
			ChangePercent: percentChangeDecimal(decimal.NewFromInt(ordersCountPrev), decimal.NewFromInt(ordersCount)),
		},
		{
			Key:           "new_members",
			Value:         decimal.NewFromInt(newMembers),
			ChangePercent: percentChangeDecimal(decimal.NewFromInt(newMembersPrev), decimal.NewFromInt(newMembers)),
		},
		{
			Key:           "revenue",
			Value:         revenue,
			ChangePercent: percentChangeDecimal(revenuePrev, revenue),
		},
	}

	salesChart, err := s.salesChart(ctx)
	if err != nil {
		return nil, err
	}

	lowStock, err := s.lowStock(ctx)
	if err != nil {
		return nil, err
	}

	recentOrders, err := s.recentOrders(ctx)
	if err != nil {
		return nil, err
	}

	topProducts, err := s.topProducts(ctx, start, end)
	if err != nil {
		return nil, err
	}

	span.AddEvent(`dashboard.svc.summary.success`)
	return &SummaryResponse{
		Cards:        cards,
		SalesChart:   salesChart,
		LowStock:     lowStock,
		RecentOrders: recentOrders,
		TopProducts:  topProducts,
	}, nil
}

func rangeWindow(rangeKey string) (time.Time, time.Time, time.Time, time.Time) {
	now := time.Now()
	end := now
	var start time.Time
	switch rangeKey {
	case "":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "today":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "week":
		start = now.AddDate(0, 0, -6)
	default:
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}
	duration := end.Sub(start)
	prevEnd := start
	prevStart := start.Add(-duration)
	return start, end, prevStart, prevEnd
}

func percentChangeDecimal(prev decimal.Decimal, current decimal.Decimal) decimal.Decimal {
	if prev.IsZero() {
		if current.IsZero() {
			return decimal.Zero
		}
		return decimal.NewFromInt(100)
	}
	return current.Sub(prev).Div(prev).Mul(decimal.NewFromInt(100))
}

func (s *Service) countOrders(ctx context.Context, start time.Time, end time.Time) (int64, error) {
	var total int64
	err := s.bunDB.DB().NewSelect().
		TableExpr("orders AS o").
		ColumnExpr("COUNT(*)").
		Where("o.created_at >= ? AND o.created_at < ?", start, end).
		Scan(ctx, &total)
	return total, err
}

func (s *Service) sumOrderNetAmount(ctx context.Context, start time.Time, end time.Time) (decimal.Decimal, error) {
	var total decimal.Decimal
	err := s.bunDB.DB().NewSelect().
		TableExpr("orders AS o").
		ColumnExpr("COALESCE(SUM(o.net_amount), 0)").
		Where("o.created_at >= ? AND o.created_at < ?", start, end).
		Scan(ctx, &total)
	return total, err
}

func (s *Service) sumOrderItemQuantity(ctx context.Context, start time.Time, end time.Time) (int64, error) {
	var total int64
	err := s.bunDB.DB().NewSelect().
		TableExpr("order_items AS oi").
		ColumnExpr("COALESCE(SUM(oi.quantity), 0)").
		Join("JOIN orders AS o ON o.id = oi.order_id").
		Where("o.created_at >= ? AND o.created_at < ?", start, end).
		Scan(ctx, &total)
	return total, err
}

func (s *Service) countNewMembers(ctx context.Context, start time.Time, end time.Time) (int64, error) {
	var total int64
	err := s.bunDB.DB().NewSelect().
		TableExpr("members AS m").
		ColumnExpr("COUNT(*)").
		Where("m.created_at >= ? AND m.created_at < ?", start, end).
		Scan(ctx, &total)
	return total, err
}

func (s *Service) salesChart(ctx context.Context) ([]SalesChartItem, error) {
	startMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location()).AddDate(0, -5, 0)
	endMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 1, 0)

	type row struct {
		Month time.Time       `bun:"month"`
		Total decimal.Decimal `bun:"total"`
	}
	var rows []row
	if err := s.bunDB.DB().NewSelect().
		TableExpr("orders AS o").
		ColumnExpr("DATE_TRUNC('month', o.created_at) AS month").
		ColumnExpr("COALESCE(SUM(o.net_amount), 0) AS total").
		Where("o.created_at >= ? AND o.created_at < ?", startMonth, endMonth).
		Group("month").
		OrderExpr("month ASC").
		Scan(ctx, &rows); err != nil {
		return nil, err
	}

	index := make(map[string]decimal.Decimal)
	for _, r := range rows {
		key := r.Month.Format("2006-01")
		index[key] = r.Total
	}

	var out []SalesChartItem
	for i := 0; i < 6; i++ {
		m := startMonth.AddDate(0, i, 0)
		key := m.Format("2006-01")
		total := index[key]
		out = append(out, SalesChartItem{
			Month: key,
			Value: total,
		})
	}
	return out, nil
}

func (s *Service) lowStock(ctx context.Context) ([]LowStockItem, error) {
	type row struct {
		ProductID   uuid.UUID `bun:"product_id"`
		NameTh      string    `bun:"name_th"`
		NameEn      string    `bun:"name_en"`
		Remaining   int       `bun:"remaining"`
		StockAmount int       `bun:"stock_amount"`
	}
	var rows []row
	if err := s.bunDB.DB().NewSelect().
		TableExpr("product_stocks AS ps").
		ColumnExpr("ps.product_id").
		ColumnExpr("p.name_th").
		ColumnExpr("p.name_en").
		ColumnExpr("ps.remaining").
		ColumnExpr("ps.stock_amount").
		Join("JOIN products AS p ON p.id = ps.product_id").
		Where("ps.deleted_at IS NULL").
		OrderExpr("ps.remaining ASC").
		Limit(10).
		Scan(ctx, &rows); err != nil {
		return nil, err
	}

	result := make([]LowStockItem, 0, len(rows))
	for _, r := range rows {
		name := r.NameTh
		if name == "" {
			name = r.NameEn
		}
		ratio := 0.0
		if r.StockAmount > 0 {
			ratio = float64(r.Remaining) / float64(r.StockAmount)
		}
		status := "Low"
		if ratio <= 0.10 {
			status = "Critical"
		} else if ratio <= 0.25 {
			status = "Warning"
		}
		result = append(result, LowStockItem{
			ProductID:   r.ProductID,
			Name:        name,
			Remaining:   r.Remaining,
			StockAmount: r.StockAmount,
			Status:      status,
			Progress:    int(ratio * 100),
		})
	}
	return result, nil
}

func (s *Service) recentOrders(ctx context.Context) ([]RecentOrderItem, error) {
	var orders []ent.OrderEntity
	if err := s.bunDB.DB().NewSelect().
		Model(&orders).
		Order("created_at DESC").
		Limit(5).
		Scan(ctx); err != nil {
		return nil, err
	}

	result := make([]RecentOrderItem, 0, len(orders))
	for _, o := range orders {
		result = append(result, RecentOrderItem{
			ID:        o.ID,
			OrderNo:   o.OrderNo,
			MemberID:  o.MemberID,
			Status:    string(o.Status),
			NetAmount: o.NetAmount,
			CreatedAt: o.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return result, nil
}

func (s *Service) topProducts(ctx context.Context, start time.Time, end time.Time) ([]TopProductItem, error) {
	type row struct {
		ProductID uuid.UUID       `bun:"product_id"`
		NameTh    string          `bun:"name_th"`
		NameEn    string          `bun:"name_en"`
		Quantity  int64           `bun:"quantity"`
		Revenue   decimal.Decimal `bun:"revenue"`
		Remaining int             `bun:"remaining"`
	}
	var rows []row
	if err := s.bunDB.DB().NewSelect().
		TableExpr("products AS p").
		ColumnExpr("p.id AS product_id").
		ColumnExpr("p.name_th").
		ColumnExpr("p.name_en").
		ColumnExpr("COALESCE(SUM(oi.quantity), 0) AS quantity").
		ColumnExpr("COALESCE(SUM(oi.total_item_amount), 0) AS revenue").
		ColumnExpr("COALESCE(ps.remaining, 0) AS remaining").
		Join("LEFT JOIN order_items AS oi ON oi.product_id = p.id").
		Join("LEFT JOIN orders AS o ON o.id = oi.order_id").
		Join("LEFT JOIN product_stocks AS ps ON ps.product_id = p.id AND ps.deleted_at IS NULL").
		Where("o.created_at >= ? AND o.created_at < ?", start, end).
		GroupExpr("p.id, p.name_th, p.name_en, ps.remaining").
		OrderExpr("revenue DESC").
		Limit(5).
		Scan(ctx, &rows); err != nil {
		return nil, err
	}

	result := make([]TopProductItem, 0, len(rows))
	for _, r := range rows {
		name := r.NameTh
		if name == "" {
			name = r.NameEn
		}
		result = append(result, TopProductItem{
			ProductID: r.ProductID,
			Name:      name,
			Quantity:  r.Quantity,
			Revenue:   r.Revenue,
			Remaining: r.Remaining,
		})
	}
	return result, nil
}
