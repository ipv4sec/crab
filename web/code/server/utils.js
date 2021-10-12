/**
 * Created by maikuraki on 2016/8/11.
 */
var crypto = require('crypto');
var request = require('request');
var config = require('../config/config.json');
var appKey = config.appConfig.key;
var appSecret = config.appConfig.secret;
var rediskey = config.webConfig.redisKey;
var i1Domain = config.webConfig.i1Host;

/**
 * 普通get请求
 * @param path
 * @param callback
 */
exports.normalGet = function (path, req, callback) {
    request({
        url: path,
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'X-Request-Id': req.headers.hasOwnProperty('X-Request-Id') ? req.headers['X-Request-Id'] : '',

            'C-Uid': req.headers.hasOwnProperty('c-uid') ? req.headers['c-uid'] : '',
            'C-Firstname': req.headers.hasOwnProperty('c-firstname') ? req.headers['c-firstname'] : '',
            'C-Lastname': req.headers.hasOwnProperty('c-lastname') ? req.headers['c-lastname'] : '',
            'C-Name': req.headers.hasOwnProperty('c-name') ? req.headers['c-name'] : '',
            'C-Domain': req.headers.hasOwnProperty('c-domain') ? req.headers['c-domain'] : '',
            'C-Ui-Domain': req.headers.hasOwnProperty('c-ui-domain') ? req.headers['c-ui-domain'] : '',
            'C-Api-Domain': req.headers.hasOwnProperty('c-api-domain') ? req.headers['c-api-domain'] : '',
            'C-Email': req.headers.hasOwnProperty('c-email') ? req.headers['c-email'] : '',

            'C-Cs-Id': req.headers.hasOwnProperty('c-cs-id') ? req.headers['c-cs-id'] : '',
            'C-Instance-Id': req.headers.hasOwnProperty('c-instance-id') ? req.headers['c-instance-id'] : '',
            'C-Groups': req.headers.hasOwnProperty('c-groups') ? req.headers['c-groups'] : '',
            'C-Roles': req.headers.hasOwnProperty('c-roles') ? req.headers['c-roles'] : '',
        }
    }, function (error, response, body) {
        if (error) {
            console.log(error)
            callback('{"status": "error","message": "服务器请求错误1"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    callback(body);
                }
            } catch (err) {
                console.log(err);
                // console.log(body);
                // console.log("Error name: " + err.name + "");
                // console.log("Error message: " + err.message);
                callback('{"status": "error","message": "服务器请求错误2"}');
            }
        }
    });
};
/**
 * 普通post请求
 * @param path
 * @param data
 * @param callback
 */
exports.normalPost = function (path, req, data, callback) {
    request({
        url: path,
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'X-Request-Id': req.headers.hasOwnProperty('X-Request-Id') ? req.headers['X-Request-Id'] : '',

            'C-Uid': req.headers.hasOwnProperty('c-uid') ? req.headers['c-uid'] : '',
            'C-Firstname': req.headers.hasOwnProperty('c-firstname') ? req.headers['c-firstname'] : '',
            'C-Lastname': req.headers.hasOwnProperty('c-lastname') ? req.headers['c-lastname'] : '',
            'C-Name': req.headers.hasOwnProperty('c-name') ? req.headers['c-name'] : '',
            'C-Domain': req.headers.hasOwnProperty('c-domain') ? req.headers['c-domain'] : '',
            'C-Ui-Domain': req.headers.hasOwnProperty('c-ui-domain') ? req.headers['c-ui-domain'] : '',
            'C-Api-Domain': req.headers.hasOwnProperty('c-api-domain') ? req.headers['c-api-domain'] : '',
            'C-Email': req.headers.hasOwnProperty('c-email') ? req.headers['c-email'] : '',

            'C-Cs-Id': req.headers.hasOwnProperty('c-cs-id') ? req.headers['c-cs-id'] : '',
            'C-Instance-Id': req.headers.hasOwnProperty('c-instance-id') ? req.headers['c-instance-id'] : '',
            'C-Groups': req.headers.hasOwnProperty('c-groups') ? req.headers['c-groups'] : '',
            'C-Roles': req.headers.hasOwnProperty('c-roles') ? req.headers['c-roles'] : '',
        },
        body: JSON.stringify(data)
    }, function (error, response, body) {
        if (error) {
            console.log(error)
            callback('{"status": "error","message": "服务器请求错误3"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    callback(body);
                }
            } catch (err) {
                console.log(err)
                console.log(body)
                // console.log("Error name: " + err.name + "");
                // console.log("Error message: " + err.message);
                callback('{"status": "error","msg": "服务器请求错误4"}');
            }
        }
    });
};
/**
 * 需要签名的的get请求
 * @param path
 * @param signature
 * @param ts
 * @param callback
 */
exports.get = function (path, req, callback) {
    var ts = Math.ceil(new Date().getTime() / 1000);
    var sig = crypto.createHash('md5').update(ts + appSecret).digest('hex');
    // console.log('sig:'+sig +'~~~~~' + 'appkey:' + appKey)
    request({
        url: path,
        method: 'GET',
        rejectUnauthorized: false,
        headers: {
            'Content-Type': 'application/json',
            'sig': sig,
            'key': appKey,
            'time': ts,
            'X-Request-Id': req.headers.hasOwnProperty('X-Request-Id') ? req.headers['X-Request-Id'] : '',
            'C-Uid': req.headers.hasOwnProperty('c-uid') ? req.headers['c-uid'] : '',
            'C-Firstname': req.headers.hasOwnProperty('c-firstname') ? req.headers['c-firstname'] : '',
            'C-Lastname': req.headers.hasOwnProperty('c-lastname') ? req.headers['c-lastname'] : '',
            'C-Name': req.headers.hasOwnProperty('c-name') ? req.headers['c-name'] : '',
            'C-Domain': req.headers.hasOwnProperty('c-domain') ? req.headers['c-domain'] : '',
            'C-Ui-Domain': req.headers.hasOwnProperty('c-ui-domain') ? req.headers['c-ui-domain'] : '',
            'C-Api-Domain': req.headers.hasOwnProperty('c-api-domain') ? req.headers['c-api-domain'] : '',
            'C-Email': req.headers.hasOwnProperty('c-email') ? req.headers['c-email'] : '',

            'C-Cs-Id': req.headers.hasOwnProperty('c-cs-id') ? req.headers['c-cs-id'] : '',
            'C-Instance-Id': req.headers.hasOwnProperty('c-instance-id') ? req.headers['c-instance-id'] : '',
            'C-Groups': req.headers.hasOwnProperty('c-groups') ? req.headers['c-groups'] : '',
            'C-Roles': req.headers.hasOwnProperty('c-roles') ? req.headers['c-roles'] : '',

        }
    }, function (error, response, body) {
        if (error) {
            console.log(error);
            callback('{"status": "error","msg": "服务器请求错误5"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    callback(body);
                }
            } catch (err) {
                console.log(err);
                // console.log(body);
                // console.log('--------------------------')
                // console.log(body)
                // console.log('--------------------------')
                // console.log("Error name: " + err.name + "");
                // console.log("Error message: " + err.message);
                // console.log("PATH: " + path);
                callback('{"status": "error","msg": "服务器请求错误6"}');
            }
        }
    });
};
/**
 * 带token的get请求
 * @param path 请求url
 * @param token 服务器分配的token
 * @param callback 回调函数
 */
exports.tokenGet = function (path, req, token, callback) {
    request({
        url: path,
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token,
            'X-Request-Id': req.headers.hasOwnProperty('X-Request-Id') ? req.headers['X-Request-Id'] : '',
            'C-Uid': req.headers.hasOwnProperty('c-uid') ? req.headers['c-uid'] : '',
            'C-Firstname': req.headers.hasOwnProperty('c-firstname') ? req.headers['c-firstname'] : '',
            'C-Lastname': req.headers.hasOwnProperty('c-lastname') ? req.headers['c-lastname'] : '',
            'C-Name': req.headers.hasOwnProperty('c-name') ? req.headers['c-name'] : '',
            'C-Domain': req.headers.hasOwnProperty('c-domain') ? req.headers['c-domain'] : '',
            'C-Ui-Domain': req.headers.hasOwnProperty('c-ui-domain') ? req.headers['c-ui-domain'] : '',
            'C-Api-Domain': req.headers.hasOwnProperty('c-api-domain') ? req.headers['c-api-domain'] : '',
            'C-Email': req.headers.hasOwnProperty('c-email') ? req.headers['c-email'] : '',

            'C-Cs-Id': req.headers.hasOwnProperty('c-cs-id') ? req.headers['c-cs-id'] : '',
            'C-Instance-Id': req.headers.hasOwnProperty('c-instance-id') ? req.headers['c-instance-id'] : '',
            'C-Groups': req.headers.hasOwnProperty('c-groups') ? req.headers['c-groups'] : '',
            'C-Roles': req.headers.hasOwnProperty('c-roles') ? req.headers['c-roles'] : '',
        }
    }, function (error, response, body) {
        if (error) {
            console.log(error)
            callback('{"status": "error","msg": "服务器请求错误7"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    callback(body);
                }
            } catch (err) {
                console.log(err)
                console.log(body)
                // console.log("Error name: " + err.name + "");
                // console.log("Error message: " + err.message);
                callback('{"status": "error","msg": "服务器请求错误8"}');
            }
        }
    });
};
/**
 * post请求
 * @param path
 * @param data
 * @param callback
 */
exports.post = function (path, req, data, callback) {
    var ts = Math.ceil(new Date().getTime() / 1000);
    var sig = crypto.createHash('md5').update(ts + appSecret).digest('hex');
    request({
        url: path,
        method: 'POST',
        rejectUnauthorized: false,
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
            'sig': sig,
            'key': appKey,
            'time': ts,
            'X-Request-Id': req.headers.hasOwnProperty('X-Request-Id') ? req.headers['X-Request-Id'] : '',

            'C-Uid': req.headers.hasOwnProperty('c-uid') ? req.headers['c-uid'] : '',
            'C-Firstname': req.headers.hasOwnProperty('c-firstname') ? req.headers['c-firstname'] : '',
            'C-Lastname': req.headers.hasOwnProperty('c-lastname') ? req.headers['c-lastname'] : '',
            'C-Name': req.headers.hasOwnProperty('c-name') ? req.headers['c-name'] : '',
            'C-Domain': req.headers.hasOwnProperty('c-domain') ? req.headers['c-domain'] : '',
            'C-Ui-Domain': req.headers.hasOwnProperty('c-ui-domain') ? req.headers['c-ui-domain'] : '',
            'C-Api-Domain': req.headers.hasOwnProperty('c-api-domain') ? req.headers['c-api-domain'] : '',
            'C-Email': req.headers.hasOwnProperty('c-email') ? req.headers['c-email'] : '',

            'C-Cs-Id': req.headers.hasOwnProperty('c-cs-id') ? req.headers['c-cs-id'] : '',
            'C-Instance-Id': req.headers.hasOwnProperty('c-instance-id') ? req.headers['c-instance-id'] : '',
            'C-Groups': req.headers.hasOwnProperty('c-groups') ? req.headers['c-groups'] : '',
            'C-Roles': req.headers.hasOwnProperty('c-roles') ? req.headers['c-roles'] : '',
        },
        formData: data
    }, function (error, response, body) {
        if (error) {
            console.log(error)
            callback('{"status": "error","msg": "服务器请求错误9"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    callback(body);
                }
            } catch (err) {
                console.log(err)
                console.log(body)
                // console.log('--------------------------')
                // console.log(body)
                // console.log('--------------------------')
                // console.log("Error name: " + err.name + "");
                // console.log("Error message: " + err.message);
                callback('{"status": "error","msg": "服务器请求错误10"}');
            }
        }
    });
};

/**
 * 带token的post请求
 * @param path 请求url
 * @param token 服务器分配的token
 * @param callback 回调函数
 */
exports.tokenPost = function (path, req, token, data, callback) {
    request({
        url: path,
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token,
            'X-Request-Id': req.headers.hasOwnProperty('X-Request-Id') ? req.headers['X-Request-Id'] : '',

            'C-Uid': req.headers.hasOwnProperty('c-uid') ? req.headers['c-uid'] : '',
            'C-Firstname': req.headers.hasOwnProperty('c-firstname') ? req.headers['c-firstname'] : '',
            'C-Lastname': req.headers.hasOwnProperty('c-lastname') ? req.headers['c-lastname'] : '',
            'C-Name': req.headers.hasOwnProperty('c-name') ? req.headers['c-name'] : '',
            'C-Domain': req.headers.hasOwnProperty('c-domain') ? req.headers['c-domain'] : '',
            'C-Ui-Domain': req.headers.hasOwnProperty('c-ui-domain') ? req.headers['c-ui-domain'] : '',
            'C-Api-Domain': req.headers.hasOwnProperty('c-api-domain') ? req.headers['c-api-domain'] : '',
            'C-Email': req.headers.hasOwnProperty('c-email') ? req.headers['c-email'] : '',

            'C-Cs-Id': req.headers.hasOwnProperty('c-cs-id') ? req.headers['c-cs-id'] : '',
            'C-Instance-Id': req.headers.hasOwnProperty('c-instance-id') ? req.headers['c-instance-id'] : '',
            'C-Groups': req.headers.hasOwnProperty('c-groups') ? req.headers['c-groups'] : '',
            'C-Roles': req.headers.hasOwnProperty('c-roles') ? req.headers['c-roles'] : '',
        },
        body: JSON.stringify(data)
    }, function (error, response, body) {
        if (error) {
            console.log(error)
            callback('{"status": "error","msg": "服务器请求错误11"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    callback(body);
                }
            } catch (err) {
                console.log(err)
                console.log(body)
                // console.log("Error name: " + err.name + "");
                // console.log("Error message: " + err.message);
                callback('{"status": "error","msg": "服务器请求错误12"}');
            }
        }
    });
};
/**
 * 带token的put请求
 * @param path
 * @param token
 * @param data
 * @param callback
 */
exports.tokenPut = function (path, req, token, data, callback) {
    request({
        url: path,
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token,
            'X-Request-Id': req.headers.hasOwnProperty('X-Request-Id') ? req.headers['X-Request-Id'] : '',

            'C-Uid': req.headers.hasOwnProperty('c-uid') ? req.headers['c-uid'] : '',
            'C-Firstname': req.headers.hasOwnProperty('c-firstname') ? req.headers['c-firstname'] : '',
            'C-Lastname': req.headers.hasOwnProperty('c-lastname') ? req.headers['c-lastname'] : '',
            'C-Name': req.headers.hasOwnProperty('c-name') ? req.headers['c-name'] : '',
            'C-Domain': req.headers.hasOwnProperty('c-domain') ? req.headers['c-domain'] : '',
            'C-Ui-Domain': req.headers.hasOwnProperty('c-ui-domain') ? req.headers['c-ui-domain'] : '',
            'C-Api-Domain': req.headers.hasOwnProperty('c-api-domain') ? req.headers['c-api-domain'] : '',
            'C-Email': req.headers.hasOwnProperty('c-email') ? req.headers['c-email'] : '',

            'C-Cs-Id': req.headers.hasOwnProperty('c-cs-id') ? req.headers['c-cs-id'] : '',
            'C-Instance-Id': req.headers.hasOwnProperty('c-instance-id') ? req.headers['c-instance-id'] : '',
            'C-Groups': req.headers.hasOwnProperty('c-groups') ? req.headers['c-groups'] : '',
            'C-Roles': req.headers.hasOwnProperty('c-roles') ? req.headers['c-roles'] : '',
        },
        body: JSON.stringify(data)
    }, function (error, response, body) {
        if (error) {
            console.log(error)
            callback('{"status": "err","msg": "服务器请求错误13"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    callback(body);
                }
            } catch (err) {
                console.log(err)
                console.log(body)
                // console.log("Error name: " + err.name + "");
                // console.log("Error message: " + err.message);
                callback('{"status": "err","msg": "服务器请求错误14"}');
            }
        }
    });
};

/**
 * 带token的delete请求
 * @param path
 * @param token
 * @param data
 * @param callback
 */
exports.tokenDelete = function (path, req, token, data, callback) {
    request({
        url: path,
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token,
            'X-Request-Id': req.headers.hasOwnProperty('X-Request-Id') ? req.headers['X-Request-Id'] : '',

            'C-Uid': req.headers.hasOwnProperty('c-uid') ? req.headers['c-uid'] : '',
            'C-Firstname': req.headers.hasOwnProperty('c-firstname') ? req.headers['c-firstname'] : '',
            'C-Lastname': req.headers.hasOwnProperty('c-lastname') ? req.headers['c-lastname'] : '',
            'C-Name': req.headers.hasOwnProperty('c-name') ? req.headers['c-name'] : '',
            'C-Domain': req.headers.hasOwnProperty('c-domain') ? req.headers['c-domain'] : '',
            'C-Ui-Domain': req.headers.hasOwnProperty('c-ui-domain') ? req.headers['c-ui-domain'] : '',
            'C-Api-Domain': req.headers.hasOwnProperty('c-api-domain') ? req.headers['c-api-domain'] : '',
            'C-Email': req.headers.hasOwnProperty('c-email') ? req.headers['c-email'] : '',

            'C-Cs-Id': req.headers.hasOwnProperty('c-cs-id') ? req.headers['c-cs-id'] : '',
            'C-Instance-Id': req.headers.hasOwnProperty('c-instance-id') ? req.headers['c-instance-id'] : '',
            'C-Groups': req.headers.hasOwnProperty('c-groups') ? req.headers['c-groups'] : '',
            'C-Roles': req.headers.hasOwnProperty('c-roles') ? req.headers['c-roles'] : '',
        },
        body: JSON.stringify(data)
    }, function (error, response, body) {
        if (error) {
            console.log(error)
            callback('{"status": "err","msg": "服务器请求错误15"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    callback(body);
                }
            } catch (err) {
                console.log(err)
                console.log(body)
                // console.log("Error name: " + err.name + "");
                // console.log("Error message: " + err.message);
                callback('{"status": "err","msg": "服务器请求错误16"}');
            }
        }
    });
};

exports.i1Get = function (path, query, req, callback) {
    let ts = new Date().toISOString();
    let app_secret = req.session.appSecret || '';
    let seq = "GET" + "\n" + path + "\n" + query + "\n" + "X-I1-Date=" + ts + "&X-I1-Host=" + req.headers['c-domain'] + "&X-I1-Key=" + req.session.appKey + "\n";
    let signature = crypto.createHmac('sha1', app_secret).update(seq).digest('base64');
    let q = query ? '?' + query : '';
    request({
        url: i1Domain + path + q,
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'X-I1-Date': ts,
            'X-I1-Host': req.headers['c-domain'],
            'X-I1-Key': req.session.appKey,
            'X-I1-Signature': signature,
        },
    }, function (error, response, body) {
        if (error) {
            console.log(error)
            callback({"status": "error", "msg": "服务器请求错误"});
        } else {
            if (typeof body === 'object') {
                callback(body);
            } else {
                try {
                    if (JSON.parse(body)) {
                        callback(JSON.parse(body));
                    }
                } catch (err) {
                    console.log(err)
                    console.log(body)
                    // console.log(body)
                    // console.log("Error name: " + err.name + "");
                    // console.log("Error message: " + err.message);
                    callback({"status": "error", "msg": "服务器请求错误"});
                }
            }
        }
    });
};

/**
 * 获取cookie
 * @param name
 * @returns {*}
 */
exports.getCookie = function (name) {
    var arrStr = document.cookie.split(";");
    for (var i = 0; i < arrStr.length; i++) {
        var temp = arrStr[i].split("=");
        var tempName = temp[0].replace(/ /, "");
        if (tempName == name) {
            return temp[1];
        }
    }
};


/**
 * 生成redis key储存登录状态
 * @param req
 */
exports.getRedisKey = function (req) {
    var key = rediskey;
    var t = new Date().toISOString();
    var token;
    try {
        token = crypto.createHmac('sha1', key).update(t).digest().toString('base64');
    } catch (err) {
        console.log(err);
    }
    req.session.m1redisKey = token;
    // console.log('redisKey:' + token);
    return token;
};