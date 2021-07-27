# D4DJ-INFO-SERVER

## Getting Started

1. Clone repository

```bash
git clone https://github.com/gangjun06/d4dj-info-server --recursive
cd d4dj-info-server
```

3. Build D4DJ-Tool

```
cd D4DJ-Tool
dotnet build --configuration Release -o ../D4DJ-Tool-bin
```

2. Add config.toml
- rename `conf.exam.toml` to `conf.toml`
- edit config file

4. Run!

```
go run main.go
```

## Credits

- [GEEKiDoS/D4DJ-Tools](https://github.com/GEEKiDoS/D4DJ-Tools): D4DJ-Tool [MIT LICENSE](https://github.com/GEEKiDoS/D4DJ-Tools/blob/master/LICENSE)

### Special Thanks

- [DCinside D4DJ Gallary](https://gall.dcinside.com/mgallery/board/lists?id=d4dj)
- [GPLNature](https://github.com/GPLNature): help to decrypt file in golang
