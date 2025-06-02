package repository

import (
	"context"
	"database/sql"
	"fmt"
	"leonid/src/internal/db"
	"strings"
	"time"
)

const (
	passExistsQuery = `
	SELECT pass, valid_by, nicknames, system_prompt
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

func (r *PassRepository) FindPass(ctx context.Context, tx *sql.Tx, pass string, byDate time.Time) (Pass, error) {
	row := r.db.QueryRowTx(ctx, tx, passExistsQuery, pass, byDate.UTC().Second())

	var p Pass
	var validBy int64
	var nicknames string
	err := row.Scan(&p.Pass, &validBy, &nicknames, &p.SystemPrompt)
	if err != nil {
		return Pass{}, fmt.Errorf("PassRepository.FindPass: %w", err)
	}

	p.ValidBy = time.Unix(validBy, 0)
	p.Nicknames = strings.Split(nicknames, ",")

	return p, nil
}

func (r *PassRepository) DeletePass(ctx context.Context, tx *sql.Tx, pass string) error {
	_, err := r.db.ExecTx(ctx, tx, passDeleteQuery, pass)
	if err != nil {
		return fmt.Errorf("PassRepository.PassDelete: %w", err)
	}
	return nil
}

type Pass struct {
	Pass         string
	ValidBy      time.Time
	Nicknames    []string
	SystemPrompt string
}
