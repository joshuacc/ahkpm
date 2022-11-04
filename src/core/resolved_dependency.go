package core

type ResolvedDependency struct {
	Name         string        `json:"name"`
	Version      string        `json:"version"`
	SHA          string        `json:"sha"`
	InstallPath  string        `json:"installPath"`
	Dependencies DependencySet `json:"dependencies"`
}

func (rd ResolvedDependency) WithDependencies(deps DependencySet) ResolvedDependency {
	rd.Dependencies = deps
	return rd
}

func (rd *ResolvedDependency) AddDependency(name string, version Version) ResolvedDependency {
	rd.Dependencies.AddDependency(NewDependency(name, version))
	return *rd
}
