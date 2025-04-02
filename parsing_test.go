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
	exts := []string{"yaml", "json", "ini"}
	for _, ext := range exts {
		t.Run(ext, func(t *testing.T) {
			p := filepath.Join(".", "test.config."+ext)
			figs := With(Options{Germinate: true})
			figs.NewString("name", "", "name")
			figs.WithValidator("name", AssureStringContains("yahuah"))
			figs.NewInt("age", 0, "age")
			figs.WithValidator("age", AssureIntInRange(17, 47))
			figs.NewString("sex", "", "sex")
			figs.WithValidator("sex", AssureStringHasSuffix("male"))
			err := figs.ParseFile(p)
			assert.NoError(t, err)
		})
	}
}
