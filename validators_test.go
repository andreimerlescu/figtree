package figtree

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTree_WithValidator(t *testing.T) {
	t.Run("CantDoDoubleValidationsOnSameKey", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewString(t.Name(), "i love yahuah", "usage")
		fig.WithValidator(t.Name(), AssureStringSubstring("love yahuah"))
		fig.WithValidator(t.Name(), AssureStringNotEmpty)
		assert.NotNil(t, fig.Parse())
	})
	t.Run("AssureStringLengthLessThan", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewString(t.Name(), "i love yahuah", "usage")
		fig.WithValidator(t.Name(), AssureStringLengthLessThan(99))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureStringLengthGreaterThan", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewString(t.Name(), "i love yahuah", "usage")
		fig.WithValidator(t.Name(), AssureStringLengthGreaterThan(3))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureStringHasSuffix", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewString(t.Name(), "i love yahuah", "usage")
		fig.WithValidator(t.Name(), AssureStringHasSuffix("yahuah"))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureStringHasPrefix", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewString(t.Name(), "i love yahuah", "usage")
		fig.WithValidator(t.Name(), AssureStringHasPrefix("i love"))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureStringLength", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewString(t.Name(), "i love yahuah", "usage")
		fig.WithValidator(t.Name(), AssureStringLength(13))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureStringSubstring", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewString(t.Name(), "i love yahuah", "usage")
		fig.WithValidator(t.Name(), AssureStringSubstring("love yahuah"))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureStringNotEmpty", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewString(t.Name(), "i love yahuah", "usage")
		fig.WithValidator(t.Name(), AssureStringNotEmpty)
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureStringContains ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewString(t.Name(), "i love yahuah", "usage")
		fig.WithValidator(t.Name(), AssureStringContains("love yahuah"))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureBoolTrue", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewBool(t.Name(), true, "usage")
		fig.WithValidator(t.Name(), AssureBoolTrue)
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureBoolFalse", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewBool(t.Name(), false, "usage")
		fig.WithValidator(t.Name(), AssureBoolFalse)
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssurePositiveInt ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewInt(t.Name(), 17, "usage")
		fig.WithValidator(t.Name(), AssurePositiveInt)
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureNegativeInt ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewInt(t.Name(), -17, "usage")
		fig.WithValidator(t.Name(), AssureNegativeInt)
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureIntGreaterThan ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewInt(t.Name(), 17, "usage")
		fig.WithValidator(t.Name(), AssureIntGreaterThan(12))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureIntLessThan ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewInt(t.Name(), 17, "usage")
		fig.WithValidator(t.Name(), AssureIntLessThan(33))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureIntInRange ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewInt(t.Name(), 47, "usage")
		fig.WithValidator(t.Name(), AssureIntInRange(17, 76))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureInt64GreaterThan", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewInt(t.Name(), 17, "usage")
		fig.WithValidator(t.Name(), AssureInt64GreaterThan(3))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureInt64LessThan ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewInt(t.Name(), 17, "usage")
		fig.WithValidator(t.Name(), AssureInt64LessThan(33))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureInt64Positive ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewInt64(t.Name(), 17, "usage")
		fig.WithValidator(t.Name(), AssureInt64Positive)
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureFloat64Positive ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewFloat64(t.Name(), 17.76, "usage")
		fig.WithValidator(t.Name(), AssureFloat64Positive)
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureFloat64InRange ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewFloat64(t.Name(), 17.76, "usage")
		fig.WithValidator(t.Name(), AssureFloat64InRange(1.0, 20.0))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureFloat64GreaterThan ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewFloat64(t.Name(), 17.76, "usage")
		fig.WithValidator(t.Name(), AssureFloat64GreaterThan(3.69))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureFloat64LessThan ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewFloat64(t.Name(), 17.76, "usage")
		fig.WithValidator(t.Name(), AssureFloat64LessThan(33.33))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureDurationGreaterThan ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewDuration(t.Name(), 30*time.Second, "usage")
		fig.WithValidator(t.Name(), AssureDurationGreaterThan(10*time.Second))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureDurationLessThan ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewDuration(t.Name(), 30*time.Second, "usage")
		fig.WithValidator(t.Name(), AssureDurationLessThan(50*time.Second))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureDurationPositive ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewDuration(t.Name(), 30*time.Second, "usage")
		fig.WithValidator(t.Name(), AssureDurationPositive)
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureDurationMax ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewDuration(t.Name(), 30*time.Second, "usage")
		fig.WithValidator(t.Name(), AssureDurationMax(1*time.Hour))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureListNotEmpty ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewList(t.Name(), []string{"yah", "i am", "yahuah"}, "usage")
		fig.WithValidator(t.Name(), AssureListNotEmpty)
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureListMinLength ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewList(t.Name(), []string{"yah", "i am", "yahuah"}, "usage")
		fig.WithValidator(t.Name(), AssureListMinLength(3))
		assert.Nil(t, fig.Parse())

	})
	t.Run("AssureMapNotEmpty ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewMap(t.Name(), map[string]string{"one": "yah", "two": "i am", "three": "yahuah"}, "usage")
		fig.WithValidator(t.Name(), AssureMapNotEmpty)
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureMapHasKey ", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewMap(t.Name(), map[string]string{"one": "yah", "two": "i am", "three": "yahuah"}, "usage")
		fig.WithValidator(t.Name(), AssureMapHasKey("three"))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureInt64InRange", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewInt64(t.Name(), 47, "usage")
		fig.WithValidator(t.Name(), AssureInt64InRange(17, 76))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureFloat64NotNaN", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewFloat64(t.Name(), 47.0, "usage")
		fig.WithValidator(t.Name(), AssureFloat64NotNaN)
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureDurationMin", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewDuration(t.Name(), 30*time.Second, "usage")
		fig.WithValidator(t.Name(), AssureDurationMin(3*time.Second))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureListContains", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewList(t.Name(), []string{"yah", "i am", "yahuah"}, "usage")
		fig.WithValidator(t.Name(), AssureListContains("yahuah"))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureMapValueMatches", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewMap(t.Name(), map[string]string{"one": "yah", "two": "i am", "three": "yahuah"}, "usage")
		fig.WithValidator(t.Name(), AssureMapValueMatches("three", "yahuah"))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureMapHasKeys", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewMap(t.Name(), map[string]string{"one": "yah", "two": "i am", "three": "yahuah"}, "usage")
		fig.WithValidator(t.Name(), AssureMapHasKeys([]string{"one", "two", "three"}))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureListContainsKey", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewList(t.Name(), []string{"yah", "i am", "yahuah"}, "usage")
		fig.WithValidator(t.Name(), AssureListContainsKey("yahuah"))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureListLength", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewList(t.Name(), []string{"yah", "i am", "yahuah"}, "usage")
		fig.WithValidator(t.Name(), AssureListLength(3))
		assert.Nil(t, fig.Parse())
	})
	t.Run("AssureMapLength", func(t *testing.T) {
		fig := With(Options{Germinate: true, Tracking: false})
		fig.NewMap(t.Name(), map[string]string{"one": "yah", "two": "i am", "three": "yahuah"}, "usage")
		fig.WithValidator(t.Name(), AssureMapLength(3))
		assert.Nil(t, fig.Parse())
	})
}
