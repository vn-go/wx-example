package core

import (
	"context"

	"github.com/vn-go/bx"
	"github.com/vn-go/dx"
)

type servicesTypes struct {
	users UserService
}

type Container struct {
	bx.Container[Container]
	userService bx.Depen[UserService, Container]
	userRepo    bx.Depen[userRepo, Container]
	pwdSvc      bx.OnlyOnce[passwordService]
	db          bx.OnlyOnce[*dx.DB]
	ctx         bx.Depen[context.Context, Container]
}

/*
Khoi tao ban dau cho container, moi conatiner bat buoc phai co ham New, dung de setup cac Resolver
The initialization for the container, each container must have a New function, used to set up the Resolvers.
*/
func (c *Container) New() {
	c.userRepo.Resolve(func() (userRepo, error) {
		db, err := c.db.Get()
		if err != nil {
			return nil, err
		}
		ctx, err := c.ctx.Get()
		if err != nil {
			return nil, err
		}
		return &userRepoSql{
			db:      db,
			context: ctx,
		}, nil

	})
	c.userService.Resolve(func() (UserService, error) {
		userRepo, err := c.userRepo.Get()
		if err != nil {
			return nil, err
		}
		pwdSvc, err := c.pwdSvc.Get()
		if err != nil {
			return nil, err
		}
		return &UserServiceSql{
			userRepo: userRepo,
			pwdSvc:   pwdSvc,
		}, nil
	})
	c.pwdSvc.Resolve(func() (passwordService, error) {
		return &bcryptPasswordService{}, nil
	})

}
