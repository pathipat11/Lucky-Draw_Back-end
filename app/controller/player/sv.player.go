package player

import (
	"app/app/model"
	"app/app/request"
	"app/app/response"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/uptrace/bun"
	"github.com/xuri/excelize/v2"
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
		Status:    req.Status,
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
	player := &model.Player{}
	err := s.db.NewSelect().Model(player).Where("id = ?", id.ID).Scan(ctx)
	if err != nil {
		return nil, false, err
	}

	if req.Prefix != "" {
		player.Prefix = req.Prefix
	}
	if req.FirstName != "" {
		player.FirstName = req.FirstName
	}
	if req.LastName != "" {
		player.LastName = req.LastName
	}
	if req.MemberID != "" {
		player.MemberID = req.MemberID
	}
	if req.Position != "" {
		player.Position = req.Position
	}
	if req.RoomID != "" {
		player.RoomID = req.RoomID
	}
	if req.Status != "" {
		player.Status = req.Status
	}

	player.IsActive = req.IsActive

	player.SetUpdateNow()

	_, err = s.db.NewUpdate().Model(player).
		Set("prefix = ?prefix, first_name = ?first_name, last_name = ?last_name, member_id = ?member_id, position = ?position, room_id = ?room_id, is_active = ?is_active, status = ?status").
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, true, errors.New("player already exists")
		}
	}

	return player, false, err
}

func (s *Service) List(ctx context.Context, req request.ListPlayer) ([]response.ListPlayer, int, error) {
	offset := (req.Page - 1) * req.Size
	m := []response.ListPlayer{}

	query := s.db.NewSelect().
		TableExpr("players AS p").
		Column("p.id", "p.prefix", "p.first_name", "p.last_name", "p.member_id", "p.position", "p.room_id", "p.is_active", "p.status", "p.created_at").
		ColumnExpr("r.name AS room_name").
		Join("LEFT JOIN rooms AS r ON r.id = p.room_id::uuid").
		Where("p.deleted_at IS NULL")

	if req.RoomID != "" {
		query.Where("p.room_id = ?", req.RoomID)
	}

	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		query = query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("LOWER(p.prefix) LIKE ?", search).
				WhereOr("LOWER(p.first_name) LIKE ?", search).
				WhereOr("LOWER(p.last_name) LIKE ?", search)
		})
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// ถ้าไม่มีการส่ง sort_by หรือส่งแบบไม่ระบุที่แน่ชัด ให้ default เป็น created_at asc
	sortBy := req.SortBy
	orderBy := req.OrderBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	if orderBy != "desc" {
		orderBy = "asc"
	}
	order := fmt.Sprintf("p.%s %s", sortBy, orderBy)

	err = query.Order(order).Limit(req.Size).Offset(offset).Scan(ctx, &m)
	if err != nil {
		return nil, 0, err
	}

	// ตรวจสอบว่ามีข้อมูลหรือไม่
	if len(m) == 0 {
		return nil, 0, fmt.Errorf("no players found")
	}

	return m, count, nil
}

func (s *Service) Get(ctx context.Context, id request.GetByIDPlayer) (*response.ListPlayer, error) {
	m := response.ListPlayer{}
	err := s.db.NewSelect().
		TableExpr("players AS p").
		Column("p.id", "p.prefix", "p.first_name", "p.last_name", "p.member_id", "p.position", "p.room_id", "p.is_active", "p.status").
		ColumnExpr("r.name AS room_name").
		Join("LEFT JOIN rooms AS r ON r.id = p.room_id::uuid").
		Where("p.id = ?::uuid", id.ID).
		Where("p.deleted_at IS NULL").
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

// new function
func (s *Service) ImportPlayersFromCSV(ctx context.Context, file io.Reader, roomID string) error {
	reader := csv.NewReader(file)
	_, err := reader.Read() // skip header
	if err != nil {
		return err
	}

	var failedLines []string
	lineNumber := 2 // เริ่มจาก 2 เพราะแถวแรกคือ header

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV on line %d: %v", lineNumber, err)
		}

		if len(record) < 7 {
			failedLines = append(failedLines, fmt.Sprintf("line %d: not enough columns", lineNumber))
			lineNumber++
			continue
		}

		isActive := false
		activeStr := strings.TrimSpace(strings.ToLower(record[6]))
		if activeStr == "true" || activeStr == "1" || activeStr == "yes" || activeStr == "เข้า" || activeStr == "เข้าร่วม" {
			isActive = true
		}

		player := &model.Player{
			Prefix:    strings.TrimSpace(record[1]),
			FirstName: strings.TrimSpace(record[2]),
			LastName:  strings.TrimSpace(record[3]),
			MemberID:  strings.TrimSpace(record[4]),
			Position:  strings.TrimSpace(record[5]),
			RoomID:    roomID,
			IsActive:  isActive,
			Status:    "not_received",
		}

		_, err = s.db.NewInsert().Model(player).Exec(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value") {
				failedLines = append(failedLines, fmt.Sprintf("line %d: duplicate member_id (%s)", lineNumber, player.MemberID))
			} else {
				failedLines = append(failedLines, fmt.Sprintf("line %d: %v", lineNumber, err))
			}
		}

		lineNumber++
	}

	if len(failedLines) > 0 {
		return fmt.Errorf("some rows failed to import:\n%s", strings.Join(failedLines, "\n"))
	}

	return nil
}

func (s *Service) ImportPlayersFromXLSX(ctx context.Context, file io.Reader, roomID string) error {
	xl, err := excelize.OpenReader(file)
	if err != nil {
		return fmt.Errorf("failed to read Excel file: %v", err)
	}

	sheetName := xl.GetSheetName(0)
	rows, err := xl.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to read rows from sheet: %v", err)
	}

	var failedLines []string

	for i, row := range rows {
		if i == 0 {
			continue // ข้าม header
		}
		if len(row) < 7 {
			failedLines = append(failedLines, fmt.Sprintf("line %d: not enough columns", i+1))
			continue
		}

		isActive := false
		activeStr := strings.TrimSpace(strings.ToLower(row[6]))
		if activeStr == "true" || activeStr == "1" || activeStr == "yes" || activeStr == "เข้า" || activeStr == "เข้าร่วม" {
			isActive = true
		}

		player := &model.Player{
			Prefix:    strings.TrimSpace(row[1]),
			FirstName: strings.TrimSpace(row[2]),
			LastName:  strings.TrimSpace(row[3]),
			MemberID:  strings.TrimSpace(row[4]),
			Position:  strings.TrimSpace(row[5]),
			RoomID:    roomID,
			IsActive:  isActive,
			Status:    "not_received",
		}

		_, err := s.db.NewInsert().Model(player).Exec(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value") {
				failedLines = append(failedLines, fmt.Sprintf("line %d: duplicate member_id (%s)", i+1, player.MemberID))
			} else {
				failedLines = append(failedLines, fmt.Sprintf("line %d: %v", i+1, err))
			}
		}
	}

	if len(failedLines) > 0 {
		return fmt.Errorf("some rows failed to import:\n%s", strings.Join(failedLines, "\n"))
	}

	return nil
}
