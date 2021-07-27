FROM golang:1.16.6-alpine3.13 as build-go
RUN mkdir app
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -o server

FROM mcr.microsoft.com/dotnet/nightly/runtime:5.0.8-alpine3.14-amd64
ENV ASSET_PATH="/app/assets"
WORKDIR /app
COPY . /app
COPY --from=build-go /app/server .
VOLUME [ ${ASSET_PATH} ]
ENTRYPOINT ["/app/server"]