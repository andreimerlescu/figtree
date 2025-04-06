package figtree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTree_MapKeys(t *testing.T) {
	figs := With(Options{Germinate: true})
	figs.NewMap(t.Name(), map[string]string{"name": "yahuah"}, "Name Map")
	assert.NoError(t, figs.Parse())
	assert.Contains(t, figs.MapKeys(t.Name()), "name")
}
