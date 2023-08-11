# Disko 后端

基于 [go-zero](https://github.com/zeromicro/go-zero) 实现 API 代码自动生成，[gorm](https://github.com/go-gorm/gorm)
实现数据库管理。

## 如何处理 JWT 在用户登出之后销毁的问题

JWT 本质上是将一些信息在服务器编码，在用户登录后将 token 发送给客户端。客户端在后续发送请求时在 `Authorization` header 带上这个 token。
服务器对 token 进行解码，若为合法数据，则认为是授权用户。

用户在登出的时候，应该将 token 失效，但是 JWT 的实现原理决定了它在过期前一定都是有效的。

解决方案：
> 黑名单校验: 创建 JWT 黑名单。
> 根据过期时间，当客户端删除其令牌时，它可能仍然有效一段时间。如果令牌生存期很短，则可能不是问题，但如果您仍希望令牌立即失效，则可以创建令牌黑名单。当后端收到注销请求时，从请求中获取 JWT 并将其存储在内存数据库中。对于每个经过身份验证的请求，您需要检查内存数据库以查看令牌是否已失效。为了保持较小的搜索空间，您可以从黑名单中删除已经过期的令牌。（根据令牌剩余有效期设置内存数据失效时间，达到自动清除的目的）