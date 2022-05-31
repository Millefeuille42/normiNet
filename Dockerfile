FROM golang:1.16.0-alpine3.13

ENV PYTHONUNBUFFERED=1
ADD ./normiNet_server /normiNet_server/
WORKDIR /normiNet_server

RUN apk add --update --no-cache python3 && ln -sf python3 /usr/bin/python
RUN python3 -m ensurepip
RUN pip3 install --no-cache --upgrade pip setuptools
RUN python3 -m pip install norminette

RUN go install .

CMD normiNet_server ./dev.env
