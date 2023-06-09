package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateCustomer(t *testing.T) {
	// prepare
	time.Local, _ = time.LoadLocation("UTC")
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		SkipDefaultTransaction: true, //agiliza os testes
		NowFunc:                time.Now().Local,
	})
	if err != nil {
		require.NoError(t, err)
	}

	err = db.AutoMigrate(&Customer{})
	assert.NoError(t, err)

	st := NewStore(db)

	// test
	c, err := st.CreateCustomer("John Doe")

	// verify
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "John Doe", c.Name)
	assert.EqualValues(t, 1, c.ID)
}
