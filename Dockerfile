FROM golang:1.20

RUN apt-get update && apt-get install -y \
    ffmpeg python3 python3-pip

RUN pip3 install scdl

# Application data
RUN mkdir -p /opt/app
ADD web /opt/app/web

# Create SCDL config
RUN mkdir -p /.config

# Make dirs writable for OCP
RUN chgrp -R 0 /opt/app && \
    chmod -R g+rwX /opt/app
RUN chgrp -R 0 /.config && \
    chmod -R g+rwX /.config

# Compile
COPY go.mod go.sum main.go /opt/app/
WORKDIR /opt/app
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

# Cleanup
RUN rm go.mod main.go

# Run
WORKDIR /opt/app/
CMD [ "./soundclown" ]