FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
RUN go get -v -d  github.com/jokimina/go-script/cmd/cron_apollo_check
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /cron_apollo_check $GOPATH/src/github.com/jokimina/go-script/cmd/cron_apollo_check

FROM scratch
COPY --from=builder /cron_apollo_check /cron_apollo_check
CMD ["/cron_apollo_check"]
