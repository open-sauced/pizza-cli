# Builds the go binary into the git ignored ./build/ dir
build:
  #!/usr/bin/env sh
  echo "Building for local arch"

  export VERSION="${RELEASE_TAG_VERSION:-dev}"
  export DATETIME=$(date -u +"%Y-%m-%dT%H:%M:%S")

  go build \
    -ldflags="-s -w" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Version=${VERSION}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Sha=$(git rev-parse HEAD)'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Sha=$(git rev-parse HEAD)'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Datetime=${DATETIME}'" \
    -o build/pizza

install: build
  sudo cp "./build/pizza" "/usr/local/bin/"

build-all: \
    build \
    build-darwin-amd64 build-darwin-arm64 \
    build-linux-amd64 build-linux-arm64 \
    build-windows-amd64 build-windows-arm64

build-darwin-amd64:
  #!/usr/bin/env sh

  echo "Building darwin amd64"

  export VERSION="${RELEASE_TAG_VERSION:-dev}"
  export CGO_ENABLED=0
  export GOOS="darwin"
  export GOARCH="amd64"

  go build \
    -ldflags="-s -w" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Version=${VERSION}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Sha=$(git rev-parse HEAD)'" \
    -o build/pizza-${GOOS}-${GOARCH}

build-darwin-arm64:
  #!/usr/bin/env sh

  echo "Building darwin arm64"

  export VERSION="${RELEASE_TAG_VERSION:-dev}"
  export CGO_ENABLED=0
  export GOOS="darwin"
  export GOARCH="arm64"

  go build \
    -ldflags="-s -w" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Version=${VERSION}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Sha=$(git rev-parse HEAD)'" \
    -o build/pizza-${GOOS}-${GOARCH}

build-linux-amd64:
  #!/usr/bin/env sh

  echo "Building linux amd64"

  export VERSION="${RELEASE_TAG_VERSION:-dev}"
  export CGO_ENABLED=0
  export GOOS="linux"
  export GOARCH="amd64"

  go build \
    -ldflags="-s -w" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Version=${VERSION}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Sha=$(git rev-parse HEAD)'" \
    -o build/pizza-${GOOS}-${GOARCH}

build-linux-arm64:
  #!/usr/bin/env sh

  echo "Building linux arm64"

  export VERSION="${RELEASE_TAG_VERSION:-dev}"
  export CGO_ENABLED=0
  export GOOS="linux"
  export GOARCH="arm64"

  go build \
    -ldflags="-s -w" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Version=${VERSION}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Sha=$(git rev-parse HEAD)'" \
    -o build/pizza-${GOOS}-${GOARCH}

build-windows-amd64:
  #!/usr/bin/env sh

  echo "Building windows amd64"

  export VERSION="${RELEASE_TAG_VERSION:-dev}"
  export CGO_ENABLED=0
  export GOOS="windows"
  export GOARCH="amd64"

  go build \
    -ldflags="-s -w" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Version=${VERSION}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Sha=$(git rev-parse HEAD)'" \
    -o build/pizza-${GOOS}-${GOARCH}

build-windows-arm64:
  #!/usr/bin/env sh

  echo "Building windows arm64"

  export VERSION="${RELEASE_TAG_VERSION:-dev}"
  export CGO_ENABLED=0
  export GOOS="windows"
  export GOARCH="arm64"

  go build \
    -ldflags="-s -w" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Version=${VERSION}'" \
    -ldflags="-X 'github.com/open-sauced/pizza-cli/pkg/utils.Sha=$(git rev-parse HEAD)'" \
    -o build/pizza-${GOOS}-${GOARCH}

# Builds the container and marks it tagged as "dev" locally
build-container:
  docker build \
    --build-arg VERSION=$(git describe --tags --always) \
    --build-arg SHA=$(git rev-parse HEAD) \
    -t pizza:dev .

clean:
  rm -rf build/

# Runs all tests
test: unit-test

# Runs all in-code, unit tests
unit-test:
  go test ./...

# Lints the Go code via golangcilint in Docker
lint:
  docker run \
    -t \
    --rm \
    -v "$(pwd)/:/app" \
    -w /app \
    golangci/golangci-lint:v1.59 \
    golangci-lint run -v

# Formats code via builtin go fmt
format:
  find . -type f -name "*.go" -exec goimports -local github.com/open-sauced/pizza-cli -w {} \;

run:
  go run main.go
