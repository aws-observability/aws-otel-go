FROM golang:1.17-alpine AS build-env

RUN apk update && apk add ca-certificates

WORKDIR /usr/src/app

COPY . .

RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"'

FROM scratch
COPY --from=build-env /usr/src/app/sampleapp /sampleapp
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/sampleapp"]
