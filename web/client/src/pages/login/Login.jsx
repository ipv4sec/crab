import React, { useState } from 'react'
import { withRouter } from 'react-router-dom'
import { connect } from 'react-redux'
import Button from '@material-ui/core/Button'
import Input from '../../components/Input'
import "../../style/sass/login.scss"
import axios from 'axios'
import store from '../../store/store'
import * as TYPE from '../../store/actions'
import Loading from '../../components/Loading'
import SnackbarCmp from '../../components/Snackbar'

const Login = (props) => {

    const [name, setName] = useState('')
    const [nameErr, setNameErr] = useState('')
    const [password, setPassword] = useState('')
    const [passwordErr, setPasswordErr] = useState('')

    function changeName(value) {
        setNameErr('')
        setName(value)
    }

    function changePassword(value) {
        setPasswordErr('')
        setPassword(value)
    }

    function login() {
        if(name.trim() === ''){
            setNameErr('请输入')
            return
        }
        if(password.trim() === '') {
            setPasswordErr('请输入')
            return
        }

        store.dispatch({
            type: TYPE.LOADING,
            val: true
        })

        axios({
            method: "GET",
            url: "/api/user/login",
            params: {
                username: name,
                password: password
            }
        }).then((res) => {
            if(res.data.code === 0) {
                window.sessionStorage.setItem('user', res.data.result.username || '')
                window.sessionStorage.setItem('curNav', '/home')
                window.location.replace('/home')
            }else {
                store.dispatch({
                    type: TYPE.SNACKBAR,
                    val: res.data.result || ''
                })
            }
            
            store.dispatch({
                type: TYPE.LOADING,
                val: false
            })
        }).catch((err) => {
            console.log('err===', err)
            store.dispatch({
                type: TYPE.LOADING,
                val: false
            })
        })

    }



    return (
        <div className="login-container">
            <div className="login-content">
                <div className="input-item">
                    <Input label="用户名：" value={name} change={changeName} inputErr={nameErr}/>
                </div>
            
                <div className="input-item">
                    <Input type="password" label="密码：" value={password} change={changePassword} inputErr={passwordErr}/>
                </div>
            
                <div className="form-btn">
                    <Button variant="contained" className="btn" color="primary" onClick={login}>登陆</Button>
                </div>
            </div>

            <Loading />
            <SnackbarCmp />
        </div>
       
    )
}

function mapPropsToState(state) {
    return state
}

export default  connect(mapPropsToState)(withRouter(Login))