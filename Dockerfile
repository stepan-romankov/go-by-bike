FROM golang:alpine as builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    GOOS=linux \
    CGO_ENABLED=0

WORKDIR /build/

COPY go.mod .
COPY go.sum .
# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

COPY . .
RUN go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o bike .

FROM alpine
RUN apk update && apk add netcat-openbsd

RUN adduser -S -D -H -h /app appuser
USER appuser

COPY --from=builder /build/bike /app/
COPY db/migrations/ /app/db/migrations/
WORKDIR /app