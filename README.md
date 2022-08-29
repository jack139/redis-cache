# Redis cache service



## Server run
```shell
go run main.go
```



## Test run
```shell
cd example
go run client.go
```



## Oracle docker test
```shell
docker run -d -p 49161:1521 -p 8080:8080 oracleinanutshell/oracle-xe-11g
```



## MS-SQL docker test
```shell
docker run -e "ACCEPT_EULA=Y" -e "SA_PASSWORD=Mssql_2019" -p 1433:1433 -d mcr.microsoft.com/mssql/server:2019-latest
```