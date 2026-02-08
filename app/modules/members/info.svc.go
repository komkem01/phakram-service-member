package members

import (
	"context"
	"database/sql"
	"log/slog"
	"phakram/app/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InfoMemberServiceResponses struct {
	ID            uuid.UUID       `json:"id"`
	MemberNo      string          `json:"member_no"`
	TierID        uuid.UUID       `json:"tier_id"`
	StatusID      uuid.UUID       `json:"status_id"`
	PrefixID      uuid.UUID       `json:"prefix_id"`
	GenderID      uuid.UUID       `json:"gender_id"`
	FirstnameTh   string          `json:"firstname_th"`
	LastnameTh    string          `json:"lastname_th"`
	FirstnameEn   string          `json:"firstname_en"`
	LastnameEn    string          `json:"lastname_en"`
	Role          string          `json:"role"`
	Phone         string          `json:"phone"`
	TotalSpent    decimal.Decimal `json:"total_spent"`
	CurrentPoints int             `json:"current_points"`
	Registration  string          `json:"registration"`
	LastLogin     string          `json:"last_login"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
	Image         *InfoMemberFile `json:"image"`
}

func (s *Service) InfoService(ctx context.Context, id uuid.UUID, memberID uuid.UUID, isAdmin bool) (*InfoMemberServiceResponses, error) {
	span, log := utils.LogSpanFromContext(ctx)
	span.AddEvent(`members.svc.info.start`)

	data, err := s.db.GetMemberByID(ctx, id)
	if err != nil {
		log.With(slog.Any(`id`, id)).Errf(`internal: %s`, err)
		return nil, err
	}
	if !isAdmin && memberID != uuid.Nil && data.ID != memberID {
		return nil, sql.ErrNoRows
	}

	var image *InfoMemberFile
	fileData, err := s.files.GetMemberFileByMemberID(ctx, data.ID)
	if err == nil {
		storageData, storageErr := s.storages.GetStorageByID(ctx, fileData.FileID)
		if storageErr != nil && storageErr != sql.ErrNoRows {
			log.With(slog.Any(`file_id`, fileData.FileID)).Errf(`internal: %s`, storageErr)
			return nil, storageErr
		}

		image = &InfoMemberFile{
			ID:       fileData.ID,
			MemberID: fileData.MemberID,
			FileID:   fileData.FileID,
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
		log.With(slog.Any(`member_id`, data.ID)).Errf(`internal: %s`, err)
		return nil, err
	}

	resp := &InfoMemberServiceResponses{
		ID:            data.ID,
		MemberNo:      data.MemberNo,
		TierID:        data.TierID,
		StatusID:      data.StatusID,
		PrefixID:      data.PrefixID,
		GenderID:      data.GenderID,
		FirstnameTh:   data.FirstnameTh,
		LastnameTh:    data.LastnameTh,
		FirstnameEn:   data.FirstnameEn,
		LastnameEn:    data.LastnameEn,
		Role:          string(data.Role),
		Phone:         data.Phone,
		TotalSpent:    data.TotalSpent,
		CurrentPoints: data.CurrentPoints,
		Registration: func() string {
			if data.Registration != nil {
				return data.Registration.Format("2006-01-02T15:04:05Z07:00")
			}
			return ""
		}(),
		LastLogin: func() string {
			if data.LastLogin != nil {
				return data.LastLogin.Format("2006-01-02T15:04:05Z07:00")
			}
			return ""
		}(),
		CreatedAt: data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: data.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Image:     image,
	}
	span.AddEvent(`members.svc.info.success`)
	return resp, nil
}
