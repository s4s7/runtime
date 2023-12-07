// Copyright 2015 go-swagger maintainers
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

package client

import (
	"net/http"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicAuth(t *testing.T) {
	r, err := newRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	writer := BasicAuth("someone", "with a password")
	err = writer.AuthenticateRequest(r, nil)
	require.NoError(t, err)

	req := new(http.Request)
	req.Header = make(http.Header)
	req.Header.Set(runtime.HeaderAuthorization, r.header.Get(runtime.HeaderAuthorization))
	usr, pw, ok := req.BasicAuth()
	require.True(t, ok)
	assert.Equal(t, "someone", usr)
	assert.Equal(t, "with a password", pw)
}

func TestAPIKeyAuth_Query(t *testing.T) {
	r, err := newRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	writer := APIKeyAuth("api_key", "query", "the-shared-key")
	err = writer.AuthenticateRequest(r, nil)
	require.NoError(t, err)

	assert.Equal(t, "the-shared-key", r.query.Get("api_key"))
}

func TestAPIKeyAuth_Header(t *testing.T) {
	r, err := newRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	writer := APIKeyAuth("x-api-token", "header", "the-shared-key")
	err = writer.AuthenticateRequest(r, nil)
	require.NoError(t, err)

	assert.Equal(t, "the-shared-key", r.header.Get("x-api-token"))
}

func TestBearerTokenAuth(t *testing.T) {
	r, err := newRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	writer := BearerToken("the-shared-token")
	err = writer.AuthenticateRequest(r, nil)
	require.NoError(t, err)

	assert.Equal(t, "Bearer the-shared-token", r.header.Get(runtime.HeaderAuthorization))
}

func TestCompose(t *testing.T) {
	r, err := newRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	writer := Compose(APIKeyAuth("x-api-key", "header", "the-api-key"), APIKeyAuth("x-secret-key", "header", "the-secret-key"))
	err = writer.AuthenticateRequest(r, nil)
	require.NoError(t, err)

	assert.Equal(t, "the-api-key", r.header.Get("x-api-key"))
	assert.Equal(t, "the-secret-key", r.header.Get("x-secret-key"))
}
