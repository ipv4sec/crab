/**
 * Created by zhengliuyang on 16/12/21.
 */
var utils = require('../server/utils');
var config = require('../config/config.json');
var requestDomain = config.webConfig.requestDomain;
var i1Domain = 'https://i1.zlibs.com/api/';
var jwt = require('jsonwebtoken');
var adminConsoleHost = config.webConfig.adminConsoleHost;
module.exports = {
    /**
     * 获取组
     * @param req
     * @param res
     */
    getGroups: (req, res) => {
        var args = req.query;
        var q = args.q || '';
        var offset = args.offset || 0;
        var size = args.size || 100;
        // var path = '/v2/groups?q=' + encodeURIComponent(q) + '&offset=' + offset + '&size=' + size;
        utils.i1Get('/v2/groups', 'offset=' + offset + '&q=' + encodeURIComponent(q)  + '&size=' + size ,req, (body) => {
            res.json(body);
        });
    },
    /**
     * 获取角色
     * @param req
     * @param res
     */
    getRole: (req, res) => {
        var args = req.query;
        var q = args.q || '';
        var offset = args.offset || 0;
        var size = args.size || 100;
        utils.i1Get('/v2/roles', 'offset=' + offset + '&q=' + encodeURIComponent(q) + '&size=' + size ,req, (body) => {
            res.json(body);
        });
    },
    /**
     * 获取用户
     * @param req
     * @param res
     */
    getPeople: (req, res) => {
        var args = req.query;
        var q = args.q || '';
        var offset = args.offset || 0;
        var size = args.size || 100;
        var order = args.order || 'desc';
        var status = args.status || '';
        utils.i1Get('/v2/users', 'offset=' + offset + '&order=' + order + '&q=' + encodeURIComponent(q) + '&size=' + size + '&status=' + status ,req, (body) => {
            res.json(body);
        });
    },
    /**
     * 获取用户权限
     * @param req
     * @param res
     */
    permission: (req, res) => {
        if (req.headers.hasOwnProperty('c-uid') && req.headers.hasOwnProperty('c-domain')) {
                let batchDelete = new Promise((resolve, reject) => {
                    utils.get(`${adminConsoleHost}/v2/api/config/dataPermissions/${encodeURI('全局')}/batchDelete/${req.headers['c-domain']}/${req.headers['c-uid']}`,req, (body) => {
                        body = JSON.parse(body);
                        resolve(body);
                        // if(body.status == 'success') {
                        //
                        // }else{
                        //     reject(body);
                        // }
                    })
                });
                let schemaManage = new Promise((resolve, reject) => {
                    // console.log(`http://119.254.147.206:2080/api/config/dataPermissions/${encodeURI('全局')}/schemaManage/${req.headers['c-domain']}/${req.headers['c-uid']}`)
                    utils.get(`${adminConsoleHost}/v2/api/config/dataPermissions/${encodeURI('全局')}/schemaManage/${req.headers['c-domain']}/${req.headers['c-uid']}`,req, (body) => {
                        body = JSON.parse(body);
                        resolve(body);
                        // if(body.status == 'success') {
                        //
                        // }else{
                        //     reject(body);
                        // }
                    })
                });
                Promise.all([batchDelete,schemaManage]).then(([b,s]) => {
                    let r = '';
                    if(b.status == 'success' && b.data) {
                        r += 'delete_article,';
                    }
                    if(s.status == 'success' && s.data) {
                        r += 'write_schema';
                    }
                    r = 'delete_article,write_schema';
                    res.json({code: 200, results: r});
                }).catch(function (err) {
                    // console.log(err)
                    res.json({code: 'error'});
                });
        } else {
            res.json({status: 'error', msg: 'get user info failed !'});
        }
    },
    apps: (req, res) => {
        var args = req.query;
        // console.log(req.session.m1redisToken)
        utils.tokenGet(i1Domain + '/v1/user/apps',req, req.session.m1redisToken, (body) => {
            res.json(JSON.parse(body));
        });
    },
};