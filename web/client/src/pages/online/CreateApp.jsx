import React, { useState, useEffect, useRef } from 'react'
import { connect } from 'react-redux'
import Button from '@material-ui/core/Button'
import '../../style/sass/online.scss'
import axios from 'axios'
import store from '../../store/store'
import * as TYPE from '../../store/actions'
import AddFile from '../../components/AddFile'
import Loading from '../../components/Loading'
import SnackbarCmp from '../../components/Snackbar'
import Editor from '../../components/Editor'

const defaultMetaHeader = `apiVersion: aam.globalsphare.com/v1alpha1
kind: Application
metadata:`

const defaultMetadata = `name: example
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

const defaultWorkloads = `- name: example
  type: webservice
  vendor: webservice
  properties:
    image: nginx:1.21
  traits:
    - type: ingress
      properties:
        k1: "v1"`



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
    const [showConfigDialog, setShowConfigDialog] = useState(false)
    const [configData, setConfigData] = useState({})


    const metaHeaderRef = useRef(null)
    const metaDataRef = useRef(null)
    const userConfigsRef = useRef(null)
    const workloadsRef = useRef(null)
    const dependenciesRef = useRef(null)
    const exportsRef = useRef(null)

    useEffect(() => {
        metaHeaderRef.current.innerText = defaultMetaHeader
        metaDataRef.current.setData(defaultMetadata)
        userConfigsRef.current.setData(defaultUserconfigs)
        workloadsRef.current.setData(defaultWorkloads)
        dependenciesRef.current.setData(defaultDependencies)
        exportsRef.current.setData(defaultExports)
        
    }, [])


    // 获取页面所有内容
    function getAllData() {
        const reg = /\n/g
        
        const pData = (
            defaultMetaHeader +
            '\n  ' + metaDataRef.current.getData().replace(reg, '\n  ') + 
            '\nspec:' +
            '\n  userconfigs:\n    ' + userConfigsRef.current.getData().replace(reg, '\n    ') + 
            '\n  workloads:\n    '+ workloadsRef.current.getData().replace(reg, '\n    ') + 
            '\n  dependencies:\n    '+ dependenciesRef.current.getData().replace(reg, '\n    ') + 
            '\n  exports:\n    '+ exportsRef.current.getData().replace(reg, '\n    ')
        ) 

        return pData
    }

    const checkRule = () => {
        if(metaDataRef.current.getData().trim() === '') {
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: 'metadata不能为空'
            })
            return false
        }

        if(workloadsRef.current.getData().trim() === '') {
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: 'workloads不能为空'
            })
            return false
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
            data: {yaml: getAllData()},
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
            data: {yaml: getAllData()},
            headers: { 'Content-Type': 'application/json'}
        }).then(res => {
            if(res.data.code == 0) {
                
                setShowConfigDialog(true)
                setConfigData(res.data.result || '')

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

        axios({
            method: "POST",
            url: `/api/app/run`,
            headers: {'Content-Type': 'application/json'},
            data: selectData
        }).then((res) => {

            if(res.data.code === 0) {
                store.dispatch({
                    type: TYPE.SNACKBAR,
                    val: '部署完成'
                })

                closeDialog()

                setTimeout(() => {
                    window.opener.postMessage('createapp', window.location.origin)
                    window.close()
                }, 1000)    
               
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
                    <div className="textarea-edit indent4"><Editor ref={metaDataRef} uniqueName="metaDataRef" /></div>
                  
                    <div className="view-text">spec:</div>
                    <div className="view-text indent6">userconfigs:</div>
                    <div className="textarea-edit indent8"><Editor ref={userConfigsRef} uniqueName="userConfigsRef"/></div>

                    <div className="view-text indent6">workloads:</div>
                    <div className="textarea-edit indent8"><Editor ref={workloadsRef}  uniqueName="workloadsRef"/></div>
                   
                    <div className="view-text indent6">dependencies:</div>
                    <div className="textarea-edit indent8"><Editor ref={dependenciesRef}  uniqueName="dependenciesRef"/></div>

                    <div className="view-text indent6">exports:</div>
                    <div className="textarea-edit indent8"><Editor ref={exportsRef}  uniqueName="exportsRef"/></div>

                    <div className="online-btns">
                        <Button className="online-btn" variant="contained" color="primary" onClick={download}>下载</Button>
                        <Button className="online-btn" variant="contained" color="primary" onClick={arrange}>部署</Button>
                    </div>
                    
                </section>
            </div>
          

            <AddFile open={showConfigDialog} title="配置实例" data={configData} close={closeDialog} submit={submitDialog}/>

            <Loading />
            <SnackbarCmp />
        </section>
    )
}

function mapStateToProps(state) {
    return state
}

export default connect(mapStateToProps)(CreateApp)