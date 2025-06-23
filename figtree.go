package figtree

import (
	"flag"
	"os"
	"sync"
	"sync/atomic"
)

// New will initialize the figTree package
// Usage:
//
//				When defining:
//				    figs := figtree.New()
//			     figs.NewInt("workers", 10, "number of workers")
//			     figs.Parse()
//		      OR err := figs.Load()
//		      OR err := figs.ParseFile("path/to/file.json")
//	       THEN workers := *figs.Int("workers") // workers is a regular int
func New() Plant {
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
// // figs implements figtree.Plant interface
func Grow() Plant {
	return With(Options{Tracking: true})
}

// With creates a new fig figTree with predefined Options
//
// Example:
//
//		figs := With(Options{
//			ConfigFilePath: "/opt/app/production.config.yaml",
//			Pollinate: true,
//			Tracking: true,
//			Harvest: 369,
//		})
//		// define your figs.New<Mutagenesis>() here...
//		figs.NewString("domain", "", "Domain name")
//		figs.WithValidator("domain", figtree.AssureStringNotEmpty)
//	 err := figs.Load()
//	 domainInProductionConfigYamlFile := *figs.String("domain")
func With(opts Options) Plant {
	angel := atomic.Bool{}
	angel.Store(true)
	fig := &figTree{
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
	fig.flagSet.Usage = fig.Usage
	angel.Store(false)
	if opts.IgnoreEnvironment {
		os.Clearenv()
	}
	return fig
}
