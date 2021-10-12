/**
 * Created by maikuraki on 2016/8/12.
 */
var redis = require('redis');
var config = require('../config/config.json');
var client = redis.createClient(config.redisConfig.port, config.redisConfig.host);

client.on('error', (error) => {
    console.log('redis连接失败：' + error);
});

client.on('connect', () => {
    console.log('--------redis连接成功--------');
});

module.exports = client;