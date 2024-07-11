# Stage 1: Build the Go binaries
FROM golang:1.22 AS builder

RUN apt-get update && DEBIAN_FRONTEND=nointeractive apt-get install --no-install-recommends --assume-yes protobuf-compiler

WORKDIR /usr/src/miku-notes-auth
COPY . .

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN mkdir -p generated
RUN protoc -I proto proto/sso.proto --go_out=./generated/ --go_opt=paths=source_relative --go-grpc_out=./generated/ --go-grpc_opt=paths=source_relative

RUN go mod download && go mod verify

# Build the main app
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/app ./cmd/sso/main.go

# Build the migrations binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/migrate ./cmd/migrations/main.go

# Stage 2: Run the migrations and start the app
FROM golang:1.22

RUN apt-get update && DEBIAN_FRONTEND=nointeractive apt-get install --no-install-recommends --assume-yes protobuf-compiler

WORKDIR /usr/src/miku-notes-auth
COPY --from=builder /usr/local/bin/app /usr/local/bin/app
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate
COPY . .

# Run migrations and start the main application
CMD /bin/sh -c "/usr/local/bin/migrate && /usr/local/bin/app"