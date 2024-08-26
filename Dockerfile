FROM --platform=$BUILDPLATFORM golang:1.22.5 AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG VERSION
ARG SHA

# Get the dependencies downloaded
WORKDIR /app
ENV CGO_ENABLED=0
COPY go.* ./
RUN go mod download
COPY . ./

# Build Go CLI binary
RUN go build -ldflags="-s -w \
    -X 'github.com/open-sauced/pizza-cli/pkg/utils.Version=${VERSION}' \
    -X 'github.com/open-sauced/pizza-cli/pkg/utils.Sha=${SHA}'" \
    -o pizza .

# Runner layer
FROM --platform=$BUILDPLATFORM golang:alpine
COPY --from=builder /app/pizza /usr/bin/
ENTRYPOINT ["/usr/bin/pizza"]
