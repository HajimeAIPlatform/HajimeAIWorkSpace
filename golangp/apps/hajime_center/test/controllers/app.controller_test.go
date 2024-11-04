package test

import (
	"github.com/magiconair/properties/assert"
	"hajime/golangp/apps/hajime_center/controllers"
	"testing"
)

func TestSmokeTest(t *testing.T) {
	_ = controllers.NewModelController(nil)

	assert.Equal(t, 0, 0)
}
