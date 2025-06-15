package figtree

import (
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

	t.Run("complex_usage", func(t *testing.T) {

		figs := With(Options{Germinate: true, Tracking: false})
		// long
		figs.NewString(cmdLong, valueLong, usage)
		figs.WithAlias(cmdLong, cmdAliasLong)
		figs.WithValidator(cmdLong, AssureStringNotEmpty)

		// short
		figs.NewString(cmdShort, valueShort, usage)
		figs.WithAlias(cmdShort, cmdAliasShort)
		figs.WithValidator(cmdShort, AssureStringNotEmpty)

		assert.NoError(t, figs.Parse())
		// long
		assert.Equal(t, valueLong, *figs.String(cmdLong))
		assert.Equal(t, valueLong, *figs.String(cmdAliasLong))
		// short
		assert.Equal(t, valueShort, *figs.String(cmdShort))
		assert.Equal(t, valueShort, *figs.String(cmdAliasShort))

		figs = nil

	})
}
