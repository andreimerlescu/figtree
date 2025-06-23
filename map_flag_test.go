package figtree

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTree_MapKeys(t *testing.T) {
	os.Args = []string{os.Args[0]}
	figs := With(Options{Germinate: true})
	figs = figs.NewMap(t.Name(), map[string]string{"name": "yahuah"}, "Name Map")
	assert.NoError(t, figs.Parse())
	mk := figs.MapKeys(t.Name())
	assert.Contains(t, mk, "name")
}

func TestMapFlag_Set(t *testing.T) {
	t.Run("PolicyMapAppend_TRUE", func(t *testing.T) {
		PolicyMapAppend = true
		defer func() { PolicyMapAppend = false }()
		os.Args = []string{os.Args[0], "-x", "name=yahuah"}
		figs := With(Options{Germinate: true})
		figs = figs.NewMap("x", map[string]string{"job": "bum"}, "Name Map")
		assert.NoError(t, figs.Parse())
		assert.Contains(t, figs.MapKeys("x"), "name")
		assert.Contains(t, figs.MapKeys("x"), "job") // Contains because of PolicyMapAppend
		os.Args = []string{os.Args[0]}
	})
	t.Run("PolicyMapAppend_DEFAULT", func(t *testing.T) {
		PolicyMapAppend = false
		os.Args = []string{os.Args[0], "-x", "name=yeshua,age=33"}
		figs := With(Options{Germinate: true})
		figs = figs.NewMap("x", map[string]string{"job": "bum"}, "Name Map")
		assert.NoError(t, figs.Parse())
		assert.Contains(t, figs.MapKeys("x"), "name")
		assert.Contains(t, figs.MapKeys("x"), "age")
		assert.NotContains(t, figs.MapKeys("x"), "job") // NotContains because no PolicyMapAppend
		os.Args = []string{os.Args[0]}
	})

}
