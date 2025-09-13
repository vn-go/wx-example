package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	
)

var dns string = "root:123456@tcp(127.0.0.1:3306)/a001"

func TestDb(t *testing.T) {
	db, err := xdb.Open("mysql", dns)
	assert.NoError(t, err)
	assert.NotNil(t, db)
}
