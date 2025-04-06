package figtree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTree_ListValues(t *testing.T) {
	figs := With(Options{Germinate: true})
	figs.NewList(t.Name(), []string{"yahuah"}, "Name List")
	assert.NoError(t, figs.Parse())
	assert.Contains(t, *figs.List(t.Name()), "yahuah")
}
