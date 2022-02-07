import React, { useState, useEffect, useRef } from 'react'
import { connect } from 'react-redux'
import '../../style/sass/online.scss'
import '../../style/sass/components.scss'
import axios from 'axios'
import store from '../../store/store'
import * as TYPE from '../../store/actions'
import Loading from '../../components/Loading'
import SnackbarCmp from '../../components/Snackbar'

const ManagerView = (props) => {

    const previewRef = useRef(null)
    const [name, setName] = useState('')

    useEffect(() => {
        setName(props.match.params.name)
        getInfo(props.match.params.id)
    }, [])

    const getInfo = (id) => {
        store.dispatch({
            type: TYPE.LOADING,
            val: true
        })

        axios({
            method: 'GET',
            url: '/api/app/detail',
            params: {id}
        }).then(res => {
            if(res.data.code == 0) {
                previewRef.current.innerText = res.data.result && res.data.result.value ? res.data.result.value : ''
            }else {
                store.dispatch({
                    type: TYPE.SNACKBAR,
                    val: res.data.result
                })
            }
           
        }).catch(err => {
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
        <section className="online-container">
            <header className="online-header">
                <div className="header-logo">Crab</div>
                {/* <div className="header-user">userinfo</div> */}
            </header>
            <div className="online-preview-content">
                <div className="oltitle">查看 {name} 描述文件</div>
                <pre className='olPreview' ref={previewRef}></pre>
            </div>
           
            <Loading />
            <SnackbarCmp />
        </section>
    )
}

function mapStateToProps(state) {
    return state
}

export default connect(mapStateToProps)(ManagerView)