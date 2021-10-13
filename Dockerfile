FROM golang:1.16.0-alpine3.13

ADD ./normiNet_server /normiNet_server/
WORKDIR /normiNet_server

RUN go install .

CMD normiNet_server ./dev.env