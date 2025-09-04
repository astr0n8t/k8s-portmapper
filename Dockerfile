# Build Stage
ARG BUILDPLATFORM
FROM --platform=${BUILDPLATFORM} golang:1.25.1 AS build-stage

LABEL app="k8s-portmapper"
LABEL REPO="https://github.com/astr0n8t/k8s-portmapper"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
# Copy all internal modules
COPY cmd/ ./cmd/
COPY pkg/ ./pkg/
COPY internal/ ./internal/
COPY version/ ./version/

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /k8s-portmapper

# Deploy the application binary into a lean image
FROM gcr.io/distroless/static-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /k8s-portmapper /k8s-portmapper

USER nonroot:nonroot

ENTRYPOINT ["/k8s-portmapper"]
