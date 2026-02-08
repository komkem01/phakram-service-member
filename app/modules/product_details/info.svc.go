package product_details

import (
	"context"
	"database/sql"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InfoProductDetailServiceResponses struct {
	ID               uuid.UUID              `json:"id"`
	ProductID        uuid.UUID              `json:"product_id"`
	Description      string                 `json:"description"`
	Material         string                 `json:"material"`
	Dimensions       string                 `json:"dimensions"`
	Weight           decimal.Decimal        `json:"weight"`
	CareInstructions string                 `json:"care_instructions"`
	Image            *InfoProductDetailFile `json:"image"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID) (*InfoProductDetailServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`product_details.svc.info.start`)

	data, err := s.db.GetProductDetailByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}

	var image *InfoProductDetailFile
	fileData, err := s.files.GetProductFileByProductID(ctx, data.ProductID)
	if err == nil {
		storageData, storageErr := s.storages.GetStorageByID(ctx, fileData.FileID)
		if storageErr != nil && storageErr != sql.ErrNoRows {
			log.With(slog.Any(`file_id`, fileData.FileID)).Errf(`internal: %s`, storageErr)
			return nil, storageErr
		}

		image = &InfoProductDetailFile{
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
		log.With(slog.Any(`product_detail_id`, data.ID)).Errf(`internal: %s`, err)
		return nil, err
	}

	resp := &InfoProductDetailServiceResponses{
		ID:               data.ID,
		ProductID:        data.ProductID,
		Description:      data.Description,
		Material:         data.Material,
		Dimensions:       data.Dimensions,
		Weight:           data.Weight,
		CareInstructions: data.CareInstructions,
		Image:            image,
	}
	span.AddEvent(`product_details.svc.info.success`)
	return resp, nil
}
