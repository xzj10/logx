# logx
可以根据日志级别将日志输出到文件或控制天的日志包

##### 日记级别: 分别为debug,info,warn,error,fatal (通过**logx.NewLogx("info")**进行设置)

> 以上日志级别依次递增, 代码内的打印的日志级别小于设定的日志级别则不打印
>
> debug级别的日志打印到控制台
>
> 大于debug级别的写入日志文件
>
> info, warn 的日志输入到  ./logs/info文件夹下; error, fatal的日志输入到 ./logs/error 文件夹下, 文件如果不存在会自动创建
>
> fatal 级别的日志会引发panic, 退出程序  

  
