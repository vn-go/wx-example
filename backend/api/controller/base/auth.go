package base

import (
	"context"
	"core"
	"fmt"
	"net/http"

	"sync"

	"core/services/jwt"

	errs "core/services/errs"

	"github.com/vn-go/dx"
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
	if authorization == "" {
		return nil, wx.Errors.NewUnauthorizedError()
	}
	a, _ := initVerifyTokenMap.LoadOrStore(authorization, &initVerifyToken{})
	i := a.(*initVerifyToken)
	i.once.Do(func() {
		ret, err := auth.Svc.AccSvc.ValidateToken(ctx, authorization)
		if err != nil {
			if err != nil {
				if auth.Svc.ErrSvc.IsForbidden(err) {
					i.err = wx.Errors.NewForbidenError("forbidden")
				} else if auth.Svc.ErrSvc.IsUnauthenticate(err) {
					i.err = wx.Errors.NewUnauthorizedError()
				} else {
					i.err = wx.Errors.NewServerError("err", err)
				}

			}
		}
		i.val = ret
	})
	if i.err != nil {
		initVerifyTokenMap.Delete(authorization)
		return nil, i.err
	}
	return i.val, i.err
}
func (auth AuthBase) GetUserDb() (*dx.DB, error) {
	return auth.Svc.TenantSvc.GetDb(auth.Data.Tenant)
}

/*
This function is used to initialize the AuthBase object.
It is called by the framework when the object is created.
Setup authentication and authorization rules.
All controller in GO need embed AuthBase to use authentication and authorization.
*/
func (auth *AuthBase) New() error {
	auth.Svc = core.Services
	// add Verify token
	auth.Authenticate.Verify(func(ctx wx.Handler) (*jwt.Indentifier, error) {
		wxCtx := ctx()
		req := wxCtx.Req
		fmt.Println(wxCtx.Req.URL.Path)
		authorization := req.Header.Get("authorization")
		// verify token and get user
		user, err := auth.Svc.AccSvc.ValidateToken(ctx().Req.Context(), authorization)
		if err != nil {
			return nil, auth.ParseError(err)
		}
		// get view path form header
		viewPath := req.Header.Get("view-path")
		if viewPath == "" && !user.IsSysAdmin {
			return nil, wx.Errors.NewForbidenError(struct {
				Msg string `json:"msg"`
			}{
				Msg: "forbidden",
			})
		}
		// return user for  controller where embed AuthBase
		return user, nil

	})
	return nil
}

func (aut *AuthBase) ParseError(err error) error {

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

	dxErr := dx.Errors.IsDbError(err)
	if dxErr != nil {
		if dxErr.IsEntryNotFoundError() {
			return &wx.HttpError{
				Code: http.StatusNotFound,
				Data: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
					//Fields  []string `json:"fields"`
				}{
					Message: "Not found",
					Code:    "not_found",
				},
			}
		}
		if dxErr.IsDuplicateEntryError() {
			return &wx.HttpError{
				Code: http.StatusConflict,
				Data: struct {
					Code    string   `json:"code"`
					Message string   `json:"message"`
					Fields  []string `json:"fields"`
				}{
					Message: "Duplicate entry",
					Fields:  dxErr.Fields,
					Code:    "duplicate_entry",
				},
			}
		}
	}

	if ok {
		if dxErr.ErrorType == dx.Errors.DUPLICATE {
			return &wx.HttpError{
				Code: http.StatusConflict,
				Data: struct {
					Fields []string `json:"fields"`
				}{
					Fields: dxErr.Fields,
				},
			}
		}
	}
	return err
}
