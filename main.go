// A generated module for GhcrBadge functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/ghcr-badge/internal/dagger"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GhcrBadge struct{}

// Returns a container with the Python build environment
func (m *GhcrBadge) BuildEnv(tag string) *dagger.Container {
	return dag.Container().
		From("python:3.11-slim").
		WithSymlink("/usr/local/bin/python3", "/usr/bin/python3").
		WithExec([]string{"apt", "update"}).
		WithExec([]string{"apt", "install", "-y", "git"}).
		WithExec([]string{"/usr/bin/python3", "-m", "venv", "/opt/venv"}).
		WithExec([]string{"/opt/venv/bin/pip", "install", "git+https://github.com/eggplants/ghcr-badge@" + tag})
}

// Builds the ghcr-badge image and pushes it to the registry
func (m *GhcrBadge) BuildAndPublish(ctx context.Context, ghUser string, ghAuthToken *dagger.Secret) error {

	latestTag, err := getGhcrBadgeLatestVersion()
	if err != nil {
		return err
	}

	buildEnv := m.BuildEnv(latestTag)

	container := dag.Container().
		From("gcr.io/distroless/python3-debian12:nonroot").
		WithLabel("org.opencontainers.image.version", latestTag).
		WithLabel("org.opencontainers.image.source", "https://github.com/kerwood/ghcr-badge-image-builder").
		WithLabel("org.opencontainers.image.created", time.Now().String()).
		WithDirectory("/opt/venv", buildEnv.Directory("/opt/venv")).
		WithEnvVariable("PATH", "/opt/venv/bin", dagger.ContainerWithEnvVariableOpts{Expand: true}).
		WithEntrypoint([]string{"ghcr-badge-server"}).
		WithRegistryAuth("ghcr.io", ghUser, ghAuthToken)

	for _, tag := range []string{"latest", latestTag} {
		_, err := container.Publish(ctx, "ghcr.io/kerwood/ghcr-badge:"+tag)
		if err != nil {
			return err
		}
	}

	return nil
}

type Tag struct {
	Name string `json:"name"`
}

func getGhcrBadgeLatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/eggplants/ghcr-badge/tags")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var tags []Tag

	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "", fmt.Errorf("No tags found on ghcr-badge repository")
	}

	return tags[0].Name, nil
}
