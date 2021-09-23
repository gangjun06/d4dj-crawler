FROM golang:1.16.6-alpine3.13 as build-go
RUN mkdir app
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -o server

FROM mcr.microsoft.com/dotnet/sdk:5.0:5.0.401-alpine3.14-amd64 as build-tools
RUN mkdir tools
COPY D4DJ-asset-extractor D4DJ-Tool /app/
WORKDIR /app
RUN cd dD4DJ-asset-extractor/UnityLive2DExtractor &&  dotnet build --configuration Release -o ../../D4DJ-assets-extractor-bin
RUN cd dD4DJ-Tool &&  dotnet build --configuration Release -o ../../D4DJ-Tool-bin

FROM mcr.microsoft.com/dotnet/runtime:5.0:5.0.10-alpine3.14-amd64
ENV ASSET_PATH="/app/assets"
WORKDIR /app
COPY --from=build-go /app/server .
COPY --from=build-tools /app/D4DJ-asset-extractor-bin /app/D4DJ-Tool-bin ./
VOLUME [ ${ASSET_PATH} ]
ENTRYPOINT ["/app/server"]