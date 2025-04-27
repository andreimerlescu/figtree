package figtree

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSomeAssurances(t *testing.T) {
	const k1, k2, k3, k4 = "k1", "k2", "k3", "k4"
	figs := With(Options{Tracking: false, Germinate: true})
	figs.NewString(k1, "default", "usage")
	figs.WithValidator(k1, AssureStringSubstring("def"))

	figs.NewInt(k2, 17, "usage")
	figs.WithValidator(k2, AssureIntPositive)
	figs.WithValidator(k2, AssureIntGreaterThan(1))
	assert.NoError(t, figs.ErrorFor(k2))

	figs.NewList(k3, []string{"yah", "i am", "yahuah"}, "usage")
	figs.WithValidator(k3, AssureListContains("yahuah"))

	figs.NewUnitDuration(k4, 33, time.Second, "usage")
	figs.WithValidator(k4, AssureDurationMin(30*time.Second))
	assert.NoError(t, figs.Parse())
}
