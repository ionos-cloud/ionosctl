package utils

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertSizeMB(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	size, err := ConvertSize("256MB", MegaBytes)
	assert.NoError(t, err)
	assert.True(t, size == 256)
	size, err = ConvertSize("4GB", MegaBytes)
	assert.NoError(t, err)
	assert.True(t, size == 4*1024)
	size, err = ConvertSize("1TB", MegaBytes)
	assert.NoError(t, err)
	assert.True(t, size == 1*1024*1024)
	size, err = ConvertSize("2PB", MegaBytes)
	assert.NoError(t, err)
	assert.True(t, size == 2*1024*1024*1024)
	size, err = ConvertSize("2PB", MegaBytes)
	assert.NoError(t, err)
	assert.True(t, size == 2*1024*1024*1024)
	err = w.Flush()
	assert.NoError(t, err)
}

func TestConvertSizeGB(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	size, err := ConvertSize("4GB", GigaBytes)
	assert.NoError(t, err)
	assert.True(t, size == 4)
	size, err = ConvertSize("1TB", GigaBytes)
	assert.NoError(t, err)
	assert.True(t, size == 1*1024)
	size, err = ConvertSize("2PB", GigaBytes)
	assert.NoError(t, err)
	assert.True(t, size == 2*1024*1024)
	size, err = ConvertSize("2 PB", GigaBytes)
	assert.NoError(t, err)
	assert.True(t, size == 2*1024*1024)
	err = w.Flush()
	assert.NoError(t, err)
}
