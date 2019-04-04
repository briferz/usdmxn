FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /project
COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -ldflags="-w -s" -o /go/bin/app ./cmd/exchangeapi/

FROM scratch

COPY --from=builder /go/bin/app /go/bin/app
ENTRYPOINT ["/go/bin/app"]