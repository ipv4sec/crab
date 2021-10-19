var request = require('request');
var config = require('../../config/config.json');

// 返回json格式数据
exports.get = function (pathName, header, callback, host='requestDomain', useRelativeHost) {
    var options = {
        // url: isStatic ? staticDomain + pathName : baseHost + pathName,
        url: useRelativeHost ? pathName : config.webConfig[host] + pathName,
        method: 'GET',
        headers: header,
    };
    // console.log('transform=url=',options.method, options.url)
    request(options, function (error, response, body) {
        if (error) {
            console.log(error);
            callback('{"status": "err","msg": "服务器请求错误"}');
        } else {
            try {
                let parseBody = JSON.parse(body)
                callback(parseBody);
            } catch (err) {
                console.log('body==', body)
                console.log(err);
                let obj = {
                    status: "error",
                    msg: "服务响应错误",
                    body: body
                };
                callback(JSON.stringify(obj));
                return
            }
        }
    });
};


// 直接返回数据
exports.getRes = function (pathName, header, callback, host='requestDomain', useRelativeHost) {
    var options = {
        // url: isStatic ? staticDomain + pathName : baseHost + pathName,
        url: useRelativeHost ? pathName : config.webConfig[host] + pathName  ,
        method: 'GET',
        headers: header,
    };
    // console.log('transform=url=',options.method, options.url)
    request(options, function (error, response, body) {
        if (error) {
            console.log(error);
            callback('{"status": "err","msg": "服务器请求错误"}');
        } else {
            callback(body);
        }
    });
};

exports.post = function (pathName, data, header, callback, host='requestDomain', useRelativeHost) {
    var options = {
        url: useRelativeHost ? pathName : config.webConfig[host] + pathName  ,
        method: 'POST',
        headers: header,
        body: JSON.stringify(data),
    };
    // console.log('transform=url=',options.method, options.url)
    request(options, function (error, response, body) {
        if (error) {
            console.log(error)
            callback('{"status": "err","msg": "服务器请求错误"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    callback(body);
                }else {
                    console.log('body==', body)
                }
            } catch (err) {
                console.log('body==', body)
                console.log(err);

                let obj = {
                    status: "error",
                    msg: "服务响应错误",
                    body: body
                };
                callback(JSON.stringify(obj));
            }
        }
    });
};

// 后台接收formdata的方式
exports.postForm = function (pathName, data, header, callback, host='requestDomain', relativeHost) {
    var options = {
        // url: isStatic ? staticDomain + pathName : baseHost + pathName,
        url: relativeHost ? relativeHost + pathName : config.webConfig[host] + pathName  ,
        method: 'POST',
        headers: header,
        formData: data,
    };
    request(options, function (error, response, body) {
       
        if (error) {
            console.log('err', error)
            callback('{"status": "err","msg": "服务器请求错误"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    let datas = JSON.parse(body)
                    callback(body);
                }else {
                    console.log('===JSON parse response error==', body)
                }
            } catch (err) {
                console.log('===response error==', body)
                console.log(err);

                let obj = {
                    status: "error",
                    msg: "服务响应错误",
                    body: body
                };
                callback(JSON.stringify(obj));
            }
        }
    });
};


exports.putForm = function (pathName, data, header, callback, host='requestDomain', relativeHost) {
    var options = {
        // url: isStatic ? staticDomain + pathName : baseHost + pathName,
        url: relativeHost ? relativeHost + pathName : config.webConfig[host] + pathName  ,
        method: 'PUT',
        headers: header,
        formData: data,
    };
    request(options, function (error, response, body) {
        if (error) {
            console.log('err', error)
            callback('{"status": "err","msg": "服务器请求错误"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    let datas = JSON.parse(body)
                    callback(body);
                }else {
                    console.log('===JSON parse response error ==', body)
                }
            } catch (err) {
                console.log('==response error ==', body)
                console.log(err);

                let obj = {
                    status: "error",
                    msg: "服务响应错误",
                    body: body
                };
                callback(JSON.stringify(obj));
            }
        }
    });
};


exports.deleteForm = function (pathName, header, callback, host='requestDomain', relativeHost) {
    var options = {
        // url: isStatic ? staticDomain + pathName : baseHost + pathName,
        url: relativeHost ? relativeHost + pathName : config.webConfig[host] + pathName ,
        method: 'DELETE',
        headers: header,
    };
    // console.log('--delete---', options)
    request(options, function (error, response, body) {
        if (error) {
            console.log('err', error)
            callback('{"status": "err","msg": "服务器请求错误"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    let datas = JSON.parse(body)
                    callback(body);
                }else {
                    console.log('===JSON parse response error ==', body)
                }
            } catch (err) {
                console.log('==response error ==', body)
                console.log(err);

                let obj = {
                    status: "error",
                    msg: "服务响应错误",
                    body: body
                };
                callback(JSON.stringify(obj));
            }
        }
    });
};



exports.postFile = function (pathName, data, header, callback, host='requestDomain', useRelativeHost) {
    // console.log('atiPost header:', header);
    var options = {
        // url: config.webConfig[host] + pathName ,
        url: useRelativeHost ? pathName : config.webConfig[host] + pathName  ,
        method: 'POST',
        headers: header,
        formData: data,
    };
    // console.log('transform=url=',options.method, options.url)
    request(options, function (error, response, body) {
       
        if (error) {
            console.log('err', error)
            callback('{"status": "err","msg": "服务器请求错误"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    let datas = JSON.parse(body)
                    callback(body);
                }else {
                    console.log('body==', body)
                }
            } catch (err) {
                console.log('body==', body)
                console.log(err);

                let obj = {
                    status: "error",
                    msg: "服务响应错误",
                    body: body
                };
                callback(JSON.stringify(obj));
            }
        }
    });
};

exports.put = function (pathName, data, header, callback, host='requestDomain', useRelativeHost) {
    var options = {
        // url: isStatic ? staticDomain + pathName : baseHost + pathName,
        url: useRelativeHost ? pathName : config.webConfig[host] + pathName  ,
        method: 'PUT',
        headers: header,
        body: JSON.stringify(data),
    };

    // console.log('put==',options)
    request(options, function (error, response, body) {
     
        if (error) {
            console.log(error);
            callback('{"status": "err","msg": "服务器请求错误"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    callback(body);
                }else {
                    console.log('body==', body)
                }
            } catch (err) {
                console.log('body==', body)
                console.log(err);

                let obj = {
                    status: "error",
                    msg: "服务响应错误",
                    body: body
                };
                callback(JSON.stringify(obj));
            }
        }
    });
};

exports.delete = function (pathName, data, header, callback, host='requestDomain', useRelativeHost) {
    var options = {
        // url: isStatic ? staticDomain + pathName : baseHost + pathName,
        url: useRelativeHost ? pathName : config.webConfig[host] + pathName  ,
        method: 'DELETE',
        headers: header,
        body: JSON.stringify(data),
    };
    // console.log('transform=url=',options.method, options.url)
    request(options, function (error, response, body) {
      
        if (error) {
            console.log(error);
            callback('{"status": "err","msg": "服务器请求错误"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    callback(body);
                }else {
                    console.log('body==', body)
                }
            } catch (err) {
                console.log('body==', body)
                console.log(err);

                let obj = {
                    status: "error",
                    msg: "服务响应错误",
                    body: body
                };
                callback(JSON.stringify(obj));
            }
        }
    });
};