package coretest

import (
	"core"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestLoadAllUsers(t *testing.T) {
	configFile := "./../../../backend/api/config.yaml"
	core.Start(configFile)
	defer core.Services.Db.Close()
	var data any
	var err error
	data, err = core.Services.Db.DslToArray("sysUsers()")
	if err != nil {
		t.Error(err)
	}
	dff, err := json.Marshal(data)
	fmt.Println(string(dff))
	for i := 0; i < 10; i++ {
		start := time.Now()
		data, err = core.Services.Db.DslToArray("sysUsers()")
		if err != nil {
			t.Error(err)
		}
		fmt.Println("load all users", time.Now().Sub(start).Microseconds())

	}
}
func BenchmarkLoadAllUsers(b *testing.B) {
	configFile := "./../../../backend/api/config.yaml"
	core.Start(configFile)
	defer core.Services.Db.Close()
	// var data any
	var err error
	_, err = core.Services.Db.DslToArray("sysUsers()")
	if err != nil {
		panic(err)
	}
	b.Run("load all users", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err = core.Services.Db.DslToArray("sysUsers()")
			if err != nil {
				panic(err)
			}

		}
	})

}
