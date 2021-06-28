个人站点[mihiru.com](https://mihiru.com)的后端接口项目, 原先采用Spring Boot开发, 感觉太占用内存上了Quarkus+GraalVM重构了一遍, 虽然内存占用降低了不少但还是挺高的. 于是干脆用GO语言又重构了一遍, 顺便把数据库从H2换成了mongodb.

# 调试运行
将[sample.yaml](./config/sample.yaml)复制一份并重命名为dev.yaml, 修改配置内容, 之后在IDE内执行main.go即可

# 编译生成部署用文件
执行buildLinux.bat可生成AMD64的linux环境下可执行文件, 需要打包其他环境的可执行文件请自行查找go编译可执行文件的相关教程.  
执行完后会在项目根路径下生成一个文件名为main的可执行文件, 将该文件拷贝到部署环境后, 为该文件添加执行权限
```shell
chmod +x main
```
在同目录下建立config目录, 在其下新建一个prod.yaml配置文件, 文件内容参考[sample.yaml](./config/sample.yaml), 按实际情况修改各项配置.  
然后使用以下命令即可启动程序
```shell
./main -e prod
```
你也可以使用如下命令在后台运行程序
```shell
nohup ./main -e prod 1>mihiru-go.log 2>&1 &
```
