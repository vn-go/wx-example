package users

import "github.com/vn-go/wx"

func (u *Users) Avatar(h struct {
	wx.Handler `route:"@/{id}" method:"get"`
	ID         string `param:"id"`
}) (any, error) {
	return nil, nil
}
