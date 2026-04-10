# Opinionated Origins: Why Figtree Is Named the Way It Is

## Words Matter

Every package makes a choice about what to name things. Most developers treat 
naming as cosmetic — a label on a box. Figtree treats naming as architectural. 
The words you use to describe a system shape the system itself. They shape what 
you build next, what feels natural to extend, and what feels like it doesn't 
belong.

Viper chose a snake. Cobra. Viper. The fangs-first family of Go configuration 
tooling. A snake is flat. It has no branches. It moves in one direction. It 
strikes. If you've ever been bitten by a race condition in a viper-backed 
application under concurrent load, the metaphor lands harder than intended.

Figtree chose a tree. Specifically, the fig tree — one of the oldest cultivated 
plants in human history, appearing in more foundational texts than almost any 
other living thing. Not because of religion. Because of what a fig tree actually 
does.

## What a Fig Tree Does

A fig tree grows from a seed. It puts down roots that draw from the environment. 
It branches. Its branches bear fruit. The fruit contains the value. You harvest 
the fruit. You don't harvest the tree.

That's the entire figtree API described in one paragraph.

- `figtree.New()` or `figtree.Grow()` — plant the seed
- `figs.NewString()`, `figs.NewInt()` — register properties, grow the branches
- `figs.Parse()` or `figs.Load()` — draw from the environment
- `*figs.String()`, `*figs.Int()` — harvest the fruit
- `figs.Fig(key)` — access the fruit itself, not just its value
- `figs.Mutations()` — observe how the fruit changes over time
- `figs.Withered` — the original state before the environment changed it
- `figs.Resurrect()` — regrow from dormant roots
- `figs.Curse()` / `figs.Recall()` — dormancy and renewal
- `figs.Pollinate()` — external forces updating the living tree
- `figs.Branch()` — the tree literally branches, each branch its own tree

None of these words were chosen arbitrarily. Each one maps to a real biological 
process that has a direct analog in what the code does. When a word fits its 
behavior this precisely, the API becomes memorable. You don't have to look up 
what `Resurrect` does. You already know.

## Why the Naming Convention Enables Functionality Viper Lacks

This isn't just aesthetics. The tree metaphor opened design space that a flat 
snake metaphor closes off.

A snake is flat. Viper's configuration is a flat map with dot-notation keys 
pretending to be hierarchy. `viper.Get("db.host")` isn't a branch — it's a 
string with a period in it. The hierarchy is a convention, not a structure. You 
can't put rules on `db`. You can't put validators on `db`. You can't watch 
`db` change independently. You can't scope callbacks to `db`. Because `db` 
doesn't exist. Only `"db.host"` exists, as a string key in a flat map.

A tree branches. `figs.NewBranch("db")` returns a real figtree. That branch 
has its own validators, its own callbacks, its own rules, its own aliases. It 
funnels mutations back to the root channel with the path recorded. It can be 
passed around, composed, and reasoned about independently. `db` exists because 
a tree has branches and a snake does not.

The naming convention didn't just describe the API. It generated the API. Once 
you commit to the metaphor honestly, the next feature becomes obvious. Trees 
have branches. Branches bear fruit. Fruit can be withered or fresh. The tree 
can be cursed or recalled. Pollination brings external changes in. The language 
tells you what to build next.

## The Concurrent Application Problem

Figtree has been backing highly concurrent Go applications in enterprise 
environments for over six years — first as `configurable`, then as `figs`, now 
as `figtree`. Not because viper was unavailable. Because viper's known race 
conditions were unacceptable in production systems that couldn't afford to lock 
an external mutex around every configuration read.

A tree is alive. It doesn't stop growing because something is reading its fruit. 
The concurrency model in figtree — `sync.RWMutex` used correctly throughout, 
read operations using `RLock` not `Lock`, `Store()` releasing locks before 
channel sends to prevent deadlock — these decisions came from six years of 
running configuration management under real concurrent load. The metaphor of a 
living tree that continues to operate while being observed isn't decorative. 
It's the design requirement.

## On Adoption

Figtree was built by someone who was adopted. Viper was not something to be 
adopted. The irony is precise and intentional.

The packages you choose for your applications are not neutral. They carry the 
assumptions of their authors. A package named after a predator assumes adversarial 
relationships — between the config and the application, between defaults and 
overrides, between what you asked for and what you got. A package named after a 
living, branching, fruit-bearing tree assumes that configuration should grow 
with your application, bear real value, and be observable as it changes.

The words you use manifest. Choose them carefully.

## On Yeshua

Figtree isn't a religious package. It doesn't require belief in anything except 
that good software is worth building carefully and naming honestly.

But the fig tree appears in contexts most developers will recognize at some level 
whether they identify as religious or not. It represents discernment — the 
ability to look at something and know whether it bears fruit or doesn't. That 
instinct is what every engineer uses when evaluating a dependency.

Does this bear fruit? Does it do what it says? Is it alive or is it dead wood?

Figtree bears fruit. That's the whole argument.

---

*Figtree is maintained by [@andreimerlescu](https://github.com/andreimerlescu).*  
*License: MIT*
