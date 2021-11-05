import React, { useState } from 'react'
import { withRouter } from 'react-router-dom'
import { connect } from 'react-redux'
import Button from '@material-ui/core/Button'
import Input from '../../components/Input'
import "../../style/sass/reset.scss"
import axios from 'axios'
import store from '../../store/store'
import * as TYPE from '../../store/actions'
import Loading from '../../components/Loading'
import SnackbarCmp from '../../components/Snackbar'

const Reset = (props) => {

    const [originPasswordErr, setOriginPasswordErr] = useState('')
    const [newPasswordErr, setNewPasswordErr] = useState('')

    let originPassword = ''
    let newPassword = ''

    function changeOriginPassword(value) {
        setOriginPasswordErr('')
        originPassword = value
    }

    function changeNewPassword(value) {
        setNewPasswordErr('')
        newPassword = value
    }

    function save() {
        if(userName.trim() === ''){
            setOriginPasswordErr('请输入原密码')
            return
        }
        if(password.trim() === '') {
            setNewPasswordErr('请输入新密码')
            return
        }

        console.log(originPassword, '-----' ,newPassword)

        store.dispatch({
            type: TYPE.LOADING,
            val: true
        })

        axios({
            method: "POST",
            url: "/api/user/reset",
            data: {
                oldPassword: originPassword,
                password: newPassword
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
            console.log('err===', err)
            store.dispatch({
                type: TYPE.LOADING,
                val: false
            })
        })

    }


    return (
        <div className="page-container reset-container">
            <div className="page-title">
                <p>密码设置</p>
            </div>
            <div className="input-item">
                <Input type="password" label="原密码：" placeholder="请输入原密码" change={changeOriginPassword} inputErr={originPasswordErr}/>
            </div>
           
            <div className="input-item">
                <Input type="password" label="新密码：" placeholder="请输入新密码" change={changeNewPassword} inputErr={newPasswordErr}/>
            </div>
           
            <div className="form-btn">
                <Button variant="contained" className="btn" color="primary" onClick={save}>保存</Button>
            </div>
        </div>
    )
}

function mapPropsToState(state) {
    return state
}

export default  connect(mapPropsToState)(withRouter(Reset))