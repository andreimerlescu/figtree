# Figtree History

## v0.0.1 

The first release of the `figtree` idea was introduced as https://github.com/andreimerlescu/configurable
and this package was `v0.0.1` offering an initial run-down of a basic command line argument parsing.

This package was in use for well over two calendar years getting tested in production live environments.

## v1.0.0

The `configurable` package was forked into https://github.com/andreimerlescu/figs because I originally
really liked the idea of doing `figs.New()` to create a new con<b>FIG</b>urable module of arguments...
but I didn't like how I was doing `cfigs := figs.New()` ; that felt off, so this release was trying
out the rename of the `configurable` to `figs`; but was discontinued.

## v1.0.1 

This version of `figs` got `.WithValidator()` added to it. 

## v2.0.0

This released version of `figs` was the same as `v1.0.1` but rebranded as `figtree` and tagged at `v2.0.0`.

## v2.0.1 

This release had minor bux fixes for the `go.mod` path of `/v2` namespacing. 

## v2.0.2 

This release introduced `.WithCallback()` to the `figtree`. 

## v2.0.3 

This release introduced `.Version()` to the `figtree` by embedding the `VERSION` file into the package.

# Roadmap

The [Road Map](/ROADMAP.md) can be viewed in the [ROADMAP.md](/ROADMAP.md). Just clicky clicky the pinky linky.