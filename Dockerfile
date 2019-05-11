FROM golang:1 as builder

WORKDIR /mockserver
COPY . .
RUN go build

FROM ulrichschreiner/cacerts
COPY --from=builder /mockserver/mockserver /mockserver
EXPOSE 9099
ENTRYPOINT [ "/mockserver" ]
