package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/shaunchuang/food-roulette-backend/internal/domain"
	"github.com/shaunchuang/food-roulette-backend/pkg/logger"
	"go.uber.org/zap"
)

// GameRepository PostgreSQL 遊戲會話資料庫操作實作
type GameRepository struct {
	db *sql.DB
}

// NewGameRepository 建立遊戲 Repository
func NewGameRepository(db *sql.DB) *GameRepository {
	return &GameRepository{
		db: db,
	}
}

// CreateSession 建立遊戲會話
func (r *GameRepository) CreateSession(ctx context.Context, session *domain.GameSession) error {
	query := `
		INSERT INTO game_sessions (id, user_id, game_type, status, started_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		session.ID,
		session.UserID,
		session.GameType,
		session.Status,
		session.StartedAt,
		now,
	)

	if err != nil {
		logger.Error("建立遊戲會話失敗", zap.Error(err), zap.String("session_id", session.ID))
		return err
	}

	logger.Info("遊戲會話建立成功", zap.String("session_id", session.ID), zap.Int("user_id", session.UserID))
	return nil
}

// GetSessionByID 根據 ID 取得遊戲會話
func (r *GameRepository) GetSessionByID(ctx context.Context, sessionID string) (*domain.GameSession, error) {
	query := `
		SELECT id, user_id, game_type, status, result_restaurant_id, started_at, completed_at, created_at
		FROM game_sessions
		WHERE id = $1`

	session := &domain.GameSession{}
	var resultRestaurantID sql.NullInt64
	var completedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.GameType,
		&session.Status,
		&resultRestaurantID,
		&session.StartedAt,
		&completedAt,
		&session.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("遊戲會話不存在")
		}
		logger.Error("取得遊戲會話失敗", zap.Error(err), zap.String("session_id", sessionID))
		return nil, err
	}

	// 處理可為空的欄位
	if resultRestaurantID.Valid {
		restaurantID := int(resultRestaurantID.Int64)
		session.ResultRestaurantID = &restaurantID
	}
	if completedAt.Valid {
		session.CompletedAt = &completedAt.Time
	}

	return session, nil
}

// UpdateSession 更新遊戲會話
func (r *GameRepository) UpdateSession(ctx context.Context, session *domain.GameSession) error {
	query := `
		UPDATE game_sessions
		SET status = $1, result_restaurant_id = $2, completed_at = $3
		WHERE id = $4`

	var resultRestaurantID interface{}
	if session.ResultRestaurantID != nil && *session.ResultRestaurantID > 0 {
		resultRestaurantID = *session.ResultRestaurantID
	}

	_, err := r.db.ExecContext(ctx, query,
		session.Status,
		resultRestaurantID,
		session.CompletedAt,
		session.ID,
	)

	if err != nil {
		logger.Error("更新遊戲會話失敗", zap.Error(err), zap.String("session_id", session.ID))
		return err
	}

	logger.Info("遊戲會話更新成功", zap.String("session_id", session.ID))
	return nil
}

// GetUserSessions 取得使用者的遊戲歷史
func (r *GameRepository) GetUserSessions(ctx context.Context, userID int, limit, offset int) ([]domain.GameSession, error) {
	query := `
		SELECT id, user_id, game_type, status, result_restaurant_id, 
		       started_at, completed_at, created_at
		FROM game_sessions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		logger.Error("取得使用者遊戲歷史失敗", zap.Error(err), zap.Int("user_id", userID))
		return nil, err
	}
	defer rows.Close()

	var sessions []domain.GameSession
	for rows.Next() {
		var session domain.GameSession
		var resultRestaurantID sql.NullInt64
		var completedAt sql.NullTime

		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.GameType,
			&session.Status,
			&resultRestaurantID,
			&session.StartedAt,
			&completedAt,
			&session.CreatedAt,
		)

		if err != nil {
			logger.Error("掃描遊戲會話資料失敗", zap.Error(err))
			continue
		}

		// 處理可為空的欄位
		if resultRestaurantID.Valid {
			restaurantID := int(resultRestaurantID.Int64)
			session.ResultRestaurantID = &restaurantID
		}
		if completedAt.Valid {
			session.CompletedAt = &completedAt.Time
		}

		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		logger.Error("處理遊戲會話查詢結果失敗", zap.Error(err))
		return nil, err
	}

	logger.Info("取得使用者遊戲歷史成功", zap.Int("user_id", userID), zap.Int("count", len(sessions)))
	return sessions, nil
}
