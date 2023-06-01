package model

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateCustomer(t *testing.T) {
	// prepare
	db, err := gorm.Open(sqlite.Open("TestCreateCustomer.db"), &gorm.Config{})
	if err != nil {
		require.NoError(t, err)
	}
	defer func() {
		err := os.Remove("TestCreateCustomer.db")
		require.NoError(t, err)
	}()

	db.AutoMigrate(&Customer{})

	st := NewStore(db)

	// test
	c, err := st.CreateCustomer("John Doe")

	// verify
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "John Doe", c.Name)
	assert.EqualValues(t, 1, c.ID)
}
