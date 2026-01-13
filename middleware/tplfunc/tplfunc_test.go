package tplfunc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
)

func TestNumberFormat(t *testing.T) {
	assert.Equal(t, "12.99", NumberFormat(12.987, 2))
	assert.Equal(t, "123.99", NumberFormat(123.987, 2))
	assert.Equal(t, "1,234.99", NumberFormat(1234.987, 2))
	assert.Equal(t, "1,234,567.99", NumberFormat(1234567.987, 2))
	assert.Equal(t, "1,234,568", NumberFormat(1234567.987, 0))
	assert.Equal(t, "234,567.99", NumberFormat(234567.987, 2))
	assert.Equal(t, "234,567.98", NumberTrim(234567.987, 2))
	assert.Equal(t, "12,234,567", NumberTrim(12234567.987, 0))
	assert.Equal(t, "234567.9", NumberTrim(234567.987, 1, ``))
}

func TestInExt(t *testing.T) {
	assert.True(t, InExt(`a/b/c/d.jpg`, `.jpg`))
	assert.False(t, InExt(`a/b/c/d.jpeg`, `.jpg`))
	assert.True(t, InExt(`a/b/c/d.jpeg`, `.jpeg`))
	var r []string
	err := com.JSONDecodeString(`null`, &r)
	assert.NoError(t, err)
}
