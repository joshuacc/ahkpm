package core

type PackagesRepository interface {
	CopyPackage(dep Dependency, path string) error
	GetPackageDependencies(dep Dependency) ([]Dependency, error)
}

type packagesRepository struct{}

func NewPackagesRepository() PackagesRepository {
	return &packagesRepository{}
}

func (pr *packagesRepository) CopyPackage(dep Dependency, path string) error {
	return nil
}

func (pr *packagesRepository) GetPackageDependencies(dep Dependency) ([]Dependency, error) {
	return []Dependency{}, nil
}
