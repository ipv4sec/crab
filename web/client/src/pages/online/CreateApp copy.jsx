import React, { useState, useEffect, useRef } from 'react'
import { connect } from 'react-redux'
import Button from '@material-ui/core/Button'
import '../../style/sass/online.scss'
import axios from 'axios'
import store from '../../store/store'
import * as TYPE from '../../store/actions'
import AddFile from '../../components/AddFile'

const defaultMetadata = `apiVersion: aam.globalsphare.com/v1alpha1
kind: Application
metadata:
    name: example
    version: 0.0.1
    description: 样例应用
    keywords:
        - 样例应用
    author: example@example.com
    maintainers:
        - email: example@example.com
          name: example
          web: https://example.com
    repositories: ["https://github.com/example/example.git"]
    bugs: https://github.com/example/example/issues
    licenses:
        - type: LGPL
          url: https://license.spec.com`

const defaultUserconfigs = `"$schema": http://json-schema.org/draft-07/schema#
"$id": http://example.com/product.schema.json
title: User
description: init user description
type: object
properties:
    username:
        type: string
    password:
        type: string
required:
    - username
    - password`

const defaultWorkloads = [`name: example
type: webservice
vendor: webservice
properties:
    image: nginx:1.21
traits:
    - type: ingress
      properties:
          k1: "v1"`]


const defaultDependencies = `- name: gitlab
  version: ">=0.0.1"
  location: user-defined(https://gitlab.com)
  items:
      /*:
          - create
          - read
          - update
          - delete`

const defaultExports = `/user:
    - create
    - read
    - update
    - delete
/admin:
    - create
    - read
    - update
    - delete`


const CreateApp = (props) => {
    const preRef = useRef(null)
    const [metadata, setMetadata] = useState(defaultMetadata)
    const [userconfig, setUserconfig] = useState(defaultUserconfigs)
    const [workloads, setWorkloads] = useState(defaultWorkloads)
    const [dependencies, setDependencies] = useState(defaultDependencies)
    const [exportsData, setExportsData] = useState(defaultExports)

    const [showConfigDialog, setShowConfigDialog] = useState(false)
    const [configData, setConfigData] = useState({})

    const [previewData, setPreviewData] = useState('')

    const changeMetadata = (e) => {
        setMetadata(e.target.value)
    }
    const changeUserconfig = (e) => {
        setUserconfig(e.target.value)
    }

    const changeWorkload = (e) => {
        const idx = e.currentTarget.dataset.index
        const newWorkloads = workloads.slice()
        newWorkloads[idx] = e.target.value
        setWorkloads(newWorkloads)
    }
    const addWorkload = () => {
        const newWorkloads = workloads.slice()
        newWorkloads.push(defaultWorkloads[0])
        setWorkloads(newWorkloads)
    }
    const removeWorkload = (e) => {
        const idx = e.currentTarget.dataset.index
        const newWorkloads = workloads.slice()
        newWorkloads.splice(idx, 1)
        setWorkloads(newWorkloads)
    }

    const changeDependencies = (e) => {
        setDependencies(e.target.value)
    }

    const changeExports = (e) => {
        setExportsData(e.target.value)
    }

    // 预览
    function preview() {
        const reg = /\n/g
        console.log()
        let pData = (
            metadata + '\nspec:' +
            '\n    userconfigs:\n        ' + userconfig.replace(reg, '\n        ') + 
            '\n    workloads:\n        '+workloads.map(item => '- ' + item.replace(reg, '\n          ')).join('\n        ') + 
            '\n    dependencies:\n        '+dependencies.replace(reg, '\n        ') + 
            '\n    exports:\n        '+ exportsData.replace(reg, '\n        ')
        ) 

        setPreviewData(pData)

        preRef.current.innerText = pData
    }

    useEffect(() => {
        preview()
        resizePreHeight()
    }, [metadata, userconfig, workloads, dependencies, exportsData])

    useEffect(() => {
        preview()
        resizePreHeight()
    }, [])

    function resizePreHeight() {
        let leftHeight = document.querySelector('.createapp-left').offsetHeight
        document.querySelector('.preview-pre').style.height = leftHeight - 50 + 'px'
    }

    const checkRule = () => {
        if(metadata.trim() === '') {
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: 'metadata不能为空'
            })
            return false
        }

        if(workloads.length === 1 && workloads[0].trim() === '') {
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: 'workloads不能为空'
            })
            return false
        }else if(workloads.length > 1) {
            let allEmpty = true
            for(let item of workloads) {
                if(item.trim() !== '') {
                    allEmpty = false
                    break
                }
            }
            if(allEmpty) {
                store.dispatch({
                    type: TYPE.SNACKBAR,
                    val: 'workloads不能为空'
                })
            }
            return !allEmpty
        }

        return true
    }


    const download = () => {
        if(!(checkRule())) { return }

        store.dispatch({
            type: TYPE.LOADING,
            val: true
        })
        axios({
            method: 'POST',
            url: '/api/online/download',
            data: {yaml: previewData},
            headers: { 'Content-Type': 'application/json'}
        }).then(res => {
            if(res.data.code == 0) {
                window.open(window.location.origin+res.data.result)
            }else {
                store.dispatch({
                    type: TYPE.SNACKBAR,
                    val: res.data.result
                })
            }
        }).catch(err => {
            console.log(err)
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: '请求错误'
            })
        }).finally(() => {
            store.dispatch({
                type: TYPE.LOADING,
                val: false
            })
        })
    }

    const arrange = () => {
        if(!(checkRule())) { return }

        store.dispatch({
            type: TYPE.LOADING,
            val: true
        })
        axios({
            method: 'POST',
            url: '/api/online/arrange',
            data: {yaml: previewData},
            headers: { 'Content-Type': 'application/json'}
        }).then(res => {
            if(res.data.code == 0) {
                
                setShowConfigDialog(true)
                setConfigData(res.data.result)

                store.dispatch({
                    type: TYPE.SNACKBAR,
                    val: '添加应用成功'
                })
            }else {
                store.dispatch({
                    type: TYPE.SNACKBAR,
                    val: res.data.result
                })
            }
        }).catch(err => {
            console.log(err)
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: '请求错误'
            })
        }).finally(() => {
            store.dispatch({
                type: TYPE.LOADING,
                val: false
            })
        })
    }

    const getAppSelect = (data) => {
        // 遍历找到所有应用的所有选择的版本
        let selectData = []
        let appConfig = {}
        const configs = (config, attr, obj) => {
            if(config) {
                if(config.type == 'object' && config.properties) {
                    obj[attr] = {}
                    Object.keys(config.properties).forEach((key) => {
                        configs(config.properties[key], key, obj[attr])
                    })
                }else {
                    obj[attr] = config.val
                }
            }
        }

        if(data && data.dependencies ) {
            Object.keys(data.dependencies).forEach((key) => {
                if(data.dependencies[key].location.selected) {
                    selectData.push({
                        "name": key,
                        "location": data.dependencies[key].location.location
                    })

                    
                }else {
                    data.dependencies[key].instances.forEach((item) => {
                        if(item.selected) {
                            selectData.push({
                                "name": item.instance.name,
                                "id": item.instance.id
                            })
                        }
                    })
                }
            })

        }
        if(data.userconfigs) {
            configs(data.userconfigs, 'userconfigs', appConfig)
        }
       
        return {
            id: data.id,
            dependencies: selectData,
            userconfigs: appConfig.userconfigs || null
        }
    }

    // 关闭弹框回调
    const closeDialog = () => {
        setShowConfigDialog(false)
    }

    // 弹框确认按钮回调
    const submitDialog = (data) => {
        if(data.notHadServe.length) {
            // 依赖中存在某些应用没有服务的情况
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: data.notHadServe.join('、') + '以上应用中不存在服务，请创建'
            })
            return 
        }
        if(data.allAppSelectServe.length) {
            // 依赖中存在没有选择服务的应用
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: data.allAppSelectServe.join('、') + '以上应用未选择服务，请选择'
            })
            return
        }

        store.dispatch({
            type: TYPE.LOADING,
            val: true
        })

        let selectData = getAppSelect(data.appInfo)

        selectData['status'] = 1

        console.log('selectData===',selectData)

        // return

        axios({
            method: "POST",
            url: `/api/app/run`,
            headers: {'Content-Type': 'application/json'},
            data: selectData
        }).then((res) => {

            console.log('res==',res)

            if(res.data.code === 0) {
                store.dispatch({
                    type: TYPE.SNACKBAR,
                    val: '部署完成'
                })

                closeDialog()
            }

            store.dispatch({
                type: TYPE.SNACKBAR,
                val: res.data.result || ''
            })

            store.dispatch({
                type: TYPE.LOADING,
                val: false
            })
            
        }).catch((err) => {
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: '请求错误'
            })
            store.dispatch({
                type: TYPE.LOADING,
                val: false
            })
        })

    }


    return (
        <section className="page-container online-container">
             <header className="online-header">
                <div className="header-logo">Crab</div>
                {/* <div className="header-user">userinfo</div> */}
            </header>
            <div className="online-content">
                <div className="oltitle">创建应用</div>
                <section className="createapp-content">
                    <div className="view-text" ref={metaHeaderRef}  ></div>
                    <AutoTextarea ref={metaDataRef} class="textarea-edit indent4" />
                        


                    <div className="createapp-left">
                        <div className="online-title"><p>metadata<span>*</span></p></div>
                        <div className="online-textarea">
                            <textarea className="textarea-input" value={metadata} onChange={changeMetadata}></textarea>
                        </div>

                        <div className="online-title"><p>userconfigs</p></div>
                        <div className="online-textarea">
                            <textarea className="textarea-input" value={userconfig} onChange={changeUserconfig}></textarea>
                        </div>

                        <div className="online-title"><p>workloads<span>*</span></p></div>
                        {
                            workloads && workloads.length ? (
                                workloads.map((item, idx) => {
                                    return (
                                        <div className="online-textarea vartextarea" key={idx} >
                                            <textarea className="textarea-input" value={item} data-index={idx} onChange={changeWorkload}></textarea>
                                            {
                                                workloads.length > 1 ? (
                                                    <Button className="createapp-removebtn" variant="contained" color="secondary" data-index={idx}  onClick={removeWorkload}>移除</Button>
                                                ) : null
                                            }
                                        </div>
                                    )
                                })
                            ) : null
                        }
                        <div className="createapp-addbtn">
                            <Button className="online-btn" variant="contained" color="primary" onClick={addWorkload}>添加</Button>
                        </div>

                        <div className="online-title"><p>dependencies</p></div>
                        <div className="online-textarea">
                            <textarea className="textarea-input" value={dependencies} onChange={changeDependencies}></textarea>
                        </div>

                        <div className="online-title"><p>exports</p></div>
                        <div className="online-textarea">
                            <textarea className="textarea-input" value={exportsData} onChange={changeExports}></textarea>
                        </div>

                    </div>
                    
                    <div className="createapp-right">
                        <div className="online-title"><p>预览</p></div>
                        <div className="createapp-preview">
                            <pre className="preview-pre" ref={preRef}></pre>
                        </div>
                    </div>
                </section>
            </div>
          

            <section className="online-btns">
                <Button className="online-btn" variant="contained" color="primary" onClick={download}>下载</Button>
                <Button className="online-btn" variant="contained" color="primary" onClick={arrange}>部署</Button>
            </section>
            


            <AddFile open={showConfigDialog} title="配置实例" data={configData} close={closeDialog} submit={submitDialog}/>

        </section>
    )
}

function mapStateToProps(state) {
    return state
}

export default connect(mapStateToProps)(CreateApp)