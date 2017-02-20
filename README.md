# hauth
	hauth是一个独立的权限控制系统。支持，域-角色-组织-用户几个维度的用户管理。
	
## 依赖库
* [beego](https://github.com/astaxie/beego)
* [jwt-go](https://github.com/dgrijalva/jwt-go)
* [dbobj](https://github.com/hzwy23/dbobj)

## 注意事项
    在使用这款系统时，首先导入必要的数据结构，默认支持mysql，mariadb数据库。
## 安装说明
1.首先将表结构与数据导入到数据库中。数据库脚本在script/init_hauth.sql。通过mysql工具导入到数据库中
> mysql -uroot -p 数据库名  < ./script/init_hauth.sql
2.编译下载下来的源代码，生成可执行文件。
3.编译完成后，请下载前端框架，下载地址：[github.com/hzwy23/devops]https://github.com/hzwy23/devops
4.将第二步中编译的可执行文件复制到第三步中下载下来的devops目录中。修改devops目录中conf目录下的system.properties文件，
```
DB.type=mysql
DB.tns = "tcp(localhost:3306)/test"
DB.user = root
DB.passwd="xzPEh+SfFL3aimN0zGNB9w=="
```
1.修改DB.tns中对应的数据库地址，端口号，数据库名称。
2.修改DB.user成相应的数据库用户名
3.修改DB.passwd成上边用户所对应的密码，系统启动后会自动加密，在此输入密码明文即可。
4.导入环境变量BIGDATA_HOME=devops，将devops修改成绝对路径。
5.在devops目录中运行可执行文件。
