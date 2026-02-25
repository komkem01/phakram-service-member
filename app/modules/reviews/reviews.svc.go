package reviews

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"phakram/app/utils/base"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const (
	reviewEditWindowDuration = 7 * 24 * time.Hour
	maxReviewImages          = 3
)

type productReviewRecord struct {
	bun.BaseModel `bun:"table:product_reviews"`

	ID          uuid.UUID `bun:"id,pk,type:uuid"`
	MemberID    uuid.UUID `bun:"member_id,type:uuid,notnull"`
	ProductID   uuid.UUID `bun:"product_id,type:uuid,notnull"`
	OrderID     uuid.UUID `bun:"order_id,type:uuid,notnull"`
	OrderItemID uuid.UUID `bun:"order_item_id,type:uuid,notnull"`
	Rating      int       `bun:"rating,notnull"`
	Comment     string    `bun:"comment,notnull"`
	IsVisible   bool      `bun:"is_visible,notnull"`
	CreatedAt   time.Time `bun:"created_at,notnull"`
	UpdatedAt   time.Time `bun:"updated_at,notnull"`
}

type productReviewImageRecord struct {
	bun.BaseModel `bun:"table:product_review_images"`

	ID        uuid.UUID `bun:"id,pk,type:uuid"`
	ReviewID  uuid.UUID `bun:"review_id,type:uuid,notnull"`
	ImageURL  string    `bun:"image_url,notnull"`
	SortOrder int       `bun:"sort_order,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull"`
	UpdatedAt time.Time `bun:"updated_at,notnull"`
}

type ListProductReviewsServiceRequest struct {
	base.RequestPaginate
	HasImages *bool
	Rating    *int
}

type ListAdminReviewsServiceRequest struct {
	base.RequestPaginate
	ProductID *uuid.UUID
	IsVisible *bool
	HasImages *bool
	Rating    *int
}

type ProductReviewItem struct {
	ID                 string   `json:"id"`
	MemberID           string   `json:"member_id"`
	MemberName         string   `json:"member_name"`
	ProductID          string   `json:"product_id"`
	OrderID            string   `json:"order_id"`
	OrderItemID        string   `json:"order_item_id"`
	OrderNo            string   `json:"order_no"`
	Rating             int      `json:"rating"`
	Comment            string   `json:"comment"`
	ImageURLs          []string `json:"image_urls"`
	IsVisible          bool     `json:"is_visible"`
	IsVerifiedPurchase bool     `json:"is_verified_purchase"`
	CreatedAt          string   `json:"created_at"`
	UpdatedAt          string   `json:"updated_at"`
}

type ProductReviewSummary struct {
	TotalReviews int     `json:"total_reviews"`
	AverageScore float64 `json:"average_score"`
}

type ListProductReviewsServiceResponse struct {
	Summary ProductReviewSummary `json:"summary"`
	Items   []*ProductReviewItem `json:"items"`
}

type CreateReviewServiceRequest struct {
	MemberID    uuid.UUID
	OrderItemID uuid.UUID
	Rating      int
	Comment     string
	ImageURLs   []string
}

type UpdateReviewServiceRequest struct {
	ReviewID  uuid.UUID
	MemberID  uuid.UUID
	Rating    int
	Comment   string
	ImageURLs []string
}

type EligibleReviewItem struct {
	OrderID     string `json:"order_id"`
	OrderNo     string `json:"order_no"`
	OrderItemID string `json:"order_item_id"`
	CreatedAt   string `json:"created_at"`
}

func sanitizeReviewImageURLs(raw []string) ([]string, error) {
	if len(raw) > maxReviewImages {
		return nil, errors.New("review images must not exceed 3 items")
	}

	urls := make([]string, 0, len(raw))
	for _, item := range raw {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		lowerValue := strings.ToLower(trimmed)
		isRemoteURL := strings.HasPrefix(lowerValue, "http://") || strings.HasPrefix(lowerValue, "https://")
		isDataImage := strings.HasPrefix(lowerValue, "data:image/") && strings.Contains(lowerValue, ";base64,")
		if !isRemoteURL && !isDataImage {
			return nil, errors.New("review image url is invalid")
		}
		urls = append(urls, trimmed)
	}

	if len(urls) > maxReviewImages {
		return nil, errors.New("review images must not exceed 3 items")
	}

	return urls, nil
}

func normalizeReviewSort(sortBy string, orderBy string) string {
	normalizedSortBy := strings.ToLower(strings.TrimSpace(sortBy))
	normalizedOrderBy := strings.ToUpper(strings.TrimSpace(orderBy))
	if normalizedOrderBy != "ASC" && normalizedOrderBy != "DESC" {
		normalizedOrderBy = "DESC"
	}

	switch normalizedSortBy {
	case "rating":
		return "pr.rating " + normalizedOrderBy + ", pr.created_at DESC"
	case "created_at":
		fallthrough
	default:
		return "pr.created_at " + normalizedOrderBy
	}
}

func (s *Service) loadReviewImageMap(ctx context.Context, reviewIDs []uuid.UUID) (map[uuid.UUID][]string, error) {
	imageMap := make(map[uuid.UUID][]string)
	if len(reviewIDs) == 0 {
		return imageMap, nil
	}

	rows := make([]*productReviewImageRecord, 0)
	if err := s.bunDB.DB().NewSelect().
		Model(&rows).
		Where("review_id IN (?)", bun.In(reviewIDs)).
		OrderExpr("sort_order ASC, created_at ASC").
		Scan(ctx); err != nil {
		return nil, err
	}

	for _, item := range rows {
		resolved := strings.TrimSpace(item.ImageURL)
		if s.railwayStorage != nil {
			resolved = s.railwayStorage.ResolveObjectURL(resolved)
		}
		imageMap[item.ReviewID] = append(imageMap[item.ReviewID], resolved)
	}

	return imageMap, nil
}

func (s *Service) toStoredReviewImageURLs(ctx context.Context, productID uuid.UUID, reviewID uuid.UUID, imageURLs []string) ([]string, error) {
	stored := make([]string, 0, len(imageURLs))
	for _, imageURL := range imageURLs {
		trimmed := strings.TrimSpace(imageURL)
		if trimmed == "" {
			continue
		}

		lowerValue := strings.ToLower(trimmed)
		if strings.HasPrefix(lowerValue, "data:image/") {
			if s.railwayStorage == nil || !s.railwayStorage.enabledForPublic() {
				return nil, errors.New("railway public storage is not configured")
			}

			uploaded, err := s.railwayStorage.UploadReviewImage(ctx, productID, reviewID, "", trimmed)
			if err != nil {
				return nil, err
			}

			stored = append(stored, uploaded.Path)
			continue
		}

		if strings.HasPrefix(lowerValue, "http://") || strings.HasPrefix(lowerValue, "https://") {
			stored = append(stored, trimmed)
			continue
		}

		stored = append(stored, trimmed)
	}

	return stored, nil
}

func (s *Service) applyPublicReviewFilters(query *bun.SelectQuery, req *ListProductReviewsServiceRequest) {
	if req.Rating != nil {
		query.Where("pr.rating = ?", *req.Rating)
	}
	if req.HasImages != nil {
		if *req.HasImages {
			query.Where("EXISTS (SELECT 1 FROM product_review_images pri WHERE pri.review_id = pr.id)")
		} else {
			query.Where("NOT EXISTS (SELECT 1 FROM product_review_images pri WHERE pri.review_id = pr.id)")
		}
	}
}

func (s *Service) ListPublicByProduct(ctx context.Context, productID uuid.UUID, req *ListProductReviewsServiceRequest) (*ListProductReviewsServiceResponse, *base.ResponsePaginate, error) {
	type summaryRow struct {
		Total   int     `bun:"total"`
		Average float64 `bun:"average"`
	}

	summary := &summaryRow{}
	summaryQuery := s.bunDB.DB().NewSelect().
		TableExpr("product_reviews AS pr").
		ColumnExpr("COUNT(*) AS total").
		ColumnExpr("COALESCE(AVG(pr.rating), 0) AS average").
		Where("pr.product_id = ?", productID).
		Where("pr.is_visible = ?", true)
	s.applyPublicReviewFilters(summaryQuery, req)
	if err := summaryQuery.Scan(ctx, summary); err != nil {
		return nil, nil, err
	}

	type row struct {
		ID          uuid.UUID `bun:"id"`
		MemberID    uuid.UUID `bun:"member_id"`
		ProductID   uuid.UUID `bun:"product_id"`
		OrderID     uuid.UUID `bun:"order_id"`
		OrderItemID uuid.UUID `bun:"order_item_id"`
		OrderNo     string    `bun:"order_no"`
		FirstnameTh string    `bun:"firstname_th"`
		LastnameTh  string    `bun:"lastname_th"`
		FirstnameEn string    `bun:"firstname_en"`
		LastnameEn  string    `bun:"lastname_en"`
		Rating      int       `bun:"rating"`
		Comment     string    `bun:"comment"`
		IsVisible   bool      `bun:"is_visible"`
		CreatedAt   time.Time `bun:"created_at"`
		UpdatedAt   time.Time `bun:"updated_at"`
	}

	rows := make([]*row, 0)
	query := s.bunDB.DB().NewSelect().
		TableExpr("product_reviews AS pr").
		Join("JOIN members AS m ON m.id = pr.member_id").
		Join("JOIN orders AS o ON o.id = pr.order_id").
		ColumnExpr("pr.id").
		ColumnExpr("pr.member_id").
		ColumnExpr("pr.product_id").
		ColumnExpr("pr.order_id").
		ColumnExpr("pr.order_item_id").
		ColumnExpr("o.order_no").
		ColumnExpr("m.firstname_th").
		ColumnExpr("m.lastname_th").
		ColumnExpr("m.firstname_en").
		ColumnExpr("m.lastname_en").
		ColumnExpr("pr.rating").
		ColumnExpr("pr.comment").
		ColumnExpr("pr.is_visible").
		ColumnExpr("pr.created_at").
		ColumnExpr("pr.updated_at").
		Where("pr.product_id = ?", productID).
		Where("pr.is_visible = ?", true).
		OrderExpr(normalizeReviewSort(req.SortBy, req.OrderBy))
	s.applyPublicReviewFilters(query, req)

	req.SetOffsetLimit(query)
	if err := query.Scan(ctx, &rows); err != nil {
		return nil, nil, err
	}

	reviewIDs := make([]uuid.UUID, 0, len(rows))
	for _, item := range rows {
		reviewIDs = append(reviewIDs, item.ID)
	}
	imageMap, err := s.loadReviewImageMap(ctx, reviewIDs)
	if err != nil {
		return nil, nil, err
	}

	items := make([]*ProductReviewItem, 0, len(rows))
	for _, item := range rows {
		memberName := strings.TrimSpace(strings.TrimSpace(item.FirstnameTh + " " + item.LastnameTh))
		if memberName == "" {
			memberName = strings.TrimSpace(strings.TrimSpace(item.FirstnameEn + " " + item.LastnameEn))
		}
		if memberName == "" {
			memberName = "ลูกค้าที่ซื้อสินค้า"
		}

		items = append(items, &ProductReviewItem{
			ID:                 item.ID.String(),
			MemberID:           item.MemberID.String(),
			MemberName:         memberName,
			ProductID:          item.ProductID.String(),
			OrderID:            item.OrderID.String(),
			OrderItemID:        item.OrderItemID.String(),
			OrderNo:            item.OrderNo,
			Rating:             item.Rating,
			Comment:            item.Comment,
			ImageURLs:          imageMap[item.ID],
			IsVisible:          item.IsVisible,
			IsVerifiedPurchase: true,
			CreatedAt:          item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:          item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &ListProductReviewsServiceResponse{
		Summary: ProductReviewSummary{
			TotalReviews: summary.Total,
			AverageScore: summary.Average,
		},
		Items: items,
	}, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(summary.Total)}, nil
}

func (s *Service) ListEligibleByProduct(ctx context.Context, memberID uuid.UUID, productID uuid.UUID) ([]*EligibleReviewItem, error) {
	type row struct {
		OrderID     uuid.UUID `bun:"order_id"`
		OrderNo     string    `bun:"order_no"`
		OrderItemID uuid.UUID `bun:"order_item_id"`
		CreatedAt   time.Time `bun:"created_at"`
	}

	rows := make([]*row, 0)
	if err := s.bunDB.DB().NewSelect().
		TableExpr("order_items AS oi").
		Join("JOIN orders AS o ON o.id = oi.order_id").
		Join("LEFT JOIN product_reviews AS pr ON pr.order_item_id = oi.id").
		ColumnExpr("o.id AS order_id").
		ColumnExpr("o.order_no").
		ColumnExpr("oi.id AS order_item_id").
		ColumnExpr("o.created_at").
		Where("o.member_id = ?", memberID).
		Where("o.status = ?", "completed").
		Where("oi.product_id = ?", productID).
		Where("pr.id IS NULL").
		OrderExpr("o.created_at DESC").
		Scan(ctx, &rows); err != nil {
		return nil, err
	}

	items := make([]*EligibleReviewItem, 0, len(rows))
	for _, item := range rows {
		items = append(items, &EligibleReviewItem{
			OrderID:     item.OrderID.String(),
			OrderNo:     item.OrderNo,
			OrderItemID: item.OrderItemID.String(),
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return items, nil
}

func (s *Service) upsertReviewImagesInTx(ctx context.Context, tx bun.Tx, reviewID uuid.UUID, imageURLs []string) error {
	if _, err := tx.NewDelete().Model((*productReviewImageRecord)(nil)).Where("review_id = ?", reviewID).Exec(ctx); err != nil {
		return err
	}

	if len(imageURLs) == 0 {
		return nil
	}

	now := time.Now().UTC()
	images := make([]*productReviewImageRecord, 0, len(imageURLs))
	for idx, imageURL := range imageURLs {
		images = append(images, &productReviewImageRecord{
			ID:        uuid.New(),
			ReviewID:  reviewID,
			ImageURL:  imageURL,
			SortOrder: idx,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	if _, err := tx.NewInsert().Model(&images).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Service) Create(ctx context.Context, req *CreateReviewServiceRequest) error {
	if req.Rating < 1 || req.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	comment := strings.TrimSpace(req.Comment)
	if comment == "" {
		return errors.New("review comment is required")
	}

	imageURLs, err := sanitizeReviewImageURLs(req.ImageURLs)
	if err != nil {
		return err
	}

	type orderItemRow struct {
		OrderID   uuid.UUID `bun:"order_id"`
		ProductID uuid.UUID `bun:"product_id"`
	}

	orderItem := &orderItemRow{}
	err = s.bunDB.DB().NewSelect().
		TableExpr("order_items AS oi").
		Join("JOIN orders AS o ON o.id = oi.order_id").
		ColumnExpr("oi.order_id").
		ColumnExpr("oi.product_id").
		Where("oi.id = ?", req.OrderItemID).
		Where("o.member_id = ?", req.MemberID).
		Where("o.status = ?", "completed").
		Limit(1).
		Scan(ctx, orderItem)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("review can be created only for completed orders")
		}
		return err
	}

	existingCount, err := s.bunDB.DB().NewSelect().
		Model((*productReviewRecord)(nil)).
		Where("order_item_id = ?", req.OrderItemID).
		Count(ctx)
	if err != nil {
		return err
	}
	if existingCount > 0 {
		return errors.New("review already exists for this order item")
	}

	now := time.Now().UTC()
	reviewID := uuid.New()
	storedImageURLs, err := s.toStoredReviewImageURLs(ctx, orderItem.ProductID, reviewID, imageURLs)
	if err != nil {
		return err
	}

	record := &productReviewRecord{
		ID:          reviewID,
		MemberID:    req.MemberID,
		ProductID:   orderItem.ProductID,
		OrderID:     orderItem.OrderID,
		OrderItemID: req.OrderItemID,
		Rating:      req.Rating,
		Comment:     comment,
		IsVisible:   true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(record).Exec(ctx); err != nil {
			return err
		}
		return s.upsertReviewImagesInTx(ctx, tx, record.ID, storedImageURLs)
	})
}

func (s *Service) Update(ctx context.Context, req *UpdateReviewServiceRequest) error {
	if req.Rating < 1 || req.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	comment := strings.TrimSpace(req.Comment)
	if comment == "" {
		return errors.New("review comment is required")
	}

	imageURLs, err := sanitizeReviewImageURLs(req.ImageURLs)
	if err != nil {
		return err
	}

	review := new(productReviewRecord)
	if err := s.bunDB.DB().NewSelect().
		Model(review).
		Where("id = ?", req.ReviewID).
		Where("member_id = ?", req.MemberID).
		Limit(1).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("review not found")
		}
		return err
	}

	if review.CreatedAt.Add(reviewEditWindowDuration).Before(time.Now().UTC()) {
		return errors.New("review edit window expired")
	}

	storedImageURLs, err := s.toStoredReviewImageURLs(ctx, review.ProductID, review.ID, imageURLs)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	return s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().
			Model((*productReviewRecord)(nil)).
			Set("rating = ?", req.Rating).
			Set("comment = ?", comment).
			Set("updated_at = ?", now).
			Where("id = ?", req.ReviewID).
			Where("member_id = ?", req.MemberID).
			Exec(ctx); err != nil {
			return err
		}

		return s.upsertReviewImagesInTx(ctx, tx, req.ReviewID, storedImageURLs)
	})
}

func (s *Service) Delete(ctx context.Context, reviewID uuid.UUID, memberID uuid.UUID) error {
	review := new(productReviewRecord)
	if err := s.bunDB.DB().NewSelect().
		Model(review).
		Where("id = ?", reviewID).
		Where("member_id = ?", memberID).
		Limit(1).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("review not found")
		}
		return err
	}

	if review.CreatedAt.Add(reviewEditWindowDuration).Before(time.Now().UTC()) {
		return errors.New("review edit window expired")
	}

	return s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewDelete().
			Model((*productReviewImageRecord)(nil)).
			Where("review_id = ?", reviewID).
			Exec(ctx); err != nil {
			return err
		}

		res, err := tx.NewDelete().
			Model((*productReviewRecord)(nil)).
			Where("id = ?", reviewID).
			Where("member_id = ?", memberID).
			Exec(ctx)
		if err != nil {
			return err
		}

		affected, _ := res.RowsAffected()
		if affected == 0 {
			return errors.New("review not found")
		}

		return nil
	})
}

func (s *Service) applyAdminReviewFilters(query *bun.SelectQuery, req *ListAdminReviewsServiceRequest) {
	if req.ProductID != nil {
		query.Where("pr.product_id = ?", *req.ProductID)
	}
	if req.IsVisible != nil {
		query.Where("pr.is_visible = ?", *req.IsVisible)
	}
	if req.Rating != nil {
		query.Where("pr.rating = ?", *req.Rating)
	}
	if req.HasImages != nil {
		if *req.HasImages {
			query.Where("EXISTS (SELECT 1 FROM product_review_images pri WHERE pri.review_id = pr.id)")
		} else {
			query.Where("NOT EXISTS (SELECT 1 FROM product_review_images pri WHERE pri.review_id = pr.id)")
		}
	}
}

func (s *Service) ListAdmin(ctx context.Context, req *ListAdminReviewsServiceRequest) ([]*ProductReviewItem, *base.ResponsePaginate, error) {
	type row struct {
		ID          uuid.UUID `bun:"id"`
		MemberID    uuid.UUID `bun:"member_id"`
		ProductID   uuid.UUID `bun:"product_id"`
		OrderID     uuid.UUID `bun:"order_id"`
		OrderItemID uuid.UUID `bun:"order_item_id"`
		OrderNo     string    `bun:"order_no"`
		FirstnameTh string    `bun:"firstname_th"`
		LastnameTh  string    `bun:"lastname_th"`
		FirstnameEn string    `bun:"firstname_en"`
		LastnameEn  string    `bun:"lastname_en"`
		Rating      int       `bun:"rating"`
		Comment     string    `bun:"comment"`
		IsVisible   bool      `bun:"is_visible"`
		CreatedAt   time.Time `bun:"created_at"`
		UpdatedAt   time.Time `bun:"updated_at"`
	}

	search := strings.TrimSpace(req.Search)

	countQuery := s.bunDB.DB().NewSelect().
		TableExpr("product_reviews AS pr").
		Join("JOIN members AS m ON m.id = pr.member_id").
		Join("JOIN orders AS o ON o.id = pr.order_id").
		Join("JOIN products AS p ON p.id = pr.product_id")
	s.applyAdminReviewFilters(countQuery, req)
	if search != "" {
		pattern := "%" + search + "%"
		countQuery = countQuery.Where("(o.order_no ILIKE ? OR p.name_th ILIKE ? OR p.name_en ILIKE ? OR pr.comment ILIKE ?)", pattern, pattern, pattern, pattern)
	}

	total, err := countQuery.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	rows := make([]*row, 0)
	query := s.bunDB.DB().NewSelect().
		TableExpr("product_reviews AS pr").
		Join("JOIN members AS m ON m.id = pr.member_id").
		Join("JOIN orders AS o ON o.id = pr.order_id").
		Join("JOIN products AS p ON p.id = pr.product_id").
		ColumnExpr("pr.id").
		ColumnExpr("pr.member_id").
		ColumnExpr("pr.product_id").
		ColumnExpr("pr.order_id").
		ColumnExpr("pr.order_item_id").
		ColumnExpr("o.order_no").
		ColumnExpr("m.firstname_th").
		ColumnExpr("m.lastname_th").
		ColumnExpr("m.firstname_en").
		ColumnExpr("m.lastname_en").
		ColumnExpr("pr.rating").
		ColumnExpr("pr.comment").
		ColumnExpr("pr.is_visible").
		ColumnExpr("pr.created_at").
		ColumnExpr("pr.updated_at").
		OrderExpr(normalizeReviewSort(req.SortBy, req.OrderBy))
	s.applyAdminReviewFilters(query, req)
	if search != "" {
		pattern := "%" + search + "%"
		query = query.Where("(o.order_no ILIKE ? OR p.name_th ILIKE ? OR p.name_en ILIKE ? OR pr.comment ILIKE ?)", pattern, pattern, pattern, pattern)
	}

	req.SetOffsetLimit(query)
	if err := query.Scan(ctx, &rows); err != nil {
		return nil, nil, err
	}

	reviewIDs := make([]uuid.UUID, 0, len(rows))
	for _, item := range rows {
		reviewIDs = append(reviewIDs, item.ID)
	}
	imageMap, err := s.loadReviewImageMap(ctx, reviewIDs)
	if err != nil {
		return nil, nil, err
	}

	items := make([]*ProductReviewItem, 0, len(rows))
	for _, item := range rows {
		memberName := strings.TrimSpace(strings.TrimSpace(item.FirstnameTh + " " + item.LastnameTh))
		if memberName == "" {
			memberName = strings.TrimSpace(strings.TrimSpace(item.FirstnameEn + " " + item.LastnameEn))
		}
		if memberName == "" {
			memberName = "ไม่ระบุชื่อ"
		}

		items = append(items, &ProductReviewItem{
			ID:                 item.ID.String(),
			MemberID:           item.MemberID.String(),
			MemberName:         memberName,
			ProductID:          item.ProductID.String(),
			OrderID:            item.OrderID.String(),
			OrderItemID:        item.OrderItemID.String(),
			OrderNo:            item.OrderNo,
			Rating:             item.Rating,
			Comment:            item.Comment,
			ImageURLs:          imageMap[item.ID],
			IsVisible:          item.IsVisible,
			IsVerifiedPurchase: true,
			CreatedAt:          item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:          item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return items, &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}, nil
}

func (s *Service) UpdateVisibility(ctx context.Context, reviewID uuid.UUID, isVisible bool) error {
	res, err := s.bunDB.DB().NewUpdate().
		Model((*productReviewRecord)(nil)).
		Set("is_visible = ?", isVisible).
		Set("updated_at = ?", time.Now().UTC()).
		Where("id = ?", reviewID).
		Exec(ctx)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("review not found")
	}

	return nil
}
