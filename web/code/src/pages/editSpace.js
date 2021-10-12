import React, { version } from 'react'
import { TextField, SelectField, MenuItem } from 'material-ui';
import SearchBox from '../components/search/index'
import BaseDialog from '../components/dialog/baseDialog'
import store from '../store/index'
import * as TYPE from '../store/actions'
import axios from 'axios';
import copy from 'copy-to-clipboard'
import { browserHistory } from 'react-router'
import JsBeautify from 'js-beautify'
const appDetail = `[{"app_id":"app-jlyN52Rm","instance_id":"ais-891ycdvw","version_id":"v0.0.6","dependences":[{"app_id":"app-Gd9sfVSa","instance_id":"ais-lnsgoa01","version_id":"v0.0.2","dependences":[{"app_id":"app-b3Dwlhzc","instance_id":"ais-nruebpei","version_id":"v0.0.2","dependences":[{"app_id":"app-jlyN52Rm","version_id":"v0.0.6"}]}]}]}]`

const styles = {
    errorStyle: {
      color: '#EC5858',
    },
    underlineFocusStyle: {
      borderColor: '#3986FF',
    },
};

export default class EditSpace extends React.Component {
    constructor(props) {
        super(props)
        this.state ={
            spaceId: this.props.params.spaceId,
            spaceInfo: {},
            appList: [
                // {app_id: '111', version_id: '1'},
                // {app_id: '222', version_id: '2'},
                // {app_id: '333', version_id: '3'},
            ],
            instanceList: [
                // {app_id: '111',instance_id: '111', version_id: '1'},
                // {app_id: '222', version_id: '2',instance_id: '222'},
                // {app_id: '333', version_id: '3',instance_id: '333'},
            ],
            open: false,
            code: ''

        }
    }

    componentDidMount() {

       this.getSpaceDetail()
       this.getInstances()
        
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
                id: this.state.spaceId,
            }
        }).then((res) => {
            let appList = []
            try{
                appList = JSON.parse(res.data.app_info)

            }catch(err) {
                console.log(err)
            }
            this.setState({
                spaceInfo: res.data,
                appList: appList
            }, () => {
                // this.getAllInstances()
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

    async getAllInstances() {
        let allInstance = []
        let data = this.state.appList
        for(let i = 0; i < data.length; i++) {
            let instance = await this.getApps(data[i].app_id)
            allInstance.concat(instance)
        }
        this.setState({
            instanceList: allInstance
        })
        // return allInstance
    }

    getApps(appId){
        let _this = this

        return new Promise((resolve, reject) => {
            axios({
                method: "GET",
                url: `${window._BASEPATH}/api/getSpaceApps`,
                params: {
                    spaceId: _this.state.spaceId,
                    appId: appId
                }
            }).then((res) => {
               
               resolve(res && res.data && res.data.list ? res.data.list : []) 
            
            }).catch((err) => {
                console.log(err)
               reject(err)
            })
        })
    }

    getInstances(){
        let _this = this

        axios({
            method: "GET",
            url: `${window._BASEPATH}/api/getSpaceApps`,
            params: {
                spaceId: _this.state.spaceId,
            }
        }).then((res) => {
            let allInstance = res && res.data && res.data.list ? res.data.list : []

            _this.setState({
                instanceList: allInstance
            },() => {
                allInstance.forEach((item) => {
                    this[item.id].value = item.branch || ''
                })
            })
        
        }).catch((err) => {
            console.log(err)
        })
    }



    addApp() {
        this.setState({
            showDialog: true
        })
    }

    closeDialog() {
        this.setState({
            showDialog: false,
            appList: []
        })
    }

    showConfig(data) {
//         this.setState({
//             open: true,
//             code:JsBeautify.js(appDetail) 
//         })
// return
        let that = this
        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })

        axios({
            method: 'GET',
            url: `${window._BASEPATH}/api/getVersionDetail`,
            params: {appId: data.app_id, version: data.version_id}
        }).then((res) => {
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
            that.setState({
                open: true,
                code: res.data ? JsBeautify.js(JSON.stringify(res.data)) : ''
            })
        }).catch((err) => {
            console.log(err)
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
        })

    }

    close(){
        this.setState({
            open: false
        })
    }

    copyData(data){
        // copy('自动构建hook')
        // https://{asb 后台 host}/app/instance/rebuild?instance_id={当前应用实例 id}
        copy(`${window._ENDHOST}/app/instance/rebuild?instance_id=${data.id}`)
        store.dispatch({
            type: TYPE.SHOW_SNACKBAR,
            val: {open: true, message: "已复制到剪切板"}
        })
    }

    saveBranch(data) {

        let _this = this
        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })
        axios({
            method: "POST",
            url: `${window._BASEPATH}/api/saveBranch`,
            data: {
                branch: _this[data.id].value,
                id: data.id
            }
        }).then((res) => {
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
            if(res.data.state == 'success') {
                store.dispatch({
                    type: TYPE.SHOW_SNACKBAR,
                    val: {open: true, message: '已更新实例分支'}
                })
            }else {
                store.dispatch({
                    type: TYPE.SHOW_SNACKBAR,
                    val: {open: true, message: '响应错误'}
                })
            }

        }).catch((err) => {
            console.log(err)
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
        })
    }


    render() {
        return (
           <div className="add-space-content">
                 <div className="space-title">
                    <span className="gray-title" onClick={() => {browserHistory.push("/space")}}>空间管理</span>
                    <span className="second-title" ><i className="iconfont icon_arrow-right"></i>详情</span>
                </div>

                <div className="space-input" style={{marginTop: '20px'}}> 
                    <label>名称</label>
                    <p className="space-text">{this.state.spaceInfo.name}</p>
                </div>
                <div className="space-input" style={{marginTop: '20px'}}> 
                    <label>ID</label>
                    <p className="space-text">
                        {this.state.spaceInfo.id} 
                        {this.state.spaceInfo.id ? (<button className="copy-btn" onClick={this.copyData.bind(this, this.state.spaceInfo.id)}>复制</button>) : null }
                    </p>
                </div>
                <div className="space-input" style={{marginTop: '20px'}}> 
                    <label>状态</label>
                    <p className="space-text">{this.state.spaceInfo.state ? '发布中' : '已发布'}</p>
                </div>

                <div className="space-app-table">
                    <p className="list-title">实例</p>

                    <table className="app-table">
                        <thead>
                            <tr>
                                <th width="27%">名称</th>
                                <th width="10%">版本</th>
                                <th width="10%">状态</th>
                                <th width="30%">自动构建</th>
                                <th width="15%">自动构建hook</th>
                                <th width="13%">应用配置</th>
                            </tr>
                        </thead>
                        <tbody style={{position: 'relative'}}>
                        {
                            this.state.instanceList.map((item, index) => {
                                return (    
                                    <tr key={index}>
                                        <td >
                                            <div className="app-td">
                                                <div style={{backgroundColor: item.color || '#54CACB'}} className="app-td-icon"><i className="iconfont icon_grey600"></i></div>
                                                {item.id || item.app_id}
                                            </div>
                                        </td>
                                        <td>{item.version_id}</td>
                                        <td>状态</td>
                                        <td><div className="td-input"><input ref={(ele) => {this[item.id] = ele}} placeholder="输入branch" /><button onClick={this.saveBranch.bind(this, item)}><i className="iconfont icon_daochu"></i></button></div></td>
                                        <td><span onClick={() => {this.copyData(item)}} className="show-pointer">复制</span></td>  
                                        <td><span onClick={this.showConfig.bind(this, item)} className="show-pointer">查看</span></td>
                                    </tr>
                                )
                            })
                        }
                        
                        </tbody>
                    </table>
                </div>

                <div className="space-app-log"> 
                    <p className="list-title">log</p>
                    <p className="space-textare">{this.state.spaceInfo.description || ''}</p>
                </div>




                <BaseDialog
                    open={this.state.open}
                    title="应用配置"
                    showCloseBtn={true}
                    close={this.close.bind(this)}
                    contentClass="middle-dialog"
                >
                   <div className="show-code">
                       <pre>
                            {this.state.code}
                       </pre>
                   </div>
                </BaseDialog>
           </div>
        )
    }
}