package core

type Dependency interface {
	Name() string
	Version() Version
}

type dependency struct {
	name    string
	version Version
}

func NewDependency(name string, version Version) Dependency {
	return dependency{
		name:    name,
		version: version,
	}
}

func (d dependency) Name() string {
	return d.name
}

func (d dependency) Version() Version {
	return d.version
}
