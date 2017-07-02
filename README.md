## hauth简介

获取完整前后台源代码,请跳转到:

[golang 版本：https://github.com/hzwy23/asofdate](https://github.com/asofdate/hauth) 

[Java 版本：https://github.com/hzwy23/asofdate-etl](https://github.com/asofdate/hauth-java) 

## 项目介绍

这部分代码是asofdate hauth项目的可执行部分,asofdate hauth包含了后台golang开发的源代码

hauth-ui是asofdate hauth编译好后的可执行部分, 这是一个完整的前端部分.

所以, 如果不想使用golang的后台,可以使用别的语言实现hauth的API, 系统所有的API,在登录系统后,

在[系统帮助]->[API文档] 可以查看hauth全部的API

## 安装介绍

A).首先导入数据库文件,数据库文件在script目录中,导入数据库文件方法

```shell
cd script
mysql -uroot -p dbname < ./init_hauth.sql
```

B).修改数据库配置信息

配置文件在conf目录中，app.conf是beego的配置文件，主要涉及到服务端口号等等，另外一个是asofdate.conf配置文件，这个里边主要是是=数据库连接信息与日志管理信息配置。

beego的配置方法，请在beego项目中查阅，请移步：beego.me。下边来讲讲asofdate.conf中数据库的配置方法。

```
DB.type=mysql
DB.tns = "tcp(localhost:3306)/dbname"
DB.user = root
DB.passwd="xzPEh+SfFL3aimN0zGNB9w=="
```

注意: 修改的文件必须保存为utf-8编码,否则可能会出现异常，DB.type=mysql，这个值请不要修改，因为当前项目中提供的数据库脚本是针对于mysql和mariadb的。

1. 修改DB.tns中对应的数据库地址，端口号，数据库名称。

2. 修改DB.user成相应的数据库用户名

3. 修改DB.passwd成上边用户所对应的密码，系统启动后会自动加密，在此输入密码明文即可。

## 启动服务

```shell
# amd64上执行
sudo ./asofdate_amd64 &
# i386上执行
sudo ./asofdate_i386 &
# window系统, 双击下边的可执行文件
asofdate.exe
```

打开浏览器,登录 https://localhost:8090

管理员 admin , 密码: hzwy23


## 温馨提示:

获取完整源代码,请移步: [https://github.com/hzwy23/asofdate](https://github.com/hzwy23/asofdate)
