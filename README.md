# ticli
Golang的TiDB客户端，基于[go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)封装，支持以下特性：

* 多TiDB连接，自动切换故障TiDB
* TiDB负载均衡
