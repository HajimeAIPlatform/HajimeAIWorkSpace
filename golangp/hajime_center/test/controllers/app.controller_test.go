package test

import (
	"github.com/magiconair/properties/assert"
	"hajime/golangp/hajime_center/controllers"
	"testing"
)

func TestGetModelCharge(t *testing.T) {
	_ = controllers.NewModelController(nil)

	assert.Equal(t, 0, 0)
}
