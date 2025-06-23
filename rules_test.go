package figtree

import (
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRules(t *testing.T) {
	kName := "name"
	os.Args = []string{os.Args[0]}
	figs := With(Options{Germinate: true, Pollinate: true})
	assert.NotNil(t, figs)
	figs = figs.NewString(kName, "", "usage of name").WithValidator(kName, AssureStringNotEmpty)
	assert.Error(t, figs.Parse())

	figs = figs.WithRule(kName, RuleNoValidations)
	assert.NoError(t, figs.Parse())

	figs = figs.WithCallback(kName, CallbackBeforeVerify, func(_ interface{}) error {
		panic("this shouldn't be called")
		return nil
	})
	assert.Panics(t, func() { _ = figs.Parse() })
	figs = figs.WithRule(kName, RuleNoCallbacks)
	assert.NoError(t, figs.Parse())

	changeEnv := func(n, m string) {
		nu := strings.ToUpper(n)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			timer := time.NewTimer(time.Second * 1)
			checker := time.NewTicker(100 * time.Millisecond)
			for {
				select {
				case <-timer.C:
					return
				case <-checker.C:
					assert.NoError(t, os.Setenv(nu, m))
				}
			}
		}()
		defer assert.NoError(t, os.Unsetenv(nu))
		wg.Wait()
	}

	// verify empty string
	assert.Equal(t, "", *figs.String(kName))

	// assign value
	figs.StoreString(kName, t.Name())

	// verify value set
	assert.Equal(t, t.Name(), *figs.String(kName))

	// assign env value
	changeEnv(kName, "Yeshua")
	// verify value was assigned
	assert.Equal(t, "Yeshua", *figs.String(kName))

	// turn off env
	figs.WithTreeRule(RuleNoEnv)
	changeEnv(kName, "Andrei")
	// value should stay Yeshua - message over messenger
	assert.Equal(t, "Yeshua", *figs.String(kName))

}
