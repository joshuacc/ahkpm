package mocks

import (
	. "ahkpm/src/core"

	"github.com/stretchr/testify/mock"
)

// TODO: Write a data-backed mock to simplify test setup
type MockPackagesRepository struct {
	mock.Mock
}

func (m *MockPackagesRepository) CopyPackage(dep ResolvedDependency, path string) error {
	args := m.Called(dep, path)
	return args.Error(1)
}

func (m *MockPackagesRepository) GetPackageDependencies(dep ResolvedDependency) (*DependencySet, error) {
	args := m.Called(dep)
	return args.Get(0).(*DependencySet), args.Error(1)
}

func (m *MockPackagesRepository) GetResolvedDependencySHA(dep Dependency) (string, error) {
	args := m.Called(dep)
	return args.String(0), args.Error(1)
}

func (m *MockPackagesRepository) ClearCache() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPackagesRepository) WithRemoveAll(removeAll func(path string) error) PackagesRepository {
	m.On("WithRemoveAll", removeAll).Return(m)
	return m
}

func (m *MockPackagesRepository) GetLatestVersion(depName string) (Version, error) {
	args := m.Called(depName)
	return args.Get(0).(Version), args.Error(1)
}
