package core

// import (
// 	"context"
// 	"core/models"
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"fmt"
// 	"sync"
// 	"time"

// 	"github.com/vn-go/bx"
// )

// type cacheUser interface {
// 	CacheUser(ctx context.Context, user models.User) error
// 	GetUserById(ctx context.Context, userId string) (*models.User, error)
// 	Sha265(content string) string
// }
// type cacheUserImpl struct {
// 	prefixKey string
// 	cache     bx.Cache
// }
// type initCacheUserImplSha265 struct {
// 	val  string
// 	once sync.Once
// }

// var cacheInitCacheUserImplSha265 sync.Map

// func (c *cacheUserImpl) Sha265(content string) string {
// 	a, _ := cacheInitCacheUserImplSha265.LoadOrStore(content, &initCacheUserImplSha265{})
// 	i := a.(*initCacheUserImplSha265)
// 	i.once.Do(func() {
// 		hash := sha256.Sum256([]byte(content))
// 		i.val = hex.EncodeToString(hash[:])
// 	})
// 	return i.val
// }
// func (c *cacheUserImpl) CacheUser(ctx context.Context, user models.User) error {
// 	return c.cache.Set(ctx, c.Sha265(fmt.Sprintf("%s@%s", user.UserId, c.prefixKey)), user, time.Hour*4)
// }
// func (c *cacheUserImpl) GetUserById(ctx context.Context, userId string) (*models.User, error) {
// 	ret := &models.User{}
// 	err := c.cache.Get(ctx, c.Sha265(fmt.Sprintf("%s@%s", userId, c.prefixKey)), ret)
// 	return ret, err
// }
// func newCacheUser(cfg *configInfo) cacheUser {
// 	return  &cacheUserImpl {
// 		prefixKey: ,
// 	}

// }
// func Test(data any) {
// 	typ:=reflect.TypeOf(data)
// 	if typ.Kind()==reflect.Kind
// }
