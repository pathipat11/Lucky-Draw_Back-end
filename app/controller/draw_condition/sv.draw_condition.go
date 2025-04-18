package draw_condition

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"app/internal/logger"
	"context"
	"errors"
	"fmt"
	"strings"
)

func (s *Service) Create(ctx context.Context, req request.CreateDrawCondition) (*model.DrawCondition, bool, error) {

	m := &model.DrawCondition{
		RoomID:         req.RoomID,
		PrizeID:        req.PrizeID,
		FilterStatus:   req.FilterStatus,
		FilterPosition: req.FilterPosition,
		Quantity:       int64(req.Quantity),
	}

	_, err := s.db.NewInsert().Model(m).Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("draw_condition already exists")
		}
	}

	return m, false, err
}

func (s *Service) Update(ctx context.Context, req request.UpdateDrawCondition, id request.GetByIDDrawCondition) (*model.DrawCondition, bool, error) {
	ex, err := s.db.NewSelect().Table("draw_conditions").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, false, err
	}

	m := &model.DrawCondition{
		ID:             id.ID,
		RoomID:         req.RoomID,
		PrizeID:        req.PrizeID,
		FilterStatus:   req.FilterStatus,
		FilterPosition: req.FilterPosition,
		Quantity:       int64(req.Quantity),
	}
	logger.Info(m)
	m.SetUpdateNow()
	_, err = s.db.NewUpdate().Model(m).
		Set("room_id = ?room_id, prize_id = ?prize_id, filter_status = ?filter_status, filter_position = ?filter_position, quantity = ?quantity").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("draw_conditions already exists")
		}
	}
	return m, false, err
}

func (s *Service) List(ctx context.Context, req request.ListDrawCondition) ([]response.ListDrawCondition, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []response.ListDrawCondition{}
	query := s.db.NewSelect().
		TableExpr("draw_conditions AS d").
		Column("d.id", "d.room_id", "d.prize_id", "d.filter_status", "d.filter_position", "d.quantity").
		Where("deleted_at IS NULL")

	if req.Search != "" {
		search := fmt.Sprintf("%" + strings.ToLower(req.Search) + "%")
		if req.SearchBy != "" {
			search := strings.ToLower(req.Search)
			query.Where(fmt.Sprintf("LOWER(d.%s) LIKE ?", req.SearchBy), search)
		} else {
			query.Where("LOWER(d.name) LIKE ?", search)
		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("d.%s %s", req.SortBy, req.OrderBy)

	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err
}

func (s *Service) Get(ctx context.Context, id request.GetByIDDrawCondition) (*response.ListDrawCondition, error) {
	m := response.ListDrawCondition{}
	err := s.db.NewSelect().
		TableExpr("draw_conditions AS d").
		Column("d.id", "d.room_id", "d.prize_id", "d.filter_status", "d.filter_position", "d.quantity").
		Where("id = ?", id.ID).
		Where("deleted_at IS NULL").
		Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIDDrawCondition) error {
	ex, err := s.db.NewSelect().Table("draw_conditions").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return err
	}

	if !ex {
		return errors.New("draw_condition not found")
	}

	// data, err := s.db.NewDelete().Table("room").Where("id = ?", id.ID).Exec(ctx)
	_, err = s.db.NewDelete().Model((*model.DrawCondition)(nil)).Where("id = ?", id.ID).Exec(ctx)
	return err
}
