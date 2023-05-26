FROM golang:1.20

RUN apt-get update && apt-get install -y \
    ffmpeg python3 python3-pip

RUN pip3 install scdl

# NodeJS shit
ENV NODE_VERSION=18.16.0
RUN apt install -y curl
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
ENV NVM_DIR=/root/.nvm
RUN . "$NVM_DIR/nvm.sh" && nvm install ${NODE_VERSION}
RUN . "$NVM_DIR/nvm.sh" && nvm use v${NODE_VERSION}
RUN . "$NVM_DIR/nvm.sh" && nvm alias default v${NODE_VERSION}
ENV PATH="/root/.nvm/versions/node/v${NODE_VERSION}/bin/:${PATH}"
RUN node --version
RUN npm --version

RUN npm install --global yarn

# Application data
RUN mkdir -p /opt/app
ADD . /opt/app/

# Compile vue thing
WORKDIR /opt/app/sc-frontend
RUN yarn build

WORKDIR /opt/app/

ENTRYPOINT [ "soundclown" ]