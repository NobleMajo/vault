### BASE
FROM golang AS base

EXPOSE 8080
WORKDIR /app

COPY go.mod* go.sum* ./
RUN go mod tidy

### LOCAL
FROM base AS local

RUN go install github.com/air-verse/air@v1

ENTRYPOINT air

### BASE DEPLOY
FROM base AS base-deploy
COPY . .
RUN make build

### DEPLOY
FROM ubuntu:24.04 AS deploy

RUN useradd -m appuser --uid 10000
USER 10000

COPY --from=base-deploy --chown=10000 /app/bin /usr/local/bin/appbin

CMD ["appbin"]
