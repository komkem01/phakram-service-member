package products

import (
	"context"
	"database/sql"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InfoProductServiceResponses struct {
	ID         uuid.UUID        `json:"id"`
	CategoryID uuid.UUID        `json:"category_id"`
	NameTh     string           `json:"name_th"`
	NameEn     string           `json:"name_en"`
	ProductNo  string           `json:"product_no"`
	Price      decimal.Decimal  `json:"price"`
	IsActive   bool             `json:"is_active"`
	CreatedAt  string           `json:"created_at"`
	UpdatedAt  string           `json:"updated_at"`
	Image      *InfoProductFile `json:"image"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*InfoProductServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`products.svc.info.start`)

	data, err := s.db.GetProductByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}

	var image *InfoProductFile
	fileData, err := s.files.GetProductFileByProductID(ctx, data.ID)
	if err == nil {
		storageData, storageErr := s.storages.GetStorageByID(ctx, fileData.FileID)
		if storageErr != nil && storageErr != sql.ErrNoRows {
			log.With(slog.Any(`file_id`, fileData.FileID)).Errf(`internal: %s`, storageErr)
			return nil, storageErr
		}

		image = &InfoProductFile{
			ID:        fileData.ID,
			ProductID: fileData.ProductID,
			FileID:    fileData.FileID,
			FilePath: func() string {
				if storageData != nil {
					return storageData.FilePath
				}
				return ""
			}(),
			CreatedAt: fileData.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: fileData.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	} else if err != sql.ErrNoRows {
		log.With(slog.Any(`product_id`, data.ID)).Errf(`internal: %s`, err)
		return nil, err
	}

	resp := &InfoProductServiceResponses{
		ID:         data.ID,
		CategoryID: data.CategoryID,
		NameTh:     data.NameTh,
		NameEn:     data.NameEn,
		ProductNo:  data.ProductNo,
		Price:      data.Price,
		IsActive:   data.IsActive,
		CreatedAt:  data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Image:      image,
	}
	span.AddEvent(`products.svc.info.success`)
	return resp, nil
}
