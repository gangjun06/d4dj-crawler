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
you should build D4DJ-Tool before build dockerfile
```bash
docker build . -t d4dj-info-server
docker run \
	-v $(pwd)/config_docker.toml:/app/config.toml \
	-v $(pwd)/assets:/app/assets \
	-p 9096:9096 \
	--name d4dj-info-server \
	d4dj-info-server
```

## Credits

- [GEEKiDoS/D4DJ-Tools](https://github.com/GEEKiDoS/D4DJ-Tools): D4DJ-Tool [MIT LICENSE](https://github.com/GEEKiDoS/D4DJ-Tools/blob/master/LICENSE)

### Special Thanks

- [DCinside D4DJ Gallary](https://gall.dcinside.com/mgallery/board/lists?id=d4dj)
- [GPLNature](https://github.com/GPLNature): help to decrypt file in golang
