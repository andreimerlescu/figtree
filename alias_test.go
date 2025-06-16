package figtree

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithAlias(t *testing.T) {
	const cmdLong, cmdAliasLong, valueLong, usage = "long", "l", "default", "usage"
	const cmdShort, cmdAliasShort, valueShort = "short", "s", "default"

	t.Run("basic_usage", func(t *testing.T) {
		figs := With(Options{Germinate: true, Tracking: false})
		figs.NewString(cmdLong, valueLong, usage)
		figs.WithAlias(cmdLong, cmdAliasLong)
		assert.NoError(t, figs.Parse())

		assert.Equal(t, valueLong, *figs.String(cmdLong))
		assert.Equal(t, valueLong, *figs.String(cmdAliasLong))
		figs = nil
	})

	t.Run("shorthand_notation", func(t *testing.T) {
		os.Args = []string{os.Args[0], "-" + cmdAliasLong, valueLong}
		figs := With(Options{Germinate: true, Tracking: false})
		figs.NewString(cmdLong, valueLong, usage)
		figs.WithAlias(cmdLong, cmdAliasLong)
		assert.NoError(t, figs.Parse())
		assert.Equal(t, valueLong, *figs.String(cmdLong))
	})

	t.Run("multiple_aliases", func(t *testing.T) {
		const k, v, u = "name", "yeshua", "the real name of god"
		ka1 := "father"
		ka2 := "son"
		ka3 := "rauch-hokadesch"
		figs := With(Options{Germinate: true, Tracking: false})
		figs.NewString(k, v, u)
		figs.WithAlias(k, ka1)
		figs.WithAlias(k, ka2)
		figs.WithAlias(k, ka3)
		assert.NoError(t, figs.Parse())

		assert.Equal(t, v, *figs.String(k))
		assert.Equal(t, v, *figs.String(ka1))
		assert.Equal(t, v, *figs.String(ka2))
		assert.Equal(t, v, *figs.String(ka3))
		figs = nil
	})

	t.Run("complex_usage", func(t *testing.T) {
		os.Args = []string{
			os.Args[0],
			"-list", "three,four,five",
			"-map", "four=4,five=5,six=6",
		}
		figs := With(Options{Germinate: true, Tracking: false})
		// long
		figs.NewString(cmdLong, valueLong, usage)
		figs.WithAlias(cmdLong, cmdAliasLong)
		figs.WithValidator(cmdLong, AssureStringNotEmpty)

		// short
		figs.NewString(cmdShort, valueShort, usage)
		figs.WithAlias(cmdShort, cmdAliasShort)
		figs.WithValidator(cmdShort, AssureStringNotEmpty)

		// list
		figs.NewList("myList", []string{"one", "two", "three"}, "usage")
		figs.WithValidator("myList", AssureListNotEmpty)
		figs.WithAlias("myList", "list")

		// map
		figs.NewMap("myMap", map[string]string{"one": "1", "two": "2", "three": "3"}, "usage")
		figs.WithValidator("myMap", AssureMapNotEmpty)
		figs.WithAlias("myMap", "map")

		assert.NoError(t, figs.Parse())

		// long
		assert.Equal(t, valueLong, *figs.String(cmdLong))
		assert.Equal(t, valueLong, *figs.String(cmdAliasLong))
		// short
		assert.Equal(t, valueShort, *figs.String(cmdShort))
		assert.Equal(t, valueShort, *figs.String(cmdAliasShort))
		// list
		assert.NotEqual(t, []string{"one", "two", "three"}, *figs.List("myList"))
		assert.Equal(t, []string{"three", "four", "five"}, *figs.List("myList"))
		// list alias
		assert.NotEqual(t, []string{"one", "two", "three"}, *figs.List("list"))
		assert.Equal(t, []string{"three", "four", "five"}, *figs.List("list"))
		// map
		assert.NotEqual(t, map[string]string{"one": "1", "two": "2", "three": "3"}, *figs.Map("myMap"))
		assert.Equal(t, map[string]string{"four": "4", "five": "5", "six": "6"}, *figs.Map("myMap"))
		// map alias
		assert.NotEqual(t, map[string]string{"one": "1", "two": "2", "three": "3"}, *figs.Map("map"))
		assert.Equal(t, map[string]string{"four": "4", "five": "5", "six": "6"}, *figs.Map("map"))

		figs = nil

	})

	t.Run("alias_with_int", func(t *testing.T) {
		figs := With(Options{Germinate: true})
		figs.NewInt("count", 42, "usage")
		figs.WithAlias("count", "c")
		assert.NoError(t, figs.Parse())
		assert.Equal(t, 42, *figs.Int("count"))
		assert.Equal(t, 42, *figs.Int("c"))
	})

	t.Run("alias_conflict", func(t *testing.T) {
		figs := With(Options{Germinate: true})
		figs.NewString("one", "value1", "usage")
		figs.NewString("two", "value2", "usage")
		figs.WithAlias("one", "x")
		figs.WithAlias("two", "x") // Should this overwrite or be ignored?
		assert.NoError(t, figs.Parse())
		assert.Equal(t, "value1", *figs.String("x")) // Clarify expected behavior
	})
}
