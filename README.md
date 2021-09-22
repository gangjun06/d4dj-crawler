# D4DJ-INFO-SERVER

## Getting Started

1. Clone repository

```bash
git clone https://github.com/gangjun06/d4dj-info-server --recursive
cd d4dj-info-server
```

2. Build D4DJ-Tool

```
cd D4DJ-Tool
dotnet build --configuration Release -o ../D4DJ-Tool-bin
```

3. Add config.toml

- rename `conf.exam.toml` to `conf.toml`
- edit config file

4. Run!

```
go run main.go
```

### Run with docker

You should build D4DJ-Tool before build dockerfile

1. Build or Pull

```bash
docker build . -t d4dj-crawler
```

OR

```bash
docker pull gangjun06/d4dj-crawler
```

2. Run

```bash

docker run \
	-v $(pwd)/config.toml:/app/config.toml \
	-v $(pwd)/assets:/app/assets \
	--name d4dj-crawler \
	{d4dj-crawler} OR {gangjun06/d4dj-crawler}
```

## Credits

### Used Tools

- [D4DJ-Tools](https://github.com/gangjun06/D4DJ-Tools): Parsing msgpack to json
  - Forked from [GEEKiDoS/D4DJ-Tools](https://github.com/GEEKiDoS/D4DJ-Tools/tree/master/D4DJ.Types)
- [D4DJ-asset-extractor](https://github.com/gangjun06/D4DJ-asset-extractor): Extract Unity AssetBundle
  - Fork from [https://github.com/Perfare/UnityLive2DExtractor](https://github.com/Perfare/UnityLive2DExtractor)
  - Using [nesrak1/AssetsTools.NET](https://github.com/nesrak1/AssetsTools.NET)

### Special Thanks

- [KJHMAGIC](https://github.com/kjhmagic): Helped with the overall game parsing.
- [GPLNature](https://github.com/GPLNature): Help to decrypt file in golang
