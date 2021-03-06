// Copyright 2018 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package name

import (
	"testing"
)

var goodStrictValidationRegistryNames = []string{
	"gcr.io",
	"gcr.io:9001",
	"index.docker.io",
	"us.gcr.io",
	"example.text",
	"localhost",
	"localhost:9090",
}

var goodWeakValidationRegistryNames = []string{
	"",
}

var badRegistryNames = []string{
	"white space",
	"gcr?com",
}

func TestNewRegistryStrictValidation(t *testing.T) {
	t.Parallel()

	for _, name := range goodStrictValidationRegistryNames {
		if registry, err := NewRegistry(name, StrictValidation); err != nil {
			t.Errorf("`%s` should be a valid Registry name, got error: %v", name, err)
		} else if registry.Name() != name {
			t.Errorf("`%v` .Name() should reproduce the original name. Wanted: %s Got: %s", registry, name, registry.Name())
		}
	}

	for _, name := range append(goodWeakValidationRegistryNames, badRegistryNames...) {
		if repo, err := NewRegistry(name, StrictValidation); err == nil {
			t.Errorf("`%s` should be an invalid Registry name, got Registry: %#v", name, repo)
		}
	}
}

func TestNewRegistry(t *testing.T) {
	t.Parallel()

	for _, name := range append(goodStrictValidationRegistryNames, goodWeakValidationRegistryNames...) {
		if _, err := NewRegistry(name, WeakValidation); err != nil {
			t.Errorf("`%s` should be a valid Registry name, got error: %v", name, err)
		}
	}

	for _, name := range badRegistryNames {
		if repo, err := NewRegistry(name, WeakValidation); err == nil {
			t.Errorf("`%s` should be an invalid Registry name, got Registry: %#v", name, repo)
		}
	}
}

func TestDefaultRegistryNames(t *testing.T) {
	testRegistries := []string{"docker.io", ""}

	for _, testRegistry := range testRegistries {
		registry, err := NewRegistry(testRegistry, WeakValidation)
		if err != nil {
			t.Fatalf("`%s` should be a valid Registry name, got error: %v", testRegistry, err)
		}

		actualRegistry := registry.RegistryStr()
		if actualRegistry != DefaultRegistry {
			t.Errorf("RegistryStr() was incorrect for %v. Wanted: `%s` Got: `%s`", registry, DefaultRegistry, actualRegistry)
		}
	}
}

func TestRegistryComponents(t *testing.T) {
	t.Parallel()
	testRegistry := "gcr.io"

	registry, err := NewRegistry(testRegistry, StrictValidation)
	if err != nil {
		t.Fatalf("`%s` should be a valid Registry name, got error: %v", testRegistry, err)
	}

	actualRegistry := registry.RegistryStr()
	if actualRegistry != testRegistry {
		t.Errorf("RegistryStr() was incorrect for %v. Wanted: `%s` Got: `%s`", registry, testRegistry, actualRegistry)
	}
}

func TestRegistryScopes(t *testing.T) {
	t.Parallel()
	testRegistry := "gcr.io"
	testAction := "whatever"

	expectedScope := "registry:catalog:*"

	registry, err := NewRegistry(testRegistry, StrictValidation)
	if err != nil {
		t.Fatalf("`%s` should be a valid Registry name, got error: %v", testRegistry, err)
	}

	actualScope := registry.Scope(testAction)
	if actualScope != expectedScope {
		t.Errorf("scope was incorrect for %v. Wanted: `%s` Got: `%s`", registry, expectedScope, actualScope)
	}
}
