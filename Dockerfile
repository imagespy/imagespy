FROM golang:1
WORKDIR /go/src/github.com/imagespy/imagespy/
COPY . .
RUN make

FROM gcr.io/distroless/base
COPY --from=0 /go/src/github.com/imagespy/imagespy/imagespy /imagespy
USER nobody
ENTRYPOINT ["/imagespy"]
