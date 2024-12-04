package interceptors

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/bbquite/go-pass-keeper/internal/utils"
	jwttoken "github.com/bbquite/go-pass-keeper/pkg/jwt_token"
	"github.com/golang-jwt/jwt/v4"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtManager    *jwttoken.JWTManager
	noAuthMethods []string
}

func NewAuthInterceptor(jwtManager *jwttoken.JWTManager, noAuthMethods []string) *AuthInterceptor {
	return &AuthInterceptor{
		jwtManager:    jwtManager,
		noAuthMethods: noAuthMethods,
	}
}

func (ai *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if slices.Contains(ai.noAuthMethods, info.FullMethod) {
			return handler(ctx, req)
		}

		userID, err := ai.authorize(ctx)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, utils.UserIDKey, userID)

		return handler(ctx, req)
	}
}

func (ai *AuthInterceptor) authorize(ctx context.Context) (uint32, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Error(codes.Unauthenticated, "Metadata not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return 0, status.Error(codes.Unauthenticated, "Invalid access token")
	}

	accessToken := values[0]
	accessToken = strings.TrimPrefix(accessToken, "Bearer ")

	authorized, err := ai.jwtManager.IsAuthorized(accessToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, status.Error(codes.Unauthenticated, "TokenExpired")
		}
		return 0, status.Error(codes.Internal, "unexpected error")
	}

	if authorized {
		userID, err := ai.jwtManager.ExtractIDFromToken(accessToken)
		if err != nil {
			return 0, status.Error(codes.Unauthenticated, "Invalid access token")
		}
		return userID, nil
	}

	return 0, status.Error(codes.Unauthenticated, "TokenExpired")
}
