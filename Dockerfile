FROM golang:1
WORKDIR /go/src/github.com/imagespy/imagespy/
COPY . .
RUN make

FROM gcr.io/distroless/base
COPY ui/ /ui/
COPY --from=0 /go/src/github.com/imagespy/imagespy/imagespy /imagespy
ENTRYPOINT ["/imagespy"]
