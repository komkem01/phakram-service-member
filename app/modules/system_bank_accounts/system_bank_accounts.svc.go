package systembankaccounts

import (
	"context"
	"errors"
	"phakram/app/modules/entities/ent"
	"phakram/app/utils"
	"phakram/app/utils/base"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ListSystemBankAccountServiceRequest struct {
	base.RequestPaginate
	BankID   uuid.UUID
	IsActive *bool
}

type SystemBankAccountServiceResponse struct {
	ID               uuid.UUID `json:"id"`
	BankID           uuid.UUID `json:"bank_id"`
	BankNameTh       string    `json:"bank_name_th"`
	BankNameEn       string    `json:"bank_name_en"`
	AccountName      string    `json:"account_name"`
	AccountNo        string    `json:"account_no"`
	Branch           string    `json:"branch"`
	IsActive         bool      `json:"is_active"`
	IsDefaultReceive bool      `json:"is_default_receive"`
	IsDefaultRefund  bool      `json:"is_default_refund"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type UpsertSystemBankAccountServiceRequest struct {
	BankID           uuid.UUID
	AccountName      string
	AccountNo        string
	Branch           string
	IsActive         bool
	IsDefaultReceive bool
	IsDefaultRefund  bool
}

type systemBankAccountRow struct {
	ID               uuid.UUID `bun:"id"`
	BankID           uuid.UUID `bun:"bank_id"`
	BankNameTh       string    `bun:"bank_name_th"`
	BankNameEn       string    `bun:"bank_name_en"`
	AccountName      string    `bun:"account_name"`
	AccountNo        string    `bun:"account_no"`
	Branch           string    `bun:"branch"`
	IsActive         bool      `bun:"is_active"`
	IsDefaultReceive bool      `bun:"is_default_receive"`
	IsDefaultRefund  bool      `bun:"is_default_refund"`
	CreatedAt        time.Time `bun:"created_at"`
	UpdatedAt        time.Time `bun:"updated_at"`
}

func (s *Service) ListService(ctx context.Context, req *ListSystemBankAccountServiceRequest) ([]*SystemBankAccountServiceResponse, *base.ResponsePaginate, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`system_bank_accounts.svc.list.start`)

	query := s.bunDB.DB().NewSelect().
		TableExpr("system_bank_accounts AS sba").
		Join("LEFT JOIN banks AS b ON b.id = sba.bank_id")

	if req.BankID != uuid.Nil {
		query = query.Where("sba.bank_id = ?", req.BankID)
	}
	if req.IsActive != nil {
		query = query.Where("sba.is_active = ?", *req.IsActive)
	}

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	offset := (req.GetPage() - 1) * req.GetSize()
	rows := make([]*systemBankAccountRow, 0)
	if err := query.
		ColumnExpr("sba.id AS id").
		ColumnExpr("sba.bank_id AS bank_id").
		ColumnExpr("b.name_th AS bank_name_th").
		ColumnExpr("b.name_en AS bank_name_en").
		ColumnExpr("sba.account_name AS account_name").
		ColumnExpr("sba.account_no AS account_no").
		ColumnExpr("sba.branch AS branch").
		ColumnExpr("sba.is_active AS is_active").
		ColumnExpr("sba.is_default_receive AS is_default_receive").
		ColumnExpr("sba.is_default_refund AS is_default_refund").
		ColumnExpr("sba.created_at AS created_at").
		ColumnExpr("sba.updated_at AS updated_at").
		OrderExpr("sba.created_at DESC").
		Offset(int(offset)).
		Limit(int(req.GetSize())).
		Scan(ctx, &rows); err != nil {
		return nil, nil, err
	}

	data := make([]*SystemBankAccountServiceResponse, 0, len(rows))
	for _, row := range rows {
		data = append(data, &SystemBankAccountServiceResponse{
			ID:               row.ID,
			BankID:           row.BankID,
			BankNameTh:       row.BankNameTh,
			BankNameEn:       row.BankNameEn,
			AccountName:      row.AccountName,
			AccountNo:        row.AccountNo,
			Branch:           row.Branch,
			IsActive:         row.IsActive,
			IsDefaultReceive: row.IsDefaultReceive,
			IsDefaultRefund:  row.IsDefaultRefund,
			CreatedAt:        row.CreatedAt,
			UpdatedAt:        row.UpdatedAt,
		})
	}

	page := &base.ResponsePaginate{Page: req.GetPage(), Size: req.GetSize(), Total: int64(total)}
	span.AddEvent(`system_bank_accounts.svc.list.success`)
	return data, page, nil
}

func (s *Service) InfoService(ctx context.Context, id string) (*SystemBankAccountServiceResponse, error) {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`system_bank_accounts.svc.info.start`)

	itemID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	row := new(systemBankAccountRow)
	err = s.bunDB.DB().NewSelect().
		TableExpr("system_bank_accounts AS sba").
		Join("LEFT JOIN banks AS b ON b.id = sba.bank_id").
		ColumnExpr("sba.id AS id").
		ColumnExpr("sba.bank_id AS bank_id").
		ColumnExpr("b.name_th AS bank_name_th").
		ColumnExpr("b.name_en AS bank_name_en").
		ColumnExpr("sba.account_name AS account_name").
		ColumnExpr("sba.account_no AS account_no").
		ColumnExpr("sba.branch AS branch").
		ColumnExpr("sba.is_active AS is_active").
		ColumnExpr("sba.is_default_receive AS is_default_receive").
		ColumnExpr("sba.is_default_refund AS is_default_refund").
		ColumnExpr("sba.created_at AS created_at").
		ColumnExpr("sba.updated_at AS updated_at").
		Where("sba.id = ?", itemID).
		Limit(1).
		Scan(ctx, row)
	if err != nil {
		return nil, err
	}

	span.AddEvent(`system_bank_accounts.svc.info.success`)
	return &SystemBankAccountServiceResponse{
		ID:               row.ID,
		BankID:           row.BankID,
		BankNameTh:       row.BankNameTh,
		BankNameEn:       row.BankNameEn,
		AccountName:      row.AccountName,
		AccountNo:        row.AccountNo,
		Branch:           row.Branch,
		IsActive:         row.IsActive,
		IsDefaultReceive: row.IsDefaultReceive,
		IsDefaultRefund:  row.IsDefaultRefund,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}, nil
}

func (s *Service) CreateService(ctx context.Context, req *UpsertSystemBankAccountServiceRequest) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`system_bank_accounts.svc.create.start`)

	if err := s.validateUpsertRequest(ctx, req); err != nil {
		return err
	}

	id := uuid.New()
	now := time.Now()
	item := &ent.SystemBankAccountEntity{
		ID:               id,
		BankID:           req.BankID,
		AccountName:      strings.TrimSpace(req.AccountName),
		AccountNo:        strings.TrimSpace(req.AccountNo),
		Branch:           strings.TrimSpace(req.Branch),
		IsActive:         req.IsActive,
		IsDefaultReceive: req.IsDefaultReceive,
		IsDefaultRefund:  req.IsDefaultRefund,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	err := s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if req.IsDefaultReceive {
			if _, err := tx.NewUpdate().Model((*ent.SystemBankAccountEntity)(nil)).Set("is_default_receive = false").Where("1 = 1").Exec(ctx); err != nil {
				return err
			}
		}
		if req.IsDefaultRefund {
			if _, err := tx.NewUpdate().Model((*ent.SystemBankAccountEntity)(nil)).Set("is_default_refund = false").Where("1 = 1").Exec(ctx); err != nil {
				return err
			}
		}

		if _, err := tx.NewInsert().Model(item).Exec(ctx); err != nil {
			return err
		}

		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionCreated,
			ActionType:   "create_system_bank_account",
			ActionID:     id,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Created system bank account with ID " + id.String(),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		_, err := tx.NewInsert().Model(auditLog).Exec(ctx)
		return err
	})
	if err != nil {
		return normalizeDBError(err)
	}

	span.AddEvent(`system_bank_accounts.svc.create.success`)
	return nil
}

func (s *Service) UpdateService(ctx context.Context, id string, req *UpsertSystemBankAccountServiceRequest) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`system_bank_accounts.svc.update.start`)

	if err := s.validateUpsertRequest(ctx, req); err != nil {
		return err
	}

	itemID, err := uuid.Parse(id)
	if err != nil {
		return normalizeDBError(err)
	}

	item := new(ent.SystemBankAccountEntity)
	if err := s.bunDB.DB().NewSelect().Model(item).Where("id = ?", itemID).Limit(1).Scan(ctx); err != nil {
		return err
	}

	now := time.Now()
	item.BankID = req.BankID
	item.AccountName = strings.TrimSpace(req.AccountName)
	item.AccountNo = strings.TrimSpace(req.AccountNo)
	item.Branch = strings.TrimSpace(req.Branch)
	item.IsActive = req.IsActive
	item.IsDefaultReceive = req.IsDefaultReceive
	item.IsDefaultRefund = req.IsDefaultRefund
	item.UpdatedAt = now

	err = s.bunDB.DB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if req.IsDefaultReceive {
			if _, err := tx.NewUpdate().Model((*ent.SystemBankAccountEntity)(nil)).Set("is_default_receive = false").Where("id <> ?", itemID).Exec(ctx); err != nil {
				return err
			}
		}
		if req.IsDefaultRefund {
			if _, err := tx.NewUpdate().Model((*ent.SystemBankAccountEntity)(nil)).Set("is_default_refund = false").Where("id <> ?", itemID).Exec(ctx); err != nil {
				return err
			}
		}

		if _, err := tx.NewUpdate().Model(item).WherePK().Exec(ctx); err != nil {
			return err
		}

		auditLog := &ent.AuditLogEntity{
			ID:           uuid.New(),
			Action:       ent.AuditActionUpdated,
			ActionType:   "update_system_bank_account",
			ActionID:     itemID,
			Status:       ent.StatusAuditSuccesses,
			ActionDetail: "Updated system bank account with ID " + itemID.String(),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		_, err := tx.NewInsert().Model(auditLog).Exec(ctx)
		return err
	})
	if err != nil {
		return err
	}

	span.AddEvent(`system_bank_accounts.svc.update.success`)
	return nil
}

func (s *Service) DeleteService(ctx context.Context, id string) error {
	span, _ := utils.LogSpanFromContext(ctx)
	span.AddEvent(`system_bank_accounts.svc.delete.start`)

	itemID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	item := new(ent.SystemBankAccountEntity)
	if err := s.bunDB.DB().NewSelect().Model(item).Where("id = ?", itemID).Limit(1).Scan(ctx); err != nil {
		return err
	}

	if item.IsDefaultReceive || item.IsDefaultRefund {
		return errors.New("default account cannot be deleted")
	}

	if _, err := s.bunDB.DB().NewDelete().Model((*ent.SystemBankAccountEntity)(nil)).Where("id = ?", itemID).Exec(ctx); err != nil {
		return normalizeDBError(err)
	}

	span.AddEvent(`system_bank_accounts.svc.delete.success`)
	return nil
}

func (s *Service) validateUpsertRequest(ctx context.Context, req *UpsertSystemBankAccountServiceRequest) error {
	if req.BankID == uuid.Nil {
		return errors.New("bank_id is required")
	}
	if strings.TrimSpace(req.AccountName) == "" {
		return errors.New("account_name is required")
	}
	if strings.TrimSpace(req.AccountNo) == "" {
		return errors.New("account_no is required")
	}

	bankCount, err := s.bunDB.DB().NewSelect().Model((*ent.BankEntity)(nil)).Where("id = ?", req.BankID).Count(ctx)
	if err != nil {
		return err
	}
	if bankCount == 0 {
		return errors.New("bank_id not found")
	}

	return nil
}

func normalizeDBError(err error) error {
	if err == nil {
		return nil
	}
	message := strings.ToLower(err.Error())
	if strings.Contains(message, "system_bank_accounts_account_no_uidx") {
		return errors.New("account_no already exists")
	}
	if strings.Contains(message, "default_receive") {
		return errors.New("default receive account already exists")
	}
	if strings.Contains(message, "default_refund") {
		return errors.New("default refund account already exists")
	}
	if strings.Contains(message, "foreign key") {
		return errors.New("bank_id not found")
	}
	return err
}
