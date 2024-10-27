package postgres

import (
	"TeamProject/api/auth/internal/storage"
	models2 "TeamProject/api/lib/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(host string, user string, password string, dbname string) (*Storage, error) {
	const op = "storage.postgres.New"

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sllmode=disable", host, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) CreateUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.CreateUser"

	var id int64

	stmt, err := s.db.Prepare("INSERT INTO users (email, passHash) VALUES (%s, %s) returning id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqlErr *pq.Error
		if errors.As(err, &sqlErr) && sqlErr.Error() == pq.ErrCouldNotDetectUsername.Error() {
		}
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (*models2.User, error) {
	const op = "storage.postgres.CheckUser"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM user WHERE email = ?")
	if err != nil {
		return &models2.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models2.User

	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models2.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return &models2.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil

}

func (s *Storage) App(ctx context.Context, appID int) (*models2.App, error) {
	op := "storage.postgres.CheckApp"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM app WHERE app_id = ?")
	if err != nil {
		return &models2.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, appID)

	var app models2.App

	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models2.App{}, storage.ErrAppNotFound
		}

		return &models2.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return &app, nil
}

func (s *Storage) DeleteUser(ctx context.Context, email string) error {
	op := "storage.postgres.DeleteUser"

	stmt, err := s.db.Prepare("DELETE FROM user WHERE email = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
