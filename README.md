# d4dj-crawler
> It only works on Windows due to multiple dependencies. 

## Getting Started

1. Clone repository

```bash
git clone https://github.com/gangjun06/d4dj-crawler --recursive
cd d4dj-crawler
```

2. Download VgmStream

Goto [vgmstream](https://vgmstream.org/downloads) and download

Unzip program and locate to [project-path]/vgmstream

3. Download Asset Extractor
Goto [gangjun06/D4DJ-asset-extractor](https://github.com/gangjun06/D4DJ-asset-extractor/releases/tag/v0.1.0) and download

Unzip program and locate to [project-path]/D4DJ-assets-extractor-bin

3. Build D4DJ-Tool

```bash
cd D4DJ-Tool
dotnet build --configuration Release -o ../D4DJ-Tool-bin
```

4. Add config.toml

- Rename `config.exam.toml` to `config.toml`
- Edit config file

5. Build
```bash
go build
```

6. Crawl Resource
> You must input assetServerPath to `config.toml`. To get the asset url, you need to analyze the game packet.
```bash
d4dj-crawler.exe -crawl
```

6-1. Parse Resource
```bash
d4dj-crawler.exe [filename]
```
or just drag and drop your file to .exe

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
