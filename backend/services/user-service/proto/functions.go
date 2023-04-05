package user_service

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"math/rand"
	customErr "userService/libs/errors"
	jwt_user "userService/services/user-service/jwt"
	"userService/services/user-service/postgres"
	user_service "userService/services/user-service/proto/user-service"
)

type Server struct {
	user_service.UnimplementedUserServiceServer
}

var db postgres.Database

func Init() error {
	var err error

	db, err = postgres.NewDatabase()
	return err
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (s Server) RandomPrivateMethod(ctx context.Context,
	r *user_service.RandomPrivateMethodRequest,
) (*user_service.RandomPrivateMethodResponse, error) {
	return &user_service.RandomPrivateMethodResponse{Msg: "ok :)"}, nil
}

func (s Server) RefreshPassword(ctx context.Context,
	r *user_service.RefreshPasswordRequest,
) (*user_service.RefreshPasswordResponse, error) {
	err := r.Validate()
	if err != nil {
		return nil, customErr.LogError(err)
	}

	err = db.CheckCode(r.Code)
	if err != nil {
		return nil, customErr.LogError(err)
	}

	user := &user_service.User{
		Username: r.Username,
		Password: r.Password,
	}

	err = db.RefreshPassword(user)
	if err != nil {
		return nil, customErr.LogError(err)
	}

	return &user_service.RefreshPasswordResponse{Msg: "password changed"}, err
}

func (s Server) GetCode(ctx context.Context,
	r *user_service.GetCodeRequest,
) (*user_service.GetCodeResponse, error) {
	err := r.Validate()
	if err != nil {
		return nil, customErr.LogError(err)
	}

	err = db.CheckUser(r.Username)
	if err != nil {
		return nil, customErr.LogError(err)
	}

	token, err := jwt_user.GetRefreshPasswordToken()
	if err != nil {
		return nil, customErr.LogError(err)
	}

	code := RandStringBytes(6)

	err = db.CreateNotification(r.Username, code)
	if err != nil {
		return nil, customErr.LogError(err)
	}

	return &user_service.GetCodeResponse{
		Code:  code,
		Token: token,
	}, err
}

func (s Server) SignUp(ctx context.Context,
	r *user_service.SignUpRequest,
) (*user_service.SignUpResponse, error) {
	err := r.Validate()
	if err != nil {
		return nil, customErr.LogError(err)
	}

	user := &user_service.User{
		Username: r.Username,
		Password: r.Password,
	}

	out, err := db.SignUp(user)
	if err != nil {
		return nil, customErr.LogError(err)
	}

	return &user_service.SignUpResponse{Id: out.Id}, err

}

func (s Server) Refresh(ctx context.Context,
	r *user_service.RefreshRequest,
) (*user_service.RefreshResponse, error) {
	err := r.Validate()
	if err != nil {
		return nil, customErr.LogError(err)
	}

	err, _ = jwt_user.CheckRefreshToken(r.Token)
	if err != nil {
		return nil, customErr.LogError(err)
	}

	token, refresh, err := jwt_user.ForwardRefresh()
	if err != nil {
		return nil, customErr.LogError(err)
	}

	tokenOut := &user_service.Token{
		Access:  token,
		Refresh: refresh,
	}

	return &user_service.RefreshResponse{
		Token: tokenOut,
	}, nil
}

func (s Server) SignIn(ctx context.Context,
	r *user_service.SignInRequest,
) (*user_service.SignInResponse, error) {
	err := r.Validate()
	if err != nil {
		return nil, customErr.LogError(err)
	}

	request := &user_service.User{
		Username: r.Username,
		Password: r.Password,
	}

	userOut, err := db.SignIn(request)
	if err != nil {
		return nil, customErr.LogError(err)
	}

	token, refresh, err := jwt_user.GetSignInToken()
	if err != nil {
		return nil, customErr.LogError(err)
	}

	tokenOut := &user_service.Token{
		Access:  token,
		Refresh: refresh,
	}

	return &user_service.SignInResponse{
		Id:       userOut.Id,
		Username: userOut.Username,
		Token:    tokenOut,
	}, err
}

func (s *Server) LogOut(ctx context.Context,
	r *user_service.LogOutRequest,
) (*user_service.LogOutResponse, error) {
	err := r.Validate()
	if err != nil {
		return nil, customErr.LogError(err)
	}

	out := &user_service.LogOutResponse{}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Printf("md: %v\n", md)
	}

	header := metadata.Pairs("authorization", "")

	err = grpc.SendHeader(ctx, header)
	if err != nil {
		return nil, customErr.LogError(err)
	}

	return out, nil
}
