package core

type DependencyArray []Dependency

func (deps DependencyArray) Equals(other DependencyArray) bool {
	if len(deps) != len(other) {
		return false
	}

	depMap := make(map[string]Dependency)
	for _, dep := range deps {
		depMap[dep.Name()] = dep
	}

	for _, otherDep := range other {
		dep, ok := depMap[otherDep.Name()]
		if !ok || !dep.Equals(otherDep) {
			return false
		}
	}

	return true
}
