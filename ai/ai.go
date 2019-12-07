package ai

import "fmt"

// Ai ...
type Ai struct {
	Name        string
	Version     int
	Family      *Family
	Description string
	FileName    string
}

// ByNameAndVersion ..
type ByNameAndVersion []*Ai

func (a ByNameAndVersion) Len() int {
	return len(a)
}

func (a ByNameAndVersion) Less(i, j int) bool {
	if a[i].Name == a[j].Name {
		return a[i].Version >= a[j].Version
	}

	return a[i].Name < a[j].Name
}

func (a ByNameAndVersion) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func actualVersion(version int, lastVersion int) int {
	if version < 0 {
		return lastVersion + version + 1
	}

	return version
}

// MatchesDescriptor checks if the ai matches a given descriptor
func (ai *Ai) MatchesDescriptor(descriptor Descriptor) bool {
	from := actualVersion(descriptor.VersionFrom, ai.Family.LastVersion.Version)
	to := actualVersion(descriptor.VersionTo, ai.Family.LastVersion.Version)

	return len(descriptor.Name) == 0 ||
		(descriptor.Name == ai.Name && from <= ai.Version && ai.Version <= to)
}

// PlayerName ...
func (ai *Ai) PlayerName() string {
	if ai.Version == 0 {
		return ai.Name
	}

	return fmt.Sprintf("%s_%d", ai.Name, ai.Version)
}

// Descriptor ...
func (ai *Ai) Descriptor() string {
	return fmt.Sprintf("%s:%d", ai.Name, ai.Version)
}
