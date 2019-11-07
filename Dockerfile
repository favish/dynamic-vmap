FROM golang as builder

# Create the necessary go directory for our source and copy it there
WORKDIR /go/src/github.com/favish/dynamic-vmap
COPY . .

# Install dep and the deps our project needs
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure

# Build the binary including OS symbols
WORKDIR /
RUN CGO_ENABLED=0 GOOS=linux go build -a github.com/favish/dynamic-vmap

# Copy the result over to a scratch image for smallest image size possible
FROM scratch
COPY --from=builder /dynamic-vmap /dynamic-vmap

ENTRYPOINT ["/dynamic-vmap"]
