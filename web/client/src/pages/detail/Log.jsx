import React, { useState, useEffect } from 'react'
import store from '../../store/store'
import * as TYPE from '../../store/actions'
import axios from 'axios'
import '../../style/sass/detail.scss'


const Log = (props) => {

    // 查看日志
    const [logList, setLogList] = useState('')

    useEffect(() => {
        readLogs()
    }, [props.id])


    const readLogs = () => {

        store.dispatch({
            type: TYPE.LOADING,
            val: true
        })
        axios({
            url: '/api/app/logs',
            method: 'GET',
            params: {id: props.id, podName: props.name}
        }).then((res) => {
            if(res.data.code === 0) {
                const logs = res.data.result
                if(logs.length) {
                    for(let item of logs) {
                        console.log('--item log--', item)
                        if(item.name === props.ctnName) {
                            setLogList(item.value)
                            break;
                        }
                    }
                }
               
            }else {
                store.dispatch({
                    type: TYPE.SNACKBAR,
                    val: res.data.result || ''
                })
            }
        }).catch((err) => {
            console.error(err)
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
        <section className="detail-log">
            <pre className="item-desc">{logList || '暂无日志'}</pre>
        </section>
    )
}

export default Log