# mercari-search-cli

too simply cui tool to fetch a search result from mercari.jp

## build

```bash
go build
```

## usage

`mercari-search-cli` has search options

|Option name                     |Description                        |
|--------------------------------|-----------------------------------|
|--page, -p                      |page number (default: 0)           |
|--keyword, -k                   |search keyword                     |
|--price-min, --min              |minimum price (default: 0)         |
|--price-max, --max              |maximum price (default: 0)         |
|--category-root                 |category root number (default: 0)  |
|--category-child                |category child number (default: 0) |
|--brand-name                    |brand name keyword                 |
|--brand-id                      |brand id number (default: 0)       |
|--desc                          |search in desc order               |
|--on-sale                       |fetch only a on-sale items         |

```bash
# Example
mercari-search-cli -k ゲームボーイ --category-root 5 --category-child 76
```

## test

unimplemented!