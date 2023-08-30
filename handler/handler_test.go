package handler

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleRequest(t *testing.T) {
	ctx := context.Background()

	// Test for the default environment
	resp, err := HandleRequest(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "Hello, this is the default environment.", resp.Body)

	// Test for the dev environment
	os.Setenv("APP_ENV", "dev")
	resp, err = HandleRequest(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "Hello, this is the dev environment.", resp.Body)

	// Test for the stg environment
	os.Setenv("APP_ENV", "stg")
	resp, err = HandleRequest(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "Hello, this is the stg environment.", resp.Body)

	// Test for the prod environment
	os.Setenv("APP_ENV", "prod")
	resp, err = HandleRequest(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "Hello, this is the prod environment.", resp.Body)
}