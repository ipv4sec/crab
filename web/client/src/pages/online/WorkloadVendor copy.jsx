import React, { useState, useEffect, useRef } from 'react'
import { connect } from 'react-redux'
import Button from '@material-ui/core/Button'
import '../../style/sass/online.scss'
import axios from 'axios'
import store from '../../store/store'
import * as TYPE from '../../store/actions'

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

const defaultMetadata = `apiVersion: aam.globalsphare.com/v1alpha1
kind: WorkloadVendor
metadata:
    name: example`

const WorkloadVendor = (props) => {
    // const preRef = useRef(null)
    
    const [yamlData, setYamlData] = useState(defaultYaml)
    const [metadata, setMetadata] = useState(defaultMetadata)
    const [systemSpec, setSystemSpec] = useState('')
    // const [cueData, setCueData] = useState('')
    const [cueTpl, setCueTpl] = useState('')

    const changeYaml = (e) => {
        setYamlData(e.target.value)
    }


    const changeMetadata = (e) => {
        setMetadata(e.target.value)
    }

    const getDefaultSystemSpec = () => {
        axios({
            method: 'GET',
            url: '/api/online/systemspec'
        }).then(res => {
            if(res.data.code == 0) {
                setSystemSpec(res.data.result || '')
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

    const changeSystemSpec = (e) => {
        setSystemSpec(e.target.value)
    }

    const changeCueTpl = (e) => {
        setCueTpl(e.target.value)
    }

    // 生成需要的数据
    function getWorkloadVendor() {
        const reg = /\n/g
        return (
            metadata +
            '\nspec: | \n    ' + systemSpec.replace(reg, '\n    ') + 
            '\n    '+cueTpl.replace(reg, '\n        ') 
        ) 

    }


    useEffect(() => {
        getDefaultSystemSpec()
    }, [])


    const checkRule = () => {
        if(metadata.trim() === '') {
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: 'metadata 不能为空'
            })
            return false
        }

        if(systemSpec.trim() === '') {
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: 'system spec 不能为空'
            })
            return false
        }

        if(cueTpl.trim() === '') {
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: 'translate spec cue 不能为空'
            })
            return false
        }

        return true
    }

    const changeYamlToCue = () => {
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
                // window.open(window.location.origin+res.data.result)
                
                // setCueData(res.data.result || '')

                // preRef.current.innerText = res.data.result || ''
                setCueTpl(res.data.result || '')

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
        if(cueTpl.trim() === '') { 
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
            data: {value: cueTpl},
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
        axios({
            method: 'POST',
            url: '/api/online/createvendor',
            data: {value: getWorkloadVendor()},
            headers: { 'Content-Type': 'application/json'}
        }).then(res => {
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: res.data.result
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


    return (
        <section className="page-container online-container">
            <div className="page-title">创建WorkloadVendor</div>
            <section className="vendor-content">

                <div className="vendor-left">
                    <div className="online-title"><p>yaml</p></div>
                    <div className="yaml-textarea">
                        <textarea className="textarea-input" value={yamlData} onChange={changeYaml}></textarea>
                    </div>
                    <div className="online-btns">
                        <Button className="online-btn" variant="contained" color="primary" onClick={changeYamlToCue}>转换yaml为cue</Button>
                    </div>
                </div>

                <div className="vendor-right">
                    <div className="online-title"><p>metadata:</p></div>
                    <div className="vendor-textarea">
                        <textarea className="textarea-input" value={metadata} onChange={changeMetadata}></textarea>
                    </div>

                    <div className="online-title"><p>system spec:</p></div>
                    <div className="vendor-textarea">
                        <textarea className="textarea-input" value={systemSpec} onChange={changeSystemSpec}></textarea>
                    </div>

                    <div className="online-title"><p>translate spec:</p></div>
                    <p>cue template</p>
                    <div className="vendor-textarea vendor-preview">
                        <textarea className="textarea-input" value={cueTpl} onChange={changeCueTpl}></textarea>
                    </div>
                    {/* <div className="vendor-preview">
                        <pre className="preview-pre" ref={preRef}></pre>
                    </div> */}
                    <div className="online-btns">
                        <Button className="online-btn" variant="contained" color="primary" onClick={checkcue}>检查</Button>
                    </div>

                </div>
                
            </section>

            <section className="online-btns">
                <Button className="online-btn" variant="contained" color="primary" onClick={save}>保存</Button>
            </section>
            
        </section>
    )
}

function mapStateToProps(state) {
    return state
}

export default connect(mapStateToProps)(WorkloadVendor)