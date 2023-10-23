# Build web
FROM docker.io/library/node:18 as node-builder
WORKDIR /build

COPY web .

RUN npm install
RUN npm run build

# Build
FROM docker.io/library/golang:1.21 as go-builder
WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
COPY --from=node-builder /build/dist web/dist
RUN \
CGO_ENABLED=0 \ 
GOOS=$TARGETOS \
GOARCH=$TARGETARCH \
go build -ldflags="-s -w" ./cmd/radiomux-demo

# Runner
FROM docker.io/library/alpine:3.18

COPY --from=go-builder /build/radiomux-demo /usr/local/bin/

ENTRYPOINT [ "radiomux-demo" ]
EXPOSE 8080
