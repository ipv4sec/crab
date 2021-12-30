import React, { useState, useEffect, useRef } from 'react'
import { connect } from 'react-redux'
import Button from '@material-ui/core/Button'
import '../../style/sass/online.scss'
import axios from 'axios'
import store from '../../store/store'
import * as TYPE from '../../store/actions'
import Loading from '../../components/Loading'
import SnackbarCmp from '../../components/Snackbar'
import AutoTextarea from '../../components/AutoTextarea'

const defaultYaml = `apiVersion: apps/v1
kind: Deployment
metadata:
    name: nginx-deployment
    labels:
        app: nginx
spec:
    replicas: 3
    selector:
        matchLabels:
            app: nginx
    template:
        metadata:
            labels:
                app: nginx
    spec:
        containers:
          - name: nginx
            image: nginx:1.14.2
            ports:
                - containerPort: 80`

const metaHeader = `apiVersion: aam.globalsphare.com/v1alpha1
kind: WorkloadVendor
metadata:`

const defaultMetadata = `name: example`

const WorkloadVendor = (props) => {
    
    const metaDataRef = useRef(null)
    const metaRef = useRef(null)
    const specRef = useRef(null)
    const cueRef = useRef(null)
    const yamlRef = useRef(null)
    const [specFold, setSpecFold] = useState(false)

    const [name, setName] = useState('')
    const [vendorInfo, setVendorInfo] = useState(null)

    const [btnDisable, setBtnDisable] = useState(false)

    const getName = () => {
        let name = ''
        if(window.location.search) {
            const params = window.location.search.substring(1, )
            if(params.indexOf('&')) {
                const kvs =  params.split('&')
            
                for(let i = 0, len = kvs.length; i < len; i++) {
                    const kv = kvs[i].split('=')
                    if(kv && kv[0] === 'name') {
                        name = kv[1]
                        break;
                    }
                }
            }
           
        }

        return name
    }

    useEffect(() => {
        const name = getName()
        metaRef.current.innerText = metaHeader
        metaDataRef.current.setData(defaultMetadata)
        if(name) {
            setName(name)
            getWorkloadVendorInfo(name)
        }else {
            yamlRef.current.setData(defaultYaml)
        }
    }, [])

    const getWorkloadVendorInfo = (name) => {
        store.dispatch({
            type: TYPE.LOADING,
            val: true
        })
        axios({
            method: 'GET',
            url: '/api/online/getworkloadvendor',
            params: {name}
        }).then(res => {
           
            if(res.data.code == 0) {
                setVendorInfo(res.data.result || {})
                yamlRef.current.setData(res.data.result.yaml || '')
                cueRef.current.setData(res.data.result.cue || '')

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


    const editWorkloadVendor = () => {
        store.dispatch({
            type: TYPE.LOADING,
            val: true
        })

        let url = `/api/cluster/editvendor?id=${vendorInfo.id || ''}`
        setBtnDisable(true)
        axios({
            url: url,
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            data: {
                metadata: metaDataRef.current.getData(),
                yaml: yamlRef.current.getData(),
                cue: cueRef.current.getData(),
                value: getWorkloadVendor()
            }
        }).then((res) => {
            if(res.data.code == 0) {
                setTimeout(() => {
                    setBtnDisable(false)  
                    window.opener.postMessage('workloadtype', window.location.origin)
                    window.close()
                }, 1000)
              
            }else {
                setBtnDisable(false)  
            }
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: res.data.result || ''
            })
        }).catch((err) => {
            console.error(err)
            setBtnDisable(false)  
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

    
    const getDefaultSystemSpec = () => {
        axios({
            method: 'GET',
            url: '/api/online/systemspec'
        }).then(res => {
            if(res.data.code == 0) {
                specRef.current.setData(res.data.result || '')
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
        })
    }


    // 生成需要的数据
    function getWorkloadVendor() {
        const reg = /\n/g
        return (
            metaHeader + 
            '\n    ' + metaDataRef.current.getData().replace(reg, '\n    ') + 
            '\nspec: | \n    ' + specRef.current.getData().replace(reg, '\n    ') +
            '\n    '+cueRef.current.getData().replace(reg, '\n        ') 
        ) 
    }

    useEffect(() => {
        getDefaultSystemSpec()
    }, [])

    const checkRule = () => {
        if(metaDataRef.current.getData().trim() === '') {
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: 'metadata 不能为空'
            })
            return false
        }

        if(specRef.current.getData().trim() === '') {
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: 'spec 不能为空'
            })
            return false
        }

        if(cueRef.current.getData().trim() === '') {
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: 'cue 不能为空'
            })
            return false
        }

        return true
    }

    const changeYamlToCue = () => {
        const yamlData = yamlRef.current.getData()
        if(yamlData.trim() === '') { 
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: '请填写yaml'
            }) 
            return 
        }

        store.dispatch({
            type: TYPE.LOADING,
            val: true
        })
        axios({
            method: 'POST',
            url: '/api/online/translateyaml',
            data: {value: yamlData},
            headers: { 'Content-Type': 'application/json'}
        }).then(res => {
            if(res.data.code == 0) {

                cueRef.current.setData(res.data.result || '')

                store.dispatch({
                    type: TYPE.SNACKBAR,
                    val: '转换完成'
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

    const checkcue = () => {
        const cueData = cueRef.current.getData()
        if(cueData.trim() === '') { 
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: 'cue 不能为空'
            }) 
            return 
        }

        store.dispatch({
            type: TYPE.LOADING,
            val: true
        })
        axios({
            method: 'POST',
            url: '/api/online/checkcue',
            data: {value: cueData},
            headers: { 'Content-Type': 'application/json'}
        }).then(res => {
        
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: res.data.code === 0 ? '检查正确' : res.data.result
            })
           
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

    const save = () => {

        if(!(checkRule())) { return }

        store.dispatch({
            type: TYPE.LOADING,
            val: true
        })
        setBtnDisable(true)  
        axios({
            method: 'POST',
            url: '/api/online/createvendor',
            data: {
                metadata: metaDataRef.current.getData(),
                yaml: yamlRef.current.getData(),
                cue: cueRef.current.getData(),
                value: getWorkloadVendor()
            },
            headers: { 'Content-Type': 'application/json'}
        }).then(res => {
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: res.data.result
            })
            if(res.data.code == 0) {
                setTimeout(() =>{
                    setBtnDisable(false)  
                    window.opener.postMessage('workloadvendor', window.location.origin)
                    window.close()
                }, 1000)
            }else {
                setBtnDisable(false)  
            }
        }).catch(err => {
            console.log(err)
            setBtnDisable(false)  
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

    const specFoldFn = () => {
        setSpecFold(!specFold)
    }


    return (
        <section className="page-container online-container">
            <header className="online-header">
                <div className="header-logo">Crab</div>
                {/* <div className="header-user">userinfo</div> */}
            </header>
            <div className="online-content">
                <div className="oltitle">{name ? '修改' : '创建'}WorkloadType</div>
                <section className="vendor-content">
                    <div className="vendor-left">
                        <div className="online-title"><p>yaml</p></div>
                        <AutoTextarea ref={yamlRef} class="yaml-textarea" />
                        <div className="online-btns">
                            <Button className="online-btn" variant="contained" color="primary" onClick={changeYamlToCue}>转换yaml为cue</Button>
                        </div>
                    </div>

                    <div className="vendor-right">
                        <div className="view-text" ref={metaRef}  ></div>
                        <AutoTextarea ref={metaDataRef} class="textarea-edit indent4" />
                        <div className="view-text" >spec: | 
                            <button className="fold-btn" onClick={specFoldFn}><span className={`iconfont ${specFold ? 'icon_navigation_combobox_down' : 'icon_navigation_combobox_up'}`}></span></button>
                        </div>
                        <AutoTextarea ref={specRef} class={`textarea-edit indent4  ${specFold ? 'hide-textarea' : ''}`} />
                       
                        <div className="view-text" >cue Template:</div>
                        <AutoTextarea ref={cueRef} class="textarea-edit indent4" />
                      
                        <div className="online-btns">
                            <Button className="online-btn" variant="contained" color="primary" onClick={checkcue}>检查</Button>
                        </div>
                        <div className="online-btns">
                            {
                                name ? (
                                    <Button disabled={btnDisable} className="online-btn" variant="contained" color="primary" onClick={editWorkloadVendor}>修改</Button>
                                ) : (
                                    <Button disabled={btnDisable} className="online-btn" variant="contained" color="primary" onClick={save}>保存</Button>
                                )
                            }
                        </div>
                    </div>
                </section>
            </div>
            
            <Loading />
            <SnackbarCmp />

        </section>
    )
}

function mapStateToProps(state) {
    return state
}

export default connect(mapStateToProps)(WorkloadVendor)