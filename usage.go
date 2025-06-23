package figtree

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/term"
)

// Usage prints a helpful table of figs in a human-readable format
func (tree *figTree) Usage() {
	fmt.Println(tree.UsageString())
}

func (tree *figTree) UsageString() string {
	termWidth := 80
	if term.IsTerminal(int(os.Stdout.Fd())) {
		if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			termWidth = width
		}
	}

	// First pass: Collect all primary flag names and aliases correctly
	// Also calculate maxFlagLen for formatting
	type flagInfo struct {
		name         string
		aliases      []string
		defValue     string
		usage        string
		mutagenesis  Mutagenesis
		isAlias      bool   // Mark if this is an alias entry itself
		originalName string // For aliases, store the original flag name
	}
	allFlagData := make(map[string]*flagInfo) // Use map to deduplicate and process by main flag name

	tree.mu.RLock() // Lock for accessing tree.figs and tree.aliases

	// Populate with main flags
	for name, fruit := range tree.figs {
		f := tree.flagSet.Lookup(name) // Get the flag.Flag object
		if f == nil {
			continue // Should not happen if figs map is consistent with flagSet
		}

		info := &flagInfo{
			name:        f.Name,
			defValue:    f.DefValue,
			usage:       f.Usage,
			mutagenesis: fruit.Mutagenesis, // Get mutagenesis from figFruit
			isAlias:     false,
		}
		allFlagData[name] = info
	}

	// Add aliases to their primary flags
	for alias, originalName := range tree.aliases {
		if info, ok := allFlagData[originalName]; ok {
			info.aliases = append(info.aliases, alias)
		}
	}

	tree.mu.RUnlock() // Release lock

	// Second pass: Re-evaluate maxFlagLen based on combined flag names
	maxFlagLen := 0
	for name := range allFlagData { // Iterate over main flag names only
		info := allFlagData[name]

		flagStr := info.name
		if len(info.aliases) > 0 {
			// Sort aliases for consistent output
			// Sort is important here because map iteration order is not guaranteed.
			// This ensures "-a|-alpha" is always consistent.
			sortedAliases := make([]string, len(info.aliases))
			copy(sortedAliases, info.aliases)
			// Apply tree.mu.RLock() and tree.mu.RUnlock() or sort.Strings outside loop
			sort.Strings(sortedAliases) // Needs sort import

			flagStr = fmt.Sprintf("%s|-%s", strings.Join(sortedAliases, "|-"), info.name)
		}

		displayValue := info.defValue
		if displayValue == `""` || displayValue == "[]" || displayValue == "{}" {
			displayValue = ""
		}
		if displayValue != "" {
			flagStr = fmt.Sprintf("%s[=%s]", flagStr, displayValue)
		}

		if len(flagStr) > maxFlagLen {
			maxFlagLen = len(flagStr)
		}
	}

	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, "Usage of %s (powered by figtree %s):\n", filepath.Base(os.Args[0]), Version())

	sortedFlagNames := make([]string, 0, len(allFlagData))
	for name := range allFlagData {
		sortedFlagNames = append(sortedFlagNames, name)
	}
	sort.Strings(sortedFlagNames)

	for _, name := range sortedFlagNames {
		info := allFlagData[name]

		flagStr := info.name
		if len(info.aliases) > 0 {
			sortedAliases := make([]string, len(info.aliases))
			copy(sortedAliases, info.aliases)
			sort.Strings(sortedAliases)
			flagStr = fmt.Sprintf("%s|-%s", strings.Join(sortedAliases, "|-"), info.name)
		}

		displayValue := info.defValue
		if displayValue == `""` || displayValue == "[]" || displayValue == "{}" {
			displayValue = ""
		}
		if displayValue != "" {
			flagStr = fmt.Sprintf("%s[=%s]", flagStr, displayValue)
		}

		typeField := fmt.Sprintf("[%s]", info.mutagenesis) // Use Mutagenesis from figFruit

		// The f.Usage is the usage string of the main flag.
		line := fmt.Sprintf("   -%-*s   %-8s   %s", maxFlagLen, flagStr, typeField, info.usage)

		// Wrap the usage text if it exceeds terminal width
		lines := wrapText(line, termWidth, maxFlagLen+8+3+3) // 8 for type field, 3 spaces each side
		for i, l := range lines {
			if i == 0 {
				_, _ = fmt.Fprintln(&sb, l)
			} else {
				// Indent wrapped lines to align with usage start
				indent := strings.Repeat(" ", maxFlagLen+8+3+3)
				_, _ = fmt.Fprintf(&sb, "%s%s\n", indent, l)
			}
		}
	}

	return sb.String()
}

// wrapText wraps a line of text to fit within the terminal width, indenting wrapped lines
func wrapText(line string, termWidth int, indentLen int) []string {
	words := strings.Fields(line)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	currentLine := words[0]
	spaceLeft := termWidth - len(words[0])

	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			lines = append(lines, currentLine)
			currentLine = word
			spaceLeft = termWidth - indentLen - len(word)
		} else {
			currentLine += " " + word
			spaceLeft -= len(word) + 1
		}
	}
	lines = append(lines, currentLine)

	// Adjust for indentation on wrapped lines
	if len(lines) > 1 {
		for i := 1; i < len(lines); i++ {
			lines[i] = strings.TrimSpace(lines[i])
		}
	}

	return lines
}
