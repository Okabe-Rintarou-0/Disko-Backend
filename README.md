# Disko 后端

基于 [go-zero](https://github.com/zeromicro/go-zero) 实现 API 代码自动生成，[gorm](https://github.com/go-gorm/gorm)
实现数据库管理。

## 利用 GORM 管理 mysql

### Soft Delete

> If your model includes a gorm.DeletedAt field (which is included in gorm.Model), it will get soft delete ability
> automatically!

## 如何处理 JWT 在用户登出之后销毁的问题

JWT 本质上是将一些信息在服务器编码，在用户登录后将 token 发送给客户端。客户端在后续发送请求时在 `Authorization` header
带上这个 token。
服务器对 token 进行解码，若为合法数据，则认为是授权用户。

用户在登出的时候，应该将 token 失效，但是 JWT 的实现原理决定了它在过期前一定都是有效的。

解决方案：
> 黑名单校验: 创建 JWT 黑名单。
> 根据过期时间，当客户端删除其令牌时，它可能仍然有效一段时间。如果令牌生存期很短，则可能不是问题，但如果您仍希望令牌立即失效，则可以创建令牌黑名单。当后端收到注销请求时，从请求中获取
> JWT 并将其存储在内存数据库中。对于每个经过身份验证的请求，您需要检查内存数据库以查看令牌是否已失效。为了保持较小的搜索空间，您可以从黑名单中删除已经过期的令牌。（根据令牌剩余有效期设置内存数据失效时间，达到自动清除的目的）

## 如何让浏览器能够下载云盘文件

参考：https://stackoverflow.com/questions/24116147/how-to-download-file-in-browser-from-go-server
> In this case for content type you can use application/octet-stream because the browser does not have to know the MIME
> type of the response.

有关 `application/octet-stream` 格式，可以参考：https://juejin.cn/post/6979224810681270309

> application/octet-stream 是应用程序文件的默认值。意思是未知的应用程序文件
> ，浏览器一般不会自动执行或询问执行。浏览器会像对待，设置了HTTP头Content-Disposition 值为 attachment
> 的文件一样来对待这类文件，即浏览器会触发下载行为。

## 限制文件上传/下载的速率

使用 [github.com/juju/ratelimit](https://github.com/juju/ratelimit):

参考：https://stackoverflow.com/questions/27187617/how-would-i-limit-upload-and-download-speed-from-the-server-in-golang

## 服务器一次性把所有的数据都推给后端？

使用 ratelimit 的时候发现前端始终没有弹开下载项，原因是服务器将数据都缓存了起来，然后一次性推给前端。

解决方法：决定不使用 ratelimit, 自己用 time.sleep + w.(http.Flusher).Flush()。

## 使用 Minikube 进行部署

+ 记得先用 minikube mount 把本地文件挂载到 minikube 中。