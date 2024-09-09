set dotenv-load

# Displays this help message
help:
  @echo "Available commands:"
  @just --list

# Builds the go binary into the git ignored ./build/ dir for the local architecture
build:
  #!/usr/bin/env sh
  echo "Building for local arch"

  export VERSION="${RELEASE_TAG_VERSION:-dev}"
  export DATETIME=$(date -u +"%Y-%m-%d %H:%M:%S")

  go build \
    -ldflags="-s -w" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Version=${VERSION}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Sha=$(git rev-parse HEAD)'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Datetime=${DATETIME}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.writeOnlyPublicPosthogKey=${POSTHOG_PUBLIC_API_KEY}'" \
    -o build/pizza

# Builds and installs the go binary for the local architecture. WARNING: requires sudo access
install: build
  sudo cp "./build/pizza" "/usr/local/bin/"

# Builds all build targets arcross all OS and architectures
build-all: \
    build \
    build-container \
    build-darwin-amd64 build-darwin-arm64 \
    build-linux-amd64 build-linux-arm64 \
    build-windows-amd64 build-windows-arm64

# Builds for Darwin linux (i.e., MacOS) on amd64 architecture
build-darwin-amd64:
  #!/usr/bin/env sh

  echo "Building darwin amd64"

  export VERSION="${RELEASE_TAG_VERSION:-dev}"
  export DATETIME=$(date -u +"%Y-%m-%d %H:%M:%S")
  export CGO_ENABLED=0
  export GOOS="darwin"
  export GOARCH="amd64"

  go build \
    -ldflags="-s -w" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Version=${VERSION}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Sha=$(git rev-parse HEAD)'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Datetime=${DATETIME}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.writeOnlyPublicPosthogKey=${POSTHOG_PUBLIC_API_KEY}'" \
    -o build/pizza-${GOOS}-${GOARCH}

# Builds for Darwin linux (i.e., MacOS) on arm64 architecture (i.e. Apple silicon)
build-darwin-arm64:
  #!/usr/bin/env sh

  echo "Building darwin arm64"

  export VERSION="${RELEASE_TAG_VERSION:-dev}"
  export DATETIME=$(date -u +"%Y-%m-%d %H:%M:%S")
  export CGO_ENABLED=0
  export GOOS="darwin"
  export GOARCH="arm64"

  go build \
    -ldflags="-s -w" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Version=${VERSION}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Sha=$(git rev-parse HEAD)'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Datetime=${DATETIME}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.writeOnlyPublicPosthogKey=${POSTHOG_PUBLIC_API_KEY}'" \
    -o build/pizza-${GOOS}-${GOARCH}

# Builds for agnostic Linux on amd64 architecture
build-linux-amd64:
  #!/usr/bin/env sh

  echo "Building linux amd64"

  export VERSION="${RELEASE_TAG_VERSION:-dev}"
  export DATETIME=$(date -u +"%Y-%m-%d %H:%M:%S")
  export CGO_ENABLED=0
  export GOOS="linux"
  export GOARCH="amd64"

  go build \
    -ldflags="-s -w" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Version=${VERSION}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Sha=$(git rev-parse HEAD)'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Datetime=${DATETIME}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.writeOnlyPublicPosthogKey=${POSTHOG_PUBLIC_API_KEY}'" \
    -o build/pizza-${GOOS}-${GOARCH}

# Builds for agnostic Linux on arm64 architecture
build-linux-arm64:
  #!/usr/bin/env sh

  echo "Building linux arm64"

  export VERSION="${RELEASE_TAG_VERSION:-dev}"
  export DATETIME=$(date -u +"%Y-%m-%d %H:%M:%S")
  export CGO_ENABLED=0
  export GOOS="linux"
  export GOARCH="arm64"

  go build \
    -ldflags="-s -w" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Version=${VERSION}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Sha=$(git rev-parse HEAD)'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Datetime=${DATETIME}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.writeOnlyPublicPosthogKey=${POSTHOG_PUBLIC_API_KEY}'" \
    -o build/pizza-${GOOS}-${GOARCH}

# Builds for Windows on amd64 architecture
build-windows-amd64:
  #!/usr/bin/env sh

  echo "Building windows amd64"

  export VERSION="${RELEASE_TAG_VERSION:-dev}"
  export DATETIME=$(date -u +"%Y-%m-%d %H:%M:%S")
  export CGO_ENABLED=0
  export GOOS="windows"
  export GOARCH="amd64"

  go build \
    -ldflags="-s -w" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Version=${VERSION}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Sha=$(git rev-parse HEAD)'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Datetime=${DATETIME}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.writeOnlyPublicPosthogKey=${POSTHOG_PUBLIC_API_KEY}'" \
    -o build/pizza-${GOOS}-${GOARCH}

# Builds for Windows on arm64 architecture
build-windows-arm64:
  #!/usr/bin/env sh

  echo "Building windows arm64"

  export VERSION="${RELEASE_TAG_VERSION:-dev}"
  export DATETIME=$(date -u +"%Y-%m-%d %H:%M:%S")
  export CGO_ENABLED=0
  export GOOS="windows"
  export GOARCH="arm64"

  go build \
    -ldflags="-s -w" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Version=${VERSION}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Sha=$(git rev-parse HEAD)'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Datetime=${DATETIME}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.writeOnlyPublicPosthogKey=${POSTHOG_PUBLIC_API_KEY}'" \
    -o build/pizza-${GOOS}-${GOARCH}

# Builds the Docker container and tags it as "dev"
build-container:
  docker build \
    --build-arg VERSION="$(git describe --tags --always)" \
    --build-arg SHA="$(git rev-parse HEAD)" \
    --build-arg DATETIME="$(date -u +'%Y-%m-%d %H:%M:%S')" \
    --build-arg POSTHOG_PUBLIC_API_KEY="${POSTHOG_PUBLIC_API_KEY}" \
    -t pizza:dev .

# Removes build artifacts
clean:
  rm -rf build/

# Runs all tests
test: unit-test

# Runs all in-code, unit tests
unit-test:
  go test ./...

# Lints Go code via golangci-lint within Docker
lint:
  docker run \
    -t \
    --rm \
    -v "$(pwd)/:/app" \
    -w /app \
    golangci/golangci-lint:v1.60 \
    golangci-lint run -v

# Formats Go code via goimports
format:
  find . -type f -name "*.go" -exec goimports -local github.com/open-sauced/pizza-cli -w {} \;

# Installs the dev tools for working with this project. Requires "go", "just", and "docker"
install-dev-tools:
  #!/usr/bin/env sh

  go install golang.org/x/tools/cmd/goimports@latest

# Runs Go code manually through the main.go
run:
  go run main.go

# Re-generates the docs from the cobra command tree
gen-docs:
  go run main.go docs ./docs/

# Runs all the dev tasks (like formatting, linting, building, etc.)
dev: format lint test build-all

# Calls the various Posthog capture events to add the Insights to the database
bootstrap-telemetry:
  #!/usr/bin/env sh
  echo "Building telemetry-oneshot"

  go build \
    -tags telemetry \
    -ldflags="-s -w" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.writeOnlyPublicPosthogKey=${POSTHOG_PUBLIC_API_KEY}'" \
    -o build/telemetry-oneshot \
    telemetry.go

  ./build/telemetry-oneshot

  rm ./build/telemetry-oneshot
