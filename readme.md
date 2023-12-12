# portsmap

golang编写的 **端口映射工具**



## 配置文件

portsmap.ini

方式一

```ini
"127.0.0.1:6666"="0.0.0.0:7777"
```

> 等号左边是原 `ip:port`，右边是映射后的 `ip:port`，默认启用
>
> 注意：等号左边必须加引号

方式二

```ini
[xrdp]
addr=127.0.0.1:3389
map=0.0.0.0:13389
active=true
```

> 在一个中括号的标签中写的配置
>
> addr 是原 `ip:port`
>
> map 是映射后的 `ip:port`
>
> active 设置是否启用
>
> 注意：如果有多个映射的话，中括号的标签不能重复

方式一和方式二可一起使用

```ini
"127.0.0.1:6666"="0.0.0.0:7777"
"127.0.0.1:8888"="0.0.0.0:9999"

[xrdp]
addr=127.0.0.1:3389
map=0.0.0.0:13389
active=true

[code]
addr=127.0.0.1:8008
map=0.0.0.0:18008
active=true
```

> 注意：方式一前必须是无中括号写的标签的，也就是说必须写在配置文件最上面



## 使用

```
./portsmap
```



## 开机启动

建议使用portsmap.sh或portsmap.cmd来使用

linux portsmap.sh

```sh
cd /home/xiao/my/portsmap
./portsmap
```

windows portsmap.cmd

```
cd /d D:\my\portsmap
portsmap.exe
```

自启动：

linux可选择新建service来自启动

> 自行查询service文件新建和使用方法

windows可选择将portsmap.cmd的快捷方式放到启动文件夹

> win+R，输入shell:startup回车，就是启动文件夹



## 适用场景

- 我比较喜欢windows的wsl，但wsl2中开启的服务无法映射到外网，只能在windows本地中以localhost来访问，所以配合这个portsmap特别好使。

- 其他单纯想换个端口的情况

- 流量加解密，当然这个没做完，有需求可以自行实现

- 配合vnet虚拟组网工具使用更香甜



## 使用的其他外部库

- gopkg.in/ini.v1  用来读取ini配置文件



## 注意事项

- 所有连接使用tcp链路，只支持tcp相关流量的转发，例如http协议也可，udp类型则不行。

- portsmap设计支持流量加密，但没有实现，有需求可在代码中的注释加密方法处自行实现

    