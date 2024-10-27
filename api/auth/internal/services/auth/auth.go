package auth

import (
	"TeamProject/api/auth/internal/storage"
	jwt "TeamProject/api/auth/lib/jwt"
	models2 "TeamProject/api/lib/models"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log          *slog.Logger
	userCreator  UserCreator
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserCreator interface {
	CreateUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (id int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (*models2.User, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (*models2.App, error)
}

func New(
	log *slog.Logger,
	userCreator UserCreator,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userCreator:  userCreator,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(
	ctx context.Context,
	email,
	password string,
	appID int,
) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("attempting to login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("user not found")

			return "", ErrInvalidCredentials
		}

		a.log.Error("failed to get user", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Warn("invalid credentials", err)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	const op = "Auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error("failed to generate password hash")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	user, err := a.userCreator.CreateUser(ctx, email, passHash)
	if err != nil {
		a.log.Error("failed to save user", err)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
