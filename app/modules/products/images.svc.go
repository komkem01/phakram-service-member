package products

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ProductImageItem struct {
	ID        uuid.UUID `json:"id"`
	FileID    uuid.UUID `json:"file_id"`
	FileName  string    `json:"file_name"`
	FilePath  string    `json:"file_path"`
	FileType  string    `json:"file_type"`
	FileSize  int64     `json:"file_size"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type UploadProductImageServiceRequest struct {
	FileName   string
	FileType   string
	FileSize   int64
	FileBase64 string
}

type productImageRow struct {
	ProductFileID uuid.UUID `bun:"product_file_id"`
	StorageID     uuid.UUID `bun:"storage_id"`
	FileName      string    `bun:"file_name"`
	FilePath      string    `bun:"file_path"`
	FileType      string    `bun:"file_type"`
	FileSize      int64     `bun:"file_size"`
	CreatedAt     time.Time `bun:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at"`
	ProductID     uuid.UUID `bun:"product_id"`
}

func isProductFilesRelationMissing(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(strings.TrimSpace(err.Error()))
	if message == "" {
		return false
	}
	return strings.Contains(message, `relation "product_files" does not exist`) || strings.Contains(message, "sqlstate 42p01")
}

func (s *Service) ListProductImagesService(ctx context.Context, productID uuid.UUID) ([]*ProductImageItem, error) {
	if _, err := s.db.GetProductByID(ctx, productID); err != nil {
		return nil, err
	}

	rows := make([]*productImageRow, 0)
	err := s.bunDB.DB().NewSelect().
		TableExpr("product_files AS pf").
		Join("JOIN storages AS st ON st.id = pf.file_id").
		ColumnExpr("pf.id AS product_file_id").
		ColumnExpr("pf.product_id AS product_id").
		ColumnExpr("st.id AS storage_id").
		ColumnExpr("st.file_name AS file_name").
		ColumnExpr("st.file_path AS file_path").
		ColumnExpr("st.file_type AS file_type").
		ColumnExpr("st.file_size AS file_size").
		ColumnExpr("st.created_at AS created_at").
		ColumnExpr("st.updated_at AS updated_at").
		Where("pf.product_id = ?", productID).
		Where("pf.deleted_at IS NULL").
		Where("st.deleted_at IS NULL").
		Where("st.is_active = true").
		OrderExpr("pf.created_at ASC").
		Scan(ctx, &rows)
	if err != nil {
		if isProductFilesRelationMissing(err) {
			return []*ProductImageItem{}, nil
		}
		return nil, err
	}

	items := make([]*ProductImageItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		resolvedPath := strings.TrimSpace(row.FilePath)
		if s.supabase != nil {
			resolvedPath = s.supabase.ResolveObjectURL(row.FilePath)
		}
		items = append(items, &ProductImageItem{
			ID:        row.StorageID,
			FileID:    row.StorageID,
			FileName:  row.FileName,
			FilePath:  resolvedPath,
			FileType:  row.FileType,
			FileSize:  row.FileSize,
			CreatedAt: row.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: row.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return items, nil
}

func (s *Service) UploadProductImageService(ctx context.Context, productID uuid.UUID, req *UploadProductImageServiceRequest) (*ProductImageItem, error) {
	if _, err := s.db.GetProductByID(ctx, productID); err != nil {
		return nil, err
	}
	if s.supabase == nil || !s.supabase.enabledPublic() {
		missing := []string{"client"}
		if s.supabase != nil {
			missing = s.supabase.missingPublicConfigFields()
		}
		return nil, fmt.Errorf("supabase public storage is not configured (missing: %s)", strings.Join(missing, ","))
	}

	uploaded, err := s.supabase.UploadProductImage(ctx, productID, strings.TrimSpace(req.FileName), strings.TrimSpace(req.FileBase64))
	if err != nil {
		return nil, err
	}

	now := time.Now()
	storageID := uuid.New()
	productFileID := uuid.New()

	storage := &ent.StorageEntity{
		ID:            storageID,
		RefID:         productID,
		FileName:      uploaded.FileName,
		FilePath:      uploaded.Path,
		FileSize:      uploaded.Size,
		FileType:      uploaded.MIMEType,
		IsActive:      true,
		RelatedEntity: ent.RelatedEntityProductFile,
		UploadedBy:    nil,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	productFile := &ent.ProductFileEntity{
		ID:        productFileID,
		ProductID: productID,
		FileID:    storageID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewInsert().Model(storage).Exec(ctx); err != nil {
			return err
		}
		if _, err := tx.NewInsert().Model(productFile).Exec(ctx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if isProductFilesRelationMissing(err) {
			return nil, errors.New("product image storage is not ready; please run database migrations")
		}
		return nil, err
	}

	return &ProductImageItem{
		ID:        storageID,
		FileID:    storageID,
		FileName:  uploaded.FileName,
		FilePath:  s.supabase.ResolveObjectURL(uploaded.Path),
		FileType:  uploaded.MIMEType,
		FileSize:  uploaded.Size,
		CreatedAt: now.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: now.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *Service) loadProductPrimaryImageMap(ctx context.Context, productIDs []uuid.UUID) (map[uuid.UUID]string, error) {
	imageMap := make(map[uuid.UUID]string)
	if len(productIDs) == 0 {
		return imageMap, nil
	}

	rows := make([]*productImageRow, 0)
	err := s.bunDB.DB().NewSelect().
		TableExpr("product_files AS pf").
		Join("JOIN storages AS st ON st.id = pf.file_id").
		ColumnExpr("pf.id AS product_file_id").
		ColumnExpr("pf.product_id AS product_id").
		ColumnExpr("st.file_path AS file_path").
		ColumnExpr("pf.created_at AS created_at").
		Where("pf.product_id IN (?)", bun.In(productIDs)).
		Where("pf.deleted_at IS NULL").
		Where("st.deleted_at IS NULL").
		Where("st.is_active = true").
		OrderExpr("pf.product_id ASC, pf.created_at ASC").
		Scan(ctx, &rows)
	if err != nil {
		if isProductFilesRelationMissing(err) {
			return imageMap, nil
		}
		return nil, err
	}

	for _, row := range rows {
		if row == nil {
			continue
		}
		if _, exists := imageMap[row.ProductID]; exists {
			continue
		}
		resolved := strings.TrimSpace(row.FilePath)
		if s.supabase != nil {
			resolved = s.supabase.ResolveObjectURL(row.FilePath)
		}
		if resolved != "" {
			imageMap[row.ProductID] = resolved
		}
	}

	return imageMap, nil
}

func (s *Service) loadProductImageURLs(ctx context.Context, productID uuid.UUID) ([]string, error) {
	rows, err := s.ListProductImagesService(ctx, productID)
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		path := strings.TrimSpace(row.FilePath)
		if path == "" {
			continue
		}
		urls = append(urls, path)
	}
	return urls, nil
}

func (s *Service) DeleteProductImageService(ctx context.Context, productID uuid.UUID, imageID uuid.UUID) error {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`products.svc.images.delete.start`)

	if _, err := s.db.GetProductByID(ctx, productID); err != nil {
		return err
	}

	row := new(productImageRow)
	err := s.bunDB.DB().NewSelect().
		TableExpr("product_files AS pf").
		Join("JOIN storages AS st ON st.id = pf.file_id").
		ColumnExpr("pf.id AS product_file_id").
		ColumnExpr("pf.product_id AS product_id").
		ColumnExpr("st.id AS storage_id").
		ColumnExpr("st.file_path AS file_path").
		Where("pf.product_id = ?", productID).
		Where("pf.file_id = ?", imageID).
		Where("pf.deleted_at IS NULL").
		Where("st.deleted_at IS NULL").
		Limit(1).
		Scan(ctx, row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("product image not found")
		}
		if isProductFilesRelationMissing(err) {
			return errors.New("product image storage is not ready; please run database migrations")
		}
		return err
	}

	now := time.Now()
	err = s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewUpdate().
			Model(&ent.ProductFileEntity{}).
			Set("deleted_at = ?", now).
			Set("updated_at = ?", now).
			Where("id = ?", row.ProductFileID).
			Exec(ctx); err != nil {
			return err
		}

		if _, err := tx.NewUpdate().
			Model(&ent.StorageEntity{}).
			Set("is_active = false").
			Set("deleted_at = ?", now).
			Set("updated_at = ?", now).
			Where("id = ?", row.StorageID).
			Exec(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if isProductFilesRelationMissing(err) {
			return errors.New("product image storage is not ready; please run database migrations")
		}
		return err
	}

	if s.supabase != nil {
		if removeErr := s.supabase.DeleteProductImageObject(ctx, strings.TrimSpace(row.FilePath)); removeErr != nil {
			log.With(slog.Any("product_id", productID), slog.Any("image_id", imageID)).Errf("products.svc.images.delete.supabase: %s", removeErr)
		}
	}

	span.AddEvent(`products.svc.images.delete.success`)
	return nil
}

func (s *Service) normalizeProductImageInput(req *UploadProductImageServiceRequest) error {
	if req == nil {
		return errors.New("image payload is required")
	}
	if strings.TrimSpace(req.FileBase64) == "" {
		return errors.New("file_base64 is required")
	}
	if req.FileSize > 0 && req.FileSize > maxProductImageFileSizeBytes {
		return fmt.Errorf("file size exceeds %d bytes", maxProductImageFileSizeBytes)
	}
	return nil
}
