package media

import "github.com/vn-go/wx"

type Media struct {
}

func (m *Media) ViewFile(ctx struct {
	FileID     string `json:"file_id"`
	wx.Handler `route:"@/{FileID};method:post"`
}) (any, error) {
	return ctx.FileID, nil
}
