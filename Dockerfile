# get a base image
FROM golang

# set the working directory at the container
WORKDIR /go/src/app

# copy the files from host to the container working directory
COPY ./app ./

# install dependencies
# RUN apk add git

# downlod all the dependecies listed in the go.mod
RUN go get -d -v
RUN go mod tidy
# RUN go mod download && go mod verify

# build the project into a binary
RUN go build -v -o app

# run the binary after build
CMD ["./app"]
