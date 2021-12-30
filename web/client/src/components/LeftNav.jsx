import React, { useState, useEffect } from 'react'
import { withRouter } from 'react-router-dom'
import { connect } from 'react-redux'
import "../style/sass/components.scss"
import store from '../store/store'
import * as TYPE from '../store/actions'
import axios from 'axios'

let list = [
                {
                    "name": "应用管理",
                    "url": "/home",
                    "sub": [],
                },
                {
                    "name": "工作负载",
                    "url": "/home/workload",
                    "sub": [],
                },
                // { // 先注释掉，后期可能会用到
                //     "name": "根域设置",
                //     "url": "/home/domain",
                //     "sub": [],
                // },
                {
                    "name": "密码设置",
                    "url": "/home/reset",
                    "sub": [],
                },
                // {
                //     "name": "在线工具",
                //     "url": "",
                //     "showChild": false,
                //     "sub": [
                //         {
                //             "name": "创建应用",
                //             "url": "/home/online/create",
                //             "sub": []
                //         },
                //         {
                //             "name": "创建Trait",
                //             "url": "/home/online/trait",
                //             "sub": []
                //         },
                //         {
                //             "name": "创建WorkloadType",
                //             "url": "/home/online/workloadtype",
                //             "sub": []
                //         },
                //         {
                //             "name": "创建WorkloadVendor",
                //             "url": "/home/online/workloadvendor",
                //             "sub": []
                //         }
                //     ],
                // }
            ]

const LeftNav = (props) => {
    if(!sessionStorage.getItem('curNav')) {
        sessionStorage.setItem('curNav', '/home')
    }
    const [curNav, setCurNav] = useState(sessionStorage.getItem('curNav'))
    const [navList, setNavList] = useState([])

    // console.log('leftNav props=',props)
    useEffect(() => {
        let menus = handleMenu(list)
        setNavList(menus)
    }, [])

    useEffect(() => {
        sessionStorage.setItem('curNav', props.common.curNav)
        setCurNav(props.common.curNav)
    }, [props.common.curNav])

    function getMenus() {
        axios({
            method: 'GET',
            url: '/api/cluster/menus'
        }).then((res) => {
            if(res.data.code === 0) {
                let menus = handleMenu(res.data.result)
                setNavList(menus)
            }else {
                store.dispatch({
                    type: TYPE.SNACKBAR,
                    val: res.data.result || '获取列表服务错误'
                })
            }
        }).catch((err) => {
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: '请求错误'
            })
            store.dispatch({
                type: TYPE.LOADING,
                val: false
            })
            console.log(err)
        })
    }

    function handleMenu(menus) {
        const deal = (list, level) => {
            if(!Array.isArray(list)) { return }
            list.forEach((item, index) => {
                item['id'] = level + '.' +index
                if(item.sub) {
                    deal(item.sub, index+1)
                }
            })
        }

        deal(menus, 1)

        return menus
    }

    function changeNav(e){
        let customData = e.currentTarget.dataset
        let id = customData.id
        let path = customData.path
        let href = 'https://' + customData.href
    
        const find = (data, id) => {
            if(data && Array.isArray(data)) {
                for(let i = 0, len = data.length; i < len; i++) {
                    if(data[i].id === id) {
                        data[i].showChild = !data[i].showChild
                        break
                    }
                    if(data[i].children) {
                        find(data[i].children, id)
                    }
                }
            }
        }

        if(path === '') {
            let newList = navList.slice()
            find(newList, id)
            setNavList(newList)
        }else {
            setCurNav(path)
            // console.log('path==',path)
            sessionStorage.setItem('curNav', path)
           
            if(path.indexOf('.') > -1) {
                props.change('/system')
                store.dispatch({
                    type: TYPE.IFRAME_SRC,
                    val: href
                })
            }else {
                props.change(path)
            }
        }
    }

    return (
        <div className="nav-container">
            <div className="logo">
                <p className="logo-title">Crab</p>
            </div>
            <div className="nav-list">
                {
                    navList.map((item, index) => {
                        return (
                            <div className="list-item" key={item.id}>
                                <div className={`item-content ${curNav == item.url ? "blueBorder" :""}`}  data-id={item.id} data-path={item.url} onClick={changeNav}>
                                    <i className={`iconfont ${item.icon || ''}`}></i>
                                    <span >{item.name}</span> 
                                </div>
                                {
                                    item.sub && item.showChild ? (
                                        item.sub.map((ele, i) => {
                                            return (
                                                <div className="list-item"  key={ele.id}>
                                                    <div className={`item-content item-content-child ${curNav == ele.url ? "blueBorder" :""}`} data-id={ele.id} data-path={ele.url} data-href={ele.url || ''} onClick={changeNav}>
                                                        <i className={`iconfont ${ele.icon || ''}`}></i>
                                                        <span >{ele.name}</span> 
                                                    </div>
                                                </div>
                                            )
                                        })
                                    ) : null
                                }
                            </div>
                        )
                    })
                }
            </div>
        </div>  
    )

}
function mapPropsToState(state) { return state }
export default connect(mapPropsToState)(withRouter(LeftNav))