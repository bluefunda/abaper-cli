FROM golang:1.25-alpine AS builder

ARG GH_PAT
RUN if [ -n "$GH_PAT" ]; then \
      git config --global url."https://${GH_PAT}@github.com/".insteadOf "https://github.com/"; \
    fi

WORKDIR /app
COPY go.mod go.sum ./
ENV GOPRIVATE=github.com/bluefunda/*
RUN go mod download

COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 go build \
    -ldflags "-X github.com/bluefunda/abaper-cli/internal/commands.version=${VERSION}" \
    -o /abaper ./cmd/abaper

FROM alpine:3.21
RUN apk --no-cache add ca-certificates
COPY --from=builder /abaper /usr/local/bin/abaper
ENTRYPOINT ["abaper"]
