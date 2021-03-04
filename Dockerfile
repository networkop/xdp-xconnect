FROM --platform=${BUILDPLATFORM} golang:1.15.6-buster as builder

WORKDIR /src

ARG LDFLAGS

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ARG TARGETOS
ARG TARGETARCH

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "${LDFLAGS}" -o xdp-xconnect main.go

FROM alpine:latest
WORKDIR /
COPY --from=builder /src/xdp-xconnect .

ENTRYPOINT ["/xdp-xconnect"]