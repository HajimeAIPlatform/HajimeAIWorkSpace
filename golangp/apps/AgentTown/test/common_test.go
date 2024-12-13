package test

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestSmokeTest(t *testing.T) {
	assert.Equal(t, 10, 9+1)
}

func TestSmokeTest2(t *testing.T) {
	assert.Equal(t, 12, 9+3)
}
