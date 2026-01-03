package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/wreckitral/distributed-backtesting-platform/internal/domain"
)

type strategyRepository struct {
	db *sql.DB
}

func NewStrategyRepository(db *sql.DB) *strategyRepository {
	return &strategyRepository{db: db}
}

func (r *strategyRepository) Create(ctx context.Context, s *domain.Strategy) error {
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now()
	}
	s.UpdatedAt = s.CreatedAt
	s.Version = 1

	query := `
		INSERT INTO strategies (
			name, description, language, code, version, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	err := r.db.QueryRowContext(
		ctx,
		query,
		s.Name,
		s.Description,
		s.Language.String(),
		s.Code,
		s.Version,
		s.CreatedAt,
		s.UpdatedAt,
	).Scan(&s.ID)

	if err != nil {
		return fmt.Errorf("failed to insert strategy: %w", err)
	}

	return nil
}

func (r *strategyRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Strategy, error) {
	query := `
		SELECT id, name, description, language, code, version, created_at, updated_at
		FROM strategies
		WHERE id = $1`

	s := &domain.Strategy{}
	var langStr string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.Name, &s.Description, &langStr,
		&s.Code, &s.Version, &s.CreatedAt, &s.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("strategy not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching strategy: %w", err)
	}

	s.Language = parseLanguage(langStr)
	return s, nil
}

func (r *strategyRepository) Update(ctx context.Context, strategy *domain.Strategy) error {
	strategy.UpdatedAt = time.Now()
	strategy.Version++

	query := `
		UPDATE strategies
		SET name = $1, description = $2, language = $3, code = $4,
		    version = $5, updated_at = $6
		WHERE id = $7`

	result, err := r.db.ExecContext(
		ctx,
		query,
		strategy.Name,
		strategy.Description,
		strategy.Language.String(),
		strategy.Code,
		strategy.Version,
		strategy.UpdatedAt,
		strategy.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update strategy: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("strategy not found")
	}

	return nil
}

func (r *strategyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM strategies WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete strategy: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("strategy not found")
	}

	return nil
}

func (r *strategyRepository) List(ctx context.Context, limit, offset int) ([]*domain.Strategy, error) {
	query := `
		SELECT id, name, description, language, code, version, created_at, updated_at
		FROM strategies
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing strategies: %w", err)
	}
	defer rows.Close()

	var strategies []*domain.Strategy
	for rows.Next() {
		s := &domain.Strategy{}
		var langStr string

		if err := rows.Scan(
			&s.ID, &s.Name, &s.Description, &langStr,
			&s.Code, &s.Version, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning strategy: %w", err)
		}

		s.Language = parseLanguage(langStr)
		strategies = append(strategies, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating strategies: %w", err)
	}

	return strategies, nil
}

func parseLanguage(l string) domain.StrategyLanguage {
	switch l {
	case "GO":
		return domain.StrategyLanguageGo
	case "PYTHON":
		return domain.StrategyLanguagePython
	default:
		log.Printf("Unknown language '%s', defaulting to Python", l)
		return domain.StrategyLanguagePython
	}
}
