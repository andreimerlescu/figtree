# Figtree Roadmap

## v2 Major Release

The package `figtree` was upgraded this year to `v2` that supports `.WithValidator()` and `.WithCallback()` handlers.

### v2.0.4 <span style="text-decoration: line-through;">Planned</span> Release

- **String**
    - [X] `figs.WithValidator(aString, figtree.AssureStringNoPrefix(prefix)`
    - [X] `figs.WithValidator(aString, figtree.AssureStringNoPrefixes([]string of prefixes)`
    - [X] `figs.WithValidator(aString, figtree.AssureStringNoSuffix(suffix)`
    - [X] `figs.WithValidator(aString, figtree.AssureStringNoSuffixes([]string of suffix)`
    - [X] `figs.WithValidator(aString, figtree.AssureStringNoPrefix(prefix)`
    - [X] `figs.WithValidator(aString, figtree.AssureStringNoSuffix(suffix)`
    - [X] `figs.WithValidator(aString, figtree.AssureStringHasPrefixes([]string of prefixes)`
    - [X] `figs.WithValidator(aString, figtree.AssureStringHasSuffixes([]string of suffix)`
    - [X] `figs.WithValidator(aString, figtree.AssureStringNoPrefixes([]string of prefixes)`
    - [X] `figs.WithValidator(aString, figtree.AssureStringNoSuffixes([]string of suffix)`
- **List**
    - [X] `figs.ListValues() []string`
- **Map**
    - [X] `figs.MapKeys() []string`
- **Callbacks**
    - [X] `figs.WithCallback(key, figtree.CallbackBeforeVerify, func(value interface{}) error {})` (before verify, run this)
    - [X] `figs.WithCallback(key, figtree.CallbackBeforeChange, func(value interface{}) error {})` (before change, run this)
    - [X] `figs.WithCallback(key, figtree.CallbackBeforeRead, func(value interface{}) error {})` (before read, run this)
- **Rules**
    - [X] `figs.WithRule(key, figtree.RulePreventChange)` (if property is changed, block the attempt)
    - [X] `figs.WithRule(key, figtree.RulePanicOnChange)` (call `panic()` when value changes)
    - [X] `figs.WithRule(key, figtree.RuleNoVerify)` (disable verification on rule))
     
### v2.1.0 Planned Release

Adding two new **Mutagenesis** types called `File` and `Directory`.

- **File**
    - [ ] `figs.NewFile(key, path, usage)` (create a new file path configurable property)
    - [ ] `figs.StoreFile(key, newPath)` (updates a file path)
    - [ ] `figs.File(key) (path string)` (retrieves a file path)
    - [ ] `figs.FileContents(key) []byte` (retrieve the file's contents)
    - [ ] `figs.FileHandler(key) *os.File` (opens a file and returns the handler)
    - [ ] `figs.FileWriteContents(key, newContents []byte) error` (writes newContents into file)
    - [ ] `figs.WithValidator(key, figtree.AssureFileExists)` (checks if the file exists)
    - [ ] `figs.WithValidator(key, figtree.AssureFileTouchIfNotExists())` (touches if the file if it does not exists)
    - [ ] `figs.WithValidator(key, figtree.AssureFileCanBeModified)` (checks if file can be modified)
    - [ ] `figs.WithValidator(key, figtree.AssureFileSizeGreaterThan(int))` (checks file size greater than value)
    - [ ] `figs.WithValidator(key, figtree.AssureFileSizeLessThan(int))` (checks file size less than value)
    - [ ] `figs.WithValidator(key, figtree.AssureFileModeIs(os.FileMode))` (checks file mode value)
    - [ ] `figs.WithValidator(key, figtree.AssureFileOwnerIs(string))` (checks file owner)
    - [ ] `figs.WithValidator(key, figtree.AssureFileGroupIs(string))` (checks file group)
- **Directory**
    - [ ] `figs.NewDirectory(key, path, usage)` (create a new directory path configurable property)
    - [ ] `figs.Directory(key) (path string)` (retrieves a directory path)
    - [ ] `figs.StoreDirectory(key, newPath)` (updates a directory path)
    - [ ] `figs.DirectoryFlushAll(key)` (recursively removes elements inside directory path)
    - [ ] `figs.WithValidator(key, figtree.AssureDirCreateIfNotExists)` (checks if directory exists and creates it)
    - [ ] `figs.WithValidator(key, figtree.AssureDirExists)` (checks if directory exists)
    - [ ] `figs.WithValidator(key, figtree.AssureDirIsReadable)` (checks for chmod value of dir)
    - [ ] `figs.WithValidator(key, figtree.AssureDirIsWritable)` (checks for chmod value of dir)
    - [ ] `figs.WithValidator(key, figtree.AssureDirOwnerIs(string))` (checks directory owner)
    - [ ] `figs.WithValidator(key, figtree.AssureDirGroupIs(string))` (checks directory group)
    - [ ] `figs.WithValidator(key, figtree.AssureDirChmod(os.FileMode)` (checks for chmod value of dir)
    - [ ] `figs.WithValidator(key, figtree.AssureDirMorePermissiveThan(os.FileMode))` (checks for chmod value of dir)
    - [ ] `figs.WithValidator(key, figtree.AssureDirLessPermissiveThan(os.FileMode))` (checks for chmod value of dir)

This update is going to incorporate the https://github.com/andreimerlescu/checkfs into the `figtree`
to support `File` and `Directory` `Mutagenesis` types. 

### v2.2.0 Planned Release

- **Define**
    - [ ] `figs.Define(property, default, usage, validators, callbacks, rules)` (all-in-one define property)
- **Short** / **Alias**
    - [ ] `figs.Short(property, propertyShortAlias)` (the irony of the 2nd arg 😹 this lets you do `-debug` as `-d`) 
    - [ ] `figs.Alias(property, alias)` (same as `.Short()` just an alias 😹)
- **Semaphore**
    - [ ] `figs.NewSemaphore(property, default, usage)`
    - [ ] `figs.Semaphore(property)` (usage `*figs.Semaphore(name).Acquire()` and `*figs.Semaphore(name).Release()`)
    - [ ] `figs.StoreSemaphore(property, newLimit)` (requires `figtree.RuleUseSmartChannels` to check if channel is empty)
    - [ ] `figs.WithValidator(property, figtree.AssureSemaphoreLessThan(1)` (tune semaphore rules by configuration policies)
    - [ ] `figs.WithValidator(property, figtree.AssureSemaphoreGreaterThan(2)` (tune semaphore rules by configuration policies)
    - [ ] `figs.WithValidator(property, figtree.AssureSemaphoreIs(3)` (tune semaphore rules by configuration policies))
- **Rules**
    - [ ] `figs.WithRule(property, figtree.RuleUseSmartChannels)`

This update is going to incorporate the https://github.com/andreimerlescu/sema package into the `figtree` to support
the `Semaphore` `Mutagenesis` type.


## v1 Major Release

The package `figtree` was introduced as `v1` for a short period of time since it was a fork and rewrite of the 
`configurable` package hosted at https://github.com/andreimerlescu/configurable. This package is over 2 years old
and has been used in a dozens of production Go applications that serve hundreds of thousands of requests efficiently. 

The name `figtree` was a derivation from the con<b>FIG</b>urable package name and the memetics of the functions 
were made to reflect what you'd expect out of nature with an actual Fig Tree.