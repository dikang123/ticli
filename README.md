# ticli
Golang的TiDB客户端，基于[go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)封装，支持以下特性：

* 多TiDB连接，自动切换故障TiDB
* TiDB负载均衡

## Usage

```go
import "github.com/sicojuy/ticli"

opt := &Option{
    Addresses: []string{"host1:port1", "host2:port2", "host3:port3"},
    User:      "user",
    Password:  "password",
    DB:        "dbname",
    Timeout:   3,
}
cli := NewClient(opt)
defer cli.Close()

db, err := cli.Open()
if err != nil {
    // open db error
}
db.Ping()
```
