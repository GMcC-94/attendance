package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gmcc94/attendance-go/internal/types"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserStore interface {
	CreateUser(ctx context.Context, username, password string) (int, error)
	AuthenticateUser(ctx context.Context, username, password string) (*types.User, error)
	UpdatePassword(ctx context.Context, userID int, newPassword string) error
}

type PostgresUserStore struct {
	DB     *sql.DB
	logger *slog.Logger
}

func NewPostgresUserStore(db *sql.DB, logger *slog.Logger) *PostgresUserStore {
	if logger == nil {
		logger = slog.Default()
	}
	return &PostgresUserStore{
		DB:     db,
		logger: logger,
	}
}

func (p *PostgresUserStore) CreateUser(ctx context.Context, username, password string) (int, error) {
	const op = "PostgresUserStore.CreateUser"

	if strings.TrimSpace(username) == "" {
		return 0, fmt.Errorf("%s: username cannot be empty", op)
	}
	if len(password) < 8 {
		return 0, fmt.Errorf("%s: password must be at least 8 characters long", op)
	}

	username = strings.ToLower(strings.TrimSpace(username))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		p.logger.Error("failed to hash password",
			slog.String("operation", op),
			slog.String("username", username),
			slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: failed to hash password: %w", op, err)
	}

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var userID int
	query := `
	INSERT INTO users (username, password_hash, created_at)
	VALUES ($1, $2, NOW())
	RETURNING id`

	err = p.DB.QueryRowContext(ctx, query, string(hashedPassword)).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				p.logger.Info("username already exists",
					slog.String("operation", op),
					slog.String("username", username))
				return 0, types.ErrUsernameTaken
			}
		}

		p.logger.Error("failed to insert user",
			slog.String("operation", op),
			slog.String("username", username),
			slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: failed to create user: %w", op, err)
	}

	p.logger.Info("user created successfully",
		slog.String("operation", op),
		slog.String("username", username),
		slog.Int("userID", userID))

	return userID, nil
}

func (p *PostgresUserStore) AuthenticateUser(ctx context.Context, username, password string) (*types.User, error) {
	const op = "PostgresUserStore.AuthenticateUser"

	// Validate input
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	username = strings.ToLower(strings.TrimSpace(username))

	// Add timeout to context if not already set
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user := &types.User{
		Username: username,
	}

	query := `SELECT id, username, password_hash, created_at, updated_at 
			  FROM users 
			  WHERE username = $1`

	err := p.DB.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p.logger.Info("authentication failed - user not found",
				slog.String("operation", op),
				slog.String("username", username))
			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		p.logger.Error("failed to query user",
			slog.String("operation", op),
			slog.String("username", username),
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: database error: %w", op, err)
	}

	// Compare password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		p.logger.Info("authentication failed - invalid password",
			slog.String("operation", op),
			slog.String("username", username))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	// Clear password hash before returning user (security best practice)
	user.PasswordHash = ""

	p.logger.Info("user authenticated successfully",
		slog.String("operation", op),
		slog.String("username", username),
		slog.Int("user_id", user.ID))

	return user, nil
}

func (p *PostgresUserStore) UpdatePassword(ctx context.Context, userID int, newPassword string) error {
	const op = "PostgresUserStore.UpdatePassword"

	if userID <= 0 {
		return fmt.Errorf("%s: invalid user ID", op)
	}
	if len(newPassword) < 8 {
		return fmt.Errorf("%s: password must be at least 8 characters", op)
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		p.logger.Error("failed to hash new password",
			slog.String("operation", op),
			slog.Int("user_id", userID),
			slog.String("error", err.Error()))
		return fmt.Errorf("%s: failed to hash password: %w", op, err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `UPDATE users 
			  SET password_hash = $1, updated_at = NOW() 
			  WHERE id = $2`

	result, err := p.DB.ExecContext(ctx, query, string(hashedPassword), userID)
	if err != nil {
		p.logger.Error("failed to update password",
			slog.String("operation", op),
			slog.Int("user_id", userID),
			slog.String("error", err.Error()))
		return fmt.Errorf("%s: failed to update password: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		p.logger.Error("failed to get rows affected",
			slog.String("operation", op),
			slog.Int("user_id", userID),
			slog.String("error", err.Error()))
		return fmt.Errorf("%s: failed to verify update: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}

	p.logger.Info("password updated successfully",
		slog.String("operation", op),
		slog.Int("user_id", userID))

	return nil
}
