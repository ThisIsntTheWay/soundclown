FROM golang:1.20

RUN apt-get update && apt-get install -y \
    ffmpeg python3 python3-pip

RUN pip3 install scdl

# Application data
RUN mkdir -p /opt/app

ADD soundclown /opt/app/
ADD web /opt/app/web

# Run
WORKDIR /opt/app/
ENTRYPOINT [ "soundclown" ]