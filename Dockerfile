FROM --platform=$BUILDPLATFORM golang:1.25 AS builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

WORKDIR /build

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOARM=${TARGETVARIANT#v} \
    go build -a -installsuffix cgo -o agent .

FROM scratch
EXPOSE 8080
WORKDIR /
COPY --from=builder /build/agent /
ENTRYPOINT ["./agent"]