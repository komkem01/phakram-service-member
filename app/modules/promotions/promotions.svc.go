package promotions

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"phakram/app/utils/base"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type promotionRecord struct {
	bun.BaseModel `bun:"table:promotions"`

	ID             uuid.UUID  `bun:"id,pk,type:uuid"`
	Code           string     `bun:"code,notnull"`
	Name           string     `bun:"name,notnull"`
	Description    string     `bun:"description"`
	DiscountType   string     `bun:"discount_type,notnull"`
	DiscountValue  float64    `bun:"discount_value,notnull"`
	MaxDiscount    *float64   `bun:"max_discount"`
	MinOrderAmount float64    `bun:"min_order_amount,notnull"`
	UsageLimit     *int       `bun:"usage_limit"`
	UsagePerMember *int       `bun:"usage_per_member"`
	UsedCount      int        `bun:"used_count,notnull"`
	StartsAt       *time.Time `bun:"starts_at"`
	EndsAt         *time.Time `bun:"ends_at"`
	IsActive       bool       `bun:"is_active,notnull"`
	CreatedAt      time.Time  `bun:"created_at,notnull"`
	UpdatedAt      time.Time  `bun:"updated_at,notnull"`
}

type promotionUsageRecord struct {
	bun.BaseModel `bun:"table:promotion_usages"`

	ID             uuid.UUID  `bun:"id,pk,type:uuid"`
	PromotionID    uuid.UUID  `bun:"promotion_id,type:uuid,notnull"`
	MemberID       uuid.UUID  `bun:"member_id,type:uuid,notnull"`
	OrderID        *uuid.UUID `bun:"order_id,type:uuid"`
	DiscountAmount float64    `bun:"discount_amount,notnull"`
	UsedAt         time.Time  `bun:"used_at,notnull"`
	CreatedAt      time.Time  `bun:"created_at,notnull"`
}

type memberPromotionCollectionRecord struct {
	bun.BaseModel `bun:"table:member_promotion_collections"`

	ID          uuid.UUID `bun:"id,pk,type:uuid"`
	MemberID    uuid.UUID `bun:"member_id,type:uuid,notnull"`
	PromotionID uuid.UUID `bun:"promotion_id,type:uuid,notnull"`
	CollectedAt time.Time `bun:"collected_at,notnull"`
	CreatedAt   time.Time `bun:"created_at,notnull"`
}

type PromotionItem struct {
	ID             string   `json:"id"`
	Code           string   `json:"code"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	DiscountType   string   `json:"discount_type"`
	DiscountValue  float64  `json:"discount_value"`
	MaxDiscount    *float64 `json:"max_discount"`
	MinOrderAmount float64  `json:"min_order_amount"`
	UsageLimit     *int     `json:"usage_limit"`
	UsagePerMember *int     `json:"usage_per_member"`
	UsedCount      int      `json:"used_count"`
	StartsAt       *string  `json:"starts_at"`
	EndsAt         *string  `json:"ends_at"`
	IsActive       bool     `json:"is_active"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

type CreatePromotionServiceRequest struct {
	Code           string
	Name           string
	Description    string
	DiscountType   string
	DiscountValue  float64
	MaxDiscount    *float64
	MinOrderAmount float64
	UsageLimit     *int
	UsagePerMember *int
	StartsAt       *string
	EndsAt         *string
	IsActive       bool
}

type UpdatePromotionServiceRequest struct {
	ID             string
	Code           string
	Name           string
	Description    string
	DiscountType   string
	DiscountValue  float64
	MaxDiscount    *float64
	MinOrderAmount float64
	UsageLimit     *int
	UsagePerMember *int
	StartsAt       *string
	EndsAt         *string
	IsActive       bool
}

type ListPromotionsServiceRequest struct {
	base.RequestPaginate
	IsActive *bool
}

type ValidatePromotionServiceRequest struct {
	Code        string
	OrderAmount float64
}

type ValidatePromotionServiceResponse struct {
	Promotion      *PromotionItem `json:"promotion"`
	IsValid        bool           `json:"is_valid"`
	Reason         string         `json:"reason"`
	DiscountAmount float64        `json:"discount_amount"`
	FinalAmount    float64        `json:"final_amount"`
}

type UsePromotionServiceRequest struct {
	PromotionID    string
	MemberID       uuid.UUID
	OrderID        *string
	DiscountAmount float64
}

type MemberPromotionItem struct {
	PromotionItem
	CollectedAt *string `json:"collected_at,omitempty"`
	IsCollected bool    `json:"is_collected"`
}

type ListMemberPromotionsServiceRequest struct {
	base.RequestPaginate
}

type PromotionReportSummary struct {
	TotalPromotions     int     `json:"total_promotions"`
	ActivePromotions    int     `json:"active_promotions"`
	CollectedCoupons    int     `json:"collected_coupons"`
	UsedCoupons         int     `json:"used_coupons"`
	TotalDiscountAmount float64 `json:"total_discount_amount"`
}

type PromotionUsageReportItem struct {
	ID             string  `json:"id"`
	PromotionID    string  `json:"promotion_id"`
	PromotionCode  string  `json:"promotion_code"`
	PromotionName  string  `json:"promotion_name"`
	MemberID       string  `json:"member_id"`
	MemberNo       string  `json:"member_no"`
	MemberName     string  `json:"member_name"`
	OrderID        string  `json:"order_id,omitempty"`
	OrderNo        string  `json:"order_no,omitempty"`
	DiscountAmount float64 `json:"discount_amount"`
	UsedAt         string  `json:"used_at"`
}

type ListPromotionUsagesServiceRequest struct {
	base.RequestPaginate
	PromotionID string
}

func normalizeDiscountType(value string) string {
	v := strings.TrimSpace(strings.ToLower(value))
	if v == "percent" || v == "amount" {
		return v
	}
	return ""
}

func parseOptionalTime(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil, nil
	}

	if parsed, err := time.Parse(time.RFC3339, trimmed); err == nil {
		utc := parsed.UTC()
		return &utc, nil
	}

	if parsed, err := time.Parse("2006-01-02T15:04", trimmed); err == nil {
		utc := parsed.UTC()
		return &utc, nil
	}

	return nil, fmt.Errorf("invalid datetime format")
}

func formatOptionalTime(value *time.Time) *string {
	if value == nil {
		return nil
	}
	v := value.Format("2006-01-02T15:04:05Z07:00")
	return &v
}

func toPromotionItem(record *promotionRecord) *PromotionItem {
	return &PromotionItem{
		ID:             record.ID.String(),
		Code:           record.Code,
		Name:           record.Name,
		Description:    record.Description,
		DiscountType:   record.DiscountType,
		DiscountValue:  record.DiscountValue,
		MaxDiscount:    record.MaxDiscount,
		MinOrderAmount: record.MinOrderAmount,
		UsageLimit:     record.UsageLimit,
		UsagePerMember: record.UsagePerMember,
		UsedCount:      record.UsedCount,
		StartsAt:       formatOptionalTime(record.StartsAt),
		EndsAt:         formatOptionalTime(record.EndsAt),
		IsActive:       record.IsActive,
		CreatedAt:      record.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      record.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *Service) List(ctx context.Context, req *ListPromotionsServiceRequest) ([]*PromotionItem, *base.ResponsePaginate, error) {
	query := s.bunDB.DB().NewSelect().Model((*promotionRecord)(nil))

	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	search := strings.TrimSpace(req.Search)
	if search != "" {
		query = query.Where("(code ILIKE ? OR name ILIKE ?)", "%"+search+"%", "%"+search+"%")
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	items := make([]*promotionRecord, 0)
	query = s.bunDB.DB().NewSelect().
		Model(&items).
		OrderExpr("created_at DESC")

	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}
	if search != "" {
		query = query.Where("(code ILIKE ? OR name ILIKE ?)", "%"+search+"%", "%"+search+"%")
	}

	req.SetOffsetLimit(query)
	if err := query.Scan(ctx); err != nil {
		return nil, nil, err
	}

	response := make([]*PromotionItem, 0, len(items))
	for _, item := range items {
		response = append(response, toPromotionItem(item))
	}

	return response, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) Info(ctx context.Context, id string) (*PromotionItem, error) {
	promotionID, err := uuid.Parse(strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}

	record := &promotionRecord{}
	if err := s.bunDB.DB().NewSelect().
		Model(record).
		Where("id = ?", promotionID).
		Limit(1).
		Scan(ctx); err != nil {
		return nil, err
	}

	return toPromotionItem(record), nil
}

func (s *Service) Create(ctx context.Context, req *CreatePromotionServiceRequest) error {
	discountType := normalizeDiscountType(req.DiscountType)
	if discountType == "" {
		return fmt.Errorf("discount type must be percent or amount")
	}
	if req.DiscountValue <= 0 {
		return fmt.Errorf("discount value must be greater than 0")
	}

	startsAt, err := parseOptionalTime(req.StartsAt)
	if err != nil {
		return err
	}
	endsAt, err := parseOptionalTime(req.EndsAt)
	if err != nil {
		return err
	}
	if startsAt != nil && endsAt != nil && endsAt.Before(*startsAt) {
		return fmt.Errorf("end date must be after start date")
	}

	now := time.Now().UTC()
	record := &promotionRecord{
		ID:             uuid.New(),
		Code:           strings.ToUpper(strings.TrimSpace(req.Code)),
		Name:           strings.TrimSpace(req.Name),
		Description:    strings.TrimSpace(req.Description),
		DiscountType:   discountType,
		DiscountValue:  req.DiscountValue,
		MaxDiscount:    req.MaxDiscount,
		MinOrderAmount: req.MinOrderAmount,
		UsageLimit:     req.UsageLimit,
		UsagePerMember: req.UsagePerMember,
		UsedCount:      0,
		StartsAt:       startsAt,
		EndsAt:         endsAt,
		IsActive:       req.IsActive,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if _, err := s.bunDB.DB().NewInsert().Model(record).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Service) Update(ctx context.Context, req *UpdatePromotionServiceRequest) error {
	promotionID, err := uuid.Parse(strings.TrimSpace(req.ID))
	if err != nil {
		return err
	}

	discountType := normalizeDiscountType(req.DiscountType)
	if discountType == "" {
		return fmt.Errorf("discount type must be percent or amount")
	}
	if req.DiscountValue <= 0 {
		return fmt.Errorf("discount value must be greater than 0")
	}

	startsAt, err := parseOptionalTime(req.StartsAt)
	if err != nil {
		return err
	}
	endsAt, err := parseOptionalTime(req.EndsAt)
	if err != nil {
		return err
	}
	if startsAt != nil && endsAt != nil && endsAt.Before(*startsAt) {
		return fmt.Errorf("end date must be after start date")
	}

	_, err = s.bunDB.DB().NewUpdate().
		Model((*promotionRecord)(nil)).
		Set("code = ?", strings.ToUpper(strings.TrimSpace(req.Code))).
		Set("name = ?", strings.TrimSpace(req.Name)).
		Set("description = ?", strings.TrimSpace(req.Description)).
		Set("discount_type = ?", discountType).
		Set("discount_value = ?", req.DiscountValue).
		Set("max_discount = ?", req.MaxDiscount).
		Set("min_order_amount = ?", req.MinOrderAmount).
		Set("usage_limit = ?", req.UsageLimit).
		Set("usage_per_member = ?", req.UsagePerMember).
		Set("starts_at = ?", startsAt).
		Set("ends_at = ?", endsAt).
		Set("is_active = ?", req.IsActive).
		Set("updated_at = ?", time.Now().UTC()).
		Where("id = ?", promotionID).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	promotionID, err := uuid.Parse(strings.TrimSpace(id))
	if err != nil {
		return err
	}

	if _, err := s.bunDB.DB().NewDelete().
		Model((*promotionRecord)(nil)).
		Where("id = ?", promotionID).
		Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Service) Validate(ctx context.Context, memberID uuid.UUID, req *ValidatePromotionServiceRequest) (*ValidatePromotionServiceResponse, error) {
	response := &ValidatePromotionServiceResponse{
		IsValid:        false,
		Reason:         "ไม่สามารถใช้โค้ดโปรโมชั่นนี้ได้",
		DiscountAmount: 0,
		FinalAmount:    req.OrderAmount,
	}

	code := strings.ToUpper(strings.TrimSpace(req.Code))
	if code == "" {
		response.Reason = "กรุณาระบุโค้ดโปรโมชั่น"
		return response, nil
	}
	if req.OrderAmount <= 0 {
		response.Reason = "ยอดคำสั่งซื้อไม่ถูกต้อง"
		return response, nil
	}

	now := time.Now().UTC()
	record := &promotionRecord{}
	if err := s.bunDB.DB().NewSelect().
		Model(record).
		Where("code = ?", code).
		Limit(1).
		Scan(ctx); err != nil {
		return response, nil
	}

	response.Promotion = toPromotionItem(record)

	if !record.IsActive {
		response.Reason = "โปรโมชั่นปิดใช้งานอยู่"
		return response, nil
	}
	if record.StartsAt != nil && now.Before(*record.StartsAt) {
		response.Reason = "โปรโมชั่นยังไม่เริ่มใช้งาน"
		return response, nil
	}
	if record.EndsAt != nil && now.After(*record.EndsAt) {
		response.Reason = "โปรโมชั่นหมดอายุแล้ว"
		return response, nil
	}
	if req.OrderAmount < record.MinOrderAmount {
		response.Reason = "ยอดสั่งซื้อไม่ถึงขั้นต่ำของโปรโมชั่น"
		return response, nil
	}

	if record.UsageLimit != nil && record.UsedCount >= *record.UsageLimit {
		response.Reason = "สิทธิ์โปรโมชั่นเต็มแล้ว"
		return response, nil
	}

	if record.UsagePerMember != nil {
		memberUsageCount, err := s.bunDB.DB().NewSelect().
			Model((*promotionUsageRecord)(nil)).
			Where("promotion_id = ?", record.ID).
			Where("member_id = ?", memberID).
			Count(ctx)
		if err != nil {
			return nil, err
		}
		if memberUsageCount >= *record.UsagePerMember {
			response.Reason = "คุณใช้โปรโมชั่นนี้ครบสิทธิ์แล้ว"
			return response, nil
		}
	}

	discount := 0.0
	if record.DiscountType == "percent" {
		discount = req.OrderAmount * (record.DiscountValue / 100)
	} else {
		discount = record.DiscountValue
	}

	if record.MaxDiscount != nil && *record.MaxDiscount > 0 && discount > *record.MaxDiscount {
		discount = *record.MaxDiscount
	}
	if discount > req.OrderAmount {
		discount = req.OrderAmount
	}
	if discount < 0 {
		discount = 0
	}

	response.IsValid = true
	response.Reason = "ใช้โปรโมชั่นได้"
	response.DiscountAmount = discount
	response.FinalAmount = req.OrderAmount - discount
	if response.FinalAmount < 0 {
		response.FinalAmount = 0
	}

	return response, nil
}

func (s *Service) Use(ctx context.Context, req *UsePromotionServiceRequest) error {
	promotionID, err := uuid.Parse(strings.TrimSpace(req.PromotionID))
	if err != nil {
		return err
	}

	usage := &promotionUsageRecord{
		ID:             uuid.New(),
		PromotionID:    promotionID,
		MemberID:       req.MemberID,
		DiscountAmount: req.DiscountAmount,
		UsedAt:         time.Now().UTC(),
		CreatedAt:      time.Now().UTC(),
	}

	if req.OrderID != nil && strings.TrimSpace(*req.OrderID) != "" {
		if orderID, parseErr := uuid.Parse(strings.TrimSpace(*req.OrderID)); parseErr == nil {
			usage.OrderID = &orderID
		}
	}

	err = s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(usage).Exec(ctx); err != nil {
			return err
		}

		if _, err := tx.NewUpdate().
			Model((*promotionRecord)(nil)).
			Set("used_count = used_count + 1").
			Set("updated_at = ?", time.Now().UTC()).
			Where("id = ?", promotionID).
			Exec(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ListAvailableForMember(ctx context.Context, memberID uuid.UUID, req *ListMemberPromotionsServiceRequest) ([]*MemberPromotionItem, *base.ResponsePaginate, error) {
	now := time.Now().UTC()
	query := s.bunDB.DB().NewSelect().
		Model((*promotionRecord)(nil)).
		Where("is_active = ?", true).
		Where("(starts_at IS NULL OR starts_at <= ?)", now).
		Where("(ends_at IS NULL OR ends_at >= ?)", now)

	search := strings.TrimSpace(req.Search)
	if search != "" {
		query = query.Where("(code ILIKE ? OR name ILIKE ?)", "%"+search+"%", "%"+search+"%")
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	records := make([]*promotionRecord, 0)
	listQuery := s.bunDB.DB().NewSelect().
		Model(&records).
		Where("is_active = ?", true).
		Where("(starts_at IS NULL OR starts_at <= ?)", now).
		Where("(ends_at IS NULL OR ends_at >= ?)", now).
		OrderExpr("created_at DESC")
	if search != "" {
		listQuery = listQuery.Where("(code ILIKE ? OR name ILIKE ?)", "%"+search+"%", "%"+search+"%")
	}

	req.SetOffsetLimit(listQuery)
	if err := listQuery.Scan(ctx); err != nil {
		return nil, nil, err
	}

	promotionIDs := make([]uuid.UUID, 0, len(records))
	for _, item := range records {
		promotionIDs = append(promotionIDs, item.ID)
	}

	collectedMap := make(map[uuid.UUID]time.Time)
	if len(promotionIDs) > 0 {
		collections := make([]*memberPromotionCollectionRecord, 0)
		if err := s.bunDB.DB().NewSelect().
			Model(&collections).
			Where("member_id = ?", memberID).
			Where("promotion_id IN (?)", bun.In(promotionIDs)).
			Scan(ctx); err != nil {
			return nil, nil, err
		}

		for _, item := range collections {
			collectedMap[item.PromotionID] = item.CollectedAt
		}
	}

	response := make([]*MemberPromotionItem, 0, len(records))
	for _, item := range records {
		baseItem := toPromotionItem(item)
		memberItem := &MemberPromotionItem{PromotionItem: *baseItem}
		if collectedAt, ok := collectedMap[item.ID]; ok {
			formatted := collectedAt.Format("2006-01-02T15:04:05Z07:00")
			memberItem.CollectedAt = &formatted
			memberItem.IsCollected = true
		}
		response = append(response, memberItem)
	}

	return response, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) ListMy(ctx context.Context, memberID uuid.UUID, req *ListMemberPromotionsServiceRequest) ([]*MemberPromotionItem, *base.ResponsePaginate, error) {
	type row struct {
		PromotionID     uuid.UUID  `bun:"promotion_id"`
		Code            string     `bun:"code"`
		Name            string     `bun:"name"`
		Description     string     `bun:"description"`
		DiscountType    string     `bun:"discount_type"`
		DiscountValue   float64    `bun:"discount_value"`
		MaxDiscount     *float64   `bun:"max_discount"`
		MinOrderAmount  float64    `bun:"min_order_amount"`
		UsageLimit      *int       `bun:"usage_limit"`
		UsagePerMember  *int       `bun:"usage_per_member"`
		UsedCount       int        `bun:"used_count"`
		StartsAt        *time.Time `bun:"starts_at"`
		EndsAt          *time.Time `bun:"ends_at"`
		IsActive        bool       `bun:"is_active"`
		PromotionCreate time.Time  `bun:"promotion_created_at"`
		PromotionUpdate time.Time  `bun:"promotion_updated_at"`
		CollectedAt     time.Time  `bun:"collected_at"`
	}

	query := s.bunDB.DB().NewSelect().
		TableExpr("member_promotion_collections AS mpc").
		Join("JOIN promotions AS p ON p.id = mpc.promotion_id").
		Where("mpc.member_id = ?", memberID)

	search := strings.TrimSpace(req.Search)
	if search != "" {
		query = query.Where("(p.code ILIKE ? OR p.name ILIKE ?)", "%"+search+"%", "%"+search+"%")
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	rows := make([]*row, 0)
	listQuery := s.bunDB.DB().NewSelect().
		TableExpr("member_promotion_collections AS mpc").
		Join("JOIN promotions AS p ON p.id = mpc.promotion_id").
		Where("mpc.member_id = ?", memberID).
		ColumnExpr("p.id AS promotion_id").
		ColumnExpr("p.code").
		ColumnExpr("p.name").
		ColumnExpr("p.description").
		ColumnExpr("p.discount_type").
		ColumnExpr("p.discount_value").
		ColumnExpr("p.max_discount").
		ColumnExpr("p.min_order_amount").
		ColumnExpr("p.usage_limit").
		ColumnExpr("p.usage_per_member").
		ColumnExpr("p.used_count").
		ColumnExpr("p.starts_at").
		ColumnExpr("p.ends_at").
		ColumnExpr("p.is_active").
		ColumnExpr("p.created_at AS promotion_created_at").
		ColumnExpr("p.updated_at AS promotion_updated_at").
		ColumnExpr("mpc.collected_at").
		OrderExpr("mpc.collected_at DESC")

	if search != "" {
		listQuery = listQuery.Where("(p.code ILIKE ? OR p.name ILIKE ?)", "%"+search+"%", "%"+search+"%")
	}

	req.SetOffsetLimit(listQuery)
	if err := listQuery.Scan(ctx, &rows); err != nil {
		return nil, nil, err
	}

	response := make([]*MemberPromotionItem, 0, len(rows))
	for _, item := range rows {
		promotion := &MemberPromotionItem{
			PromotionItem: PromotionItem{
				ID:             item.PromotionID.String(),
				Code:           item.Code,
				Name:           item.Name,
				Description:    item.Description,
				DiscountType:   item.DiscountType,
				DiscountValue:  item.DiscountValue,
				MaxDiscount:    item.MaxDiscount,
				MinOrderAmount: item.MinOrderAmount,
				UsageLimit:     item.UsageLimit,
				UsagePerMember: item.UsagePerMember,
				UsedCount:      item.UsedCount,
				StartsAt:       formatOptionalTime(item.StartsAt),
				EndsAt:         formatOptionalTime(item.EndsAt),
				IsActive:       item.IsActive,
				CreatedAt:      item.PromotionCreate.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:      item.PromotionUpdate.Format("2006-01-02T15:04:05Z07:00"),
			},
			CollectedAt: func() *string {
				formatted := item.CollectedAt.Format("2006-01-02T15:04:05Z07:00")
				return &formatted
			}(),
			IsCollected: true,
		}
		response = append(response, promotion)
	}

	return response, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) Collect(ctx context.Context, memberID uuid.UUID, promotionID string) error {
	parsedPromotionID, err := uuid.Parse(strings.TrimSpace(promotionID))
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	promotion := &promotionRecord{}
	if err := s.bunDB.DB().NewSelect().
		Model(promotion).
		Where("id = ?", parsedPromotionID).
		Limit(1).
		Scan(ctx); err != nil {
		return err
	}

	if !promotion.IsActive {
		return fmt.Errorf("promotion is inactive")
	}
	if promotion.StartsAt != nil && now.Before(*promotion.StartsAt) {
		return fmt.Errorf("promotion is not active yet")
	}
	if promotion.EndsAt != nil && now.After(*promotion.EndsAt) {
		return fmt.Errorf("promotion has expired")
	}

	record := &memberPromotionCollectionRecord{
		ID:          uuid.New(),
		MemberID:    memberID,
		PromotionID: parsedPromotionID,
		CollectedAt: now,
		CreatedAt:   now,
	}

	if _, err := s.bunDB.DB().NewInsert().
		Model(record).
		On("CONFLICT (member_id, promotion_id) DO NOTHING").
		Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Service) ReportSummary(ctx context.Context) (*PromotionReportSummary, error) {
	totalPromotions, err := s.bunDB.DB().NewSelect().Model((*promotionRecord)(nil)).Count(ctx)
	if err != nil {
		return nil, err
	}

	activePromotions, err := s.bunDB.DB().NewSelect().
		Model((*promotionRecord)(nil)).
		Where("is_active = ?", true).
		Count(ctx)
	if err != nil {
		return nil, err
	}

	collectedCoupons, err := s.bunDB.DB().NewSelect().Model((*memberPromotionCollectionRecord)(nil)).Count(ctx)
	if err != nil {
		return nil, err
	}

	usedCoupons, err := s.bunDB.DB().NewSelect().Model((*promotionUsageRecord)(nil)).Count(ctx)
	if err != nil {
		return nil, err
	}

	var totalDiscount sql.NullFloat64
	if err := s.bunDB.DB().NewSelect().
		Model((*promotionUsageRecord)(nil)).
		ColumnExpr("COALESCE(SUM(discount_amount), 0)").
		Scan(ctx, &totalDiscount); err != nil {
		return nil, err
	}

	return &PromotionReportSummary{
		TotalPromotions:     totalPromotions,
		ActivePromotions:    activePromotions,
		CollectedCoupons:    collectedCoupons,
		UsedCoupons:         usedCoupons,
		TotalDiscountAmount: totalDiscount.Float64,
	}, nil
}

func (s *Service) ListUsages(ctx context.Context, req *ListPromotionUsagesServiceRequest) ([]*PromotionUsageReportItem, *base.ResponsePaginate, error) {
	type usageRow struct {
		ID             uuid.UUID  `bun:"id"`
		PromotionID    uuid.UUID  `bun:"promotion_id"`
		PromotionCode  string     `bun:"promotion_code"`
		PromotionName  string     `bun:"promotion_name"`
		MemberID       uuid.UUID  `bun:"member_id"`
		MemberNo       string     `bun:"member_no"`
		FirstnameTH    string     `bun:"firstname_th"`
		LastnameTH     string     `bun:"lastname_th"`
		FirstnameEN    string     `bun:"firstname_en"`
		LastnameEN     string     `bun:"lastname_en"`
		OrderID        *uuid.UUID `bun:"order_id"`
		OrderNo        string     `bun:"order_no"`
		DiscountAmount float64    `bun:"discount_amount"`
		UsedAt         time.Time  `bun:"used_at"`
	}

	baseQuery := s.bunDB.DB().NewSelect().
		TableExpr("promotion_usages AS pu").
		Join("JOIN promotions AS p ON p.id = pu.promotion_id").
		Join("LEFT JOIN members AS m ON m.id = pu.member_id").
		Join("LEFT JOIN orders AS o ON o.id = pu.order_id")

	promotionID := strings.TrimSpace(req.PromotionID)
	if promotionID != "" {
		parsedPromotionID, err := uuid.Parse(promotionID)
		if err != nil {
			return nil, nil, err
		}
		baseQuery = baseQuery.Where("pu.promotion_id = ?", parsedPromotionID)
	}

	search := strings.TrimSpace(req.Search)
	if search != "" {
		like := "%" + search + "%"
		baseQuery = baseQuery.Where("(p.code ILIKE ? OR p.name ILIKE ? OR m.member_no ILIKE ? OR o.order_no ILIKE ?)", like, like, like, like)
	}

	total, err := baseQuery.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	rows := make([]*usageRow, 0)
	listQuery := s.bunDB.DB().NewSelect().
		TableExpr("promotion_usages AS pu").
		Join("JOIN promotions AS p ON p.id = pu.promotion_id").
		Join("LEFT JOIN members AS m ON m.id = pu.member_id").
		Join("LEFT JOIN orders AS o ON o.id = pu.order_id").
		ColumnExpr("pu.id").
		ColumnExpr("pu.promotion_id").
		ColumnExpr("p.code AS promotion_code").
		ColumnExpr("p.name AS promotion_name").
		ColumnExpr("pu.member_id").
		ColumnExpr("m.member_no").
		ColumnExpr("m.firstname_th").
		ColumnExpr("m.lastname_th").
		ColumnExpr("m.firstname_en").
		ColumnExpr("m.lastname_en").
		ColumnExpr("pu.order_id").
		ColumnExpr("COALESCE(o.order_no, '') AS order_no").
		ColumnExpr("pu.discount_amount").
		ColumnExpr("pu.used_at").
		OrderExpr("pu.used_at DESC")

	if promotionID != "" {
		parsedPromotionID, _ := uuid.Parse(promotionID)
		listQuery = listQuery.Where("pu.promotion_id = ?", parsedPromotionID)
	}
	if search != "" {
		like := "%" + search + "%"
		listQuery = listQuery.Where("(p.code ILIKE ? OR p.name ILIKE ? OR m.member_no ILIKE ? OR o.order_no ILIKE ?)", like, like, like, like)
	}

	req.SetOffsetLimit(listQuery)
	if err := listQuery.Scan(ctx, &rows); err != nil {
		return nil, nil, err
	}

	result := make([]*PromotionUsageReportItem, 0, len(rows))
	for _, row := range rows {
		memberName := strings.TrimSpace(strings.TrimSpace(row.FirstnameTH + " " + row.LastnameTH))
		if memberName == "" {
			memberName = strings.TrimSpace(strings.TrimSpace(row.FirstnameEN + " " + row.LastnameEN))
		}

		item := &PromotionUsageReportItem{
			ID:             row.ID.String(),
			PromotionID:    row.PromotionID.String(),
			PromotionCode:  row.PromotionCode,
			PromotionName:  row.PromotionName,
			MemberID:       row.MemberID.String(),
			MemberNo:       row.MemberNo,
			MemberName:     memberName,
			DiscountAmount: row.DiscountAmount,
			UsedAt:         row.UsedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		if row.OrderID != nil {
			item.OrderID = row.OrderID.String()
			item.OrderNo = row.OrderNo
		}

		result = append(result, item)
	}

	return result, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}
