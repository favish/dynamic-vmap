FROM golang as builder

WORKDIR /workspace

RUN go get github.com/favish/dynamic-vmap
RUN CGO_ENABLED=0 GOOS=linux go build -a github.com/favish/dynamic-vmap

FROM scratch
COPY --from=builder /workspace/dynamic-vmap /dynamic-vmap

ENTRYPOINT ["/dynamic-vmap"]
