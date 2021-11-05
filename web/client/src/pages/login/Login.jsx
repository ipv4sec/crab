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

    const [nameErr, setNameErr] = useState(false)
    const [passwordErr, setPasswordErr] = useState(false)

    let userName = ''
    let password = ''


    function changeName(value) {
        setNameErr(false)
        userName = value
    }

    function changePassword(value) {
        setPasswordErr(false)
        password = value
    }

    function login() {
        if(userName.trim() === ''){
            setNameErr(true)
            return
        }
        if(password.trim() === '') {
            setPasswordErr(true)
            return
        }

        console.log(userName, '-----' ,password)

        store.dispatch({
            type: TYPE.LOADING,
            val: true
        })

        axios({
            method: "POST",
            url: "/api/user/login",
            data: {
                userName: userName,
                password: password
            }
        }).then((res) => {
            let data = res.data
            if(data.code === 0) {
                window.sessionStorage.setItem("token", res.data.result.token)
                this.props.history.redirect('/')
            }else {
                store.dispatch({
                    type: TYPE.SNACKBAR,
                    val: data.result || ''
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
            <div className="input-item">
                <Input label="用户名：" change={changeName} inputErr={nameErr}/>
            </div>
           
            <div className="input-item">
                <Input type="password" label="密码：" change={changePassword} inputErr={passwordErr}/>
            </div>
           
            <div className="form-btn">
                <Button variant="contained" className="btn" color="primary" onClick={login}>登陆</Button>
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