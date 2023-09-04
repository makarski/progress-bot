FROM golang:alpine3.18

WORKDIR /progress-bot
ADD . .

RUN go install

ENTRYPOINT [ "progress-bot" ]
