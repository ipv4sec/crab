import React, {useState, useEffect, lazy, Suspense} from 'react'
import { Route, Switch, useHistory } from 'react-router-dom'
import { connect } from 'react-redux'
import SuspenseLoading from '../../components/SuspenseLoading'
import axios from 'axios'
import '../../style/sass/detail.scss'
import store from '../../store/store'
import * as TYPE from '../../store/actions'
import moment from 'moment'
import {js_beautify} from 'js-beautify'

import DetailNav from './DetailNav'
import List from './List'
import Log from './Log'
import Tail from './Tail'
import Crumbs from './DetailCrumb'
import Desc from './Desc'

// const Tail = lazy(() => import('./Tail'))
// const Log = lazy(() => import('./Log'))


const defaultCrumbs = [
    'workload', 'pods'
]

const defaultNavList = [
    {
        id: 1,
        name: 'workload',
        children: [
            {name: 'pods1'},{name: 'pods2'},{name: 'pods3'},{name: 'pods4'},
        ]
    },
    {
        id: 2,
        name: 'workload2',
        children: [
            {name: 'pods1'},{name: 'pods2pods2pods2pods2pods2pods2pods2pods2'},{name: 'pods3'},{name: 'pods4'},
        ]
    }
]

const navlist = {
    "cronJob": [{
        "metadata": {
          "name": "charlie-ndgfb1"
        }
      },{
        "metadata": {
          "name": "charlie-ndgfb2"
        }
      }],
      "daemonSet": [{
        "metadata": {
          "name": "charlie-ndgfb3"
        }
      },{
        "metadata": {
          "name": "charlie-ndgfb4"
        }
      }]
}

const listKeys = {
    'cronJob': []
}

const Detail = (props) => {
    const history = useHistory()
    const [type, setType] = useState('list')
    const [instanceInfo, setInstanceInfo] = useState({})
    const [detailData, setDetailData] = useState({})
    const [detailDesc, setDetailDesc] = useState('')
    const [instanceList, setInstanceList] = useState([])
    const [crumbs, setCrumbs] = useState([])
    const [navList, setNavList] = useState([])
    const [curNav, setCurNav] = useState('')
    const [tailData, setTailData] = useState({})
    const [listData, setListData] = useState([])
    const [descName, setDescName] = useState('')
    const [logName, setLogName] = useState('')

    useEffect(() => {
        store.dispatch({
            type: TYPE.SET_INSTANCE_INFO,
            val: props.match.params
        })
        setInstanceInfo(props.match.params)
        getDetailData()
    },[])

    const parseDetailData = (data) => {
        // console.log('--parseDetailData--', data)
        if(!data || !data.details || !(Object.keys(data.details).length)) { return }

        let detailInfo = []
        Object.keys(data.details).forEach(key => {
            let tmp = {
                id: key,
                name: key,
                header: [],
                body: [],
                navList: [],
            }

            switch(key){
                case 'cronJob':
                    tmp.header = ['名称', '创建时间', '最后执行时间']
                    data.details[key].forEach((el) => {
                        tmp.navList.push({id: el.metadata.uid || '', name: el.metadata.name || ''})
                        tmp.body.push([el.metadata.name || '', moment(el.metadata.creationTimestamp || '').format('YYYY-MM-DD hh:mm:ss') || '', moment(el.status.lastScheduleTime || '').format('YYYY-MM-DD hh:mm:ss') || ''])
                    })
                    break;
                
                case 'daemonSet':
                    tmp.header = ['名称', '创建时间', '当前可用数']
                    data.details[key].forEach((el) => {
                        tmp.navList.push({id: el.metadata.uid || '', name: el.metadata.name || ''})
                        tmp.body.push([el.metadata.name || '',moment(el.metadata.creationTimestamp || '').format('YYYY-MM-DD hh:mm:ss') || '', el.status.numberAvailable || 0])
                    })
                    break;
                case 'deployment':
                    tmp.header = ['名称', '创建时间', '当前可用数']
                    data.details[key].forEach((el) => {
                        tmp.navList.push({id: el.metadata.uid || '', name: el.metadata.name || ''})
                        tmp.body.push([el.metadata.name || '',moment(el.metadata.creationTimestamp || '').format('YYYY-MM-DD hh:mm:ss') || '', el.status.numberAvailable || 0])
                    })
                    break
                case 'job':
                    tmp.header = ['名称', '创建时间', '完成时间']
                    data.details[key].forEach((el) => {
                        tmp.navList.push({id: el.metadata.uid || '', name: el.metadata.name || ''})
                        tmp.body.push([el.metadata.name || '', moment(el.metadata.creationTimestamp || '').format('YYYY-MM-DD hh:mm:ss') || '', moment(el.status.completionTime || '').format('YYYY-MM-DD hh:mm:ss') || ''])
                    })
                    break;

                case 'pod':
                    tmp.header = ['名称', '创建时间', '状态', '启动时间', '操作']
                    data.details[key].forEach((el) => {
                        tmp.navList.push({id: el.metadata.uid || '', name: el.metadata.name || ''})
                        tmp.body.push([el.metadata.name || '', moment(el.metadata.creationTimestamp || '').format('YYYY-MM-DD hh:mm:ss') || '',el.status.phase ||'',  moment(el.status.startTime || '').format('YYYY-MM-DD hh:mm:ss') || '', 'log'])
                    })
                    break;
                case 'replicaSet':
                    tmp.header = ['名称', '创建时间', '当前可用数']
                    data.details[key].forEach((el) => {
                        tmp.navList.push({id: el.metadata.uid || '', name: el.metadata.name || ''})
                        tmp.body.push([el.metadata.name || '',moment(el.metadata.creationTimestamp || '').format('YYYY-MM-DD hh:mm:ss') || '', el.status.numberAvailable || 0])
                    })
                    break;
                case 'replicationController':
                    tmp.header = ['名称', '创建时间', '当前副本数']
                    data.details[key].forEach((el) => {
                        tmp.navList.push({id: el.metadata.uid || '', name: el.metadata.name || ''})
                        tmp.body.push([el.metadata.name || '',moment(el.metadata.creationTimestamp || '').format('YYYY-MM-DD hh:mm:ss') || '', el.status.replicas || 0])
                    })
                    break;
                case 'statefulSet':
                    tmp.header = ['名称', '创建时间', '当前副本数']
                    data.details[key].forEach((el) => {
                        tmp.navList.push({id: el.metadata.uid || '', name: el.metadata.name || ''})
                        tmp.body.push([el.metadata.name || '',moment(el.metadata.creationTimestamp || '').format('YYYY-MM-DD hh:mm:ss') || '', el.status.replicas || 0])
                    })
                    break;
                case 'service':
                    tmp.header = ['名称', '创建时间']
                    data.details[key].forEach((el) => {
                        tmp.navList.push({id: el.metadata.uid || '', name: el.metadata.name || ''})
                        tmp.body.push([el.metadata.name || '', moment(el.metadata.creationTimestamp || '').format('YYYY-MM-DD hh:mm:ss') || ''])
                    })
                    break;
                case 'configMap':
                    tmp.header = ['名称', '创建时间']
                    data.details[key].forEach((el) => {
                        tmp.navList.push({id: el.metadata.uid || '', name: el.metadata.name || ''})
                        tmp.body.push([el.metadata.name || '', moment(el.metadata.creationTimestamp || '').format('YYYY-MM-DD hh:mm:ss') || ''])
                    })
                    break;
                case 'pvc':
                    tmp.header = ['名称', '创建时间', '状态']
                    data.details[key].forEach((el) => {
                        tmp.navList.push({id: el.metadata.uid || '', name: el.metadata.name || ''})
                        tmp.body.push([el.metadata.name || '', moment(el.metadata.creationTimestamp || '').format('YYYY-MM-DD hh:mm:ss') || '', el.status.phase || '' ])
                    })
                    break;
                case 'secret':
                    tmp.header = ['名称', '创建时间']
                    data.details[key].forEach((el) => {
                        tmp.navList.push({id: el.metadata.uid || '', name: el.metadata.name || ''})
                        tmp.body.push([el.metadata.name || '', moment(el.metadata.creationTimestamp || '').format('YYYY-MM-DD hh:mm:ss') || ''])
                    })
                    break;
                case 'roleBinding':
                    tmp.header = ['名称', '创建时间']
                    data.details[key].forEach((el) => {
                        tmp.navList.push({id: el.metadata.uid || '', name: el.metadata.name || ''})
                        tmp.body.push([el.metadata.name || '', moment(el.metadata.creationTimestamp || '').format('YYYY-MM-DD hh:mm:ss') || ''])
                    })
                    break;
                case 'role':
                    tmp.header = ['名称', '创建时间']
                    data.details[key].forEach((el) => {
                        tmp.navList.push({id: el.metadata.uid || '', name: el.metadata.name || ''})
                        tmp.body.push([el.metadata.name || '', moment(el.metadata.creationTimestamp || '').format('YYYY-MM-DD hh:mm:ss') || ''])
                    })
                    break;
                case 'serviceAccount':
                    tmp.header = ['名称', '创建时间']
                    data.details[key].forEach((el) => {
                        tmp.navList.push({id: el.metadata.uid || '', name: el.metadata.name || ''})
                        tmp.body.push([el.metadata.name || '', moment(el.metadata.creationTimestamp || '').format('YYYY-MM-DD hh:mm:ss') || ''])
                    })
                    break;
              
            }

            detailInfo.push(tmp)
        })

        // console.log('detailInfo===',detailInfo)

        setInstanceList(detailInfo)
        if(detailInfo.length) {
            const crumbs = [{id: detailInfo[0].id, name: detailInfo[0].name}]
            setCrumbs(crumbs)
            setCurNav(detailInfo[0].id)
            setListData(detailInfo[0])
        }
    }

    const getDetailData = () => {

        store.dispatch({
            type: TYPE.LOADING,
            val: true
        })

        axios({
            method: "GET",
            url: '/api/app/detail',
            params: {id: props.match.params.id || ''}
        }).then((res) => {
            if(res.data.code === 0) {
                setDetailData(res.data.result)
                parseDetailData(res.data.result)
            }
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: res.data.result || ''
            })
            
        }).catch((err) => {
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

    const changeCrumb = (curNav, id) => {
        // console.log('--changecrumb---', id, curNav)
        // console.log(crumbs)
        if(crumbs.length === 2 && id === crumbs[0].id) {
            let newCrumbs = crumbs.slice()
            newCrumbs.splice(1,1)
            setCrumbs(newCrumbs)

            let data = instanceList.find((item) => item.id == id)
            setType('list')
            setListData(data || [])

        }

        return 

        // console.log(e.currentTarget.dataset)
        // const curNav = e.currentTarget.dataset.name
        // const idx = parseInt(e.currentTarget.dataset.index)
        console.log('---change crumbs---',curNav, idx)
        let newCrumbs = crumbs.slice(0,idx+1)
        setCrumbs(newCrumbs)
        setCurNav(curNav)
        store.dispatch({
            type: TYPE.SET_SERVICE,
            val: curNav
        })
        if(newCrumbs.length === 3) {
            history.push('/detail/tail')
        }else {
            history.push('/detail')
        }
    }

   
    
    const changeNav = (data) => {
        let newCrumbs = []
        newCrumbs.push(data.parent)
        if(data.child) {
            newCrumbs.push(data.child)
            setCurNav(data.child.id)
            store.dispatch({
                type: TYPE.SET_SERVICE,
                val: data.child
            })
        }else {
            setCurNav(data.parent.id)
            store.dispatch({
                type: TYPE.SET_SERVICE,
                val: data.parent
            })
            setListData(instanceList[data.index])
        }

        setCrumbs(newCrumbs)
        setType('list')
    }

    const goTail = (data) => {
        setTailData(data)
        setType('tail')
        let newCrumbs = crumbs.slice()
        newCrumbs.push(data.name)
        setCrumbs(newCrumbs)
        setCurNav(data.name)
    }


    const goLog = (name) => {
        setType('log')
        let newCrumbs = crumbs.slice()
        newCrumbs.push({id: name, name: name})
        setCrumbs(newCrumbs)
        setLogName(name)
    }

    const showDesc = (name) => {
        setType('desc')
        let newCrumbs = crumbs.slice()
        newCrumbs.push({id: name, name: name})
        setCrumbs(newCrumbs)
        // setCurNav(name)

        setDescName(name)


        store.dispatch({
            type: TYPE.LOADING,
            val: true
        })

        axios({
            method: "GET",
            url: '/api/app/detailDesc',
            params: {instanceId: props.match.params.id || '', resourceType: curNav, resourceName: name}
        }).then((res) => {
            if(res.data.code === 0) {
                console.log(typeof res.data.result)
                try {
                    const resData = JSON.stringify(res.data.result)
                    const d = js_beautify(resData)
                    setDetailDesc(d)
                    
                }catch(err) {
                    store.dispatch({
                        type: TYPE.SNACKBAR,
                        val: res.data.result || ''
                    })
                }
               
            }
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: res.data.result || ''
            })
            
        }).catch((err) => {
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
        <section className="detail-container">
            <header className="detail-header">
                <div className="header-logo">Crab</div>
                {/* <div className="header-user">userinfo</div> */}
            </header>
            <Crumbs name={instanceInfo.name + ' ( ' + instanceInfo.id + ' ) ' } data={crumbs} change={changeCrumb} />
            <section className="detail-content">
                <section className="detail-left">
                    <DetailNav change={changeNav} data={instanceList} curNav={curNav}/>
                </section>
                <section className="detail-right">
                    {
                        type === 'list' ? (
                            <List data={listData} goTail={goTail} toDesc={showDesc} goLog={goLog} />
                        ) : null
                    } 
                     {
                        type === 'tail' ? (
                            <Tail data={tailData} />
                        ) : null
                    }
                     {
                        type === 'log' ? (
                            <Log id={detailData.id || ''} name={logName} />
                        ) : null
                    }
                     {
                        type === 'desc' ? (
                            <Desc name={descName} data={detailDesc || ''} />
                        ) : null
                    }

                {/* 
                    <Suspense fallback={SuspenseLoading()}>
                        <Switch>
                            <Route path="/detail/log" component={Log}/>
                            <Route path="/detail/tail" component={Tail}/>
                            <Route path="/detail/list" component={List} />
                        </Switch>
                    </Suspense> */}
                </section>
            </section>

           
        </section>
    )
}

function mapStateToProps(state) { return state }
export default connect(mapStateToProps)(Detail)