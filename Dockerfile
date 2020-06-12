# Start by building the application.
FROM golang:1.14-buster as build

WORKDIR /go/src/greeter
ADD . /go/src/greeter

RUN go get -d -v ./...

RUN go build -o /go/bin/greeter

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/greeter /
CMD ["/greeter"]
