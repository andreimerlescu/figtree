package figtree

import (
	"flag"
	"os"
	"sync"
	"sync/atomic"
)

// New will initialize the Tree package
// Usage:
//
//				When defining:
//				    figs := figtree.New()
//			     figs.NewInt("workers", 10, "number of workers")
//			     figs.Parse()
//		      OR err := figs.Load()
//		      OR err := figs.ParseFile("path/to/file.json")
//	       THEN workers := *figs.Int("workers") // workers is a regular int
func New() Fruit {
	return With(Options{Tracking: false})
}

// Grow is a memetic alias of New
// Example:
//
//	 ctx, cancel := context.WithCancel(context.Background())
//	 defer cancel()
//		figs := figtree.Grow()
//	 go func() {
//	   for {
//	     select {
//	       case <-ctx.Done():
//	         return
//	       case mutation, ok := <-figs.Mutations():
//	         if ok {
//	           log.Println(mutation)
//	         }
//	     }
//	   }
//	 }()
//
// // figs implements figtree.Fruit interface
func Grow() Fruit {
	return With(Options{Tracking: true})
}

func With(opts Options) Fruit {
	angel := atomic.Bool{}
	angel.Store(true)
	fig := &Tree{
		ConfigFilePath: opts.ConfigFile,
		filterTests:    opts.Germinate,
		pollinate:      opts.Pollinate,
		tracking:       opts.Tracking,
		harvest:        opts.Harvest,
		angel:          &angel,
		figs:           make(map[string]*Fig),
		withered:       make(map[string]Fig),
		mu:             sync.RWMutex{},
		mutationsCh:    make(chan Mutation),
		flagSet:        flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
	}
	angel.Store(false)
	return fig
}
