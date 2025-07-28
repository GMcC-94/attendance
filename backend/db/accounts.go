package db

import (
	"database/sql"
	"fmt"

	"github.com/gmcc94/attendance-go/helpers"
	"github.com/gmcc94/attendance-go/types"
)

type AccountsStore interface {
	AddAccountEntries(entries []types.AccountEntry, category string) error
	GetAccounts() ([]types.AccountEntry, error)
}

type PostgresAccountsStore struct {
	DB *sql.DB
}

func (p *PostgresAccountsStore) AddAccountEntries(entries []types.AccountEntry, category string) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	for _, e := range entries {
		if _, err := tx.Exec(`
		INSERT INTO club_accounts (description, amount, category, created_at)
		VALUES ($1, $2, $3)`, e.Description, e.Amount, category); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert account entry: %w", err)
		}
	}
	return tx.Commit()
}

func (p *PostgresAccountsStore) GetAccounts() ([]types.AccountEntry, error) {
	rows, err := p.DB.Query(`
	SELECT id, description, amount, category, created_at
	FROM club_accounts
	ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []types.AccountEntry
	for rows.Next() {
		var e types.AccountEntry
		if err := rows.Scan(
			&e.ID,
			&e.Description,
			&e.Amount,
			&e.Category,
			&e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, nil
}

func (s *PostgresAccountsStore) GetGroupedAccounts() ([]types.GroupedAccounts, error) {
	entries, err := s.GetAccounts()
	if err != nil {
		return nil, err
	}
	grouped := helpers.GroupedAccounts(entries)
	return helpers.BuildGroupedResponse(grouped), nil
}
