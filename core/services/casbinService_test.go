package services

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vn-go/dx"
)

var pgDsn = "postgres://postgres:123456@localhost:5432/hrm?sslmode=disable"

func TestCasbinService(t *testing.T) {
	db, err := dx.Open("postgres", pgDsn)
	if err != nil {
		t.Fatal(err)
	}
	a, err := NewCasbinEnforcer(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, a)
	// Thêm policy trực tiếp
	ok, err := a.AddPolicy("alice", "data1", "read")
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Println("policy already exists")
	}
	policies, err := a.GetFilteredPolicy(0, "alice")
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range policies {
		fmt.Println(p) // ["alice", "data1", "read"]
	}

	// Lấy policy với user = "alice" và resource = "data1"
	policies2, _ := a.GetFilteredPolicy(0, "alice", "data1")
	for _, p := range policies2 {
		fmt.Println(p) // ["alice", "data1", "read"]
	}

}
