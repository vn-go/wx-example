package base

import (
	"context"
	"core"

	"sync"

	"core/services/jwt"

	"github.com/vn-go/wx"
)

type AuthBase struct {
	Base
	wx.Authenticate[jwt.Indentifier]
}
type initVerifyToken struct {
	val  *jwt.Indentifier
	err  error
	once sync.Once
}

var initVerifyTokenMap sync.Map

func (auth *AuthBase) verifyToken(ctx context.Context, authorization string) (*jwt.Indentifier, error) {
	a, _ := initVerifyTokenMap.LoadOrStore(authorization, &initVerifyToken{})
	i := a.(*initVerifyToken)
	i.once.Do(func() {
		ret, err := auth.Svc.AccSvc.ValidateToken(ctx, authorization)
		if err != nil {
			if err != nil {
				if auth.Svc.ErrSvc.IsForbidden(err) {
					i.err = wx.Errors.NewForbidenError("forbidden")
				}
				i.err = wx.Errors.NewServerError("err", err)
			}
		}
		i.val = ret
	})
	return i.val, i.err
}
func (auth *AuthBase) New() error {
	auth.Svc = core.Services
	auth.Authenticate.Verify(func(ctx wx.Handler) (*jwt.Indentifier, error) {
		req := ctx().Req
		authorization := req.Header.Get("authorization")
		return auth.verifyToken(ctx().Req.Context(), authorization)

	})
	return nil
}
