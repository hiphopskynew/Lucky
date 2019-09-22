package cache

import (
	"time"

	"github.com/bluele/gcache"
)

var c gcache.Cache = gcache.New(99999).Build()
var result = make(map[string]interface{})

func Set(key string, value map[string]interface{}) error {
	return c.Set(key, value)
}

func SetWithExpire(key string, value map[string]interface{}, expire time.Duration) error {
	return c.SetWithExpire(key, value, expire)
}

func Get(key string) (map[string]interface{}, error) {
	v1, e := c.Get(key)
	return v1.(map[string]interface{}), e
}

func GetWithoutErr(key string) (result map[string]interface{}) {
	defer func() {
		recover()
	}()
	v1, _ := c.Get(key)
	result = v1.(map[string]interface{})
	return result
}
