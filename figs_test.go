package figtree

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWith(t *testing.T) {
	type args struct {
		opts Options
	}
	tests := []struct {
		name      string
		args      args
		wantNil   bool
		wantTrack bool
		wantFile  string
	}{
		{
			name:      "With tracking enabled",
			args:      args{opts: Options{Tracking: true}},
			wantNil:   false,
			wantTrack: true,
			wantFile:  "",
		},
		{
			name:      "With tracking disabled and custom config file",
			args:      args{opts: Options{Tracking: false, ConfigFile: "custom.yaml"}},
			wantNil:   false,
			wantTrack: false,
			wantFile:  "custom.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment to ensure clean slate
			os.Clearenv()

			figs := With(tt.args.opts)
			assert.Equal(t, tt.wantNil, figs == nil, "With(%v) should return a non-nil Fruit", tt.args.opts)
			if !tt.wantNil {
				tree, ok := figs.(*Tree)
				assert.True(t, ok, "With() should return a *Tree")
				assert.Equal(t, tt.wantTrack, tree.tracking, "With() should set tracking to %v", tt.wantTrack)
				assert.NotNil(t, tree.figs, "With() should initialize figs map")
				assert.NotNil(t, tree.withered, "With() should initialize withered map")
				assert.NotNil(t, tree.mutationsCh, "With() should initialize mutationsCh")
				assert.Equal(t, tt.wantFile, tree.ConfigFilePath, "With() should set ConfigFilePath to %v", tt.wantFile)
				// Verify mutationsCh behavior based on tracking
				select {
				case _, ok := <-figs.Mutations():
					assert.True(t, ok, "Mutations channel should be open after With")
				default:
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		wantNil   bool
		wantTrack bool
	}{
		{
			name:      "New creates tree without tracking",
			wantNil:   false,
			wantTrack: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment to ensure clean slate
			os.Clearenv()

			figs := New()
			assert.Equal(t, tt.wantNil, figs == nil, "New() should return a non-nil Fruit")
			if !tt.wantNil {
				tree, ok := figs.(*Tree)
				assert.True(t, ok, "New() should return a *Tree")
				assert.Equal(t, tt.wantTrack, tree.tracking, "New() should disable tracking by default")
				assert.NotNil(t, tree.figs, "New() should initialize figs map")
				assert.NotNil(t, tree.withered, "New() should initialize withered map")
				assert.NotNil(t, tree.mutationsCh, "New() should initialize mutationsCh")
				// Verify mutationsCh is open (even if not tracking, channel should exist)
				select {
				case _, ok := <-figs.Mutations():
					assert.True(t, ok, "Mutations channel should be open after New")
				default:
				}
			}
		})
	}
}

func TestGrow(t *testing.T) {
	tests := []struct {
		name      string
		wantNil   bool
		wantTrack bool
	}{
		{
			name:      "Grow creates tree with tracking",
			wantNil:   false,
			wantTrack: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment to ensure clean slate
			os.Clearenv()

			figs := Grow()
			assert.Equal(t, tt.wantNil, figs == nil, "Grow() should return a non-nil Fruit")
			if !tt.wantNil {
				tree, ok := figs.(*Tree)
				assert.True(t, ok, "Grow() should return a *Tree")
				assert.True(t, tree.tracking, "Grow() should enable tracking")
				assert.NotNil(t, tree.figs, "Grow() should initialize figs map")
				assert.NotNil(t, tree.withered, "Grow() should initialize withered map")
				assert.NotNil(t, tree.mutationsCh, "Grow() should initialize mutationsCh")
				// Verify mutationsCh is open
				select {
				case _, ok := <-figs.Mutations():
					assert.True(t, ok, "Mutations channel should be open after Grow")
				default:
				}
			}
		})
	}
}

func TestIsTracking(t *testing.T) {
	// Add a timeout to ensure the test fails if it runs too long
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute+34*time.Second)
	defer cancel()

	// Run the test with a timeout
	t.Run("IsTracking", func(t *testing.T) {
		// Create a figtree with Tracking and Germinate enabled
		figs := With(Options{Tracking: true, Germinate: true})
		var k, d, u = "name", "yahuah", "usage"
		figs.NewString(k, d, u)
		assert.Nil(t, figs.Parse())

		// Use a WaitGroup to synchronize goroutines
		var wg sync.WaitGroup

		// Use atomic counters to collect errors and count writes
		var errorCount int32
		var writeCount int32

		// Collect mutations from the Mutations channel
		mutations := make([]Mutation, 0, 500)
		var mutationsMu sync.Mutex
		done := make(chan struct{})
		wg.Add(1)
		go func() {
			defer wg.Done()
			for mutation := range figs.Mutations() {
				mutationsMu.Lock()
				mutations = append(mutations, mutation)
				mutationsMu.Unlock()
				// t.Logf("Received mutation: %v", mutation)
			}
			close(done)
		}()

		runFor := time.NewTimer(7 * time.Second)
		readEvery := time.NewTicker(3 * time.Millisecond)
		writeEvery := time.NewTicker(77 * time.Millisecond)

		// Channel to collect errors
		errChan := make(chan error, 500)

		// Main loop to handle reads and writes
		wg.Add(1)
		go func() {
			defer wg.Done()
			// start := time.Now()
			for {
				select {
				case <-ctx.Done():
					t.Error("Test timed out after 2 minutes")
					return
				case <-runFor.C:
					// case n := <-readEvery.C:
					// t.Logf("Timer fired at %v after %v", n, time.Since(start))
					readEvery.Stop()
					writeEvery.Stop()
					// Close the mutations channel to allow the collector goroutine to exit
					// t.Log("cursing the fig tree...")
					figs.Curse()
					return
				case <-readEvery.C:
					// case n := <-runFor.C:
					// t.Logf("Read ticker fired at %v", n)
					val := figs.String(k)
					if val == nil {
						errChan <- fmt.Errorf("String(%q) returned nil", k)
						atomic.AddInt32(&errorCount, 1)
						continue
					}
					// t.Logf("Read: %s", *val)
				case <-writeEvery.C:
					// case n := <-writeEvery.C:
					// t.Logf("Write ticker fired at %v", n)
					// Toggle the value between "yahuah" and "andrei"
					if d == "yahuah" {
						d = "andrei"
					} else {
						d = "yahuah"
					}
					figs.StoreString(k, d)
					atomic.AddInt32(&writeCount, 1)
					// t.Logf("Wrote: %s (write #%d)", d, writeCount)
				}
			}
		}()

		// Wait for goroutines to finish
		wg.Wait()
		<-done

		// Collect and log any errors
		errors := make([]error, 0, len(errChan))
		for len(errChan) > 0 {
			err := <-errChan
			if err != nil {
				errors = append(errors, err)
			}
		}

		// Log errors for debugging
		for i, err := range errors {
			t.Logf("Error %d: %v", i+1, err)
		}

		// Verify that no errors occurred
		assert.Equal(t, int32(0), errorCount, "Test failed with %d errors", errorCount)

		actualWrites := int(atomic.LoadInt32(&writeCount))
		assert.Equal(t, actualWrites, len(mutations), "Number of mutations does not match number of writes")

		// Verify the final value
		finalVal := figs.String(k)
		assert.NotNil(t, finalVal)
	})
}
