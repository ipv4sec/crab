import React, { version } from 'react'
import axios from 'axios'
import store from '../../store/index'
import * as TYPE from '../../store/actions'
import BaseDialog from '../../components/dialog/baseDialog'
// 首先获取根应用的详情，和版本，并且每个应用下都有个外部实例的选项，进行填写key：value
// 如某个节点的子节点依赖的应用中存在于同级或者父级的实例，那么这么子节点的版本列表中就要有一个即有实例的选项


export default class AppTree extends React.Component{
    constructor(props){
        super(props)
        this.state = {
            curAppId: this.props.appId,
            curVersion: '',
            curNodeId: '',
            curInstanceId: '',
            curAppConfig: [], // 选择的应用的config配置
            selectAppKey: '', // 选择的应用的key
            appInfo: {
                // "instanceid": "fdencpyq",
                // "dependencies": {
                //     "demo-app1": { //#应用名
                //         "instances": [], //#应用实例 id 列表
                //         "location": { location: "https://www.huanqiu.com", selected: false},
                //         "resources": [
                //             {
                //                 "actions": [
                //                     "GET"
                //                 ],
                //                 "uri": "/app"
                //             }
                //         ],
                //         "type": "mutable",
                //         "userconfig": {
                //             "properties": {
                //                 "param1": {
                //                     "type": "string"
                //                 },
                //                 "param2": {
                //                     "type": "number"
                //                 },
                //                 "param3": {
                //                     "type": "object",
                //                     "properties": {
                //                         "param3_1":{
                //                             "type": "string"
                //                         },
                //                         "param3_2":{
                //                             "type": "number"
                //                         }
                //                     }
                //                 }
                //             },
                //             "required": [
                //                 "param1"
                //             ],
                //             "type": "object"
                //         }
                //     }
                // }
            },
            appList: [
                // {
                //     app_name: 'app1', 
                //     app_id: '1', 
                //     value: '1',
                //     app_versions: [{version: 'external', value: 'external',selected: false, config: false},{version: 'v0.0.1', value: '0.0.1',selected: false, config: false},{version: 'v0.0.2', value: '0.0.2',selected: false, config: false}],
                //     dependences: [
                //         {app_name: 'app2', app_id: '2', value: '2', app_versions: [{version: 'v0.0.1', value: '0.0.1',selected: false, config: false},{version: 'v0.0.2', value: '0.0.2',selected: false, config: false}],},
                //         {
                //             app_name: 'app3', 
                //             app_id: '3',
                //             value: '3',
                //             app_versions: [{version: 'v0.0.1', value: '0.0.1',selected: false, config: false},{version: 'v0.0.2', value: '0.0.2',selected: false, config: false}],
                //             dependences: [
                //                 {app_name: 'app4', app_id: '4', value: '4', app_versions: [{version: 'v0.0.1', value: '0.0.1',selected: false, config: false},{version: 'v0.0.2', value: '0.0.2',selected: false, config: false}],},
                //                 {app_name: 'app5', app_id: '5', value: '5',app_versions:[{version: 'v0.0.1', value: '0.0.1',selected: false, config: false},{version: 'v0.0.2', value: '0.0.2',selected: false, config: false}],},
                //                 {
                //                     app_name: 'app1', app_id: '1', value: '6',
                //                     app_versions: [{version: 'v0.0.1', value: '0.0.1',selected: false, config: false},{version: 'v0.0.2', value: '0.0.2',selected: false, config: false}],
                //                     dependences: [
                //                     ]
                //                 }
                //             ]
                //         }
                //     ]
                // }
            ],
            notHadServe: [], // 依赖的应用中是否都存在服务，如果有不存在的，不能添加应用
            allAppSelectServe: [], // 是否所有的依赖app都选择了服务
            showDialog: false, // 添加依赖时打开某个应用的配置文件
            appConfig: {
                name: {
                    type: 'string',
                    required: true,
                    value: '',
                    error: '请输入'
                },
                age: {
                    type: 'number',
                    required: false,
                    value: ''
                },

            }
        }

        this.setData = this.setData.bind(this)
        this.getData = this.getData.bind(this)
        this.getAllApps = this.getAllApps.bind(this)
    }

    componentWillMount(){
       
    }

    componentDidMount() {

        return
         // 暂时不需要回显编辑数据
        if(this.props.isEdit){
            // console.log('---appInfo---', sessionStorage.getItem('appInfo'))
            try{
                let list = sessionStorage.getItem('appInfo') && sessionStorage.getItem('appInfo') !== '' ? JSON.parse(sessionStorage.getItem('appInfo')) : null
                this.setState({
                    // appList: list ? [list] : []
                    appList: this.state.list
                },() => {
                    this.renderTree(this.state.appList)
                    this.props.checkedAppList(this.state.appList)
                    let selected = {}
                    const findLastExternal = (nodes) => {
                        if(nodes && nodes.length) {
                            nodes.forEach((item) => {
                                console.log(item)
                                if(item.dependences && item.dependences.length){
                                    findLastExternal(item.dependences)
                                }else {
                                    if(item.version_id == 'external') {
                                        selected['node_id'] = item.node_id
                                        selected['appName'] = item.app_name
                                        selected['appId'] = item.app_id
                                        selected['version'] = item.version_id
                                        selected['instanceId'] = item.instance_id
                                        selected['conf'] = item.conf || []

                                    }
                                }
                            })
                        }
                    }
                   
                    findLastExternal(this.state.appList)
                    this.props.versionConfig(selected)
                })
            }catch(err) {
                console.log(err)
            }
           
        }else {
            this.getAppInfo();
        }
       
       
    }


    // 返回所有应用、没有服务的应用，是否选全都选了服务的应用
    getData() {
        // 判断应用是否都选择了服务
        let hadNoAppSelectServe = this.checkAllAppSelectServe(this.state.appInfo.dependencies || {})
        
        return {
            // appList: this.state.appList, 
            appInfo: this.state.appInfo, 
            notHadServe: this.state.notHadServe,
            allAppSelectServe: hadNoAppSelectServe
        }
    }

    setData(data) {
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

        this.setState({
            appInfo: data,
        })

        // 获取所有config配置字段
        let allConfigs = []
        let level = -1
        const configs = (config, attr, required, level) => {
            if(config) {
                if(config.type == 'object' && config.properties) {
                    if(attr !== 'userconfig') { 
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

        configs(data.userconfig, 'userconfig', false, level)

        this.setState({
            showDialog: true,
            curAppConfig: allConfigs, //userconfig
            // selectAppKey: key 
        })
    }

    // 判断所有依赖app中是否都有对应的服务
    allAppHadServe(dependencies) {
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


    getAllApps() {
        return this.state.appList
    }

    generateInstanceId(randomFlag=true, min=8, max=8){
        let str = "",
            range = min,
            arr = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'];
        
        // 随机产生
        if(randomFlag){
            range = Math.round(Math.random() * (max-min)) + min;
        }
        for(var i=0; i<range; i++){
            let pos = Math.round(Math.random() * (arr.length-1));
            str += arr[pos];
        }
        return 'ais-'+str;
    }

    getAppDetail(type, param){
        let params = {}
        if(type == 'name') {
            params = {name: param}
        }else {
            params = {appId: param}
        }
        return new Promise((resolve, reject)=>{
            axios({
                method: 'GET',
                url: `${window._BASEPATH}/api/appDetail`,
                params: params
            }).then((res) => {
                if(res && res.data && res.data.id) {
                    resolve(res.data)
                }else {
                    resolve(false)
                }
            }).catch((err) => {
                console.log(err)
                reject(false)
            })
        })
    }

    getAppVersions(appId, willHadInstance){
        return new Promise((resolve, reject)=>{
            axios({
                method: 'GET',
                url: `${window._BASEPATH}/api/getAppVersions`,
                params: {
                    appId: appId, //this.props.params.appId
                }
            }).then((res) => {
                let data = [{selected: false, version: 'external'}]
                if(willHadInstance && willHadInstance.length) {
                    data.push({selected: false, version: 'willHadInstance'})
                }
                if(res && res.data && res.data.list && res.data.list.length){
                    res.data.list.forEach((item) => {
                        data.push({
                            selected: false,
                            version: item
                        })
                    })
                    resolve(data)
                }else {
                    resolve(false)
                }
            }).catch((err) => {
                reject(false)
            })
        })
    }

    getAppInfo(){
        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })
        Promise.all([this.getAppDetail('appId',this.state.curAppId), this.getAppVersions(this.state.curAppId)]).then((res) => {
            let detail = res[0]
            let versions = res[1]
            let appList = []
            let tmp = {
                node_id: this.generateInstanceId(), // 节点唯一ID
                app_name: detail.name,
                app_versions: versions,
                app_id: detail.id,
                version_id: "",
                instance_id: this.generateInstanceId(),
                external: 0, // 外部配置
                conf: {}, // 外部配置参数
                selectWillInstance: null, // 选了即有实例
                willInstances: [], // 既有实例列表
                dependences: []
            }
            appList.push(tmp)
            this.setState({
                appList: appList
            })

            this.renderTree(appList)
            this.props.checkedAppList(appList)

            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })

        }).catch((err) => {
            console.log(err)
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
        })
    }


    // 获取某个版本的详情（包含依赖）
    getDependApp() {
        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })
        axios({
            method: "GET",
            url: `${window._BASEPATH}/api/getVersionDetail`,
            params: {
                appId: this.state.curAppId,
                version: this.state.curVersion 
            }
        }).then(async (res) => {
            let dependences = []
            if(res.data.dependencies && Object.keys(res.data.dependencies).length){
                let keys = Object.keys(res.data.dependencies)
                if(keys.length) {
                    for(let i = 0; i < keys.length; i++){
                        let detail = await this.getAppDetail('name',keys[i]);
                        if(!detail) {
                            continue
                        }
                        let willHadInstance = this.getHadInstance(detail.id)
                        let versions = await this.getAppVersions(detail.id, willHadInstance);
                        if(!versions) {
                            continue
                        }
                        dependences.push({
                            node_id: this.generateInstanceId(), // 节点唯一id
                            app_name: keys[i],
                            app_id: detail.id,
                            app_versions: versions,
                            version_id: '',
                            instance_id: this.generateInstanceId(),
                            external: 0, // 外部配置
                            conf: {}, // 外部配置参数
                            dependences: [],
                            selectWillInstance: null, // 选了即有实例
                            willInstances: willHadInstance, // 既有实例
                        })
                    }
                }
            }

            let newAppList = this.addDependence(dependences)
            this.setState({
                appList: newAppList
            })
            this.props.checkedAppList(newAppList)
            this.renderTree(newAppList)

            // 添加后向下滑动显示出一个依赖版本
            this.treeContainerRef.scrollTop = this.treeContainerRef.scrollTop + 120

            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
        }).catch((err) => {
            console.log(err)
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
        })
    }

    getHadInstance(appId){
        let selectedInstances = []
        const findSameInstance = (allInstance) => {
            if(allInstance && allInstance.length) {
                allInstance.forEach((item) => {
                    // 后期 可能会增加 集群中已存在的实例 字段，如果是已存在的，那么必须要排除出去这个实例，因为他不算既有实例
                    if(item.app_id == appId && item.version_id && item.version_id !== ''){
                        selectedInstances.push({
                            node_id: item.node_id,
                            app_name: item.app_name,
                            app_id: item.app_id,
                            version_id: item.version_id,
                            instance_id: item.instance_id
                        })
                    }else {
                        if(item.dependences){
                            findSameInstance(item.dependences)
                        }
                    }
                })
            }
        }

        

        findSameInstance(this.state.appList)
        return selectedInstances
    }

    

    addDependence(data){
        const findNode = (nodes) => {
            if(nodes && nodes.length) {
                for(let i = 0; i < nodes.length; i ++) {
                    if(nodes[i].node_id == this.state.curNodeId) {
                        nodes[i].dependences = data
                        break;
                    }else {
                        findNode(nodes[i].dependences)
                    }
                }
            }
        }

        let allApps = this.state.appList.slice()
        
        findNode(allApps)
        
        return allApps
    }

    changeInput(key) {
        let newData = Object.assign({}, this.state.appInfo)
        newData.dependencies[key].location.location = event.target.value
        this.setState({
            appInfo: newData
        })
    }


    renderTree(data) {
        let dom =[]
        let left = -34
      
        const tree = (dependences) => {
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
                                                <span className="label-icon version-icon"  onClick={this.checkVersion.bind(this, 'instance', key, i, dependences[key].userconfig)}>
                                                    <i style={{color: dependences[key].type === 'mutable' ? '#54CACB' : '#e0e0e0'}} className={`${v.selected ? "iconfont icon_d-pass" : ""}`}></i> 
                                                </span>
                                                <span className="label-name" style={{color: dependences[key].type === 'mutable' ? '#262626' : 'gray'}} >
                                                    实例：{v.instance.instanceid || ''} {v.instance.version}
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
                                        <span className="label-icon version-icon"  onClick={this.checkVersion.bind(this, 'location',  key, '', dependences[key].userconfig )}>
                                            <i style={{color: dependences[key].type === 'mutable' ? '#54CACB' : '#e0e0e0'}} className={`${dependences[key].location.selected ? "iconfont icon_d-pass" : ""}`}></i> 
                                        </span>
                                        <span className="label-name" style={{color: dependences[key].type === 'mutable' ? '#262626' : 'gray'}} >
                                            服务地址：<input  onChange={this.changeInput.bind(this, key)} value={dependences[key].location.location} disabled={dependences[key].type === 'immutable'} style={{color: dependences[key].type === 'mutable' ? '#262626' : 'gray'}} />
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

        
        tree(data.dependencies || {})

        return dom
    }

    checkVersion_bak(curVersion) {
        if(curVersion.selected) { return } // 已经选择时点击无效
        console.log('curData===',app, curVersion)
        // let selected = {}

        const findVersion = (id, curV, list) => {
            if(list && list.length) {
                list.forEach((item) => {
                    if(id == item.app_id) { 
                        if(item.app_versions && item.app_versions.length) {
                            item.app_versions.forEach((ele) => {
                                // console.log(ele.version , curV)
                                if(ele.version == curV) {
                                    ele.selected = true
                                   
                                }else {
                                    ele.selected = false
                                }
                            })
                        }
                    }else if(item.dependences) {
                        findVersion(id, curV, item.dependences)
                    }
                })
            }
        }

        let newList = this.state.appList
        findVersion(app.node_id, curVersion.version, newList)


        this.setState({
            appList: newList,
            // allAppSelectServe: hadNoSelected,
            // curAppId: selected['appId'] || '',
            // curVersion: selected['version'] || '',
            // curInstanceId: selected['instanceId'] || '',
            // curNodeId: selected['node_id'] || '',
        },() => {
            // if(!(/external|willHadInstance/.test(selected.version))){
            //     this.getDependApp()
            // }else {
            //     this.addDependence([])
            // }
            
            // 点击后返回所有的app， 其实应该返回所有选中的app和该app的版本
            // this.props.checkedAppList(newList)
        })
    }


    checkVersion(type, key, idx, userconfig) {

        // console.log('checkVersion===',userconfig)

        if(this.state.appInfo.dependencies[key].type === 'immutable' && type == 'instance' ){
            return
        }
        let selected = false
        let newData = Object.assign({}, this.state.appInfo)
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
        this.setState({
            appInfo: newData
        })

        // if(selected) {

        //     let allConfigs = []
        //     const configs = (config, attr) => {
        //         if(config) {
        //             if(config.type == 'object' && config.properties) {
        //                 Object.keys(config.properties).forEach((key) => {
        //                     configs(config.properties[key], key)
        //                 })
        //             }else {
        //                 console.log(222)
        //                 allConfigs.push({
        //                     key: attr,
        //                     type: config.type,
        //                     val: '',
        //                     required: userconfig.required.indexOf(attr) !== -1, //config.required,
        //                     error: ''
        //                 })
        //             }
        //         }
        //     }

        //     configs(userconfig, 'userconfig')

        //     // console.log('--allConfigs', allConfigs )
        //     this.setState({
        //         showDialog: true,
        //         curAppConfig: allConfigs, //userconfig
        //         selectAppKey: key 
        //     })
        // }
    }

    checkAllAppSelectServe(dependencies) {
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


    setConfig(type, appInfo, config){
        const changeConfig = (nodeId, conf, list) => {
            if(list && list.length) {
                for(let i = 0; i < list.length; i++){
                    if(nodeId == list[i].node_id){
                        list[i].conf = conf
                        break
                    }else {
                        if(list[i].dependences) {
                            changeConfig(nodeId, conf, list[i].dependences)
                        }
                    }
                }
            }
        }

        const changeWillInstance = (instanceId, selectWillInstance, list) => {
            if(list && list.length) {
                for(let i = 0; i < list.length; i++){
                    if(instanceId == list[i].instance_id){
                        list[i].selectWillInstance = selectWillInstance
                        break
                    }else {
                        if(list[i].dependences) {
                            changeWillInstance(instanceId, selectWillInstance, list[i].dependences)
                        }
                    }
                }
            }
        }

        let newList = this.state.appList
       

        if(type == 'external'){
            changeConfig(appInfo.node_id, config, newList)
        }else if(type === 'willHadInstance'){
            let findInstance = appInfo.willHadInstance.find((val) => val.instance_id == config)
            let selectedInstance = {app_id: findInstance.app_id, version_id: findInstance.version_id, instance_id: findInstance.instance_id}
            changeWillInstance(appInfo.instanceId, selectedInstance, newList)
        }

        this.setState({
            appList: newList
        },() => {
            // 点击后返回所有的app， 其实应该返回所有选中的app和该app的版本
            this.props.checkedAppList(newList)
        })
    }

    toConfig(appInfo) {
        event.stopPropagation();
        this.props.versionConfig(appInfo)
    }

    closeDialog() {
        this.setState({
            showDialog: false
        })
    }

    changeKey(key, idx) {
        let newConfig = this.state.curAppConfig.slice()
        newConfig[idx].val = event.target.value //newConfig[idx].type == 'number' ? Number(event.target.value) :  event.target.value
        newConfig[idx].error = ''
        this.setState({
            curAppConfig: newConfig
        })
    }

    confirmDialog() {
        let success = true
        let configs = this.state.curAppConfig.slice()

        configs.forEach((item) => {
            if(item.type !== 'object' && item.required) {
                if(item.val.trim() == '') {
                    item.error = '请输入'
                    success = false
                }
            }
        })

        if(success) {
            let newAppInfo = Object.assign({}, this.state.appInfo)
            const configs = (config, attr) => {
                if(config) {
                    // console.log('---renderconfig--', config.type , config.properties, config.type == 'object' && config.properties)
                    if(config.type == 'object' && config.properties) {
                        Object.keys(config.properties).forEach((key) => {
                            configs(config.properties[key], key)
                        })
                    }else {
                        console.log(222)
                        let idx = this.state.curAppConfig.findIndex((val) => val.key === attr)
                        config.val = this.state.curAppConfig[idx].type === 'number' ? Number(this.state.curAppConfig[idx].val) : this.state.curAppConfig[idx].val
                    }
                }
            }
            configs(newAppInfo.userconfig, 'userconfig') 

            console.log('--newAppInfo--', newAppInfo)

            this.setState({
                appInfo: newAppInfo
            })

            this.closeDialog()
        }else {
            this.setState({
                curAppConfig: configs
            })
        }

    }

    renderConfig(userconfig, attribute) {

        let doms = []
        const configs = (config, attr) => {
            if(config) {
                // console.log('---renderconfig--', config.type , config.properties, config.type == 'object' && config.properties)
                if(config.type == 'object' && config.properties) {
                    Object.keys(config.properties).forEach((key) => {
                        configs(config.properties[key], key)
                    })
                }else {
                    doms.push(
                        <div key={attr} className="config-item">
                            <div className="item-input">
                                <label>{config.required ? <span>* </span> : null}{attr}</label>
                                <input className={config.error ? 'red-border' : ''} onChange={this.changeKey.bind(this, attr)} value={config.value} />
                            </div>
                            {
                                config.error ? (
                                    <p className="item-error">{config.error}</p>
                                ) : null
                            }
                            
                        </div>
                    )
                }
            }
        }


        configs(userconfig, attribute)

        return doms
    }

    render(){
        return(
            <div className="app-tree-container" ref={(ele) => {this.treeContainerRef = ele}}
                style={{border: this.state.appInfo && this.state.appInfo.dependencies && Object.keys(this.state.appInfo.dependencies).length ? '1px solid #dcdcdc' : 'none'}}    
            >
                {
                    this.renderTree(this.state.appInfo).map((item, index) => item)
                }

                <BaseDialog
                    open={this.state.showDialog}
                    title={"应用配置"}
                    showCloseBtn={false}
                    // close={this.closeDialog.bind(this)}
                    confirm={this.confirmDialog.bind(this)}
                    contentClass="middle-dialog"
                >
                    <div className="appconfig-content">
                        {
                            this.state.curAppConfig.map((item, index) => {
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
                                                    onChange={this.changeKey.bind(this, item.key, index )} 
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

                </BaseDialog>

            </div> 
        )
    }
}