import React, { useState, useEffect } from 'react'
import Dialog from '@material-ui/core/Dialog'
import DialogTitle from '@material-ui/core/DialogTitle'
import DialogContent from '@material-ui/core/DialogContent'
import DialogActions from '@material-ui/core/DialogActions'
import Button from '@material-ui/core/Button'
import '../style/sass/addfile.scss'

const AddFile = (props) => {
    const [openConfig, setOpenConfig] = useState(false)
    const [curAppConfig, setCurAppConfig] = useState([])
    const [appInfo, setAppInfo] = useState()

    useEffect(() => {
        setData(props.data)
    }, [props.data])

    const getData = () => {
        // 判断应用是否都选择了服务
        let hadNoAppSelectServe = checkAllAppSelectServe(appInfo.dependencies || {})
        
        console.log('--getData-0', appInfo)
        return { 
            appInfo: appInfo, 
            notHadServe: allAppHadServe(appInfo.dependencies || []) || [], //notHadServe,
            allAppSelectServe: hadNoAppSelectServe
        }
    }

    const checkAllAppSelectServe = (dependencies) => {
        let hadNoAppSelectServe = []
        const findNoSelect = (app) => {
           if(Object.keys(app).length) {
               Object.keys(app).forEach((key) => {
                   let selected = false
                   if(app[key].location.selected) {
                    selected = true
                   }else {
                       app[key].instances.forEach((v) => {
                           if(v.selected) {
                               selected = true
                           }
                       })
                   }
                   if(!selected){
                       hadNoAppSelectServe.push(key)
                   }
               })
           }
        }

        findNoSelect(dependencies)

        return hadNoAppSelectServe
    }

    // 判断所有依赖app中是否都有对应的服务
    const allAppHadServe = (dependencies) => {
        let notHadServe = []
        const hadServe = (app) => {
            if(Object.keys(app).length){
                Object.keys(app).forEach((key) => {
                    if(!app[key].instances.length && !app[key].location) {
                        notHadServe.push(key)
                    }
                })
            }
        }

        hadServe(dependencies)

        return notHadServe
    }

    const setData = (data) => {
        if(!data || !data.dependencies) { return }
        // 默认选中第一个
        Object.keys(data.dependencies).forEach((key) => {
            let tmp = []

            data.dependencies[key].instances.forEach((its, index) => {
                if(index == 0 && data.dependencies[key].type === 'mutable') {
                    tmp.push({
                        instance: its,
                        selected: true  // 如果需要默认选中第一个的话，设为true
                    })
                }else {
                    tmp.push({
                        instance: its,
                        selected: false
                    })
                }
            })

            if(!data.dependencies[key].instances || !data.dependencies[key].instances.length || data.dependencies[key].type === 'immutable'){
                data.dependencies[key].location = {
                    location: data.dependencies[key].location,
                    selected: true // 如果需要默认选中的话，设为true
                }
            }else {
                data.dependencies[key].location = {
                    location: data.dependencies[key].location,
                    selected: false
                }
            }
          
            data.dependencies[key].instances = tmp

        })

        console.log('data===',data)
        setAppInfo(data)


        if(data.userconfigs && Object.keys(data.userconfigs).length) {
            // 获取所有config配置字段
            let allConfigs = []
            let level = -1
            const configs = (config, attr, required, level) => {
                if(config) {
                    if(config.type == 'object' && config.properties) {
                        if(attr !== 'userconfigs') { 
                            allConfigs.push({
                                key: attr,
                                type: config.type,
                                val: '',
                                required: false,
                                error: '',
                                level: level
                            })
                        }
                        
                        level += 1
                        Object.keys(config.properties).forEach((key) => {
                            configs(config.properties[key], key, config.required, level)
                        })
                    }else {
                        allConfigs.push({
                            key: attr,
                            type: config.type,
                            val: '',
                            required: required ? required.indexOf(attr) !== -1 : false,
                            error: '',
                            level: level
                        })
                    }
                }
            }

            configs(data.userconfigs, 'userconfigs', false, level)


            setOpenConfig(true)
            setCurAppConfig(allConfigs)
        }
       
        
    }

    function closeDialog() {
        props.close()
    }

    function submitDialog() {
        props.submit(getData())
    }

    function checkVersion(type, key, idx) {

        if(appInfo.dependencies[key].type === 'immutable' && type == 'instance' ){
            return
        }
        let selected = false
        let newData = Object.assign({}, appInfo)
        if(type === 'instance'){
            newData['dependencies'][key]['location'].selected = false
            newData['dependencies'][key].instances.forEach((item, index) => {
                if(index == idx) {
                    item.selected = !item.selected
                    selected = item.selected
                    // newData['dependencies'][key][idx].selected = !newData['dependencies'][key]['instance[idx].selected
                }else {
                    item.selected = false
                    // newData[key]['dependencies'][idx].selected = false
                }
            })

        }else if(type === 'location') {
            newData['dependencies'][key]['instances'].forEach((item, index) => {
               item.selected = false
            })
            newData['dependencies'][key]['location'].selected = !newData['dependencies'][key]['location'].selected
            selected = newData['dependencies'][key]['location'].selected
        }

        setAppInfo(newData)
    }

    function changeInput() {
        const key = event.target.dataset.key
        let newData = Object.assign({}, appInfo)
        newData.dependencies[key].location.location = event.target.value
        setAppInfo(newData)
    }
    
    function changeKey() {
        const idx = event.target.dataset.index
        let newConfig = curAppConfig.slice()
        newConfig[idx].val = event.target.value //newConfig[idx].type == 'number' ? Number(event.target.value) :  event.target.value
        newConfig[idx].error = ''
        setCurAppConfig(newConfig)
    }

    function renderTree(data) {
        let dom =[]
        let left = -34
      
        const tree = (dependences) => {
            if(!dependences) { return }
            left += 34
            Object.keys(dependences).forEach((key, index) => {
                dom.push(
                    <div style={{paddingLeft: left+'px' }} className="app-item" key={index}>
                        <div className="app-label">
                            <span style={{backgroundColor: '#54CACB'}} className="label-icon app-icon"><i className="iconfont icon_grey600"></i> </span>
                            <span className="label-name">{key}</span>
                        </div>

                        {
                            (dependences[key].instances || []).map((v,i) => {
                                return (
                                        <div key={i} className="app-item-versions">
                                            <div className="app-label" key={i} >
                                                <span className="label-icon version-icon"  onClick={() => {checkVersion('instance', key, i)}}>
                                                    <i style={{color: dependences[key].type === 'mutable' ? '#54CACB' : '#e0e0e0'}} className={`${v.selected ? "iconfont icon_d-pass" : ""}`}></i> 
                                                </span>
                                                <span className="label-name" style={{color: dependences[key].type === 'mutable' ? '#262626' : 'gray'}} >
                                                    实例：{v.instance && v.instance.id ? v.instance.id : ''}  {v.instance && v.instance.name ? v.instance.name : ''}
                                                </span>
                                            </div>
                                        </div>
                                )
                            })  
                        }
                        {
                            dependences[key].location ? (
                                <div className="app-item-versions">
                                    <div className="app-label" >
                                        <span className="label-icon version-icon"  onClick={() => {checkVersion('location',  key, '')}}>
                                            <i style={{color: dependences[key].type === 'mutable' ? '#54CACB' : '#e0e0e0'}} className={`${dependences[key].location.selected ? "iconfont icon_d-pass" : ""}`}></i> 
                                        </span>
                                        <span className="label-name" style={{color: dependences[key].type === 'mutable' ? '#262626' : 'gray'}} >
                                            服务地址：<input data-key={key}  onChange={changeInput} value={dependences[key].location.location} disabled={dependences[key].type === 'immutable'} style={{color: dependences[key].type === 'mutable' ? '#262626' : 'gray'}} />
                                        </span>
                                    </div>
                                </div>
                            ) : null
                        }

                        {
                            !dependences[key].instances.length && !dependences[key].location ? (
                            <div  className="app-item-versions">
                                <div className="app-label" >
                                    <span className="label-icon version-icon" >
                                        {/* <i style={{color: item.color ? item.color : '#54CACB'}} className={`${v.selected ? "iconfont icon_d-pass" : ""}`}></i>  */}
                                    </span>
                                    <span className="label-name">
                                        暂无，需创建
                                    </span>                                            
                                </div>
                            </div>) : null
                        }
                    </div>
                )
            })

        }

        
        tree(data.dependencies || null)

        return dom
    }

    function confirmConfig() {
        let success = true
        let configs = curAppConfig.slice()

        configs.forEach((item) => {
            if(item.type !== 'object' && item.required) {
                if(item.val.trim() == '') {
                    item.error = '请输入'
                    success = false
                }
            }
        })

        if(success) {
            let newAppInfo = Object.assign({}, appInfo)
            const configs = (config, attr) => {
                if(config) {
                    // console.log('---renderconfig--', config.type , config.properties, config.type == 'object' && config.properties)
                    if(config.type == 'object' && config.properties) {
                        Object.keys(config.properties).forEach((key) => {
                            configs(config.properties[key], key)
                        })
                    }else {
                        console.log(222)
                        let idx = curAppConfig.findIndex((val) => val.key === attr)
                        config.val = curAppConfig[idx].type === 'number' ? Number(curAppConfig[idx].val) : curAppConfig[idx].val
                    }
                }
            }
            configs(newAppInfo.userconfigs, 'userconfigs') 

            console.log('--newAppInfo--', newAppInfo)

            setAppInfo(newAppInfo)

            setOpenConfig(false)

        }else {
            setCurAppConfig(configs)
          
        }

    }

    return (
        <Dialog
            open={props.open}
            onClose={closeDialog}
            aria-labelledby="upload-file-title"
        >
            <DialogTitle id="upload-file-title">{props.title}</DialogTitle>
            <DialogContent>

                <div className="instance-content">
                    {
                        renderTree(props.data).map((item, index) => item)
                    }
                </div>

                <Dialog
                    open={openConfig}
                    aria-labelledby="config-title"
                >
                    <DialogTitle id="config-title">实例配置</DialogTitle>
                    <DialogContent>
                        <div className="appconfig-content">
                            {
                                curAppConfig.map((item, index) => {
                                    if(item.type === 'object') {
                                        return (
                                            <div key={item.key} className="config-item" style={{marginLeft: item.level * 40 + 'px'}}>
                                                <div className="item-input">
                                                    <label>{item.key}:</label>
                                                </div>
                                            </div>
                                        )
                                    }else {
                                        return (
                                            <div key={item.key} className="config-item" style={{marginLeft: item.level * 40 + 'px'}}>
                                                <div className="item-input">
                                                    <label>{item.required ? <span>* </span> : null}{item.key}:</label>
                                                    <input 
                                                        type={item.type == 'number' ? 'number' : 'text'}  
                                                        className={item.error ? 'red-border' : ''} 
                                                        placeholder={item.type === 'number' ? '请输入number类型' : '请输入string类型'}
                                                        data-index={index}
                                                        onChange={changeKey} 
                                                        value={item.val} />
                                                </div>
                                                {
                                                    item.error ? (
                                                        <p className="item-error">{item.error}</p>
                                                    ) : null
                                                }
                                                
                                            </div>
                                        )
                                    }
                                    
                                })
                            }
                        </div>
                    </DialogContent>
                    <DialogActions>
                        <Button className="common-btn" color="primary" onClick={confirmConfig}>确定</Button>
                    </DialogActions>
                </Dialog>

            </DialogContent>
            <DialogActions>
                <Button className="common-btn" onClick={closeDialog}>取消</Button>
                <Button className="common-btn" color="primary" onClick={submitDialog}>确定</Button>
            </DialogActions>

        </Dialog>
    )
}

export default AddFile