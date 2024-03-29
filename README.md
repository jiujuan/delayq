## delayq
延迟队列，基于redis的zset，list数据结构开发的。写的一个小应用。

## 项目依赖库
- Redis内存数据库 - [github.com/redis/redis](https://github.com/redis/redis)
- redis的go库 - [github.com/gomodule/redigo](https://github.com/gomodule/redigo)
- web框架Gin - [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
- 日志zap - [go.uber.org/zap](https://github.com/uber-go/zap)
- 配置viper - [github.com/spf13/viper](https://github.com/spf13/viper)
- json解析库json-iterator - [github.com/json-iterator/go](https://github.com/json-iterator/go)

## 参考
https://tech.youzan.com/queuing_delay/  有赞的延迟队列

参考有赞这篇文章来编写程序。

## 需求背景：

	1. 用户下单发送短信服务
	2. 订单多少分钟未支付需要通知客户支付
	3. 订单未支付需要关闭订单，并退还库存
	4. 店铺信息快要到期时候发送通知
	5. 订单完成后通知用户评价
	6. 红包 24 小时未被查收，需要延迟执退还业务

等等需求

## 解决的方法有：

	1. 扫表
业务少的时候，可以扫表来解决，数据量大了，扫表肯定会出现时间的误差，效率会很低。
每个业务也需要维护自己的一套扫表逻辑。业务越来越多，扫表的业务也会越来越多。但是这部分逻辑又是重复的

	2. 延迟队列
延迟队列功能解决上面的需求


