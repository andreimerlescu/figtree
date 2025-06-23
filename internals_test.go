package figtree

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTree_checkAndSetFromEnv(t *testing.T) {
	const k, u = "workers-check-and-set-from-env", "usage"

	// create a new fig tree internally to test the func checkAndSetFromEnv
	angel := atomic.Bool{}
	angel.Store(true)
	opts := Options{Germinate: true}
	figs := &figTree{
		ConfigFilePath: opts.ConfigFile,
		ignoreEnv:      opts.IgnoreEnvironment,
		filterTests:    opts.Germinate,
		pollinate:      opts.Pollinate,
		tracking:       opts.Tracking,
		harvest:        opts.Harvest,
		angel:          &angel,
		problems:       make([]error, 0),
		aliases:        make(map[string]string),
		figs:           make(map[string]*figFruit),
		values:         &sync.Map{},
		withered:       make(map[string]witheredFig),
		mu:             sync.RWMutex{},
		mutationsCh:    make(chan Mutation),
		flagSet:        flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
	}
	figs.flagSet.Usage = figs.Usage
	angel.Store(false)
	if opts.IgnoreEnvironment {
		os.Clearenv()
	}

	figs = figs.NewInt(k, 10, "Number").(*figTree)

	// assign an int to k
	assert.NoError(t, figs.Parse())

	// verify assignment
	assert.Equal(t, 10, *figs.Int(k))

	// set env for k to 17
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
				assert.NoError(t, os.Setenv(strings.ToUpper(k), "17"))
			}
		}
	}()
	defer assert.NoError(t, os.Unsetenv(k))
	wg.Wait()
	figs.checkAndSetFromEnv(k)
	assert.Equal(t, 17, *figs.Int(k))
}

func TestTree_setValue(t *testing.T) {
	type fields struct {
		ConfigFilePath string
		figs           map[string]*figFruit
		withered       map[string]witheredFig
		aliases        map[string]string
		mu             sync.RWMutex
		tracking       bool
		mutationsCh    chan Mutation
	}
	type args struct {
		flagVal interface{}
		value   interface{}
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantValue interface{}
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "Set int value",
			fields: fields{
				figs:        make(map[string]*figFruit),
				withered:    make(map[string]witheredFig),
				aliases:     make(map[string]string),
				mutationsCh: make(chan Mutation, 1),
			},
			args: args{
				flagVal: new(int),
				value:   42,
			},
			wantValue: 42,
			wantErr:   assert.NoError,
		},
		{
			name: "Set string value",
			fields: fields{
				figs:        make(map[string]*figFruit),
				withered:    make(map[string]witheredFig),
				aliases:     make(map[string]string),
				mutationsCh: make(chan Mutation, 1),
			},
			args: args{
				flagVal: new(string),
				value:   "hello",
			},
			wantValue: "hello",
			wantErr:   assert.NoError,
		},
		{
			name: "Invalid type",
			fields: fields{
				figs:        make(map[string]*figFruit),
				withered:    make(map[string]witheredFig),
				aliases:     make(map[string]string),
				mutationsCh: make(chan Mutation, 1),
			},
			args: args{
				flagVal: new(float32), // Unsupported type
				value:   3.14,
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fig := &figTree{
				ConfigFilePath: tt.fields.ConfigFilePath,
				figs:           tt.fields.figs,
				withered:       tt.fields.withered,
				mu:             tt.fields.mu,
				values:         &sync.Map{},
				tracking:       tt.fields.tracking,
				mutationsCh:    tt.fields.mutationsCh,
				flagSet:        flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
				filterTests:    true,
			}
			err := fig.setValue(tt.args.flagVal, tt.args.value)
			if tt.wantErr(t, err, fmt.Sprintf("setValue(%v, %v)", tt.args.flagVal, tt.args.value)) {
				return
			}

			// Verify value only if no error is expected
			if tt.wantErr(t, nil) {
				switch v := tt.args.flagVal.(type) {
				case *int:
					assert.Equal(t, tt.wantValue, *v, "setValue should set int value")
				case *string:
					assert.Equal(t, tt.wantValue, *v, "setValue should set string value")
				}
			}
		})
	}
}

func TestTree_setValuesFromMap(t *testing.T) {
	tree := &figTree{
		figs:        make(map[string]*figFruit),
		withered:    make(map[string]witheredFig),
		aliases:     make(map[string]string),
		mu:          sync.RWMutex{},
		tracking:    false,
		values:      &sync.Map{},
		mutationsCh: make(chan Mutation, 1),
		flagSet:     flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
		filterTests: true,
	}
	m := map[string]interface{}{
		"name": "yahuah",
		"age":  33,
		"sex":  "male",
	}
	err := tree.setValuesFromMap(m)
	assert.NoError(t, err)
	assert.Equal(t, "yahuah", m["name"])
	assert.Equal(t, "male", m["sex"])
}
