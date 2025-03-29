package figtree

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestTree_Parse(t *testing.T) {
	figs := With(Options{Germinate: true})
	assert.Nil(t, figs.Parse())
}

func TestTree_ParseFile(t *testing.T) {
	p := filepath.Join(".", "test.config.yaml")
	figs := With(Options{Germinate: true})
	figs.NewString("name", "", "name")
	figs.NewInt("age", 0, "age")
	figs.NewString("sex", "", "sex")
	err := figs.ParseFile(p)
	assert.Nil(t, err)
}
