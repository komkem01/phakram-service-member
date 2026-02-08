package entities

import (
	"context"
	entitiesdto "phakram/app/modules/entities/dto"
	"phakram/app/modules/entities/ent"
	entitiesinf "phakram/app/modules/entities/inf"
	"phakram/app/utils/base"

	"github.com/google/uuid"
)

var _ entitiesinf.AuditLogEntity = (*Service)(nil)

func (s *Service) ListAuditLogs(ctx context.Context, req *entitiesdto.ListAuditLogsRequest) ([]*ent.AuditLogEntity, *base.ResponsePaginate, error) {
	data := make([]*ent.AuditLogEntity, 0)

	_, page, err := base.NewInstant(s.db).GetList(
		ctx,
		&data,
		&req.RequestPaginate,
		[]string{"action", "action_type", "status", "action_by"},
		[]string{"created_at", "action_type"},
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	return data, page, nil
}

func (s *Service) GetAuditLogByID(ctx context.Context, id uuid.UUID) (*ent.AuditLogEntity, error) {
	data := new(ent.AuditLogEntity)
	err := s.db.NewSelect().
		Model(data).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) CreateAuditLog(ctx context.Context, log *ent.AuditLogEntity) error {
	_, err := s.db.NewInsert().
		Model(log).
		Exec(ctx)
	return err
}

func (s *Service) UpdateAuditLog(ctx context.Context, log *ent.AuditLogEntity) error {
	_, err := s.db.NewUpdate().
		Model(log).
		Where("id = ?", log.ID).
		Exec(ctx)
	return err
}

func (s *Service) DeleteAuditLog(ctx context.Context, logID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&ent.AuditLogEntity{}).
		Where("id = ?", logID).
		Exec(ctx)
	return err
}
