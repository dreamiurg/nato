package alphabet

import (
	"fmt"
	"sort"
	"strings"
)

var registry = make(map[string]*Alphabet)

// Register adds an alphabet to the registry.
func Register(a *Alphabet) {
	registry[strings.ToLower(a.Name)] = a
}

// Get returns an alphabet by name (case-insensitive).
func Get(name string) (*Alphabet, error) {
	a, ok := registry[strings.ToLower(name)]
	if !ok {
		return nil, fmt.Errorf("unknown alphabet %q", name)
	}
	return a, nil
}

// List returns all registered alphabets sorted by name.
func List() []*Alphabet {
	alphabets := make([]*Alphabet, 0, len(registry))
	for _, a := range registry {
		alphabets = append(alphabets, a)
	}
	sort.Slice(alphabets, func(i, j int) bool {
		return alphabets[i].Name < alphabets[j].Name
	})
	return alphabets
}

// Names returns all registered alphabet names sorted.
func Names() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
