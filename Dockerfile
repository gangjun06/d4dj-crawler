FROM golang:1.16.6-alpine3.13 as build-go
RUN mkdir app
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -o server

FROM mcr.microsoft.com/dotnet/nightly/sdk:5.0.103-alpine3.12 as build-tools
RUN mkdir app
COPY . /app/
WORKDIR /app/D4DJ-Tool
RUN dotnet build D4DJ-Tool.csproj --configuration Release -o ../D4DJ-Tool-bin
WORKDIR /app/D4DJ-asset-extractor/UnityLive2DExtractor
RUN dotnet build D4DJAssetExtractor.csproj --configuration Release -o ../../D4DJ-asset-extractor-bin

FROM mcr.microsoft.com/dotnet/nightly/runtime:5.0-alpine3.11
ENV ASSET_PATH="/app/assets"
WORKDIR /app
COPY --from=build-go /app/server .
COPY --from=build-tools /app/D4DJ-asset-extractor-bin/ /app/D4DJ-Tool-bin/ ./
VOLUME [ ${ASSET_PATH} ]
ENTRYPOINT ["/app/server", "-crawl"]