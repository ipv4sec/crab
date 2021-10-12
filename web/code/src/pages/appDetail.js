import React from 'react'
import {FlatButton, MenuItem, SelectField} from 'material-ui'
import HelpTip from '../components/helptip/index'
// import SVGComp from '../components/appdetail/svg'
import Graph from '../components/appdetail/graph'
import AppTree from '../components/appdetail/appTree'
import axios from 'axios'
import store from '../store/index'
import * as TYPE from '../store/actions'
import LoadingComp from '../components/showLoading/index'
import Snack from '../components/snackbar/index'
import BaseDialog from '../components/dialog/baseDialog'

const styles = {
    errorStyle: {
      color: '#EC5858',
    },
    underlineFocusStyle: {
      display: 'none',
    },
    underlineStyle: {
        display: 'none',
    },
    labelStyle:{
        fontSize: '14px'
    }
};
export default class AppDetail extends React.Component{
    constructor(props) {
        super(props)
        // console.log('params=',this.props.params)
        // console.log('appInfo===',sessionStorage.getItem('appInfo'))
        this.state={
            config: '',
            appId: this.props.params.appId,
            isEdit: this.props.params.isEdit == '1' ? true : false,
            version: '',
            allApps: [],
            configs: [
                // {label: '配置型1', code: "sdfdsfds", tooltip: "这是一个提示项", value: '10' },
                // {label: '配置型1', code: "sdfdsfds", tooltip: "这是一个提示项", value: '11' },
                // {label: '配置型1', code: "sdfdsfds", tooltip: "这是一个提示项" , value: '12'},
                // {label: '配置型1', code: "sdfdsfds", tooltip: "这是一个提示项", value: '13' },
                // {label: '配置型1', code: "sdfdsfds", tooltip: "这是一个提示项" , value: '14'},
                // {label: '配置型1', code: "sdfdsfds", tooltip: "这是一个提示项" , value: '15'},
                // {label: '配置型1', code: "sdfdsfds", tooltip: "这是一个提示项" , value: '16'},
                // {label: '配置型1', code: "sdfdsfds", tooltip: "这是一个提示项" , value: '17'},
            ],
            nameSpaceDetail: {},
            externalInfo: null,
            curAppInfo: null,
            configType: '',
            willInstances: [],
            showDialog: false
        }
        this.getedApps = false
        this.svgLoaded = false
    }

    componentWillMount(){
        // console.log('params=',this.props.params)
        // console.log('appInfo===',sessionStorage.getItem('appInfo'))
    }

    componentDidMount() {
      
       
    }

    getSpaceDetail(){
        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })
        axios({
            method: "GET",
            url: `${window._BASEPATH}/api/getSpaceDetail`,
            params: {
                id: this.props.params.appId,
            }
        }).then((res) => {
            
            this.setState({
                nameSpaceDetail: res.data
            })

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

    // 选择app版本后触发
    checkedAppList(data) {
        // 选择某个版本后会返回所有app树
        this.setState({
            allApps: data
        },() => {
            this.getedApps = true
            // 更新svg部分
            if(this.svgRef && this.svgLoaded){
                this.svgRef.setData(this.state.allApps)
            }
        })
    }

    // 点击请配置后触发
    versionConfig(data){
        if(data.version == 'external') {
            this.setState({
                configType: '',
                configs: [{key: '', value: ''}]
            },() => {
                let keyValues = []
                if(data.conf && Object.keys(data.conf).length) {
                    Object.keys(data.conf).forEach((key) => {
                        keyValues.push({
                            key: key,
                            value: data['conf'][key]
                        })
                    })
                }else {
                    keyValues = [{key: '', value: ''}]
                }
                // console.log('keyValues====',keyValues)
                this.setState({
                    curAppInfo: data,
                    configType: 'external',
                    // externalInfo: data,
                    configs:  keyValues,//[{key: '', value: ''}],
                    willInstances: [],
                    config:''
                })
            })
           
        }else if(data.version == 'willHadInstance'){
            this.setState({
                curAppInfo: data,
                configType: 'willHadInstance',
                configs: [],
                // externalInfo: {},
                willInstances: data.willHadInstance,
                config:  data.willHadInstance && data.willHadInstance.length && data.willHadInstance[0].instance_id ? data.willHadInstance[0].instance_id : ''
            },() => {
                this.appTreeRef.setConfig('willHadInstance',this.state.curAppInfo, this.state.config)
            })
            
        }else {
            this.setState({
                curAppInfo: {},
                configType: '',
                externalInfo: {},
                configs: [],
                willInstances: [],
                config: ''
            })
        }
    //    this.getAppConfig(data.appId, data.version)
    }

    addConfig() {
        let newConfigs = this.state.configs.slice()
        newConfigs.push({key: '', value: ''})
        this.setState({
            configs: newConfigs
        })
    }

    saveConfig() {
        let newConfigs = {}
        this.state.configs.forEach((item, index) => {
            if(this['config'+index] && this['config'+index].getData){
                let con = this['config'+index].getData()
                newConfigs[con.key] = con.value
            }
        })
        
        this.appTreeRef.setConfig('external',this.state.curAppInfo, newConfigs)
        store.dispatch({
            type: TYPE.SHOW_SNACKBAR,
            val: {open: true, message: '已保存配置'}
        })


    }

    svgLoadedFn(flag) {
        this.svgLoaded = flag
        if(this.getedApps && this.svgRef){
            this.svgRef.setData(this.state.allApps)
        }
    }
   
    // 获取该app某版本的配置信息
    getAppConfig(appId, version) {
        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })
        axios({
            method: "GET",
            url: `${window._BASEPATH}/api/getAppConfig`,
            params: {
                appId: appId,
                version: version
            }
        }).then((res) => {
            
            this.setState({
                configs: res.data.data
            })

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

    checkAllHasVersion(allApps) {
        let noHadVersion = {}
        const checkedVersion = (apps) => {
            if(apps && apps.length) {
                for(let i = 0; i < apps.length; i++) {
                    let versionNoCheck = false
                    if(apps[i].app_versions && apps[i].app_versions.length) {
                        for(let j = 0; j < apps[i].app_versions.length; j++){
                            if(apps[i].app_versions[j].selected){
                                versionNoCheck = true
                            }
                        }

                    }
                    if(versionNoCheck) {
                        if(apps[i].dependences) {
                            checkedVersion(apps[i].dependences)
                        }
                    }else {
                        noHadVersion = {
                            had: true,
                            appName: apps[i].app_name
                        }
                        break
                    }
                }
            }
        }
        checkedVersion(allApps)
        return noHadVersion
    }

    submit() {
        let _this = this
        let appList = this.appTreeRef.getAllApps()
        // 组装选择的应用及其版本，但是这个还是放到最后总体提交的时候进行组装，这个时候先带着所有的版本和数据，可以进行回显
        let checkApps = this.checkAllHasVersion(appList)
        if(!this.state.allApps.length || (checkApps && checkApps.had)){
            store.dispatch({
                type: TYPE.SHOW_SNACKBAR,
                val: {open: true, message: (checkApps.appName || '') +'未选择版本，请选择！'}
            })
        }else {
            window.opener.postMessage(JSON.stringify({action: 'refreshAppList',refreshPage: true, appData: appList, isEdit: this.state.isEdit}), '*')
        
            store.dispatch({
                type: TYPE.SHOW_SNACKBAR,
                val: {open: true, message: this.state.isEdit ? '修改成功' : '添加成功'}
            })
            setTimeout(() => {
                _this.confirmDialog()
            }, 1000)
        } 
       
    }

    close() {
        // window.close()
        this.setState({
            showDialog: true
        })
    }

    closeDialog() {
        this.setState({
            showDialog: false
        })
    }

    confirmDialog() {
        window.close()
    }



    render(){
        return(
            <div className="app-detail-container">
                <div className="app-detail-top">
                    <FlatButton
                        label="取消"
                        style={{color: '#6A6A6A'}}
                        onClick={this.close.bind(this)}
                        icon={<i style={{fontSize: '14px'}} className="iconfont icon_navigation_close"></i>}
                    />
                    <FlatButton
                        label="确认"
                        style={{color: '#4A90E2'}}
                        onClick={this.submit.bind(this)}
                        icon={<i style={{fontSize: '18px'}} className="iconfont icon_g-save"></i>}
                    />
                </div>

                <div className="app-detail-bottom">
                    <div className="detail-left">
                        <p className="detail-bottom-title" >应用</p>
                        <AppTree ref={(ele) => {this.appTreeRef = ele}} appId={this.state.appId} isEdit={this.state.isEdit} checkedAppList={this.checkedAppList.bind(this)} versionConfig={this.versionConfig.bind(this)}/>

                    </div>
                    <div className="detail-middle">
                        <p className="detail-bottom-title">配置</p>
                        <div className="detail-middle-content">
                            {
                                this.state.configType === 'external' ? (this.state.configs.map((item, index) => {
                                    return (
                                        <div className="detail-config" key={index}>
                                            <KeyValue data={item} ref={(ele) => {this['config'+index] = ele}} />
                                        </div>
                                    )
                                }) ): null
                            }
                            {
                                this.state.configType === 'external' ? (
                                    <div className="config-btns">
                                        <button onClick={this.addConfig.bind(this)}>添加</button>
                                        <button onClick={this.saveConfig.bind(this)}>保存</button>
                                    </div>
                                ) : null
                            }

                            

                            {/* {
                                this.state.configs.map((item, index) => {
                                    return (
                                        <div className="detail-config" key={index}>
                                            <p className="config-title">
                                                {item.label}
                                                <HelpTip  
                                                    size="14px"
                                                    content={`<p>${item.tooltip}</p>`}
                                                />
                                            </p>
                                            <pre>{item.code}</pre>
                                        </div>
                                    )
                                })
                            } */}

                            {
                                this.state.configType === 'willHadInstance' ? (
                                    <div className="detail-config">
                                        <p className="config-title">
                                            选择配置项
                                            <HelpTip  
                                                size="14px"
                                                content={`<p>请选择即将创建的实例</p>`}
                                            />
                                        </p>
                                        <div className="config-select">
                                            <SelectField
                                                value={this.state.config}
                                                labelStyle={styles.labelStyle}
                                                underlineFocusStyle={styles.underlineFocusStyle}
                                                underlineStyle={styles.underlineStyle}
                                                onChange={(event,index,value) => {
                                                    this.setState({
                                                        config: value
                                                    })
                                                    this.appTreeRef.setConfig('willHadInstance',this.state.curAppInfo, value)
                                                }}
                                                fullWidth={true}
                                            >
                                                {
                                                    this.state.willInstances.map((item, index) => {
                                                        return(
                                                            <MenuItem innerDivStyle={{color: '#4a4a4a'}} key={index} value={item.instance_id} primaryText={item.app_name+' : '+item.version_id} />
                                                        )
                                                    })
                                                }
                                            </SelectField>
                                        </div>
                                    

                                    </div>
                                ) : null
                            }
                            

                        </div>
                       
                    </div>
                    <div className="detail-right">
                        {/* <SVGComp appId={this.state.appId} allApps={this.state.allApps} ref={(ele) => {this.svgRef = ele}} svgLoaded={this.svgLoadedFn.bind(this)}/> */}
                        <Graph appId={this.state.appId} allApps={this.state.allApps} ref={(ele) => {this.svgRef = ele}} svgLoaded={this.svgLoadedFn.bind(this)}/>
                    </div>
                </div>

                <BaseDialog
                    open={this.state.showDialog}
                    title="关闭"
                    showCloseBtn={false}
                    close={this.closeDialog.bind(this)}
                    confirm={this.confirmDialog.bind(this)}
                    contentClass="small-dialog"
                >
                   <div className="confirm-dialog-content">
                        是否放弃编辑应用？
                   </div>
                </BaseDialog>

                <LoadingComp />
                <Snack />
            </div>
        )
    }
}


class KeyValue extends React.Component{
    constructor(props) {
        super(props)
        this.state = {
            key: this.props.data.key || '',
            value: this.props.data.value || ''
        }
    }

    // componentDidUpdate() {
    //     if(this.props.data && Object.keys(this.props.data).length) {
    //         this.setState({
    //             key: this.props.key,
    //             value: this.props.value
    //         })
    //     }
    // }

    getData() {
        return {
            key: this.state.key,
            value: this.state.value
        }
    }

    render() {
        return (
            <div className="config-input-content">
                <input className="config-input" value={this.state.key} onChange={(e) => {this.setState({key: e.target.value})}} /> :
                <input className="config-input" value={this.state.value} onChange={(e) => {this.setState({value: e.target.value})}} />
            </div>
        )
       
    }
}