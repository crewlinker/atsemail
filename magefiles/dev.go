// Package main provides repo automation using mage.
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// init performs some sanity checks before running anything.
func init() {
	mustBeInRootIfNotTest()
}

// Dev groups commands for local development.
type Dev mg.Namespace

// Lint the codebase through static code analysis.
func (Dev) Lint() error {
	if err := sh.Run("golangci-lint", "run"); err != nil {
		return fmt.Errorf("failed to run golang-ci: %w", err)
	}

	return nil
}

// Serve will create or replace containers used for development.
func (Dev) Serve() error {
	// to make queue names unique but consistent between containers we provide an environment
	// variable that uses a millisecond timestamp.
	os.Setenv("STAMP", strconv.FormatInt(time.Now().UnixMilli(), 10))

	if err := sh.RunWith(map[string]string{}, "docker", "compose",
		"-f", "docker-compose.yml",
		"up",
		"-d", "--build", "--remove-orphans", "--force-recreate",
	); err != nil {
		return fmt.Errorf("failed to run: %w", err)
	}

	return nil
}

// Test tests all the code using Gingo, with an empty label filter.
func (Dev) Test() error {
	return (Dev{}).TestSome("!e2e")
}

// TestE2e will run the e2e tests.
func (Dev) TestE2e(env string) error {
	if err := godotenv.Load("e2e." + env + ".env"); err != nil {
		return fmt.Errorf("failed to load e2e env: %w", err)
	}

	return (Dev{}).testSome("", "./e2e")
}

// TestSome tests the whole repo using Ginkgo test runner with label filters applied.
func (Dev) TestSome(labelFilter string) error {
	if err := (Dev{}).testSome(labelFilter, "./..."); err != nil {
		return fmt.Errorf("failed to run ginkgo: %w", err)
	}

	return nil
}

// testSome the whole repo using Ginkgo test runner with label filters applied.
func (Dev) testSome(labelFilter, dir string) error {
	if err := sh.Run(
		"go", "run", "-mod=readonly", "github.com/onsi/ginkgo/v2/ginkgo",
		"-p", "-randomize-all", "--fail-on-pending", "--race", "--trace",
		"--junit-report=test-report.xml",
		"--label-filter", labelFilter,
		dir,
	); err != nil {
		return fmt.Errorf("failed to run ginkgo: %w", err)
	}

	return nil
}

// Generate generates code across the repository.
func (Dev) Generate() error {
	// run std go generators
	if err := sh.Run("go", "generate", "./..."); err != nil {
		return fmt.Errorf("failed to go generate: %w", err)
	}

	// generate protobuf in subdirs
	if err := sh.Run("buf", "generate"); err != nil {
		return fmt.Errorf("failed to generate protobuf: %w", err)
	}

	// generate mock types
	if err := sh.Run("go", "run", "-mod=readonly", "github.com/vektra/mockery/v2"); err != nil {
		return fmt.Errorf("failed to run mockery: %w", err)
	}

	return nil
}

// mustBeInRootIfNotTest checks that the command is run in the project root.
func mustBeInRootIfNotTest() {
	if _, err := os.ReadFile("go.mod"); err != nil && !strings.Contains(strings.Join(os.Args, ""), "-test.") {
		panic("must be in project root, couldn't stat go.mod file: " + err.Error())
	}
}
