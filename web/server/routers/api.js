const express = require('express')
const router = express.Router()
const request = require('../utils/request.js')
const fs = require('fs')
const path = require('path')
const multiparty = require('multiparty')
const FormData = require('form-data')
const archiver = require('archiver')

// 登录
router.get('/user/login', (req, res) => {
    request.get('/user/'+req.query.username, '', req.headers, function(response) {
        if(response.data.code === 0 && response.data.result.password === req.query.password) {
            // req.session['user'] = req.query.username 
            res.send({
                code: 0,
                result: {
                    username: req.query.username,
                    message: '登陆成功'
                }
            })
        }else {
            res.send({code: 500, result: '登陆失败'})
        }
    })
    
})

// 退出登陆
router.get('/user/logout', (req, res) => {
    // delete req.session['user']
    res.send({
        code: 40404,
        result: '退出登陆'
    })
})



// 获取管理员信息
router.get('/user/root', (req, res) => {
    request.get('/user/root', '', req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})


// 修改管理员密码
router.post('/user/reset', (req, res) => {
    request.put('/user/root', req.body, req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 添加实例
router.post('/app/upload', (req, res) => {
    
    const dirPath = path.resolve(__dirname,'../tempfiles')

    if (!fs.existsSync(dirPath)) {
        fs.mkdirSync(dirPath,{ recursive: true })
    }

    var form = new multiparty.Form({uploadDir: dirPath});

    form.parse(req, function (err, fields, files) {
        let filePath = files.file[0].path,
        fileName = files.file[0].originalFilename;
         if(fields.folder_name) {
             let fileInfo = filePath.split('/')
             let names =  fileName.split('/')
            fileName = names[names.length - 1]
            filePath = 'tempfiles/'+fileInfo[fileInfo.length - 1]
        }
        const newPath = path.join(path.dirname(filePath), fileName) // 得到newPath新地址用于创建读取流
        fs.renameSync(filePath, newPath)
        let file = fs.createReadStream(newPath)
        // console.log(file)
        let formData = new FormData()
        formData.append('file', file)
        let headers = formData.getHeaders()
        let header = Object.assign({}, headers)
        request.postForm('/app', formData, header, function(response) {
            if (fs.existsSync(newPath)) {
                fs.unlink(newPath, (err) => {})
            }
            // // res.set(response.headers)
            res.send(response.data)
        })

    })
})


// 获取应用列表
router.get('/app/list', (req, res) => {
    request.get('/app', {limit: req.query.limit, offset: req.query.offset}, req.headers, function(response) {
       
        // // res.set(response.headers)
        res.send(response.data)
    })
})


// 运行实例
router.post('/app/run', (req, res) => {
    let newData = {
        dependencies: req.body.dependencies,
        userconfigs: req.body.userconfigs,
        status: req.body.status
    }
    request.put('/app/'+req.body.id, newData, req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 查看日志
router.get('/app/logs_bak', (req, res) => {
    request.get('/app/'+req.query.id+'/status','', req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})


// 查看日志
router.get('/app/logs', (req, res) => {
    request.get('/app/'+req.query.id+'/logs','', req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})


// 删除实例
router.get('/delete/instance', (req, res) => {
    request.del('/app/'+req.query.id, '', req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 获取实例详情
router.get('/app/detail', (req, res) => {
    request.get('/app/'+req.query.id,'', req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})


// 获取实例详情
router.get('/app/detail_list', (req, res) => {
    request.get('/app/'+req.query.id,'', req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})



// 导出实例配置文件
router.get('/app/output', (req, res) => {
    request.get('/app/'+req.query.id,'', req.headers, function(response) {
        // console.log('----',response.data)
        if(response.data.code === 0) {
            const dirPath = path.resolve(__dirname,'../tempfiles')
            const filePath = path.resolve(__dirname,'../tempfiles/'+response.data.result.id+'.yaml')
            if (!fs.existsSync(dirPath)) {
                fs.mkdirSync(dirPath,{ recursive: true })
            }
    
            fs.writeFileSync(filePath, response.data.result.deployment, 'utf8')

            res.download(filePath, response.data.result.id+'.yaml', (err) => {
                console.log(err)
            })
           
            fs.rm(filePath, (err) => {
                console.log('---remove file error ---', err)
            })
        }else {
            res.send(response)
        }

    })
})

// 获取工作负载源
router.get('/cluster/mirror', (req, res) => {
    request.get('/cluster/mirror', '', req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 设置工作负载源
router.post('/cluster/mirror', (req, res) => {
    request.put('/cluster/mirror',req.body, req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})


// 获取traitlist 
router.get('/cluster/traitlist', (req, res) => {
    request.get(`/trait?limit=${req.query.limit}&offset=${req.query.offset}`,'',req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 获取workloadTypelist 
router.get('/cluster/workloadlist', (req, res) => {
    request.get(`/workloadType?limit=${req.query.limit}&offset=${req.query.offset}`,'',req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})
// 获取workloadVendorlist 
router.get('/cluster/vendorlist', (req, res) => {
    request.get(`/workloadVendor?limit=${req.query.limit}&offset=${req.query.offset}`,'',req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})


// 修改trait
router.post('/cluster/edittrait', (req, res) => {
    request.put('/trait/'+req.query.id, req.body, req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 修改workloadtype
router.post('/cluster/editworkload', (req, res) => {
    request.put('/workloadType/'+req.query.id, req.body,req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 修改trait
router.post('/cluster/editvendor', (req, res) => {
    request.put('/workloadVendor/'+req.query.id, req.body,req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})


// 删除trait
router.get('/cluster/deletetrait', (req, res) => {
    request.del('/trait/'+req.query.id, req.body,req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 删除workloadtype
router.get('/cluster/deleteworkload', (req, res) => {
    request.del('/workloadType/'+req.query.id, '',req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 删除workloadvendor
router.get('/cluster/deletevendor', (req, res) => {
    request.del('/workloadVendor/'+req.query.id, '',req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})


// 获取节点地址
router.get('/cluster/addrs', (req, res) => {
    request.get('/cluster/addrs','', req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 获取根域
router.get('/cluster/domain', (req, res) => {
    request.get('/cluster/domain','', req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 修改根域
router.post('/cluster/domain', (req, res) => {
    request.put('/cluster/domain',req.body,req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})



// 获取磁盘列表
router.get('/cluster/storage', (req, res) => {
    request.get('/cluster/storage','',req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 修改磁盘
router.post('/cluster/storage', (req, res) => {
    request.put('/cluster/storage',req.body,req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 获取菜单
router.get('/cluster/menus', (req, res) => {
    request.get('/cluster/menus','',req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 下载manifest.yaml文件
router.post('/online/download', (req, res) => {
    const text = req.body.yaml

    const dir = path.join(__dirname, '../tempfiles')
    if(!fs.existsSync(dir)) {
        fs.mkdirSync(dir)
    }

    const filepath = path.join(dir, 'manifest.yaml')

    fs.writeFile(filepath, text, 'utf8', err => {
        if(err) {
            console.log('write manifest.yaml error: ')
            console.log(err)
            res.send({
                code: 50001,
                result: err
            })
            return
        }

        res.send({
            code: 0,
            result: '/api/online/downloadyaml'
        })
    })

})

router.get('/online/downloadyaml', (req, res) => {
    const filepath = path.join(__dirname, '../tempfiles/manifest.yaml')
    res.download(filepath, 'manifest.yaml', err => {
        if(err) {
            if(err) {
                console.log('download manifest.yaml error: ')
                console.log(err)
                res.send({
                    code: 50001,
                    result: err
                })
                return
            }
    
            res.send({
                code: 0,
                result: '下载成功'
            })
        }
    })
})


// 部署
router.post('/online/arrange', async (req, res) => {
    const text = req.body.yaml

    const dir = path.join(__dirname, '../tempfiles')
    if(!fs.existsSync(dir)) {
        fs.mkdirSync(dir)
    }

    const filepath = path.join(dir, 'manifest.yaml')

    fs.writeFileSync(filepath, text, 'utf8')

    let zipRes = await zipYaml()

    if(zipRes === 'success') {
        const newPath = path.join(__dirname, '../tempfiles/manifest.zip')
        let file = fs.createReadStream(newPath)
        // console.log(file)
        let formData = new FormData()
        formData.append('file', file)
        let headers = formData.getHeaders()
        let header = Object.assign({}, headers)
        request.postForm('/app', formData, header, function(response) {
            if (fs.existsSync(newPath)) {
                fs.unlink(newPath, (err) => {})
            }
            // res.set(response.headers)
            res.send(response.data)
        })

    }else {
        res.send({
            code: 50001,
            result: '压缩zip文件失败'
        })
    }

})


// 创建trait
router.post('/online/createtrait', (req, res) => {
    request.post('/trait',req.body,req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})
// 获取trait详情
router.get('/online/gettrait', (req, res) => {
    request.get('/trait/'+req.query.name, '',req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})


// 创建workloadtype
router.post('/online/createworkloadtype', (req, res) => {
    request.post('/workloadType',req.body,req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 获取workloadtype 详情
router.get('/online/getworkloadtype', (req, res) => {
    request.get('/workloadType/'+req.query.name, '',req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})


// 获取workloadVendor的system spec默认数据
router.get('/online/systemspec', (req, res) => {
    request.get('/tool/systemTemplate','',req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 获取getworkloadvendor详情
router.get('/online/getworkloadvendor', (req, res) => {
    request.get('/workloadVendor/'+req.query.name, '',req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})



// workloadvendor 转换yaml 为 cue 格式 
router.post('/online/translateyaml', (req, res) => {
    request.post('/tool/convertion',req.body,req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 检查CUE语法
router.post('/online/checkcue', (req, res) => {
    request.post('/tool/spelling',req.body,req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 创建workloadvendor
router.post('/online/createvendor', (req, res) => {
    request.post('/workloadVendor',req.body,req.headers, function(response) {
        // res.set(response.headers)
        res.send(response.data)
    })
})

// 流水线接口
router.post('/deployment/:id', (req, res) => {
    const dirPath = path.resolve(__dirname,'../tempfiles')
    if (!fs.existsSync(dirPath)) {
        fs.mkdirSync(dirPath, { recursive: true })
    }
    var form = new multiparty.Form({uploadDir: dirPath});
    form.parse(req, function (err, fields, files) {
        let filePath = files.file[0].path,
        fileName = files.file[0].originalFilename;
         if(fields.folder_name) {
             let fileInfo = filePath.split('/')
             let names =  fileName.split('/')
            fileName = names[names.length - 1]
            filePath = 'tempfiles/'+fileInfo[fileInfo.length - 1]
        }
        const newPath = path.join(path.dirname(filePath), fileName) // 得到newPath新地址用于创建读取流
        fs.renameSync(filePath, newPath)
        let file = fs.createReadStream(newPath)
        let formData = new FormData()
        formData.append('file', file)
        let headers = formData.getHeaders()
        let header = Object.assign({}, headers)
        request.putForm('/deployment/'+req.params.id, formData, header, function(response) {
            if (fs.existsSync(newPath)) {
                fs.unlink(newPath, (err) => {})
            }
            res.send(response.data)
        })

    })
})


function zipYaml () {
    return new Promise((resolve, reject) => {
        
        const output = fs.createWriteStream(path.join(__dirname, '../tempfiles/manifest.zip'))

        const archive = archiver('zip', { zlib: { level: 9 } })
    
        
        output.on('end', () => {
            console.log('Data has been drained');
        })

        output.on('close', () => {
            console.log(archive.pointer() + ' total bytes');
            console.log('archiver has been finalized and the output file descriptor has closed.');

            resolve('success')

        })
        
        archive.on('error', function(err) {
            if(err) {
                reject(err)
                console.log(err)
            }
        });
    
        archive.pipe(output);
        archive.file(path.join(__dirname, '../tempfiles/manifest.yaml'), {name: 'manifest.yaml'})
        archive.finalize()
    
    })
}



router.get('/app/testdata', (req, res) => {
    fs.readFile(path.join(__dirname, '../tempfiles/detail.json'), 'utf8', (err, data) => {
        if(err) {
            console.log(err)
            res.send({
                code: 50001,
                result: '获取数据失败'
            })

        }else {
            res.send(data)
        }
    })
})



module.exports = router