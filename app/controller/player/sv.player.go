package player

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

func (s *Service) Create(ctx context.Context, req request.CreatePlayer) (*model.Player, bool, error) {

	m := &model.Player{
		Prefix:    req.Prefix,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		MemberID:  req.MemberID,
		Position:  req.Position,
		RoomID:    req.RoomID,
		IsActive:  req.IsActive,
	}

	_, err := s.db.NewInsert().Model(m).Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("player already exists")
		}
	}

	return m, false, err
}

func (s *Service) Update(ctx context.Context, req request.UpdatePlayer, id request.GetByIDPlayer) (*model.Player, bool, error) {
	ex, err := s.db.NewSelect().Table("players").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, false, err
	}

	m := &model.Player{
		ID:        id.ID,
		Prefix:    req.Prefix,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		MemberID:  req.MemberID,
		Position:  req.Position,
		RoomID:    req.RoomID,
		IsActive:  req.IsActive,
	}
	logger.Info(m)
	m.SetUpdateNow()
	_, err = s.db.NewUpdate().Model(m).
		Set("prefix = ?prefix, first_name = ?first_name, last_name = ?last_name, member_id = ?member_id, position = ?position, room_id = ?room_id, is_active = ?is_active").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("player already exists")
		}
	}
	return m, false, err
}

func (s *Service) List(ctx context.Context, req request.ListPlayer) ([]response.ListPlayer, int, error) {
	offset := (req.Page - 1) * req.Size

	m := []response.ListPlayer{}
	query := s.db.NewSelect().
		TableExpr("players AS p").
		Column("p.id", "p.prefix", "p.first_name", "p.last_name", "p.member_id", "p.position", "p.room_id", "p.is_active").
		Where("deleted_at IS NULL")

	if req.Search != "" {
		search := fmt.Sprintf("%" + strings.ToLower(req.Search) + "%")
		if req.SearchBy != "" {
			search := strings.ToLower(req.Search)
			query.Where(fmt.Sprintf("LOWER(p.%s) LIKE ?", req.SearchBy), search)
		} else {
			query.Where("LOWER(p.first_name) LIKE ?", search)
		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("p.%s %s", req.SortBy, req.OrderBy)

	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}
	return m, count, err
}

func (s *Service) Get(ctx context.Context, id request.GetByIDPlayer) (*response.ListPlayer, error) {
	m := response.ListPlayer{}
	err := s.db.NewSelect().
		TableExpr("players AS p").
		Column("p.id", "p.prefix", "p.first_name", "p.last_name", "p.member_id", "p.position", "p.room_id", "p.is_active").
		Where("id = ?", id.ID).
		Where("deleted_at IS NULL").
		Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIDPlayer) error {
	ex, err := s.db.NewSelect().Table("players").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return err
	}

	if !ex {
		return errors.New("player not found")
	}

	// data, err := s.db.NewDelete().Table("room").Where("id = ?", id.ID).Exec(ctx)
	_, err = s.db.NewDelete().Model((*model.Player)(nil)).Where("id = ?", id.ID).Exec(ctx)
	return err
}
