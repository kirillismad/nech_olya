FROM golang:1.26.4-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN GOWORK=off go mod download

COPY main.go ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/

RUN GOWORK=off CGO_ENABLED=0 go build -mod=readonly -trimpath -ldflags="-s -w" -o /out/notes .

FROM alpine:3.23

RUN addgroup -S -g 10001 notes \
	&& adduser -S -D -H -u 10001 -G notes notes

COPY --from=builder /out/notes /usr/local/bin/notes

USER notes:notes

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/notes"]
CMD ["serve"]
