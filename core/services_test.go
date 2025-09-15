package core

import (
	"context"
	"core/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vn-go/dx"
)

var mysqlDsn = "root:123456@tcp(127.0.0.1:3306)/hrm"

func TestUserservice(t *testing.T) {

	c, err := (&Container{}).NewAndSetValue(func(c *Container) error {
		c.db.Resolve(func() (*dx.DB, error) {

			ret, err := dx.Open("mysql", mysqlDsn)
			if err != nil {
				return nil, err
			}

			return ret, nil

		})
		c.ctx.Resolve(func() (context.Context, error) {
			return t.Context(), nil
		})
		return nil
	})
	assert.NoError(t, err)
	userservice, err := c.userService.Get()
	assert.NoError(t, err)
	user, _ := dx.NewThenSetDefaultValues(func() (*models.User, error) {
		return &models.User{
			Username: "x-000002",
		}, nil
	})
	err = userservice.CreateUser(user)
	assert.NoError(t, err)
}
func BenchmarkUserservice(t *testing.B) {
	for i := 0; i < t.N; i++ {

		c, err := (&Container{}).NewAndSetValue(func(c *Container) error {
			c.db.Resolve(func() (*dx.DB, error) {

				ret, err := dx.Open("mysql", mysqlDsn)
				if err != nil {
					return nil, err
				}

				return ret, nil

			})
			c.ctx.Resolve(func() (context.Context, error) {
				return t.Context(), nil
			})
			return nil
		})
		assert.NoError(t, err)
		userservice, err := c.userService.Get()
		assert.NoError(t, err)
		user, _ := dx.NewThenSetDefaultValues(func() (*models.User, error) {
			return &models.User{
				Username:     fmt.Sprintf("bench_user_di_%d", i),
				HashPassword: "password123",
			}, nil
		})
		userservice.CreateUser(user)

	}
}
