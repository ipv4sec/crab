const express = require('express')
const router = express.Router()
const request = require('../utils/request.js')
const fs = require('fs')
const path = require('path')

// 登录 废弃
// router.post('/user/login', (req, res) => {
//     request.post('/user/login',req.body, req.headers, function(response) {
//         res.set(response.headers)
//         res.send(response.data)
//     })
// })

// 登录
router.post('/user/login', (req, res) => {
    request.get('/api/user/'+req.body.userName, '', req.headers, function(response) {
        if(response.data.code === 0 && response.data.result.password === req.body.password) {
            res.send({
                code: 0,
                result: '登陆成功'
            })
        }else {
            res.send(response)
        }
    })
})

// 获取管理员信息
router.get('/user/root', (req, res) => {
    request.get('/user/root', '', req.headers, function(response) {
        res.set(response.headers)
        res.send(response.data)
    })
})


// 修改管理员密码
router.post('/user/root', (req, res) => {
    request.put('/user/root', req.body, req.headers, function(response) {
        res.set(response.headers)
        res.send(response.data)
    })
})


// 获取应用列表
router.get('/app/list', (req, res) => {
    request.get('/app', {limit: req.query.limit, offset: req.query.offset}, req.headers, function(response) {
        res.set(response.headers)
        res.send(response.data)
    })
})


// 运行实例
router.post('/app/run', (req, res) => {
    let newData = {
        dependencies: req.body.dependencies,
        userconfig: req.body.userconfig,
        status: req.body.status
    }
    request.put('/app/'+req.body.instanceid, newData, req.headers, function(response) {
        res.set(response.headers)
        res.send(response.data)
    })
})

// 查看日志
router.get('/app/logs', (req, res) => {
    request.get('/app/'+req.query.id+'/status','', req.headers, function(response) {
        res.set(response.headers)
        res.send(response.data)
    })
})


// 删除实例
router.get('/delete/instance', (req, res) => {
    request.del('/app/'+req.query.id, '', req.headers, function(response) {
        res.set(response.headers)
        res.send(response.data)
    })
})

// 获取实例详情
router.get('/app/detail', (req, res) => {
    request.get('/app/'+req.query.id,'', req.headers, function(response) {
        res.set(response.headers)
        res.send(response.data)
    })
})

// 导出实例配置文件
router.get('/app/output', (req, res) => {
    request.get('/app/'+req.query.id,'', req.headers, function(response) {
        const dirPath = path.resolve(__dirname,'../tempfiles')
        const filePath = path.resolve(__dirname,'../tempfiles/'+response.data.id+'.yaml')
        fs.access(dirPath, (err) => {
            console.log('=====',err)
            if(err) {
                fs.mkdirSync(dirPath,{ recursive: true })
            }
    
            fs.writeFileSync(filePath, response.data.deployment, 'utf8')
    
            res.download(filePath)
        
            fs.rm(filePath, (err) => {
                console.log('---remove file error ---', err)
            })
        })
    })
})

// 获取工作负载源
router.get('/cluster/mirror', (req, res) => {
    request.get('/cluster/mirror', '', req.headers, function(response) {
        res.set(response.headers)
        res.send(response.data)
    })
})

// 设置工作负载源
router.post('/cluster/mirror', (req, res) => {
    request.put('/cluster/domain',req.body, req.headers, function(response) {
        res.set(response.headers)
        res.send(response.data)
    })
})



// 获取节点地址
router.get('/cluster/addrs', (req, res) => {
    request.get('/cluster/addrs','', req.headers, function(response) {
        res.set(response.headers)
        res.send(response.data)
    })
})

// 获取根域
router.get('/cluster/domain', (req, res) => {
    request.get('/cluster/domain','', req.headers, function(response) {
        res.set(response.headers)
        res.send(response.data)
    })
})

// 修改根域
router.post('/cluster/domain', (req, res) => {
    request.put('/cluster/domain',req.body,req.headers, function(response) {
        res.set(response.headers)
        res.send(response.data)
    })
})





// 获取磁盘列表
router.get('/cluster/storage', (req, res) => {
    request.get('/cluster/storage','',req.headers, function(response) {
        res.set(response.headers)
        res.send(response.data)
    })
})

// 修改磁盘
router.post('/cluster/storage', (req, res) => {
    request.put('/cluster/storage',req.body,req.headers, function(response) {
        res.set(response.headers)
        res.send(response.data)
    })
})

// 获取菜单
router.get('/cluster/menus', (req, res) => {
    request.get('/cluster/menus','',req.headers, function(response) {
        res.set(response.headers)
        res.send(response.data)
    })
})



module.exports = router