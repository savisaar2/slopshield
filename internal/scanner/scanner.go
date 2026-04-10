package scanner

type Dependency struct {
	Name    string
	Version string
	Source  string // e.g., "package.json", "requirements.txt"
}

type Scanner interface {
	Scan(path string) ([]Dependency, error)
}
