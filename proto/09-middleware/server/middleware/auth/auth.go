package auth

import (
	"context"
	"errors"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Token 用户信息.
type TokenInfo struct {
	ID    string
	Roles []string
}

// AuthInterceptor 认证拦截器，对以 authorization 为头部，
// 形式为 `bearer token` 的 Token 进行验证.
func AuthInterceptor(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	tokenInfo, err := parseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, " %v", err)
	}
	// 使用 context.WithValue 添加了值后，可以使用 Value(key) 方法获取值.
	newCtx := context.WithValue(ctx, tokenInfo.ID, tokenInfo)
	return newCtx, nil
}

// 解析 token，并进行验证.
func parseToken(token string) (TokenInfo, error) {
	var tokenInfo TokenInfo
	if token == "grpc.auth.token" {
		tokenInfo.ID = "1"
		tokenInfo.Roles = []string{"admin"}
		return tokenInfo, nil
	}
	return tokenInfo, errors.New("Token 无效: bearer " + token)
}

// 从 token 中获取用户唯一标识.
func userClaimsFromToken(tokenInfo TokenInfo) string {
	return tokenInfo.ID
}
