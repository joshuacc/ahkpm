package mocks

import (
	. "ahkpm/src/core"

	"github.com/stretchr/testify/mock"
)

type MockPackagesRepository struct {
	mock.Mock
}

func (m *MockPackagesRepository) CopyPackage(dep Dependency, path string) error {
	args := m.Called(dep, path)
	return args.Error(1)
}

func (m *MockPackagesRepository) GetPackageDependencies(dep Dependency) ([]Dependency, error) {
	args := m.Called(dep)
	return args.Get(0).([]Dependency), args.Error(1)
}
