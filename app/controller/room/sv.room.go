package room

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"app/internal/logger"
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Create(ctx context.Context, req request.CreateRoom) (*model.Room, bool, error) {
	var hashedPassword string
	var err error

	// Hash password if provided
	if req.Password != "" {
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, false, errors.New("failed to hash password")
		}
		hashedPassword = string(hashedBytes)
	}

	m := &model.Room{
		Name:     req.Name,
		Password: hashedPassword,
	}

	_, err = s.db.NewInsert().Model(m).Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("room already exists")
		}
		return nil, false, err
	}

	return m, false, nil
}

func (s *Service) Update(ctx context.Context, req request.UpdateRoom, id request.GetByIDRoom) (*model.Room, bool, error) {
	ex, err := s.db.NewSelect().Table("rooms").Where("id = ?", id.ID).Exists(ctx)
	if err != nil {
		return nil, false, err
	}

	if !ex {
		return nil, false, err
	}

	m := &model.Room{
		ID:   id.ID,
		Name: req.Name,
	}
	logger.Info(m)
	m.SetUpdateNow()
	_, err = s.db.NewUpdate().Model(m).
		Set("name = ?name").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("room already exists")
		}
	}
	return m, false, err
}

func (s *Service) List(ctx context.Context, req request.ListRoom) ([]response.ListRoom, int, error) {
	offset := (req.Page - 1) * req.Size

	type roomRaw struct {
		ID       string  `bun:"id"`
		Name     string  `bun:"name"`
		Password *string `bun:"password"` // ใช้ pointer เพื่อเช็ค null
	}

	rawRooms := []roomRaw{}

	query := s.db.NewSelect().
		TableExpr("rooms AS r").
		Column("r.id", "r.name", "r.password").
		Where("r.deleted_at IS NULL").
		Where("EXISTS (?)", s.db.NewSelect().
			Table("players").
			Where("players.room_id = r.id::text").
			Where("players.deleted_at IS NULL")).
		Where("EXISTS (?)", s.db.NewSelect().
			Table("prizes").
			Where("prizes.room_id = r.id::text").
			Where("prizes.deleted_at IS NULL"))

	if req.Search != "" {
		search := fmt.Sprintf("%%%s%%", strings.ToLower(req.Search))
		if req.SearchBy != "" {
			query.Where(fmt.Sprintf("LOWER(r.%s) LIKE ?", req.SearchBy), search)
		} else {
			query.Where("LOWER(r.name) LIKE ?", search)
		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	order := fmt.Sprintf("r.%s %s", req.SortBy, req.OrderBy)
	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &rawRooms)
	if err != nil {
		return nil, 0, err
	}

	// แปลงผลลัพธ์เป็น response.ListRoom
	result := make([]response.ListRoom, 0, len(rawRooms))
	for _, r := range rawRooms {
		hasPassword := r.Password != nil && *r.Password != ""
		result = append(result, response.ListRoom{
			ID:          r.ID,
			Name:        r.Name,
			HasPassword: hasPassword,
		})
	}

	return result, count, nil
}

func (s *Service) Get(ctx context.Context, id request.GetByIDRoom) (*response.ListRoom, error) {
	m := response.ListRoom{}
	err := s.db.NewSelect().
		TableExpr("rooms AS u").
		Column("u.id", "u.name").
		Where("id = ?", id.ID).Where("deleted_at IS NULL").Scan(ctx, &m)
	return &m, err
}

func (s *Service) Delete(ctx context.Context, id request.GetByIDRoom) error {
	ex, err := s.db.NewSelect().Table("rooms").Where("id = ?", id.ID).Where("deleted_at IS NULL").Exists(ctx)
	if err != nil {
		return err
	}

	if !ex {
		return errors.New("room not found")
	}

	// data, err := s.db.NewDelete().Table("room").Where("id = ?", id.ID).Exec(ctx)
	_, err = s.db.NewDelete().Model((*model.Room)(nil)).Where("id = ?", id.ID).Exec(ctx)
	return err
}

// new function

func (s *Service) Login(ctx context.Context, req request.LoginRoom) (*model.Room, error) {
	var room model.Room

	err := s.db.NewSelect().
		Model(&room).
		Where("id = ?", req.ID).
		Where("deleted_at IS NULL").
		Scan(ctx)

	if err != nil {
		return nil, errors.New("room not found")
	}

	// เช็คว่าห้องนี้ตั้งรหัสไว้ไหม
	if room.Password == "" || room.Password == "null" {
		return &room, nil
	}

	// ตรวจสอบรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(room.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return &room, nil
}
