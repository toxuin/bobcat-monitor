FROM golang:1.17-alpine AS build_base
RUN apk add ca-certificates

WORKDIR /tmp/app

COPY . .
RUN go get -d ./... && CGO_ENABLED=0 go build -ldflags="-w -s" -o ./out/bobcatmonitor

FROM scratch
COPY --from=build_base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --from=build_base /tmp/app/out/bobcatmonitor /bobcatmonitor

ENTRYPOINT ["/bobcatmonitor"]