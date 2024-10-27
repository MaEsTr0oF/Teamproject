package authgrpc

import (
	"context"

	sso "github.com/AlzoabiAF/protos/gen/go/TeamProject/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
}

type serverAPI struct {
	sso.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *sso.LoginRequest) (*sso.LoginResponse, error) {
	if err := validateLogin(req.Email, req.Password, req.GetAppId()); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.Email, req.Password, int(req.AppId))
	if err != nil {
		//TODO: implement verification logic
		return nil, err
	}

	return &sso.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	if err := validateRegister(req.Email, req.Password); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.Email, req.Password)
	if err != nil {
		//TODO: implement verfication logic
		return nil, err
	}

	return &sso.RegisterResponse{
		UserId: userID,
	}, nil
}

func validateLogin(email string, password string, appID int32) error {
	if err := validateEmailAndPassword(email, password); err != nil {
		return err
	}

	if appID == 0 {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}

func validateRegister(email string, password string) error {
	return validateEmailAndPassword(email, password)
}

func validateEmailAndPassword(email string, password string) error {
	if email == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if password == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}
