package figtree

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFlesh(t *testing.T) {
	assert.NotNil(t, NewFlesh(t.Name()))
	assert.Equal(t, t.Name(), NewFlesh(t.Name()).ToString())
	assert.Equal(t, 0, NewFlesh(t.Name()).ToInt())
	assert.Equal(t, map[string]string{}, NewFlesh(t.Name()).ToMap())
	assert.Equal(t, []string{t.Name()}, NewFlesh(t.Name()).ToList())
	var x interface{}
	assert.NotNil(t, NewFlesh(x))
	assert.NotNil(t, NewFlesh(x).ToString())
}

func TestFleshInterface(t *testing.T) {
	t.Run("Is", func(t *testing.T) {
		t.Run("Map", func(t *testing.T) {
			os.Args = []string{os.Args[0]}
			figs := With(Options{Germinate: true, IgnoreEnvironment: true})
			figs.NewMap(t.Name(), map[string]string{"name": "yahuah"}, t.Name())
			assert.NoError(t, figs.Parse())
			var flesh Flesh
			flesh = figs.FigFlesh(t.Name())
			assert.NotNil(t, flesh)
			assert.True(t, flesh.Is(tMap))
			assert.False(t, flesh.Is(tBool))
		})
		t.Run("List", func(t *testing.T) {
			os.Args = []string{os.Args[0]}
			figs := With(Options{Germinate: true, IgnoreEnvironment: true})
			figs.NewList(t.Name(), []string{"yahuah"}, t.Name())
			assert.NoError(t, figs.Parse())
			var flesh Flesh
			flesh = figs.FigFlesh(t.Name())
			assert.NotNil(t, flesh)
			assert.True(t, flesh.Is(tList))
			assert.False(t, flesh.Is(tBool))
		})
		t.Run("UnitDuration", func(t *testing.T) {
			os.Args = []string{os.Args[0]}
			figs := With(Options{Germinate: true, IgnoreEnvironment: true})
			figs.NewUnitDuration(t.Name(), 17, time.Second, t.Name())
			assert.NoError(t, figs.Parse())
			var flesh Flesh
			flesh = figs.FigFlesh(t.Name())
			assert.NotNil(t, flesh)
			assert.True(t, flesh.Is(tUnitDuration))
			assert.False(t, flesh.Is(tBool))
		})
		t.Run("Duration", func(t *testing.T) {
			os.Args = []string{os.Args[0]}
			figs := With(Options{Germinate: true, IgnoreEnvironment: true})
			figs.NewDuration(t.Name(), time.Second, t.Name())
			assert.NoError(t, figs.Parse())
			var flesh Flesh
			flesh = figs.FigFlesh(t.Name())
			assert.NotNil(t, flesh)
			assert.True(t, flesh.Is(tDuration))
			assert.False(t, flesh.Is(tBool))
		})
		t.Run("Float64", func(t *testing.T) {
			os.Args = []string{os.Args[0]}
			figs := With(Options{Germinate: true, IgnoreEnvironment: true})
			figs.NewFloat64(t.Name(), 1.0, t.Name())
			assert.NoError(t, figs.Parse())
			var flesh Flesh
			flesh = figs.FigFlesh(t.Name())
			assert.NotNil(t, flesh)
			assert.True(t, flesh.Is(tFloat64))
			assert.False(t, flesh.Is(tInt))
		})
		t.Run("Int64", func(t *testing.T) {
			os.Args = []string{os.Args[0]}
			figs := With(Options{Germinate: true, IgnoreEnvironment: true})
			figs.NewInt64(t.Name(), 1, t.Name())
			assert.NoError(t, figs.Parse())
			var flesh Flesh
			flesh = figs.FigFlesh(t.Name())
			assert.NotNil(t, flesh)
			assert.True(t, flesh.Is(tInt64))
			assert.False(t, flesh.Is(tFloat64))
		})
		t.Run("Int", func(t *testing.T) {
			os.Args = []string{os.Args[0]}
			figs := With(Options{Germinate: true, IgnoreEnvironment: true})
			figs.NewInt(t.Name(), 1, t.Name())
			assert.NoError(t, figs.Parse())
			var flesh Flesh
			flesh = figs.FigFlesh(t.Name())
			assert.NotNil(t, flesh)
			assert.True(t, flesh.Is(tInt))
			assert.False(t, flesh.Is(tFloat64))
		})
		t.Run("String", func(t *testing.T) {
			os.Args = []string{os.Args[0]}
			figs := With(Options{Germinate: true, IgnoreEnvironment: true})
			figs.NewString(t.Name(), t.Name(), t.Name())
			assert.NoError(t, figs.Parse())
			var flesh Flesh
			flesh = figs.FigFlesh(t.Name())
			assert.NotNil(t, flesh)
			assert.True(t, flesh.Is(tString))
			assert.False(t, flesh.Is(tFloat64))
		})

	})
	t.Run("ToMap", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		figs := With(Options{Germinate: true, IgnoreEnvironment: true})
		figs.NewMap(t.Name(), map[string]string{"name": "yahuah"}, t.Name())
		assert.NoError(t, figs.Parse())
		var flesh Flesh
		flesh = figs.FigFlesh(t.Name())
		assert.NotNil(t, flesh)
		assert.Contains(t, flesh.ToMap(), "name")
	})
	t.Run("ToList", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		figs := With(Options{Germinate: true, IgnoreEnvironment: true})
		figs = figs.NewList(t.Name(), []string{"yahuah"}, t.Name())
		assert.NoError(t, figs.Parse())
		var flesh Flesh
		flesh = figs.FigFlesh(t.Name())
		assert.NotNil(t, flesh)
		assert.Contains(t, flesh.ToList(), "yahuah")
	})
	t.Run("ToUnitDuration", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		figs := With(Options{Germinate: true, IgnoreEnvironment: true})
		figs.NewUnitDuration(t.Name(), 1, time.Second, t.Name())
		assert.NoError(t, figs.Parse())
		var flesh Flesh
		flesh = figs.FigFlesh(t.Name())
		assert.NotNil(t, flesh)
		assert.Equal(t, flesh.ToUnitDuration(), time.Second)
	})
	t.Run("ToDuration", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		figs := With(Options{Germinate: true, IgnoreEnvironment: true})
		figs.NewDuration(t.Name(), 1, t.Name())
		assert.NoError(t, figs.Parse())
		var flesh Flesh
		flesh = figs.FigFlesh(t.Name())
		assert.NotNil(t, flesh)
		assert.Equal(t, time.Duration(1), flesh.ToDuration())
	})
	t.Run("ToFloat64", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		figs := With(Options{Germinate: true, IgnoreEnvironment: true})
		figs.NewFloat64(t.Name(), 1, t.Name())
		assert.NoError(t, figs.Parse())
		var flesh Flesh
		flesh = figs.FigFlesh(t.Name())
		assert.NotNil(t, flesh)
		assert.Equal(t, 1.0, flesh.ToFloat64())
	})
	t.Run("ToInt64", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		figs := With(Options{Germinate: true, IgnoreEnvironment: true})
		figs.NewInt64(t.Name(), 1, t.Name())
		assert.NoError(t, figs.Parse())
		var flesh Flesh
		flesh = figs.FigFlesh(t.Name())
		assert.NotNil(t, flesh)
		assert.Equal(t, int64(1), flesh.ToInt64())
	})
	t.Run("ToInt", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		figs := With(Options{Germinate: true, IgnoreEnvironment: true})
		figs.NewInt(t.Name(), 1, t.Name())
		assert.NoError(t, figs.Parse())
		var flesh Flesh
		flesh = figs.FigFlesh(t.Name())
		assert.NotNil(t, flesh)
		assert.Equal(t, 1, flesh.ToInt())
	})
	t.Run("ToString", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		figs := With(Options{Germinate: true, IgnoreEnvironment: true})
		figs.NewString(t.Name(), t.Name(), t.Name())
		assert.NoError(t, figs.Parse())
		var flesh Flesh
		flesh = figs.FigFlesh(t.Name())
		assert.NotNil(t, flesh)
		assert.Equal(t, t.Name(), flesh.ToString())
	})
}
