package ai

// Family ...
type Family struct {
	Name        string
	LastVersion *Ai
	Ais         []*Ai
}

// Add ...
func (f *Family) Add(ais ...*Ai) {
	if len(ais) > 0 {
		f.Name = ais[0].Name
	}

	for _, ai := range ais {
		ai.Family = f

		f.Ais = append(f.Ais, ai)

		if f.LastVersion == nil || ai.Version > f.LastVersion.Version {
			f.LastVersion = ai
		}
	}
}
