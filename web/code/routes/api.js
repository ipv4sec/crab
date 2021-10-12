/**
 * Created by zhengliuyang on 16/11/7.
 */
var express = require('express');
var router = express.Router();
var utils = require('../server/utils');
var token = require('../server/token');
var config = require('../config/config.json');
var sso = require('../server/sso');
var search = require('../server/search');
var approval = require('../server/approval');
var schema = require('../server/schema');
// var i1 = require('../server/i1');

router.use((req, res, next) => {
    next();
    // var mtoken = '';
    // var pathname = req['_parsedUrl'].pathname;
    // if(pathname == '/oidc') {
    //     next();
    // }else{
    //     if (req.method == 'POST') {
    //         mtoken = req.body.token || '';
    //         if (req.header('token')) {
    //             mtoken = req.header('token');
    //         }
    //     } else if (req.method == 'GET') {
    //         mtoken = req.query.token || '';
    //     }
    //     if (mtoken) {
    //         mtoken = mtoken.replace(/\s/g, '\+');
    //         mtoken = decodeURIComponent(mtoken);
    //     }
    //     if (token.checkToken(req, mtoken)) {
    //         next();
    //     } else {
    //         var _token = token.getToken(req);
    //         res.cookie('m1_token', _token, {
    //             maxAge: config.webConfig.cookieMaxAge,
    //             httpOnly: false,
    //             path: '/',
    //             secure: false
    //         });
    //         res.json({'status': 'error', 'msg': '非法请求，请刷新页面重试'});
    //     }
    // }
});

/**
 * OIDC
 */
router.get('/oidc', (req, res) => {
    sso.oidc(req, res);
});

router.get('/loginout', (req, res) => {
    sso.loginlout(req, res);
});

/**
 * 搜索
 */
router.get('/search', (req, res) => {
    search.search(req, res);
});
/**
 * 获取以保存的条件
 */
router.get('/search/condition', (req, res) => {
    search.getCondition(req, res);
});
/**
 * 删除以保存的条件
 */
router.get('/search/condition/delete', (req, res) => {
    search.deleteCondition(req, res);
});
/**
 * 保存条件
 */
router.post('/search/condition/save', (req, res) => {
    search.saveCondition(req, res);
});
router.post('/search/condition/update', (req, res) => {
   search.updateCondition(req, res);
});

router.post('/search/delete', (req, res) => {
   search.deleteArticle(req, res);
});

router.get('/search/export', (req, res) => {
    search.exportResult(req, res);
});


/*******************审稿*******************/
router.get('/approval/list', (req, res) => {
    approval.list(req, res);
});
router.get('/approval/articleDetail', (req, res) => {
    approval.articleDetail(req, res);
});
router.get('/approval/nearPage', (req, res) => {
    approval.nearPage(req, res);
});
router.post('/approval/comment', (req, res) => {
    approval.comment(req, res);
});
router.post('/approval/pass', (req, res) => {
    approval.pass(req, res);
});
router.post('/approval/del', (req, res) => {
    approval.del(req, res);
});
router.post('/approval/comment/edit', (req, res) => {
    approval.commentEdit(req, res);
});
router.post('/approval/comment/del', (req, res) => {
    approval.commentDel(req, res);
});
router.post('/approval/marker', (req, res) => {
    approval.marker(req, res);
});
/*******************schema*******************/
/**
 * 搜索schema
 */
router.get('/schema/search', (req, res) => {
    schema.search(req, res);
});
/**
 * 获取schema详情
 */
router.get('/schema/detail', (req, res) => {
    schema.detail(req, res);
});
/**
 * 修改单个schema状态
 */
router.get('/schema/status', (req, res) => {
    schema.status(req, res);
});
/**
 * 修改多个schema状态
 */
router.post('/schema/status', (req, res) => {
    schema.multiStatus(req, res);
});
/**
 * 删除单个schema状态
 */
router.get('/schema/delete', (req, res) => {
    schema.delete(req, res);
});
/**
 * 删除多个schema状态
 */
router.post('/schema/delete', (req, res) => {
    schema.multiDelete(req, res);
});
/**
 * 添加schema
 */
router.post('/schema/add', (req, res) => {
    schema.add(req, res);
});
/**
 * 更新schema
 */
router.post('/schema/update', (req, res) => {
    schema.update(req, res);
});
/**
 * 获取当前对应的表名字
 */
router.get('/schema/table', (req, res) => {
    schema.table(req, res);
});
/**
 * 获取schema浏览记录
 */
router.get('/schema/history', (req,res) => {
    schema.getHistory(req, res);
});
/**
 * 保存schema浏览记录
 */
router.post('/schema/history/save', (req, res) => {
    schema.saveHistory(req, res);
});
/**
 * 删除schema浏览记录
 */
router.get('/schema/history/delete', (req, res) => {
   schema.deleteHistory(req, res);
});





// router.get('/i1/user', (req, res) => {
//     i1.getPeople(req, res);
// });


// router.get('/permission', (req, res) => {
//     i1.permission(req, res);
// });

// router.get('/i1/apps', (req, res) => {
//     i1.apps(req, res);
// });

/******************* 获取用户有哪些应用权限 ******************/
router.get('/userAppsInfo', (req, res) => {
    utils.tokenGet('http://api.localdomain/portal/getSuiteEntrance?cs_id='+req.headers['c-cs-id']+'&uid='+req.headers['c-uid']+'&i1_domain='+req.headers['c-domain'],req, '', (body) => {
        res.json(JSON.parse(body));
    })
});

module.exports = router;
