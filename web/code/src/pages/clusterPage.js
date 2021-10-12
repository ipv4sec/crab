import React from 'react';
import Popover from 'material-ui/Popover'
import Menu from 'material-ui/Menu'
import MenuItem from 'material-ui/MenuItem';
import store from '../store/index'
import s from '../store/index'
import * as TYPE from '../store/actions'
// import * as events from '../../libs/events';
import axios from 'axios'
import copy from 'copy-to-clipboard'
import Moment from 'moment';
// import SelectSearch from '../../components/selectSearch/index'
// import Pagination from '../../components/globalComponent/globalSearch/pagination/onebyone'
// import StaticUser from './StaticUser'
import Pagination from '../components/pagination/index'
import BaseDialog from '../components/dialog/baseDialog'
import AddCluster from './addCluster';
import Search from '../components/search/index'



export default class ClusterPage extends React.Component {
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
            clusterData: [
                // {name: '测试1', id: 'clst-foo1', status: '1'},
                // {name: '测试1', id: 'clst-foo2', status: '1'},
                // {name: '测试1', id: 'clst-foo3', status: '0'},
                // {name: '测试1', id: 'clst-foo4', status: '1'},
                // {name: '测试1', id: 'clst-foo5', status: '0'},
            ],
            errorNameTip: '',
            createDesc: '',
            createName: '',
           
            showTipPop: false,
            tableTipEl: null,
            tableTipDesc: '',

            addCluster: false,
            dialogTitle: '添加集群',
            type: 'add'
        };
    }

    componentDidMount() {
        this.getClusterData();
    }

    componentWillUnmount() {

    }
    
    getClusterData() {
        let that = this

        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })
        let data = {
            offset: this.state.current_page * this.state.pageSize || 0,
            limit: this.state.pageSize
        }
 
        axios({
            method: 'GET',
            url: `${window._BASEPATH}/api/getAllCluster`,
            params: data
        }).then((res) => {
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
            that.setState({
                clusterData: res.data.list || [],
                total: res.data.total ? Math.ceil(res.data.total/that.state.pageSize) : 0,
                totalRes: res.data.total || 0
            })
           
            // else {
            //     that.setState({
            //         clusterData: [],
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
            this.getClusterData()
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
                this.getClusterData()
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
            addCluster: false
        })
    }

    confirmDialog() {
        if(this.addClusterRef) {
            let data = this.addClusterRef.getData()
            if(data) {

                store.dispatch({
                    type: TYPE.SHOW_LOADING,
                    val: true
                })

                let formData = new FormData()
                formData.append("name", data.name)
                formData.append("host", data.host)
                formData.append("other", data.other)
                formData.append("file", data.file)

                // let postData = {
                //     name: data.name,
                //     host: data.host,
                //     certificate_url: data.file,
                //     description: data.other
                // }

                if(this.state.type === 'add') {
                    axios({
                        method: "POST",
                        url: `${window._BASEPATH}/api/addCluster`,
                        // headers: {"Content-Type": "application/json"},
                        headers: {"Content-Type": "application/x-www-form-urlencoded"},
                        data: formData
                    }).then((res) => {
                        if(res.data.id){
                            this.closeDialog()
                            this.getClusterData()
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
                    // postData['id'] = this.state.currentItem.id || ''
                    formData.append("id", this.state.currentItem.id || '') 
                    axios({
                        method: "POST",
                        url: `${window._BASEPATH}/api/editCluster`,
                        // headers: {'Content-Type': 'application/json'},
                        headers: {'Content-Type': 'application/x-www-form-urlencoded'},
                        data: formData
                    }).then((res) => {
                        if(res.data.state == 'success'){
                            this.closeDialog()
                            this.getClusterData()
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
    }

    editCluster() {
        this.setState({
            addCluster: true,
            dialogTitle: '编辑集群',
            openMenu: false,
            type: 'edit'
        })
    }


    render() {
        const { 
            clusterData, openMenu, anchorEl, showTipPop, tableTipEl, tableTipDesc, 
        } = this.state

        return (
            <div className="space-container" style={{minWidth: '1000px'}}>
                 <div className="space-title">
                    <span onClick={() => {browserHistory.push("/page")}}>集群管理</span>
                </div>
                <div className="top" >
                    <div className="searchBtn" style={{display: 'flex', marginRight: '10px'}}>
                        <button className="staticSearchBtn" onClick={() => {this.setState({type: 'add',addCluster: true,dialogTitle: '添加集群'})}}>添加</button>
                    </div>
                    {/* <div className="staticSearch">
                        <Search placeholder="搜索" callback={this.search.bind(this)} />
                    </div> */}
                </div>
                <div className="bottom">
                    <table className="data bottom-table staticTable">
                        <thead>
                            <tr>
                                <th width="45%">名称</th>
                                <th width="45%">ID</th>
                                <th width="10%">操作</th>
                            </tr>
                        </thead>
                        <tbody style={{position: 'relative'}}>
                        {
                            clusterData.map((item, index) => {
                                return (    
                                    <tr key={index}
                                        onClick={() => {}}>
                                        <td>{item.name}</td>
                                        <td>{item.id}</td>
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
                        <MenuItem style={{minHeight: '40px', lineHeight: '40px'}} onClick={this.editCluster.bind(this)}><div className="staticPopoverMenu"><i className="iconfont icon_baseline_edit"></i>编辑</div></MenuItem>
                    </Menu>
                </Popover>


                {
                    this.state.addCluster ? (
                        <BaseDialog
                            open={this.state.addCluster}
                            title={this.state.dialogTitle}
                            close={this.closeDialog.bind(this)}
                            confirm={this.confirmDialog.bind(this)}
                            contentClass="middle-dialog"
                        >
                            <AddCluster type={this.state.type} ref={(ele) => {this.addClusterRef = ele}} data={this.state.currentItem}/>
                        </BaseDialog>
                    ) : null
                }



            </div> 
        );
    }
}