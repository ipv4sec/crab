/**
 * Created by maikuraki on 2016/8/12.
 */
var crypto = require('crypto');
var config = require('../config/config.json');
/**
 * 生成token值
 */
exports.getToken = function (req) {
    var key = config.webConfig.tokenKey;
    var t = new Date().toISOString();
    var token = '';
    try {
        token = crypto.createHmac('sha1', key).update(t).digest().toString('base64');
    } catch (err) {
        console.log(err);
    }
    req.session.token = token;
    return token;
};
/**
 * 校检请求中的token是否和session中的token相同
 */
exports.checkToken = function (req, token) {
    if (req.session.token && req.session.token == token) {
        return true;
    } else {
        return false;
    }
};