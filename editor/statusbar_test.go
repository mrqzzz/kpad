package editor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatusBar(t *testing.T) {

	// test with large width

	txt := fitText(10, "123", "abc")
	assert.Equal(t, "123    abc", txt)

	txt = fitText(10, "123", "")
	assert.Equal(t, "123       ", txt)

	txt = fitText(10, "123", "abcdef")
	assert.Equal(t, "123 abcdef", txt)

	txt = fitText(10, "123", "abcdefg")
	assert.Equal(t, "123abcdefg", txt)

	txt = fitText(10, "123", "abcdefghijkl")
	assert.Equal(t, "123abcdefg", txt)

	txt = fitText(10, "", "abcdefghijkl")
	assert.Equal(t, "abcdefghij", txt)

	// test with smaller width

	txt = fitText(3, "123", "abcdefghijkl")
	assert.Equal(t, "123", txt)

	txt = fitText(2, "123", "abcdefghijkl")
	assert.Equal(t, "12", txt)

	// test with width = 1

	txt = fitText(1, "123", "abcdefghijkl")
	assert.Equal(t, "1", txt)

	txt = fitText(1, "123", "")
	assert.Equal(t, "1", txt)

	txt = fitText(1, "", "abcde")
	assert.Equal(t, "a", txt)

	txt = fitText(1, "", "")
	assert.Equal(t, " ", txt)

	// test with width = 0

	txt = fitText(0, "123", "abc")
	assert.Equal(t, "", txt)

	txt = fitText(0, "", "")
	assert.Equal(t, "", txt)
}
