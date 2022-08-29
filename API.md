## 接口设计



### 协议

 http/2 + TLS证书验证



### 提交格式
post form，字段名为 key



### 伪代码：

```golang
post_url = "https://localhost:8443/redis/cache"

client.PostForm(post_url, url.Values{"key": {key}})
```
