var express = require('express');
var router = express.Router();
var config = require('../config/config.json');
var back = require('../server/util/back');
var headers = require('../server/util/header');
 
var fs = require('fs');
var multiparty = require('multiparty');
var path = require('path')
var FormData = require('form-data')
var request = require('request');


// 上传form表单，或者以表单形式上传单文件
//添加应用文件
router.post('/addFile', function (req, res, next) {
    let header = headers.getHeader(req);

    if (!fs.existsSync('./tempFiles')) {
        fs.mkdir('./tempFiles', () => {});
    }

    var form = new multiparty.Form({uploadDir: './tempFiles'});
    form.parse(req, function (err, fields, files) {
        let filePath = files.file[0].path,
        fileName = files.file[0].originalFilename;
         if(fields.folder_name) {
             let fileInfo = filePath.split('/')
             let names =  fileName.split('/')
            fileName = names[names.length - 1]
            filePath = 'tempFiles/'+fileInfo[fileInfo.length - 1]
        }
        const newPath = path.join(path.dirname(filePath), fileName) // 得到newPath新地址用于创建读取流
        fs.renameSync(filePath, newPath)
        let file = fs.createReadStream(newPath)
        let data = {
           file: file
        }
        header = Object.assign({}, header, {'Content-Type': 'multipart/form-data'})

    
        back.postForm('/api/app', data, header, function (body) {
            if (fs.existsSync(newPath)) {
                fs.unlink(newPath, (err) => {})
            }
            res.send(body);
        });

    })
});

// 运行应用
router.post('/run', function (req, res, next) {
    let header = headers.getHeader(req);

    // header = Object.assign({}, header, {'Content-Type': 'multipart/form-data'})
    back.put('/api/app', req.body, header, function (body) {
        res.send(body);
    });
});

// 卸载应用
router.post('/uninstallApp', function (req, res, next) {
    let header = headers.getHeader(req);

    // header = Object.assign({}, header, {'Content-Type': 'multipart/form-data'})
    back.put('/api/app', req.body, header, function (body) {
        res.send(body);
    });
});


// 删除应用
router.get('/deleteApp', function (req, res, next) {
    let header = headers.getHeader(req);

    back.deleteForm('/api/app?instanceid='+req.query.instanceid, header, function (body) {
        res.send(body);
    });
});



//获取应用列表
router.get('/getAllApp', function (req, res, next) {
    let header = headers.getHeader(req);
    let queries = ''
    if(req.query.limit) {
        queries = `?offset=${req.query.offset}&limit=${req.query.limit}`
    }
    back.get(`/api/app${queries}`, header, function (body) {
        res.send(body);
    });
});


// //获取所有空间
// router.get('/getAllSpace', function (req, res, next) {
//     let header = headers.getHeader(req);
//     let queries = ''
//     if(req.query.limit) {
//         queries = `?offset=${req.query.offset}&limit=${req.query.limit}`
//     }
//     back.get(`/namespace${queries}`, header, function (body) {
//         res.send(body);
//     });
// });

// // 添加空间
// router.post('/addSpace', function (req, res, next) {
//     let header = headers.getHeader(req);
//     header = Object.assign({}, header, {'Content-Type': 'multipart/form-data'})

//     back.postForm('/namespace', req.body, header, function (body) {
//         res.send(body);
//     });
// });

// // 空间详情
// router.get('/getSpaceDetail', function (req, res, next) {
//     let header = headers.getHeader(req);

//     back.get(`/namespace?id=${req.query.id}`, header, function (body) {
//         res.send(body);
//     });
// });


// // 空间详情
// router.post('/saveBranch', function (req, res, next) {
//     let header = headers.getHeader(req);
//     header = Object.assign({}, header, {'Content-Type': 'multipart/form-data'})
//     back.putForm('/app/instance', req.body, header, function (body) {
//         res.send(body);
//     });
// });


// // 导出空间配置
// router.get('/output', function (req, res, next) {
//     let header = headers.getHeader(req);
//     back.getRes(`/namespace/deployment_yml?namespace_id=${req.query.spaceId}`, header, function (body) {
//        try{
//             fs.writeFile(path.join(__dirname, '../tempFiles/deployment_yml'), body, (err) => {
             
//                 if(err) {
//                     console.log(err)
//                     return
//                 }
//                 res.download(path.join(__dirname, '../tempFiles/deployment_yml'));
//                 // res.end();
//             })
           
//         }catch(err) {
//             res.send({ 
//                 status: "error",
//                 msg: "服务响应错误",
//                 body: body}
//             );
//         }
      
//     });
// });


// // 获取空间下应用
// router.get('/getSpaceApps', function (req, res, next) {
//     let header = headers.getHeader(req);
//     back.get(`/app/instances?namespace_id=${req.query.spaceId}`, header, function (body) {
//         res.send(body);
//     });
// });


// // 获取所有应用
// router.get('/getAppList', function (req, res, next) {
//     let header = headers.getHeader(req);

//     back.get('/app/list', header, function (body) {
//         res.send(body);
//     });
// });


// // 获取app某个版本的配置信息
// router.get('/getAppConfig', function (req, res, next) {
//     let header = headers.getHeader(req);
//     back.get(`/getAppConfig?appId=${req.query.appId}&version=${req.query.version}`, header, function (body) {
//         res.send(body);
//     });
// });


// //删除空间
// router.post('/deleteSpace', function (req, res, next) {
//     let header = headers.getHeader(req);
//     header = Object.assign({}, header, {'Content-Type': 'multipart/form-data'})

//     back.postForm('/deleteSpace/', req.body, header, function (body) {
//         res.send(body);
//     });
// });






// //添加应用
// router.post('/addApp', function (req, res, next) {
//     let header = headers.getHeader(req);
//     header = Object.assign({}, header, {'Content-Type': 'multipart/form-data'})
  
//     back.postForm('/app', req.body, header, function (body) {
//         res.send(body);
//     });
// });

// //修改应用
// router.post('/editApp', function (req, res, next) {
//     let header = headers.getHeader(req);
//     header = Object.assign({}, header, {'Content-Type': 'multipart/form-data'})
  
//     back.putForm('/app', req.body, header, function (body) {
//         res.send(body);
//     });
// });


// //获取应用下所有版本
// router.get('/getAppVersions', function (req, res, next) {
//     let header = headers.getHeader(req);
//     back.get(`/app/versions?app_id=${req.query.appId}`, header, function (body) {
//         res.send(body);
//     });
// });

// //获取应用详情
// router.get('/appDetail', function (req, res, next) {
//     let header = headers.getHeader(req);
//     let query = ''
//     if(req.query.name){
//         query = `?name=${req.query.name}`
//     }else if(req.query.appId){
//         query = `?id=${req.query.appId}`
//     }
//     back.get('/app'+query, header, function (body) {
//         res.send(body);
//     });
// });

// //获取应用下某个版本的详情
// router.get('/getVersionDetail', function (req, res, next) {
//     let header = headers.getHeader(req);
//     back.get(`/app/version_info?app_id=${req.query.appId}&version_id=${req.query.version}`, header, function (body) {
//         res.send(body);
//     });
// });


// //获取集群列表
// router.get('/getAllCluster', function (req, res, next) {
//     let header = headers.getHeader(req);
//     let queries = ''
//     if(req.query.limit) {
//         queries = `?offset=${req.query.offset}&limit=${req.query.limit}`
//     }
//     back.get(`/cluster${queries}`, header, function (body) {
//         res.send(body);
//     });
// });



// //添加集群文件
// router.post('/addFile', function (req, res, next) {
//     let header = headers.getHeader(req);

//     if (!fs.existsSync('./tempFiles')) {
//         fs.mkdir('./tempFiles', () => {});
//     }

//     var form = new multiparty.Form({uploadDir: './tempFiles'});
//     form.parse(req, function (err, fields, files) {
//         let filePath = files.file[0].path,
//         fileName = files.file[0].originalFilename;
//          if(fields.folder_name) {
//              let fileInfo = filePath.split('/')
//              let names =  fileName.split('/')
//             fileName = names[names.length - 1]
//             filePath = 'tempFiles/'+fileInfo[fileInfo.length - 1]
//         }
//         const newPath = path.join(path.dirname(filePath), fileName) // 得到newPath新地址用于创建读取流
//         fs.renameSync(filePath, newPath)
//         let file = fs.createReadStream(newPath)
//         let formData = new FormData()
//         formData.append('file', file, encodeURI(fileName))
//         formData.append('name', fields.name[0])
//         formData.append('host', fields.host[0])
//         formData.append('other', fields.other[0])
//         header = Object.assign({}, header, formData.getHeaders(), {'Content-Type': false, processData: false})
//         back.postFile('/cluster', formData, header, function (body) {
//             if (fs.existsSync(newPath)) {
//                 fs.unlink(newPath, (err) => {})
//             }
//             res.send(body);
//         });

//     })
// });

// //添加集群
// router.post('/addCluster', function (req, res, next) {
//     let header = headers.getHeader(req);

//     if (!fs.existsSync('./tempFiles')) {
//         fs.mkdir('./tempFiles', () => {});
//     }

//     var form = new multiparty.Form({uploadDir: './tempFiles'});
//     form.parse(req, function (err, fields, files) {
//         let filePath = files.file[0].path,
//         fileName = files.file[0].originalFilename;
//          if(fields.folder_name) {
//              let fileInfo = filePath.split('/')
//              let names =  fileName.split('/')
//             fileName = names[names.length - 1]
//             filePath = 'tempFiles/'+fileInfo[fileInfo.length - 1]
//         }
//         const newPath = path.join(path.dirname(filePath), fileName) // 得到newPath新地址用于创建读取流
//         fs.renameSync(filePath, newPath)
//         let file = fs.createReadStream(newPath)
//         let data = {
//             name: fields.name[0],
//             // host: fields.host[0],
//             url: fields.host[0],
//             description: fields.other[0],
//             // certificate_url: file
//             token_file: file
//         }
//         header = Object.assign({}, header, {'Content-Type': 'multipart/form-data'})
//         back.postForm('/cluster', data, header, function (body) {
//             if (fs.existsSync(newPath)) {
//                 fs.unlink(newPath, (err) => {})
//             }
//             res.send(body);
//         });

//     })
// });

// //修改集群
// router.post('/editCluster', function (req, res, next) {
//     let header = headers.getHeader(req);

//     if (!fs.existsSync('./tempFiles')) {
//         fs.mkdir('./tempFiles', () => {});
//     }

//     var form = new multiparty.Form({uploadDir: './tempFiles'});
//     form.parse(req, function (err, fields, files) {
//         let filePath = files.file[0].path,
//         fileName = files.file[0].originalFilename;
//          if(fields.folder_name) {
//              let fileInfo = filePath.split('/')
//              let names =  fileName.split('/')
//             fileName = names[names.length - 1]
//             filePath = 'tempFiles/'+fileInfo[fileInfo.length - 1]
//         }
//         const newPath = path.join(path.dirname(filePath), fileName) // 得到newPath新地址用于创建读取流
//         fs.renameSync(filePath, newPath)
//         let file = fs.createReadStream(newPath)
//         let data = {
//             id: fields.id[0],
//             name: fields.name[0],
//             // host: fields.host[0],
//             url: fields.host[0],
//             description: fields.other[0],
//             // certificate_url: file
//             token_file: file
//         }
//         header = Object.assign({}, header, {'Content-Type': 'multipart/form-data'})
       
//         back.putForm('/cluster', data, header, function (body) {
//             if (fs.existsSync(newPath)) {
//                 fs.unlink(newPath, (err) => {})
//             }
//             res.send(body);
//         });

//     })
// });


// 登录
router.post('/login', function (req, res, next) {
    // console.log('login==',req.body)
    let header = headers.getHeader(req);
  
    var options = {
        url:config.webConfig.requestDomain +'/api/user/'+ req.body.userName,
        method: 'GET',
        headers: header,
    };

    request(options, function (error, response, body) {
        if (error) {
            console.log(error);
            callback('{"status": "err","msg": "服务器请求错误"}');
        } else {

            console.log('body====',body)
            try {
                let parseBody = JSON.parse(body)
                if(parseBody.password === req.body.password) {
                    res.send(JSON.stringify({
                        status: 'success',
                        msg: '登录成功'
                    }));
                }else {
                    res.send(JSON.stringify({
                        status: 'error',
                        msg: '登录失败'
                    }));
                }
            } catch (err) {
                console.log('body==', body)
                console.log(err);
                let obj = {
                    status: "error",
                    msg: "服务响应错误",
                    body: body
                };
                res.send(JSON.stringify(obj));
                return
            }
        }
    });
});



module.exports = router;
