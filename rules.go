package figtree

type RuleKind int

const (
	RuleUndefined                 RuleKind = iota // RuleUndefined is default and does no action
	RulePreventChange             RuleKind = iota // RulePreventChange blocks Mutagensis Store methods
	RulePanicOnChange             RuleKind = iota // RulePanicOnChange will throw a panic on the Mutagenesis Store methods
	RuleNoValidations             RuleKind = iota // RuleNoValidations will skip over all WithValidator assignments
	RuleNoCallbacks               RuleKind = iota // RuleNoCallbacks will skip over all WithCallback assignments
	RuleCondemnedFromResurrection RuleKind = iota // RuleCondemnedFromResurrection will panic if there is an attempt to resurrect a condemned fig
	RuleNoMaps                    RuleKind = iota // RuleNoMaps blocks NewMap, StoreMap, and Map from being called on the Tree
	RuleNoLists                   RuleKind = iota // RuleNoLists blocks NewList, StoreList, and List from being called on the Tree
	RuleNoFlags                   RuleKind = iota // RuleNoFlags disables the flag package from the Tree
	RuleNoEnv                     RuleKind = iota // RuleNoEnv skips over all os.Getenv related logic
)

func (tree *figTree) HasRule(rule RuleKind) bool {
	if rule == RuleUndefined {
		return false
	}
	for _, r := range tree.GlobalRules {
		if r == rule {
			return true
		}
	}
	return false
}

func (fig *figFruit) HasRule(rule RuleKind) bool {
	for _, r := range fig.Rules {
		if r == rule {
			return true
		}
	}
	return false
}

func (tree *figTree) WithTreeRule(rule RuleKind) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.GlobalRules = append(tree.GlobalRules, rule)
	return tree
}

// WithRule attaches a Rule to to the Fig
func (tree *figTree) WithRule(name string, rule RuleKind) Plant {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	fruit, exists := tree.figs[name]
	if !exists || fruit == nil {
		tree.mu.Unlock()
		tree.Resurrect(name)
		tree.mu.Lock()
		fruit = tree.figs[name]
	}
	if fruit == nil {
		return tree
	}
	fruit.Rules = append(fruit.Rules, rule)
	tree.figs[name] = fruit
	return tree
}
