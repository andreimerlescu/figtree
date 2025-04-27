package figtree

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTree_MapKeys(t *testing.T) {
	figs := With(Options{Germinate: true})
	figs.NewMap(t.Name(), map[string]string{"name": "yahuah"}, "Name Map")
	assert.NoError(t, figs.Parse())
	assert.Contains(t, figs.MapKeys(t.Name()), "name")
}

func TestMapFlag_Set(t *testing.T) {
	t.Run("PolicyMapAppend_TRUE", func(t *testing.T) {
		PolicyMapAppend = true
		os.Args = []string{os.Args[0], "-x", "name=yahuah"}
		figs := With(Options{Germinate: true})
		figs.NewMap("x", map[string]string{"job": "bum"}, "Name Map")
		assert.NoError(t, figs.Parse())
		assert.Equal(t, "", figs.Fig("x").ToString())
		assert.Contains(t, figs.MapKeys("x"), "name")
		assert.Contains(t, figs.MapKeys("x"), "job") // Contains because of PolicyMapAppend
		os.Args = []string{os.Args[0]}
	})
	t.Run("PolicyMapAppend_DEFAULT", func(t *testing.T) {
		PolicyMapAppend = false
		os.Args = []string{os.Args[0], "-x", "name=yahuah"}
		figs := With(Options{Germinate: true})
		figs.NewMap("x", map[string]string{"job": "bum"}, "Name Map")
		assert.NoError(t, figs.Parse())
		assert.Equal(t, "", figs.Fig("x").ToString())
		assert.Contains(t, figs.MapKeys("x"), "name")
		assert.NotContains(t, figs.MapKeys("x"), "job") // NotContains because no PolicyMapAppend
		os.Args = []string{os.Args[0]}
	})

}
