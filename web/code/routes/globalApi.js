var express = require('express');
var request = require('request');
var router = express.Router();
var config = require('../config/config.json');

// var global = config.globalDomain;
var global = config.webConfig.globalDomain;
var e1 = config.e1Api;
var headers = {
    'Content-Type': 'application/json',
    'C-Uid': config.appConfig.dev_header.c_uid || '',
    'C-Name': config.appConfig.dev_header.c_name || '',
    'C-Firstname': config.appConfig.dev_header.c_firstname ||  '',
    'C-Lastname': config.appConfig.dev_header.c_lastname || '',
    'C-Avatar': config.appConfig.dev_header.c_avatar || '',
    'C-Email': config.appConfig.dev_header.c_email || '',
    'C-Domain': config.appConfig.dev_header.c_domain || '',
    'C-Cs-Id': config.appConfig.dev_header.c_cs_id || '',
    'C-Groups': config.appConfig.dev_header.c_groups || [],
    'C-Roles': config.appConfig.dev_header.c_roles || [],
    'C-Instance-Id': config.appConfig.dev_header.c_instance_id || '',
};


var getGlobal = function (path, req, callback) {
    request({
        url: path,
        method: 'GET',
        headers: {
            // "Cookie": `PHPSESSID=${config.appConfig.dev_header.phpsessid}`
            "Cookie": `PHPSESSID=${config.webConfig.phpsessid}`

            
        }
    }, function (error, response, body) {
        if (error) {
            console.log(error);
            callback('{"status": "error","msg": "服务器请求错误"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    callback(body);
                }
            } catch (err) {
                console.log(body);
                callback('{"status": "error","msg": "服务器请求错误"}');
            }
        }
    });
};
var deleteGlobal = function (path, req, data, callback) {
    request({
        url: path,
        method: 'delete',
        headers: {
            "Cookie": `PHPSESSID=${config.appConfig.dev_header.phpsessid}`
        },
        body: JSON.stringify(data)
    }, function (error, response, body) {
        if (error) {
            console.log(error);
            callback('{"status": "error","msg": "服务器请求错误"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    callback(body);
                }
            } catch (err) {
                console.log(body);
                callback('{"status": "error","msg": "服务器请求错误"}');
            }
        }
    });
};
var postGlobal = function (path, req, data, callback) {
    request({
        url: path,
        method: 'post',
        headers: {
            'Content-Type': 'application/json',
            "Cookie": `PHPSESSID=${config.appConfig.dev_header.phpsessid}`
        },
        body: JSON.stringify(data)
    }, function (error, response, body) {
        if (error) {
            console.log(error);
            callback('{"status": "error","msg": "服务器请求错误"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    callback(body);
                }
            } catch (err) {
                console.log(body);
                callback('{"status": "error","msg": "服务器请求错误"}');
            }
        }
    });
};

var get = function (path, req, callback) {
    request({
        url: path,
        method: 'GET',
        headers: headers,
    }, function (error, response, body) {
        if (error) {
            console.log(error);
            callback('{"status": "error","msg": "服务器请求错误"}');
        } else {
            try {
                if (JSON.parse(body)) {
                    callback(body);
                }
            } catch (err) {
                console.log(body);
                callback('{"status": "error","msg": "服务器请求错误"}');
            }
        }
    });
};

router.get('/global/instance', (req, res) => {
    var admin = '';
    if(req.query.admin){
        admin = `?admin=${req.query.admin}`;
    }
    getGlobal(`${global}/global/instance${admin}`, req, (data) => {
        res.json(JSON.parse(data));
    })
});
router.get('/global/user', (req, res) => {
    getGlobal(`${global}/global/user`, req, (data) => {
        res.json(JSON.parse(data));
    })
});
router.get('/global/tree', (req, res) => {
    getGlobal(`${global}/global/tree`, req, (data) => {
        res.json(JSON.parse(data));
    })
});
router.get('/global/field', (req, res) => {
    let type = req.query.type || '';
    let params = '';

    if(Array.isArray ? Array.isArray(type) : Object.prototype.toString.call(type) === '[object Array]'){
        type.map((d) => {
            params += `&type[]=${d}`;
        });
        params = params.substr(1);
    }else{
        params = `type=${type}`;
    }

    getGlobal(`${global}/global/field?${params}`, req, (data) => {
        res.json(JSON.parse(data));
    })
});
router.get('/global/attr', (req, res) => {
    var match = '', node = '', path = '';
    if(req.query.match){
        match = `match=${req.query.match}`;
    }
    if(req.query.node){
        node = `node=${req.query.node}`;
    }
    if(match && node){
        path = `?node=${node}&match=${match}`;
    }else if(match && !node){
        path = `?match=${match}`;
    }else if(node && !match){
        path = `?node=${node}`;
    }
    getGlobal(`${global}/global/attr${path}`, req, (data) => {
        res.json(JSON.parse(data));
    })
});
router.get('/global/plan/attr', (req, res) => {
    var plan_id = req.query.plan_id || '';
    getGlobal(`${global}/global/plan/attr?plan_id=${plan_id}`, req, (data) => {
        res.json(JSON.parse(data));
    })
});
router.get('/global/plan/field', (req, res) => {
    var plan_id = req.query.plan_id || '';
    getGlobal(`${global}/global/plan/field?plan_id=${plan_id}`, req, (data) => {
        res.json(JSON.parse(data));
    })
});
router.get('/global/plan/flow', (req, res) => {
    var plan_id = req.query.plan_id || '';
    getGlobal(`${global}/global/plan/flow?plan_id=${plan_id}`, req, (data) => {
        res.json(JSON.parse(data));
    })
});
router.get('/global/user/query', (req, res) => {
    getGlobal(`${global}/global/user/query`, req, (data) => {
        res.json(JSON.parse(data));
    })
});
router.delete('/global/user/query', (req, res) => {
    deleteGlobal(`${global}/global/user/query`, req, req.body, (data) => {
        res.json(JSON.parse(data));
    })
});
router.post('/global/user/query', (req, res) => {
    postGlobal(`${global}/global/user/query`, req, req.body, (data) => {
        res.json(JSON.parse(data));
    })
});
router.get('/global/flow', (req, res) => {
    getGlobal(`${global}/global/flow?plan_id=${req.query.plan_id}`, req, (data) => {
        res.json(JSON.parse(data));
    })
});
router.get('/global/steps', (req, res) => {
    getGlobal(`${global}/global/steps`, req, (data) => {
        res.json(JSON.parse(data));
    })
});
router.delete('/global/article', (req, res) => {
    var aid = req.query.aid,
        ver = req.query.ver;
    var param = '';
    if(aid && ver){
       param = `?aid=${aid}&ver=${ver}`;
    }else if(aid && !ver){
        param = `?aid=${aid}`;
    }
    deleteGlobal(`${global}/global/article${param}`, req, '', (data) => {
        res.json(JSON.parse(data));
    })
});
router.post('/global/article', (req, res) => {
    postGlobal(`${global}/global/article`, req, req.body, (data) => {
        res.json(JSON.parse(data));
    })
});
router.get('/global/article', (req, res) => {
    var aid = req.query.aid ? req.query.aid : '';
    var ver = req.query.ver ? req.query.ver : '';
    var sql = req.query.sql ? req.query.sql : '';
    var self = req.query.self ? req.query.self : 0;
    var group = req.query.group ? req.query.group : 0;
    var params = '';
    if(sql){
        params = `?sql=${sql}&self=${self}&group=${group}`;
    }else{
        if(aid && !ver){
            params = `?aid=${aid}`;
        }else if(aid && ver){
            params = `?aid=${aid}&ver=${ver}`;
        }
    }
    getGlobal(`${global}/global/article${params}`, req, (data) => {
        res.json(JSON.parse(data));
    })
});

router.get('/global/field', (req, res) => {
    getGlobal(`${global}/global/field`, req, (data) => {
        res.json(JSON.parse(data));
    })
});

router.get('/global/messages', (req, res) => {
    let page = req.query.page || 1;
    getGlobal(`${global}/global/messages?page=${page}`, req, (data) => {
        res.json(JSON.parse(data));
    })
});
router.post('/global/message/read', (req, res) => {
    let data = {id: req.body.id};
    postGlobal(`${global}/global/message/read`, req, data, (data) => {
        res.json(JSON.parse(data));
    })
});
router.get('/global/languages', (req, res) => {
    getGlobal(`${global}/global/languages`, req, (data) => {
        res.json(JSON.parse(data));
    })
});

router.get('/global/field', (req, res) => {
    getGlobal(`${global}/global/field`, req, (data) => {
        res.json(JSON.parse(data));
    })
});


module.exports = router;
