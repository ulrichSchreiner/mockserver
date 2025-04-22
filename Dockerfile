FROM golang:1 as builder

WORKDIR /mockserver
COPY . .
RUN CGO_ENABLED=0 go build

FROM registry.gitlab.com/ulrichschreiner/cacerts:latest
#FROM alpine
COPY --from=builder /mockserver/mockserver /mockserver
EXPOSE 9099
ENTRYPOINT [ "/mockserver" ]
