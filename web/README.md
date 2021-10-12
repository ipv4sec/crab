# appManager 相关说明

> 前端部分使用的是 react@16.3 + react-router + redux + sass ，使用webpack进行项目打包

> 后端部分使用NodeJs 的 express框架。


## 目录

* [项目目录]
    * code: 项目源码目录
    * config: 项目配置文件目录
    * docker: 项目的Dockerfile,用于docker镜像构建


## 运行
### 在docker中运行

> 项目目录下进行执行： docker build -f ./docker/Dockerfile -t imageName:version .  创建镜像


### 在服务器中使用pm2启动项目运行

1、安装
```
npm install -g pm2
```

2、建立软连接到全局运行环境
```
ln -s /usr/local/src/node-v6.5.0-linux-x64/bin/pm2 /usr/local/bin/pm2
```

3、进入项目路径(以 /opt/sites/appManager 为例),进行启动:
```
cd /opt/sites/appManager
pm2 start ./code/pm2.json
```

4、如需重启或者停止应用,可在进入项目目录后使用如下命令:
```
pm2 restart all
pm2 stop all
```
