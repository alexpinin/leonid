package repository

import (
	"context"
	"database/sql"
	"fmt"
	"leonid/src/internal/db"
	"time"
)

const (
	passExistsQuery = `
	SELECT COUNT(*) > 0
	FROM pass
	WHERE pass = $1
	  AND valid_by >= $2
`
	passDeleteQuery = `
	DELETE FROM pass 
	WHERE pass = $1
`
)

type PassRepository struct {
	db *db.DB
}

func NewPassRepository(db *db.DB) *PassRepository {
	return &PassRepository{
		db: db,
	}
}

func (r *PassRepository) PassExists(ctx context.Context, tx *sql.Tx, pass string, byDate time.Time) (bool, error) {
	row := r.db.QueryRowTx(ctx, tx, passExistsQuery, pass, byDate.UTC().Second())
	res := false
	err := row.Scan(&res)
	if err != nil {
		return false, fmt.Errorf("PassRepository.PassExists: %w", err)
	}
	return res, nil
}

func (r *PassRepository) DeletePass(ctx context.Context, tx *sql.Tx, pass string) error {
	_, err := r.db.ExecTx(ctx, tx, passDeleteQuery, pass)
	if err != nil {
		return fmt.Errorf("PassRepository.PassDelete: %w", err)
	}
	return nil
}
