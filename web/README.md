# crab-front

## 项目结构
```
client // 前端部分
    |-- build // 打包文件
    |-- config // 配置文件
    |-- src // 源码
    |-- utils // 公共方法
    |
docker // docker配置文件
    |-- config.json // 配置文件
    |-- Dockerfile 
    |
node_modules // 依赖
    |
public // 前端打包后的文件，并包含图片、第三方插件 
    |-- css // 打包后的css文件
    |-- font // 字体文件
    |-- images // 图片
    |-- js // 打包后的js文件
    |-- static // 放置第三方插件
    |-- index.html // 入口html文件
    |
server // node服务端代码
    |-- routers // 所有路由文件
    |-- utils // 服务端公共方法
    |-- views // ejs模版文件
    |-- index.js // 入口文件
    |
babel.config.json // babel配置文件
    |
package.json 
    |
package-lock.json
    |
README.md
```


## 主要的技术栈
1. react@17 、 react-dom@17
2. react-router-dom@5
3. redux@4
4. sass@1
5. webpack@5
6. express@4
7. axios


## 开发流程
```
1. 拉取代码，执行npm install
2. npm run dev 开启webpack-dev-server前端开发环境
3. npm run serve 开启node服务端
4. npm run build 执行前端打包
5. 执行完打包后，将public>index.html 中script和link元素复制到 server>views>index.ejs对应的位置中 （该操作以后优化）
```

## 部署方法
### 仅部署前端静态代码
```
将public文件夹下放到nginx的相应位置
```

### 部署node服务
```
将server文件夹和public文件夹放到同一个文件目录下，可以使用pm2去管理进程

或者使用docker去部署：在项目根目录下执行 docker build -f docker/Dockerfile -t imageName:version .  来创建镜像
```