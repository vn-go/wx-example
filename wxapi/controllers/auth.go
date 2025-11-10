package controllers

import (
	"context"
	"core"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/vn-go/wx"
)

type Auth struct {
}

// Business logic core
func (auth *Auth) LoginCore(ctx context.Context, username, password string) (*core.OAuthResponse, error) {
	if !core.Services.Config.Tenant.IsMulti {
		oauth, err := core.Services.AuthSvc.Login("", ctx, username, password)
		if err != nil || oauth == nil {
			return nil, wx.Errors.NewUnauthorizedError()
		}
		return oauth, nil
	}
	if strings.Contains(username, "/") {
		parts := strings.SplitN(username, "/", 2)
		tenant, user := parts[0], parts[1]
		oauth, err := core.Services.AuthSvc.Login(tenant, ctx, user, password)
		if err != nil || oauth == nil {
			return nil, wx.Errors.NewUnauthorizedError()
		}
		return oauth, nil
	}

	oauth, err := core.Services.TenantSvc.Login(ctx, username, password)
	if err != nil || oauth == nil {
		return nil, wx.Errors.NewUnauthorizedError()
	}
	return oauth, nil
}
func (auth *Auth) Login(handler wx.Handler, data wx.Form[struct {
	Username string `check:"range:[3:50]"`
	Password string `check:"range:[3:50]"`
}]) (*core.OAuthResponse, error) {

	ctx := handler().Req.Context()
	ret, err := auth.LoginCore(ctx, data.Data.Username, data.Data.Password)
	if err != nil {
		return nil, err
	}

	// -------------------------
	// 1. Tạo refresh token an toàn
	// Có thể mã hóa hoặc dùng reference token
	refreshToken := ret.RefreshToken // nếu dùng reference token, lưu vào DB
	expires := time.Second * time.Duration(ret.ExpiresIn)

	// -------------------------
	// 2. Set HttpOnly, Secure cookie
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,                  // ngăn JS đọc
		Secure:   true,                  // HTTPS
		SameSite: http.SameSiteNoneMode, // chặt chẽ hơn Lax
		Path:     "/",                   // toàn app
	}
	http.SetCookie(handler().Res, cookie)

	// -------------------------
	// 3. Tạo CSRF token cho SPA
	csrfToken := auth.GenerateCSRFToken() // random, dài, lưu server-side
	csrfCookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(expires),
		HttpOnly: false, // JS cần đọc để gửi header
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
	}
	http.SetCookie(handler().Res, csrfCookie)

	// Trả response chứa access token (thường JWT ngắn hạn)
	return ret, nil
}

// GenerateCSRFToken tạo CSRF token ngẫu nhiên, base64 URL-safe
func (auth *Auth) GenerateCSRFToken() string {
	// độ dài 32 byte → ~43 ký tự base64, đủ an toàn
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		// fallback, ít khi xảy ra
		panic("cannot generate CSRF token: " + err.Error())
	}
	// base64 URL-safe, không có + / =
	return base64.RawURLEncoding.EncodeToString(b)
}
