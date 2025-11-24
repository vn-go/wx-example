/*
this file support function create SHA256 key from key
Some Cache Key too log we need zip in Sha256
*/
package cacher

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
)

// Memoization Sha256 calculation struct
type initCacheServiceImplSha265 struct {
	val  string
	once sync.Once
}

var cacheInitCacheServiceImplSha265 sync.Map

/*
this function will calcuate SHA256 content and cache for next access
*/
func (c *CacheServiceImpl) sha265(content string) string {
	a, _ := cacheInitCacheServiceImplSha265.LoadOrStore(content, &initCacheServiceImplSha265{})
	i := a.(*initCacheServiceImplSha265)
	i.once.Do(func() {
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s://%s", c.prefixKey, content)))
		i.val = hex.EncodeToString(hash[:])
	})
	return i.val
}
