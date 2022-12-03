/*
Copyright 2020 The Flux authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package notifier

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGiteaBasic(t *testing.T) {
	g, err := NewGitea("https://try.gitea.io/foo/bar", "foobar", nil)
	assert.Nil(t, err)
	assert.Equal(t, g.Owner, "foo")
	assert.Equal(t, g.Repo, "bar")
	assert.Equal(t, g.BaseURL, "https://try.gitea.io")
}

func TestNewGiteaInvalidUrl(t *testing.T) {
	_, err := NewGitea("https://try.gitea.io/foo/bar/baz", "foobar", nil)
	assert.NotNil(t, err)
}

func TestNewGiteaEmptyToken(t *testing.T) {
	_, err := NewGitea("https://try.gitea.io/foo/bar", "", nil)
	assert.NotNil(t, err)
}
