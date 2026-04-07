package figtree

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithAlias_ConflictsWithExistingFig(t *testing.T) {
	os.Args = []string{os.Args[0]}
	figs := With(Options{Germinate: true})
	figs = figs.NewString("long", "default", "usage")
	figs = figs.NewString("short", "default", "usage")

	// Attempt to register "short" as an alias for "long" — but "short" is
	// already a registered fig name, so this should record a problem and
	// not panic.
	figs = figs.WithAlias("long", "short")

	assert.NoError(t, figs.Parse())
	problems := figs.(*figTree).Problems()
	assert.Len(t, problems, 1, "expected one problem recorded for alias conflict")
	assert.Contains(t, problems[0].Error(), "conflicts with existing fig name")
}

func TestConcurrentPollinateReads(t *testing.T) {
	os.Args = []string{os.Args[0]}
	figs := With(Options{Pollinate: true, Germinate: true, Tracking: false})
	figs.NewString("concurrent_key", "initial", "usage")
	assert.NoError(t, figs.Parse())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	defer os.Unsetenv("CONCURRENT_KEY")

	go func() {
		vals := []string{"alpha", "beta", "gamma"}
		i := 0
		ticker := time.NewTicker(5 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				os.Setenv("CONCURRENT_KEY", vals[i%3])
				i++
			}
		}
	}()

	var wg sync.WaitGroup
	for n := 0; n < 10; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ticker := time.NewTicker(10 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					_ = figs.String("concurrent_key")
				}
			}
		}()
	}
	wg.Wait()
}

func TestWithAlias(t *testing.T) {
	const cmdLong, cmdAliasLong, valueLong, usage = "long", "l", "default", "usage"
	const cmdShort, cmdAliasShort, valueShort = "short", "s", "default"

	t.Run("basic_usage", func(t *testing.T) {
		os.Args = []string{os.Args[0], "-l", t.Name()}
		figs := With(Options{Germinate: true, Tracking: false})
		figs = figs.NewString(cmdLong, valueLong, usage)
		figs = figs.WithAlias(cmdLong, cmdAliasLong)
		assert.NoError(t, figs.Parse())

		assert.NotEqual(t, valueLong, *figs.String(cmdLong))
		assert.Equal(t, t.Name(), *figs.String(cmdLong))
		assert.NotEqual(t, valueLong, *figs.String(cmdAliasLong))
		assert.Equal(t, t.Name(), *figs.String(cmdAliasLong))
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
		os.Args = []string{os.Args[0]}
		const k, v, u = "name", "yeshua", "the real name of god"
		ka1 := "father"
		ka2 := "son"
		ka3 := "rauch-hokadesch"
		figs := With(Options{Germinate: true, Tracking: false})
		figs = figs.NewString(k, v, u)
		figs = figs.WithAlias(k, ka1)
		figs = figs.WithAlias(k, ka2)
		figs = figs.WithAlias(k, ka3)
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
		figs = figs.NewString(cmdLong, valueLong, usage)
		figs = figs.WithAlias(cmdLong, cmdAliasLong)
		figs = figs.WithValidator(cmdLong, AssureStringNotEmpty)

		// short
		figs = figs.NewString(cmdShort, valueShort, usage)
		figs = figs.WithAlias(cmdShort, cmdAliasShort)
		figs = figs.WithValidator(cmdShort, AssureStringNotEmpty)

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
		assert.Equal(t, []string{"five", "four", "three"}, *figs.List("myList"))

		// list alias
		assert.NotEqual(t, []string{"one", "two", "three"}, *figs.List("list"))
		assert.Equal(t, []string{"five", "four", "three"}, *figs.List("list"))
		// map
		assert.NotEqual(t, map[string]string{"one": "1", "two": "2", "three": "3"}, *figs.Map("myMap"))
		assert.Equal(t, map[string]string{"four": "4", "five": "5", "six": "6"}, *figs.Map("myMap"))
		// map alias
		assert.NotEqual(t, map[string]string{"one": "1", "two": "2", "three": "3"}, *figs.Map("map"))
		assert.Equal(t, map[string]string{"four": "4", "five": "5", "six": "6"}, *figs.Map("map"))

		figs = nil
	})

	t.Run("alias_with_int", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		figs := With(Options{Germinate: true})
		figs = figs.NewInt("count", 42, "usage")
		figs = figs.WithAlias("count", "c")
		assert.NoError(t, figs.Parse())
		assert.Equal(t, 42, *figs.Int("count"))
		assert.Equal(t, 42, *figs.Int("c"))
	})

	t.Run("alias_conflict", func(t *testing.T) {
		os.Args = []string{os.Args[0]}
		figs := With(Options{Germinate: true})
		figs = figs.NewString("one", "value1", "usage")
		figs = figs.NewString("two", "value2", "usage")
		figs = figs.WithAlias("one", "x")
		figs = figs.WithAlias("two", "x") // Should this overwrite or be ignored?
		assert.NoError(t, figs.Parse())
		assert.Equal(t, "value1", *figs.String("x")) // Clarify expected behavior
		us := figs.UsageString()
		assert.NotEmpty(t, us)
	})
}
