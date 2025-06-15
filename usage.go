package figtree

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

// Usage prints a helpful table of figs in a human-readable format
func (tree *figTree) Usage() string {
	termWidth := 80
	if term.IsTerminal(int(os.Stdout.Fd())) {
		if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			termWidth = width
		}
	}

	maxFlagLen := 0
	flag.VisitAll(func(f *flag.Flag) {
		flagStr := f.Name
		// Handle special cases for default values
		defValue := f.DefValue
		if defValue == `""` || defValue == "[]" || defValue == "{}" {
			defValue = ""
		}
		if defValue != "" {
			flagStr = fmt.Sprintf("%s[=%s]", f.Name, defValue)
		}
		if len(flagStr) > maxFlagLen {
			maxFlagLen = len(flagStr)
		}
	})

	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, "Usage of %s (powered by figree %s):\n", filepath.Base(os.Args[0]), Version())
	flag.VisitAll(func(f *flag.Flag) {
		flagStr := f.Name
		defValue := f.DefValue
		if defValue == `""` || defValue == "[]" || defValue == "{}" {
			defValue = ""
		}
		if defValue != "" {
			flagStr = fmt.Sprintf("%s[=%s]", f.Name, defValue)
		}

		typeStr := "Unknown"
		tree.mu.RLock()
		if fig, ok := tree.figs[f.Name]; ok && fig != nil {
			typeStr = string(fig.Mutagenesis)
		}
		tree.mu.RUnlock()
		typeField := fmt.Sprintf("[%s]", typeStr)
		var aliasList []string
		for alias, name := range tree.aliases {
			if name == f.Name {
				aliasList = append(aliasList, alias)
			}
		}
		if len(aliasList) > 0 {
			flagStr = fmt.Sprintf("%s|-%s", strings.Join(aliasList, "|-"), flagStr)
		}

		line := fmt.Sprintf("   -%-*s   %-8s   %s", maxFlagLen, flagStr, typeField, f.Usage)

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
	})

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
