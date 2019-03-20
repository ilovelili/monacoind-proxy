FROM golang:1.11 as builder

LABEL maintainer="<m_ju@indiesquare.me>"

ENV SRC_DIR=/go/src/monacoind-proxy/
WORKDIR $SRC_DIR

COPY . $SRC_DIR

# Go dep
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

WORKDIR $SRC_DIR
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o proxy .

FROM alpine

LABEL maintainer="<m_ju@indiesquare.me>"

ENV BUILDER_DIR=/go/src/monacoind-proxy
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder $BUILDER_DIR/config.json $BUILDER_DIR/proxy /root/