FROM golang:1.12.6

LABEL maintainer="Vishal Wadhwa <vishal.wadhwa@zomato.com>"

WORKDIR $GOPATH/src/github.com/vishalwadhwa13/restaurant-service

COPY . .

ENV ENV docker

# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

RUN make build

EXPOSE 8080/tcp

# Run the executable
CMD ["out/restaurant-service"]
