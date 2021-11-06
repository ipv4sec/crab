FROM harbor1.zlibs.com/foo/node:13.13.0
MAINTAINER huangyulong@huanqiu.com
WORKDIR /crab-front
COPY    ./public             ./public
COPY    ./server             ./server
COPY    ./docker/config.json        ./server/config.json
COPY    ./package.json       ./package.json
COPY    ./package-lock.json  ./package-lock.json
RUN     npm install --production --registry=https://registry.npm.taobao.org
EXPOSE  80
CMD     ["node", "./server/index.js"]