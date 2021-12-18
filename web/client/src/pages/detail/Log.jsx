import React, { useState, useEffect } from 'react'
import store from '../../store/store'
import * as TYPE from '../../store/actions'
import axios from 'axios'
import '../../style/sass/detail.scss'


const Log = (props) => {

    // 查看日志
    const [logList, setLogList] = useState([])

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
            params: {id: props.id}
        }).then((res) => {
            if(res.data.code === 0) {
                setLogList(res.data.result || [])
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
             <div className="log-list">
                {
                    logList.map((item, index) => {
                        return (
                            <div key={index} className="log-item">
                                <p>{item.name}：</p>
                                <p className="item-desc">{item.value}</p>
                            </div>
                        )
                    })
                }
            </div>
        </section>
    )
}

export default Log