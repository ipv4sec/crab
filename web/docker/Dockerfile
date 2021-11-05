FROM harbor1.zlibs.com/foo/node:13.13.0
MAINTAINER huangyulong@huanqiu.com
WORKDIR /island-install-front
COPY    ./public             ./public
COPY    ./server             ./server
COPY    ./docker/config.json        ./server/config.json
COPY    ./package.json       ./package.json
COPY    ./package-lock.json  ./package-lock.json
RUN     npm install --production
EXPOSE  80
CMD     ["node", "./server/index.js"]