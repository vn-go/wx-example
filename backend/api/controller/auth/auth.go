package auth

import (
	"apicore/controller/base"
	"net/http"
	"time"

	"core/services/errs"
	

	"github.com/vn-go/wx"
)

type Auth struct {
	base.Base
}

func (aut *Auth) Login(h struct {
	wx.Handler `route:"/api/auth/login"`
}, data wx.Form[struct {
	Username string `json:"username"`
	Password string `json:"password"`
}]) (any, error) {
	ret, err := aut.Svc.AccSvc.LoginAndGetJWT(h.Handler().Req.Context(), "", data.Data.Username, data.Data.Password)
	if err != nil {
		return nil, aut.ParseError(err)
	}
	refreshToken := ret.RefreshToken // nếu dùng reference token, lưu vào DB

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
	http.SetCookie(h.Handler().Res, cookie)
	return ret, aut.ParseError(err)
}

func (aut *Auth) ParseError(err error) error {

	cErr, ok := err.(*errs.Err)
	if ok {
		switch cErr.Typ {
		case errs.ErrUnautheticate:
			return wx.Errors.NewUnauthorizedError()
		case errs.ErrForbidden:
			return wx.Errors.NewForbidenError("forbidden")
		case errs.ErrBadRequest:
			return &wx.HttpError{
				Data: "Bad request",
				Code: http.StatusBadRequest,
			}
		case errs.ErrSytem:
			return &wx.HttpError{
				Code: http.StatusInternalServerError,
				Data: "Internal server error",
			}
		default:
			return &wx.HttpError{
				Code: http.StatusInternalServerError,
				Data: "Internal server error",
			}
		}
	}
	return err
}
