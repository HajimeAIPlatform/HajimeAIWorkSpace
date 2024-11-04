package test

import (
	"github.com/magiconair/properties/assert"
	"hajime/golangp/common/utils"
	"testing"
)

func TestSmokeTest(t *testing.T) {
	genStr := utils.GenerateRandomString(10)
	assert.Equal(t, 10, len(genStr))
}
