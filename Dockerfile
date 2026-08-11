### BASE
FROM golang:latest AS base

WORKDIR /app

### LOCAL
FROM base AS local

ENV HOME /home/tester
RUN groupadd -g 1000 tester \
    && useradd -u 1000 -g 1000 -m tester

RUN go install github.com/air-verse/air@v1
RUN echo 'export PS1="\u@go-container:\w\$ "' >> /etc/bash.bashrc

ENTRYPOINT air

### BASE DEPLOY
FROM base AS base-deploy
COPY . .
RUN make build

### DEPLOY
FROM ubuntu:latest AS deploy

RUN mkdir -p /app/data \
	&& chown -R 1000:1000 /app \
	&& chmod -R 755 /app

WORKDIR /app
USER 1000:1000

COPY --from=base-deploy --chown=1000:1000 \
	/app/bin /usr/local/bin/appbin

CMD ["appbin"]
