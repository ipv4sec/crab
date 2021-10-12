import React from 'react';
import Popover from 'material-ui/Popover'
import Menu from 'material-ui/Menu'
import MenuItem from 'material-ui/MenuItem';
import store from '../store'
import * as TYPE from '../store/actions'
import copy from 'copy-to-clipboard'
import Moment from 'moment';
import Pagination from '../components/pagination/index'
import BaseDialog from '../components/dialog/baseDialog'
import Search from '../components/search/index'
import AddApp from './addApp'
import Versions from './versions'
import axios from 'axios'
import { select } from 'underscore';


export default class AppPage extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            navList: [],
            name: '',
            current_page: 0,
            pageSize: 20,
            total: 0,
            totalRes: 0,
            openMenu: false,
            anchorEl: null,
            currentItem: {}, // 当前碎片信息
            search: '',
            appList: [
                // {name: '测试1', id: 'clst-foo1', status: '1', color: 'blue', git: 'dsfsfsfds', description: 'dsfdsfdsfdsf'},
                // {name: '测试1', id: 'clst-foo2', status: '1', color: 'red', git: 'dsfsfsfds', description: 'dsfdsfdsfdsf'},
                // {name: '测试1', id: 'clst-foo3', status: '0', color: 'yellow', git: 'dsfsfsfds', description: 'dsfdsfdsfdsf'},
                // {name: '测试1', id: 'clst-foo4', status: '1', color: 'green', git: 'dsfsfsfds', description: 'dsfdsfdsfdsf'},
                // {name: '测试1', id: 'clst-foo5', status: '0', color: 'black', git: 'dsfsfsfds', description: 'dsfdsfdsfdsf'},
            ],
            errorNameTip: '',
            createDesc: '',
            createName: '',
           
            showTipPop: false,
            tableTipEl: null,
            tableTipDesc: '',

            type: 'add',
            addApp: false,
            dialogTitle: '添加应用',

            showVersions: false,
            versionTitle: '版本列表',


        };
    }

    componentDidMount() {
        this.getAppList();
    }

    componentWillUnmount() {

    }
    
    getAppList() {
        let that = this

        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })
        let data = {
            offset: this.state.current_page * this.state.pageSize,
            limit: this.state.pageSize,
            // search: this.state.search
        }
        axios({
            method: 'GET',
            url: `${window._BASEPATH}/manager/getAllApp`,
            params: data,
        }).then((res) => {
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
            that.setState({
                appList: res.data.rows || [],
                total: res.data.total ? (Math.ceil(res.data.total / that.state.pageSize)) : 0,
                totalRes: res.data.total || 0
            })


            // else {
            //     that.setState({
            //         appList: [],
            //         total: 0,
            //         totalRes: 0
            //     })
            // }
        }).catch((err) => {
            console.log(err)
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
        })
    }


    search(data) {
        this.setState({
            current_page: 0,
            search: data
        }, () => {
            this.getAppList()
        })
       
    }

    create () {
      
    }

    closePopover(auto) {
        if(auto) {
            this.setState({
                currentItem: {}
            })
        }
        this.setState({
            anchorEl: null,
            openMenu: false,
        })
    }

    copyRoute() {
        const {currentItem} = this.state
        copy(currentItem.route)
        // events.customEvent.emit(events.SHOW_SNACK, '已复制到剪切板!');
        this.setState({
            currentItem: {}
        })
        this.closePopover()
    }

     // change事件
     changeFragmentName(e) {
        this.setState({
            errorNameTip: '',
            createName: e.target.value
        })
    }


    // 失去焦点验证
    linkBlur() {
        if(this.state.createName !== '' && this.state.createName.length < 30) {
            this.setState({
                errorNameTip: ''
            })
        }else {
            this.setState({
                errorNameTip: '名称不能为空且长度不超30个字符'
            })
        }
    }

  
   
    showTip(desc) {
        if(desc && desc.length && (desc.length * 12 > event.target.offsetWidth)) {
            this.setState({
                tableTipEl: {x: event.pageX, y: event.pageY},
                showTipPop: true,
                tableTipDesc: desc
            })
        }
       
    }

    hideTip(e) {
        this.setState({
            tableTipEl: null,
            showTipPop: false,
            tableTipDesc: ''
        })
    }

    addDialogClose() {
        this.setState({
            addDialogTitle: '',
            addDialogOpen: false,
            currentItem: {}
        })
    }


    handlePageChanged(newPage) {
        this.setState(
            {current_page: newPage},
            () => {
                this.getAppList()
            }
        );
    }

    closeUser() {
        this.setState({
            openUesr: false
        })
    }

    closeDialog() {
        this.setState({
            addApp: false
        })
    }

    uninstallApp() {
        // this.state.currentItem
        this.closePopover()
        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })

        axios({
            method: "POST",
            url: `${window._BASEPATH}/manager/uninstallApp`,
            headers: {'Content-Type': 'application/json'},
            data: {instanceid: this.state.currentItem.id, status: 3}
        }).then((res) => {
            this.closeDialog()
            this.getAppList()

            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })

            store.dispatch({
                type: TYPE.SHOW_SNACKBAR,
                val: {open: true, message: res.data.result}
            })
            
        }).catch((err) => {
            store.dispatch({
                type: TYPE.SHOW_SNACKBAR,
                val: {open: true, message: '请求错误'}
            })
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
        })
    }

    deleteApp() {
        // this.state.currentItem
        this.closePopover()
        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })

        axios({
            method: "GET",
            url: `${window._BASEPATH}/manager/deleteApp`,
            headers: {'Content-Type': 'application/json'},
            params: {instanceid: this.state.currentItem.id}
        }).then((res) => {
            this.closeDialog()
            this.getAppList()

            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })

            store.dispatch({
                type: TYPE.SHOW_SNACKBAR,
                val: {open: true, message: res.data.result}
            })
            
        }).catch((err) => {
            store.dispatch({
                type: TYPE.SHOW_SNACKBAR,
                val: {open: true, message: '请求错误'}
            })
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
        })
    }

    confirmDialog() {
        let data = this.addAppRef.getData()
        if(data.notHadServe.length) {
            // 依赖中存在某些应用没有服务的情况
            store.dispatch({
                type: TYPE.SHOW_SNACKBAR,
                val: {open: true, message: data.notHadServe.join('、') + '已上应用中不存在服务，请创建'}
            })
            return 
        }
        if(data.allAppSelectServe.length) {
            // 依赖中存在没有选择服务的应用
            store.dispatch({
                type: TYPE.SHOW_SNACKBAR,
                val: {open: true, message: data.allAppSelectServe.join('、') + '已上应用未选择服务，请选择'}
            })
            return
        }

        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })

        let selectData = this.getAppSelect(data.appInfo)

        selectData['status'] = 1

        axios({
            method: "POST",
            url: `${window._BASEPATH}/manager/run`,
            headers: {'Content-Type': 'application/json'},
            data: selectData
        }).then((res) => {

            console.log('res==',res)
            store.dispatch({
                type: TYPE.SHOW_SNACKBAR,
                val: {open: true, message: res.data.result}
            })
            this.getAppList()
           
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
           this.closeDialog()
            
        }).catch((err) => {
            store.dispatch({
                type: TYPE.SHOW_SNACKBAR,
                val: {open: true, message: '请求错误'}
            })
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
        })


        return

        if(data) {
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: true
            })
            if(this.state.type === 'add') {
                axios({
                    method: "POST",
                    url: `${window._BASEPATH}/api/addApp`,
                    headers: {'Content-Type': 'application/json'},
                    data: data
                }).then((res) => {
                    if(res.data.id){
                        this.closeDialog()
                        this.getAppList()
                    }
                    store.dispatch({
                        type: TYPE.SHOW_LOADING,
                        val: false
                    })
                    
                }).catch((err) => {
                    store.dispatch({
                        type: TYPE.SHOW_SNACKBAR,
                        val: {open: true, message: '请求错误'}
                    })
                    store.dispatch({
                        type: TYPE.SHOW_LOADING,
                        val: false
                    })
                })
            }else if(this.state.type === 'edit') {
                data['id'] = this.state.currentItem.id || '' // app_id需要
                axios({
                    method: "POST",
                    url: `${window._BASEPATH}/api/editApp`,
                    headers: {'Content-Type': 'application/json'},
                    data: data
                }).then((res) => {
                    if(res.data.state == 'success'){
                        this.closeDialog()
                        this.getAppList()
                    }
                    store.dispatch({
                        type: TYPE.SHOW_LOADING,
                        val: false
                    })
                    
                }).catch((err) => {
                    store.dispatch({
                        type: TYPE.SHOW_SNACKBAR,
                        val: {open: true, message: '请求错误'}
                    })
                    store.dispatch({
                        type: TYPE.SHOW_LOADING,
                        val: false
                    })
                })
            }
        }
    }

    getAppSelect(data) {
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
                                "name": key,
                                "instanceid": item.instance.instanceid
                            })
                        }
                    })
                }
            })

        }
        if(data.userconfig) {
            configs(data.userconfig, 'userconfig', appConfig)
        }
       
        return {
            instanceid: data.instanceid,
            dependencies: JSON.stringify(selectData),
            userconfig: JSON.stringify(appConfig.userconfig ? appConfig.userconfig : {})
        }
    }

    closeVersionDialog() {
        this.setState({
            showVersions: false
        })
    }


    render() {
  
        const { 
            appList, openMenu, anchorEl, showTipPop, tableTipEl, tableTipDesc, 
        } = this.state

        return (
            <div className="space-container" style={{minWidth: '1000px'}}>
                 <div className="space-title">
                    {/* <span onClick={() => {browserHistory.push("/")}}>应用管理</span> */}
                    <span>应用管理</span>
                </div>
                <div className="top" >
                    <div className="searchBtn" style={{display: 'flex', marginRight: '10px'}}>
                        <button className="staticSearchBtn" onClick={() => {this.setState({addApp: true, dialogTitle: '添加应用', type: 'add'})}}>添加</button>
                    </div>
                    {/* <div className="staticSearch">
                        <Search placeholder="搜索" callback={this.search.bind(this)} />
                    </div> */}
                </div>
                <div className="bottom">
                    <table className="data bottom-table staticTable">
                        <thead>
                            <tr>
                                <th width="20%">名称</th>
                                <th width="20%">状态</th>
                                <th width="20%">版本</th>
                                <th width="30%">链接</th>
                                <th width="10%">操作</th>
                            </tr>
                        </thead>
                        <tbody style={{position: 'relative'}}>
                        {
                            appList.map((item, index) => {
                                return (    
                                    <tr key={item.instance_id}
                                        onClick={() => {}}>
                                      
                                        <td >
                                            <div className="app-td">
                                                {item.name}
                                            </div>
                                        </td>
                                        <td>{item.status}</td>
                                        <td>{item.version}</td>
                                        <td className="list-entry"><a href={item.entry} target="_blank">{item.entry}</a></td>
                                        <td onClick={(e) => {
                                            e.stopPropagation()
                                            this.setState({
                                                anchorEl: e.target,
                                                openMenu: true,
                                                currentItem: item
                                            })
                                        }}><i className="iconfont icon_navigation_more" style={{cursor: "pointer"}}></i></td>
                                    </tr>
                                )
                            })
                        }
                         {
                                showTipPop ? (
                                    <div
                                        className="showTablePop"
                                        style={{
                                            left: tableTipEl.x + 10 + 'px',
                                            top: tableTipEl.y - 20 + 'px',
                                        }}
                                    >
                                    {tableTipDesc}
                                    </div>
                                ) : null
                            }
                        </tbody>
                    </table>
                     <div className="page-pagination" style={{marginBottom: '20px'}}>
                        <Pagination
                            total={this.state.total}
                            totalRes={this.state.totalRes}
                            current={this.state.current_page}
                            className="pagination-box"
                            handleChange={this.handlePageChanged.bind(this)}
                        />
                    </div>
                </div> 
               
                        
                <Popover
                    open={openMenu}
                    anchorEl={anchorEl}
                    anchorOrigin={{horizontal: 'left', vertical: 'bottom'}}
                    targetOrigin={{horizontal: 'right', vertical: 'top'}}
                    onRequestClose={this.closePopover.bind(this, true)}
                >
                    <Menu>
                        {
                            this.state.currentItem.status === "卸载完成" ? (
                                <MenuItem key='1' style={{minHeight: '40px', lineHeight: '40px'}} onClick={this.deleteApp.bind(this)}><div className="staticPopoverMenu"><i className="iconfont icon_baseline_delete"></i>删除</div></MenuItem>
                            ) : (
                                <MenuItem key='1' style={{minHeight: '40px', lineHeight: '40px'}} onClick={this.uninstallApp.bind(this)}><div className="staticPopoverMenu"><i className="iconfont icon_baseline_delete"></i>卸载</div></MenuItem>
                            )
                        }
                        {/* <MenuItem key='1' style={{minHeight: '40px', lineHeight: '40px'}} onClick={this.deleteApp.bind(this)}><div className="staticPopoverMenu"><i className="iconfont icon_baseline_delete"></i>卸载</div></MenuItem> */}
                        {/* <MenuItem key='0' style={{minHeight: '40px', lineHeight: '40px'}} onClick={() => {this.setState({addApp: true, dialogTitle: '编辑应用', type: 'edit',  openMenu: false})}}><div className="staticPopoverMenu"><i className="iconfont icon_baseline_edit"></i>编辑</div></MenuItem> */}
                    </Menu>
                </Popover>


                {
                    this.state.addApp ? (
                        <BaseDialog
                            open={this.state.addApp}
                            title={this.state.dialogTitle}
                            close={this.closeDialog.bind(this)}
                            confirm={this.confirmDialog.bind(this)}
                            contentClass="middle-dialog"
                        >
                            <AddApp  
                                type={this.state.type} 
                                ref={(ele) => {this.addAppRef = ele}}
                                data={this.state.currentItem}
                            />
                        </BaseDialog>
                    ) : null
                }

                {
                    this.state.showVersions ? (
                        <BaseDialog
                            open={this.state.showVersions}
                            title={this.state.versionTitle}
                            showCloseBtn={true}
                            close={this.closeVersionDialog.bind(this)}
                            contentClass="middle-dialog"
                        >
                            <Versions  
                                ref={(ele) => {this.versionRef = ele}}
                                data={this.state.currentItem}

                            />
                        </BaseDialog>
                    ) : null
                }

            </div> 
        );
    }
}

