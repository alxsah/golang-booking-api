FROM golang
# Create app directory
WORKDIR /usr/src/app

# Get the repo
RUN apt-get update && apt-get install -y netcat
RUN go get github.com/alxsah/golang-booking-api
WORKDIR $GOPATH/src/github.com/alxsah/golang-booking-api
# Build app
RUN go build app.go

# Expose port and run
EXPOSE 3001
CMD [ "sh", "start.sh"]