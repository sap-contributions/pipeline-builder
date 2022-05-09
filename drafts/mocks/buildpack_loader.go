// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	drafts "github.com/paketo-buildpacks/pipeline-builder/drafts"
	mock "github.com/stretchr/testify/mock"
)

// BuildpackLoader is an autogenerated mock type for the BuildpackLoader type
type BuildpackLoader struct {
	mock.Mock
}

// LoadBuildpack provides a mock function with given fields: id
func (_m *BuildpackLoader) LoadBuildpack(id string) (drafts.Buildpack, error) {
	ret := _m.Called(id)

	var r0 drafts.Buildpack
	if rf, ok := ret.Get(0).(func(string) drafts.Buildpack); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(drafts.Buildpack)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadBuildpacks provides a mock function with given fields: uris
func (_m *BuildpackLoader) LoadBuildpacks(uris []string) ([]drafts.Buildpack, error) {
	ret := _m.Called(uris)

	var r0 []drafts.Buildpack
	if rf, ok := ret.Get(0).(func([]string) []drafts.Buildpack); ok {
		r0 = rf(uris)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]drafts.Buildpack)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(uris)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}