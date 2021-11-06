import React, { useState, useEffect } from 'react'
import Input from '../../components/Input'
import Button from '@material-ui/core/Button'
import store from '../../store/store'
import * as TYPE from '../../store/actions'
import axios from 'axios'
import { connect } from 'react-redux'
import '../../style/sass/workload.scss'

const WorkLoad = (props) => {
    const [host, setHost] = useState('')
    const [inputErr, setInputErr] = useState('')

    useEffect(() => {
        getClusterMirror()
    }, [])


    const getClusterMirror = () => {
        store.dispatch({
            type: TYPE.LOADING,
            val: true
        })

        axios({
            method: 'GET',
            url: '/api/cluster/mirror'
        }).then((res) => {
            if(res.data.code === 0) {
                setHost(res.data.result)
            }else {
                store.dispatch({
                    type: TYPE.SNACKBAR,
                    val: res.data.result
                })
            }
            store.dispatch({
                type: TYPE.LOADING,
                val: false
            })
        }).catch((err) => {
            console.log(err)
            store.dispatch({
                type: TYPE.LOADING,
                val: false
            })
        })
    }

    function changeHost(data) {
        setInputErr('')
        setHost(data)
    }

    function changeOrigin(){
        if(host.trim() == '') {
            setInputErr('请输入')
            return
        }

        store.dispatch({
            type: TYPE.LOADING,
            val: true
        })

        axios({
            method: 'POST',
            url: '/api/cluster/mirror',
            headers: {"Content-Type": "application/json"},
            data: {
                mirror: host
            }
        }).then((res) => {
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: res.data.result || ''
            })
            store.dispatch({
                type: TYPE.LOADING,
                val: false
            })
        }).catch((err) => {
            console.log(err)
            store.dispatch({
                type: TYPE.LOADING,
                val: false
            })
        })
    }

    return (
        <div className="page-container workload-container">
            <div className="page-title">
                <p>设置</p>
            </div>

            <div className="workload-content">
                <div className="host-input">
                    <Input change={changeHost} inputErr={inputErr} value={host} />
                </div> 
                <div className="host-btn">
                    <Button variant="contained" color="primary" className="btn-item" onClick={changeOrigin}>保存</Button>
                </div> 
               
            </div>
        </div> 
    )


}

function mapPropsToState(state) { return state }

export default connect(mapPropsToState)(WorkLoad)