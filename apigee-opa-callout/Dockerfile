FROM golang:latest as builder
ADD ./api /go/src/apigee-opa-callout/api
ADD ./cmd/server /go/src/apigee-opa-callout/cmd/server
ADD ./pkg /go/src/apigee-opa-callout/pkg
ADD go.mod /go/src/apigee-opa-callout
ADD go.sum /go/src/apigee-opa-callout

WORKDIR /go/src/apigee-opa-callout

RUN groupadd -r -g 20000 app && useradd -M -u 20001 -g 0 -r -c "Default app user" app && chown -R 20001:0 /go
ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -a -ldflags='-s -w -extldflags "-static"' -o /go/bin/apigee-opa-callout /go/src/apigee-opa-callout/cmd/server/main.go

#without these certificates, we cannot verify the JWT token
FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
WORKDIR /
COPY --from=builder /go/bin/apigee-opa-callout .
COPY --from=builder /etc/passwd /etc/group /etc/shadow /etc/
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
USER 20001
EXPOSE 50051
ENTRYPOINT ["/apigee-opa-callout"]