FROM golang as builder

COPY . /app

WORKDIR /app

RUN go build -o /tmp/snapshot 

FROM alpine

RUN apk --no-cache add ca-certificates

COPY --from=builder /tmp/snapshot /

WORKDIR /

ENTRYPOINT ["/snapshot"]
