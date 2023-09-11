# get Go version 1.18
FROM golang:1.18-alpine

# create new directory
WORKDIR /app

# copy everything in the root directory into the new app directory
COPY . .

# copy production env
COPY .env .


# download all dependencies
RUN go mod download

# build the Go application
RUN go build -o /staff-scheduling

# run the application inside docker container
CMD [ "/staff-scheduling" ]