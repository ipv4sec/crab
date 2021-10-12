import React from 'react';
import Popover from 'material-ui/Popover'
import Menu from 'material-ui/Menu'
import MenuItem from 'material-ui/MenuItem';
import { browserHistory } from 'react-router'
import store from '../store/index'
import * as TYPE from '../store/actions'
import copy from 'copy-to-clipboard'
import Moment from 'moment';
import Pagination from '../components/pagination/index'
import Search from '../components/search/index'
import axios from 'axios'

export default class SpacePage extends React.Component {
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
            searchObj: {},
            spaceData: [
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
        };
    }

    componentDidMount() {
        this.getSpaceData();
    }

    componentWillUnmount() {

    }
    
    getSpaceData() {
        let that = this

        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })
        let data = {
            offset: this.state.current_page * this.state.pageSize, 
            limit: this.state.pageSize
        }

        axios({
            method: 'GET',
            url: `${window._BASEPATH}/api/getAllSpace`,
            params: data
        }).then((res) => {
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
            that.setState({
                spaceData: res.data.list || [],
                total:res.data.total ? Math.ceil(res.data.total/this.state.pageSize) : 0,
                totalRes: res.data.total || 0
            })
            // else {
            //     that.setState({
            //         spaceData: [],
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
            this.getSpaceData()
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

    delete() {
        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })
        this.closePopover()

        axios({
            method: "POST",
            url: `${window._BASEPATH}/api/deleteSpace`,
            data: {
                spaceId: this.state.currentItem.id
            },
            headers: {'Content-Type': 'Application/json'}
        }).then((res) => {
            if(res.data.code == 200) {
                store.dispatch({
                    type: TYPE.SHOW_SNACKBAR,
                    val: {open: true, message: '删除成功'}
                })
                this.getSpaceData()
            }else  {
                store.dispatch({
                    type: TYPE.SHOW_SNACKBAR,
                    val: {open: true, message: res.data.message}
                })
            }
            
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
        }).catch((err) => {
            console.log(err)
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

    sysout() {
        // store.dispatch({
        //     type: TYPE.SHOW_LOADING,
        //     val: true
        // })
        this.closePopover()
        // axios({
        //     method: "GET",
        //     url: `${window._BASEPATH}/api/output`,
        //     params: {
        //         spaceId: this.state.currentItem.id
        //     }
        // }).then((res) => {
        //     store.dispatch({
        //         type: TYPE.SHOW_LOADING,
        //         val: false
        //     })
        // }).catch((err) => {
        //     console.log(err)
        //     store.dispatch({
        //         type: TYPE.SHOW_SNACKBAR,
        //         val: {open: true, message: '请求错误'}
        //     })
        //     store.dispatch({
        //         type: TYPE.SHOW_LOADING,
        //         val: false
        //     })
        // })
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
                this.getSpaceData()
            }
        );
    }

    closeUser() {
        this.setState({
            openUesr: false
        })
    }

    edit() {
        this.closePopover()
        // window.open(`/appdetail/${this.state.namespace_id.id}/1`, "_blank")
        browserHistory.push(`/space/edit/${this.state.currentItem && this.state.currentItem.id ?this.state.currentItem.id : 0}`)
    }


    render() {
  
        const { 
            spaceData, openMenu, anchorEl, showTipPop, tableTipEl, tableTipDesc, 
        } = this.state

        return (
            <div className="space-container" style={{minWidth: '1000px'}}>
                 <div className="space-title">
                    <span onClick={() => {browserHistory.push("/page")}}>空间管理</span>                    
                </div>
                <div className="top" >
                    <div className="searchBtn" style={{display: 'flex', marginRight: '10px'}}>
                        {/* <button className="staticSearchBtn" onClick={() => {browserHistory.push('/space/add/1/0')}}>添加</button> */}
                        <button className="staticSearchBtn" onClick={() => {browserHistory.push('/space/add')}}>添加</button>
                    </div>
                    {/* <div className="staticSearch">
                        <Search placeholder="搜索" callback={this.search.bind(this)} />
                    </div> */}
                </div>
                <div className="bottom">
                    <table className="data bottom-table staticTable">
                        <thead>
                            <tr>
                                <th width="30%">名称</th>
                                <th width="30%">ID</th>
                                <th width="30%">部署状态</th>
                                <th width="10%">操作</th>
                            </tr>
                        </thead>
                        <tbody style={{position: 'relative'}}>
                        {
                            spaceData.map((item, index) => {
                                return (    
                                    <tr key={index}
                                        onClick={() => {}}>
                                       
                                        <td>{item.name || ''}</td>
                                        <td>{item.id || ''}</td>
                                        <td>{
                                            item.state == 1 ? (
                                                <p className="published"><i className="iconfont icon_d-pass"></i>已完成</p>
                                            ) : (
                                                <p className="publishing"><i className="iconfont icon_assembly_details"></i>发布中...</p>
                                            )}
                                        </td>
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
                        <MenuItem style={{minHeight: '40px', lineHeight: '40px'}} onClick={this.edit.bind(this)}><div className="staticPopoverMenu"><i className="iconfont icon_menu2"></i>详情</div></MenuItem>
                        <MenuItem style={{minHeight: '40px', lineHeight: '40px'}} onClick={this.sysout.bind(this)} >
                            <a href={`/api/output?spaceId=${this.state.currentItem.id}`}><div className="staticPopoverMenu"><i className="iconfont icon_daochu"></i>导出配置</div></a>
                        </MenuItem>
                        <MenuItem style={{minHeight: '40px', lineHeight: '40px'}} onClick={this.delete.bind(this)}><div className="staticPopoverMenu"><i className="iconfont icon_baseline_delete"></i>删除</div></MenuItem>
                    </Menu>
                </Popover>
                
            </div> 
        );
    }
}