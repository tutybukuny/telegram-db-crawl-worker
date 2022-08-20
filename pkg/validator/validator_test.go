package utilities

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var vld IValidator

func TestMain(m *testing.M) {
	vld = New()

	code := m.Run()
	os.Exit(code)
}

func TestValidatorImpl_Bool(t *testing.T) {
	str := ""
	_, err := vld.Bool(str, "required")
	assert.Error(t, err)

	str = "tru"
	v, err := vld.Bool(str, "required")
	assert.NoError(t, err)
	assert.Equal(t, false, v)

	str = "true"
	v, err = vld.Bool(str, "required")
	assert.NoError(t, err)
	assert.Equal(t, true, v)
}

func TestValidatorImpl_Float(t *testing.T) {
	str := ""
	_, err := vld.Float(str, "required")
	assert.Error(t, err)

	str = "1"
	v, err := vld.Float(str, "required")
	assert.NoError(t, err)
	assert.Equal(t, 1.0, v)

	str = "1.12"
	v, err = vld.Float(str, "required")
	assert.NoError(t, err)
	assert.Equal(t, 1.12, v)
}

func TestValidatorImpl_Integer(t *testing.T) {
	str := ""
	_, err := vld.Integer(str, "required")
	assert.Error(t, err)

	str = "1"
	v, err := vld.Integer(str, "required")
	assert.NoError(t, err)
	assert.Equal(t, 1, v)
}

func TestValidatorImpl_String(t *testing.T) {
	str := ""
	_, err := vld.String(str, "required")
	assert.Error(t, err)

	str = "hello"
	v, err := vld.String(str, "required")
	assert.NoError(t, err)
	assert.Equal(t, "hello", v)
}
