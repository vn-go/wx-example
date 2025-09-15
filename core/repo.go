package core

import (
	"context"

	"github.com/vn-go/bx"
	"github.com/vn-go/dx"
)

type modelType struct {
}

var Models = &modelType{}

type repoType struct {
}

func (repo *repoType) User(db *dx.DB, ctx context.Context) userRepo {
	ret, _ := bx.OnceCall[repoType]("User@"+db.DbName+"/"+db.DriverName, func() (userRepo, error) {
		return newUserRepoSql(db, ctx), nil
	})
	return ret

}

var Repo = &repoType{}
