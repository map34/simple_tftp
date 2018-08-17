FROM golang

ENV APP_LOCATION=$GOPATH/src/github.com/map34/simple_tftp

WORKDIR ${APP_LOCATION}

COPY . ${APP_LOCATION}

RUN curl https://glide.sh/get | sh && cd ${APP_LOCATION} && glide install