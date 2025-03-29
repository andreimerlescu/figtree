package figtree

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkMutations(b *testing.B) {
	const k, v, u = "name", "yahuah", "usage"
	figs := With(Options{Tracking: false, Germinate: true})
	b.Run("new string", func(b *testing.B) {
		figs.NewString(k+strconv.Itoa(b.N), v, u)
	})
	figs.NewString(k, v, u)
	b.Run("parse", func(b *testing.B) {
		if e := figs.Parse(); e != nil {
			b.Fatal(e)
		}
	})
	b.Run("reads", func(b *testing.B) {
		l := *figs.String(k)
		if l != v {
			b.Fatalf("expected %s, got %s", v, l)
		}
	})
	b.Run("writes", func(b *testing.B) {
		figs.StoreString(k, v)
	})
}

func TestTree_String(t *testing.T) {
	figs := With(Options{Germinate: true})
	figs.NewString("test", "default", "usage")
	assert.Nil(t, figs.Parse())
	assert.Equal(t, *figs.String("test"), "default")
}

func TestTree_Bool(t *testing.T) {
	figs := With(Options{Germinate: true})
	figs.NewBool("test1", true, "usage")
	figs.NewBool("test2", false, "usage")
	assert.Nil(t, figs.Parse())
	assert.Equal(t, *figs.Bool("test1"), true)
	assert.Equal(t, *figs.Bool("test2"), false)
}

func TestTree_Int(t *testing.T) {
	figs := With(Options{Germinate: true})
	figs.NewInt("test-int", 42, "usage")
	assert.Nil(t, figs.Parse())
	assert.Equal(t, *figs.Int("test-int"), 42)
}

func TestTree_Int64(t *testing.T) {
	figs := With(Options{Germinate: true})
	figs.NewInt64("test-int64", 42, "usage")
	assert.Nil(t, figs.Parse())
	assert.Equal(t, *figs.Int64("test-int64"), int64(42))
}

func TestTree_Float64(t *testing.T) {
	figs := With(Options{Germinate: true})
	figs.NewFloat64("test-float64", 1.23, "usage")
	assert.Nil(t, figs.Parse())
	assert.Equal(t, *figs.Float64("test-float64"), 1.23)
}

func TestTree_Duration(t *testing.T) {
	figs := With(Options{Germinate: true})
	figs.NewDuration("test-duration", 42*time.Millisecond, "usage")
	assert.Nil(t, figs.Parse())
	assert.Equal(t, *figs.Duration("test-duration"), 42*time.Millisecond)
}

func TestTree_UnitDuration(t *testing.T) {
	figs := With(Options{Germinate: true})
	figs.NewUnitDuration("test-unit-duration", 42, time.Millisecond, "usage")
	assert.Nil(t, figs.Parse())
	assert.Equal(t, *figs.UnitDuration("test-unit-duration"), 42.0*time.Millisecond)
}

func TestTree_List(t *testing.T) {
	figs := With(Options{Germinate: true})
	figs.NewList("test-list", []string{"one", "two", "three"}, "usage")
	assert.Nil(t, figs.Parse())
	l := *figs.List("test-list")
	assert.Equal(t, 3, len(l))
	assert.Equal(t, l[0], "one")
	assert.Equal(t, l[1], "two")
	assert.Equal(t, l[2], "three")
}

func TestTree_Map(t *testing.T) {
	figs := With(Options{Germinate: true})
	figs.NewMap("test-map", map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"}, "usage")
	assert.Nil(t, figs.Parse())
	m := *figs.Map("test-map")
	k1, k1ok := m["key1"]
	assert.True(t, k1ok)
	assert.Equal(t, k1, "value1")
	k2, k2ok := m["key2"]
	assert.True(t, k2ok)
	assert.Equal(t, k2, "value2")
	k3, k3ok := m["key3"]
	assert.True(t, k3ok)
	assert.Equal(t, k3, "value3")

}
