# Build the manager binary
FROM golang:1.19 as builder

ARG TARGETOS
ARG TARGETARCH
ARG GIT_HEAD_COMMIT
ARG GIT_TAG_COMMIT
ARG GIT_LAST_TAG
ARG GIT_MODIFIED
ARG GIT_REPO
ARG BUILD_DATE

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY version.go version.go
COPY api/ api/
COPY controllers/ controllers/
COPY pkg/ pkg/


# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH GO111MODULE=on go build \
        -gcflags "-N -l" \
        -ldflags "-X main.GitRepo=$GIT_REPO -X main.GitTag=$GIT_LAST_TAG -X main.GitCommit=$GIT_HEAD_COMMIT -X main.GitDirty=$GIT_MODIFIED -X main.BuildTime=$BUILD_DATE" \
        -o manager

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
