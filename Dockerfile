FROM golang:1.20

RUN apt-get update && apt-get install -y \
    ffmpeg python3 python3-pip

RUN pip3 install scdl

# Application data
RUN mkdir -p /opt/app
ADD web /opt/app/web

# Compile
COPY go.mod main.go /opt/app/
WORKDIR /opt/app
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

# Run
WORKDIR /opt/app/
CMD [ "./soundclown" ]