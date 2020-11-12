###############
# BUILDER     #
###############
FROM golang:1.13.14-alpine3.12 AS build_base

# File Author / Maintainer
LABEL maintainer="Lucas Costa <lucas@itguru.com.br>"

# Force the go compiler to use modules
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://proxy.golang.org

# Set Workdir
WORKDIR /build

# Copy Go-Modules files to download and cache dependencies 
COPY go.mod .
COPY go.sum .
# Download and cache dependencies
RUN go mod download

# Copy the project container
COPY . .

# Build with static flags
RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .

############### Build final layer ###############
FROM scratch
# Copy certificate from the main image. This is important to make http requests
COPY --from=build_base /dist/main /
COPY --from=build_base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT [ "/main" ]