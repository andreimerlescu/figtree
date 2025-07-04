package figtree

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTree_ListValues(t *testing.T) {
	figs := With(Options{Germinate: true})
	figs.NewList(t.Name(), []string{"yahuah"}, "Name List")
	assert.NoError(t, figs.Parse())
	assert.Contains(t, *figs.List(t.Name()), "yahuah")
}

func TestListFlag_Set(t *testing.T) {
	t.Run("PolicyListAppend_TRUE", func(t *testing.T) {
		PolicyListAppend = true
		defer func() { PolicyListAppend = false }()
		os.Args = []string{os.Args[0], "-x", "yeshua"}
		figs := With(Options{Germinate: true})
		figs.NewList("x", []string{"bum"}, "Name List")
		assert.NoError(t, figs.Parse())
		assert.Equal(t, "bum,yeshua", figs.FigFlesh("x").ToString())
		assert.Contains(t, *figs.List("x"), "yeshua")
		assert.Contains(t, *figs.List("x"), "bum")
		os.Args = []string{os.Args[0]}
	})
	t.Run("PolicyListAppend_DEFAULT", func(t *testing.T) {
		PolicyListAppend = false
		os.Args = []string{os.Args[0], "-x", "yahuah"}
		figs := With(Options{Germinate: true})
		figs.NewList("x", []string{"bum"}, "Name List")
		assert.NoError(t, figs.Parse())
		assert.Equal(t, "yahuah", figs.FigFlesh("x").ToString())
		assert.Contains(t, *figs.List("x"), "yahuah")
		assert.NotContains(t, *figs.List("x"), "bum") // NotContains because of PolicyListAppend
		os.Args = []string{os.Args[0]}
	})
}
