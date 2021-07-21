# D4DJ-INFO-SERVER

## Getting Started

1. Clone repository

```bash
git clone https://github.com/gangjun06/d4dj-info-server
cd d4dj-info-server
```

3. build D4DJ-Tool

```
cd D4DJ-Tool
dotnet build --configuration Release
```

2. create .env file

```
SERVER_PORT=9096
ASSET_PATH=./assets
TOOL_PATH=/PATH/TO/D4DJ-Tool/EXECUTEABLE/FILE
```

4. Run!

```
go run main.go
```

## How it works
Crawling D4DJ Groovy-mix asset server every 1 hour.
The Asset server has a file called iOSResourceList.msgpack. If this file is changed, the added files are downloaded.


## Credits

- [GEEKiDoS/D4DJ-Tools](https://github.com/GEEKiDoS/D4DJ-Tools): D4DJ-Tool [MIT LICENSE](https://github.com/GEEKiDoS/D4DJ-Tools/blob/master/LICENSE)

### Special Thanks

- [DCinside D4DJ Gallary](https://gall.dcinside.com/mgallery/board/lists?id=d4dj)
- [GPLNature](https://github.com/GPLNature): help to decrypt file in golang
