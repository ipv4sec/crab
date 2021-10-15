FROM harbor1.zlibs.com/foo/node:13.13.0
MAINTAINER huangyulong@huanqiu.com
WORKDIR /appManager
COPY    ./code/bin                ./bin
COPY    ./code/public             ./public
COPY    ./code/routes             ./routes
COPY    ./code/server             ./server
COPY    ./code/views              ./views
COPY    ./code/app.js             ./app.js
COPY    ./code/package.json       ./package.json
COPY    ./code/package-lock.json  ./package-lock.json
COPY    ./config/config.json      ./config/config.json
RUN     npm install --production --registry=https://registry.npm.taobao.org
EXPOSE  80
CMD     ["node", "./bin/www"]