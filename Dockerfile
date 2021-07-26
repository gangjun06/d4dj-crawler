FROM golang:1.16

ENV ASSET_PATH="/app/assets"

RUN mkdir app

COPY . /app

WORKDIR /app

RUN go build -o server

VOLUME [ ${ASSET_PATH} ]

ENTRYPOINT ["/app/server"]
