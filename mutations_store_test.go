package figtree

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTree_StoreString(t *testing.T) {
	const k, u = "test-store-string", "usage"

	// new fig tree with string
	figs := With(Options{Germinate: true})
	figs.NewString(k, "default", u)
	assert.Nil(t, figs.Parse())

	// get a string k
	s := *figs.String(k)
	assert.Equal(t, "default", s)

	// store a new string in k
	figs.StoreString(k, "new")

	// verify new string in k
	s = *figs.String(k)
	assert.Equal(t, "new", s)
}

func TestTree_StoreBool(t *testing.T) {
	const k, u = "test-store-bool", "usage"

	// new fig tree with bool
	figs := With(Options{Germinate: true})
	figs.NewBool(k, false, u)
	assert.Nil(t, figs.Parse())

	// get a bool s
	s := *figs.Bool(k)
	assert.Equal(t, false, s)

	// store a new bool in k
	figs.StoreBool(k, true)

	// verify new bool in k
	s = *figs.Bool(k)
	assert.Equal(t, true, s)
}

func TestTree_StoreInt(t *testing.T) {
	const k, u = "test-store-int", "usage"

	// new fig tree with int
	figs := With(Options{Germinate: true})
	figs.NewInt(k, 17, u)
	assert.Nil(t, figs.Parse())

	// get an int s
	s := *figs.Int(k)
	assert.Equal(t, 17, s)

	// store a new int in k
	figs.StoreInt(k, 76)

	// verify new int in k
	s = *figs.Int(k)
	assert.Equal(t, 76, s)
}

func TestTree_StoreInt64(t *testing.T) {
	const k, u = "test-store-int64", "usage"

	// new fig tree with int64
	figs := With(Options{Germinate: true})
	figs.NewInt64(k, int64(17), u)
	assert.Nil(t, figs.Parse())

	// get an int64 s
	s := *figs.Int64(k)
	assert.Equal(t, int64(17), s)

	// store a new int64 in k
	figs.StoreInt64(k, 76)

	// verify new int64 in k
	s = *figs.Int64(k)
	assert.Equal(t, int64(76), s)
}

func TestTree_StoreFloat64(t *testing.T) {
	const k, u = "test-store-float64", "usage"

	// new fig tree with float64
	figs := With(Options{Germinate: true})
	figs.NewFloat64(k, float64(17), u)
	assert.Nil(t, figs.Parse())

	// get a float64 s
	s := *figs.Float64(k)
	assert.Equal(t, float64(17), s)

	// store a new float64 in k
	figs.StoreFloat64(k, 76)

	// verify new float64 in k
	s = *figs.Float64(k)
	assert.Equal(t, float64(76), s)
}

func TestTree_StoreDuration(t *testing.T) {
	const k, u = "test-stoe-duration", "usage"

	// new fig tree with time.Duration
	figs := With(Options{Germinate: true})
	figs.NewDuration(k, time.Duration(17), u)
	assert.Nil(t, figs.Parse())

	// get a time.Duration s from k
	s := *figs.Duration(k)
	assert.Equal(t, time.Duration(17), s)

	// store a new time.Duration in k
	figs.StoreDuration(k, time.Duration(76))

	// verify time.Duration in k
	s = *figs.Duration(k)
	assert.Equal(t, time.Duration(76), s)
}

func TestTree_StoreUnitDuration(t *testing.T) {
	const k, u = "test-store-unit-duration", "usage"

	// new fig tree with unit duration
	figs := With(Options{Germinate: true})
	figs.NewUnitDuration(k, time.Duration(17), time.Second, u)
	assert.Nil(t, figs.Parse())

	// get a time.Duration with units s from k
	s := *figs.UnitDuration(k)
	assert.Equal(t, 17*time.Second, s)

	// store a new time.Duration with units in k
	figs.StoreUnitDuration(k, time.Duration(76), time.Minute)

	// verify new time.Duration in k
	s = *figs.UnitDuration(k)
	assert.Equal(t, 76*time.Minute, s)
}

func TestTree_StoreList(t *testing.T) {
	const k, u = "test-store-list", "usage"

	// new fig tree with a list
	figs := With(Options{Germinate: true})
	figs = figs.NewList(k, []string{"yah", "i am", "yahuah"}, u)
	assert.Nil(t, figs.Parse())

	// get the list from k as s
	s := *figs.List(k)
	assert.Equal(t, 3, len(s))
	assert.Equal(t, []string{"i am", "yah", "yahuah"}, s)

	// store a new list in k
	figs = figs.StoreList(k, []string{"yah", "its", "true", "he", "is"})

	// verify the new list
	s = *figs.List(k)
	assert.Equal(t, 5, len(s))
	assert.Equal(t, []string{"he", "is", "its", "true", "yah"}, s)
}

func TestTree_StoreMap(t *testing.T) {
	const k, u = "test-store-map", "usage"
	d := map[string]string{"name": "yahuah", "sex": "male"}

	// new fig tree with a map as k
	figs := With(Options{Germinate: true})
	figs.NewMap(k, d, u)
	assert.Nil(t, figs.Parse())

	// get map k as s
	s := *figs.Map(k)
	assert.Equal(t, 2, len(s))

	// verify name in map
	n, nok := s["name"]
	assert.True(t, nok)
	assert.Equal(t, "yahuah", n)

	// create a new map
	a := map[string]string{"name": "andrei", "sex": "male"}

	// replace the map of k that is d with a new map of a in k
	figs.StoreMap(k, a)

	// verify the map replacement
	s = *figs.Map(k)

	// check the name again for a change
	n, nok = s["name"]
	assert.True(t, nok)
	assert.Equal(t, "andrei", n)
}
