package common

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestValidatePositiveInt(t *testing.T) {
	t.Run("valid positive", func(t *testing.T) {
		id, err := ValidatePositiveInt(42)
		assert.NoError(t, err)
		assert.Equal(t, int64(42), id)
	})
	t.Run("zero", func(t *testing.T) {
		_, err := ValidatePositiveInt(0)
		assert.Error(t, err)
	})
	t.Run("negative", func(t *testing.T) {
		_, err := ValidatePositiveInt(-5)
		assert.Error(t, err)
	})
}

func TestValidatePositiveIntString(t *testing.T) {
	t.Run("valid string", func(t *testing.T) {
		id, err := ValidatePositiveIntString("123")
		assert.NoError(t, err)
		assert.Equal(t, int64(123), id)
	})
	t.Run("empty string", func(t *testing.T) {
		_, err := ValidatePositiveIntString("")
		assert.Error(t, err)
	})
	t.Run("invalid int", func(t *testing.T) {
		_, err := ValidatePositiveIntString("abc")
		assert.Error(t, err)
	})
}

func TestValidateIDList(t *testing.T) {
	t.Run("all positive", func(t *testing.T) {
		err := ValidateIDList([]int64{1, 2, 3})
		assert.NoError(t, err)
	})
	t.Run("contains zero", func(t *testing.T) {
		err := ValidateIDList([]int64{1, 0, 3})
		assert.Error(t, err)
	})
	t.Run("contains negative", func(t *testing.T) {
		err := ValidateIDList([]int64{1, -2, 3})
		assert.Error(t, err)
	})
}

func TestValidateIDListFromString(t *testing.T) {
	t.Run("valid list", func(t *testing.T) {
		ids, err := ValidateIDListFromString("1,2,3")
		assert.NoError(t, err)
		assert.Equal(t, []int64{1, 2, 3}, ids)
	})
	t.Run("invalid int in list", func(t *testing.T) {
		_, err := ValidateIDListFromString("1,a,3")
		assert.Error(t, err)
	})
	t.Run("zero in list", func(t *testing.T) {
		_, err := ValidateIDListFromString("1,0,3")
		assert.Error(t, err)
	})
}

func TestIsValidPassword(t *testing.T) {
	t.Run("valid password", func(t *testing.T) {
		assert.True(t, IsValidPassword("Test@123"))
	})
	t.Run("short password", func(t *testing.T) {
		assert.False(t, IsValidPassword("T@1a"))
	})
	t.Run("no uppercase", func(t *testing.T) {
		assert.False(t, IsValidPassword("test@123"))
	})
	t.Run("no special", func(t *testing.T) {
		assert.False(t, IsValidPassword("Test1234"))
	})
}
