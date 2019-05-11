FROM golang:1 as builder

WORKDIR /mockserver
COPY . .
RUN CGO_ENABLED=0 go build

FROM ulrichschreiner/cacerts
#FROM alpine
COPY --from=builder /mockserver/mockserver /mockserver
EXPOSE 9099
ENTRYPOINT [ "/mockserver" ]
