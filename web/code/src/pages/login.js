import React from 'react'
import {browserHistory} from 'react-router';
import { connect } from 'react-redux'
// import Input from '../components/form/input'
import axios from 'axios'
import store from '../store/index'
import * as TYPE from '../store/actions'
import LoadingComp from '../components/showLoading/index'
import Snack from '../components/snackbar/index'

import {TextField, FlatButton} from 'material-ui'

const styles = {
    errorStyle: {
      color: '#EC5858',
    },
    underlineFocusStyle: {
      borderColor: '#3986FF',
    },
  };

class Login extends React.Component {
    constructor(props){
        super(props)
        this.state = {
            userName: '',
            password: '',
            nameErr: '',
            passwordErr: ''
        }
        this.login = this.login.bind(this)
        this.changeName = this.changeName.bind(this)
        this.changePassword = this.changePassword.bind(this)
    }

    componentDidMount() {

    }


    changeName() {
        this.setState({
            userName: event.target.value,
            nameErr: ''
        })
    }
    nameblur() {

    }

    changePassword() {
        this.setState({
            password: event.target.value,
            passwordErr: ''
        })
    }
    passwordblur() {

    }

    login() {
        if(this.state.userName.trim() === ''){
            this.setState({
                nameErr: '请输入用户名'
            })
            return
        }
        if(this.state.password.trim() === '') {
            this.setState({
                passwordErr: '请输入密码'
            })
            return
        }


        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })

        axios({
            method: "POST",
            url: "/manager/login",
            headers: {"Content-Type": "application/json"},
            data:{
                userName: this.state.userName,
                password: this.state.password
            }
        }).then((res) => {
            console.log('res==',res)
            let data = res.data
            if(data.status === 'success') {
                window.sessionStorage.setItem('user', this.state.userName)
              window.location.replace('/')
            }else {
                store.dispatch({
                    type: TYPE.SHOW_SNACKBAR,
                    val: {
                        open: true,
                        message: data.msg
                    }
                })
            }
            
            store.dispatch({
                type: TYPE.SHOW_LOADING,
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

    render(){
        return (
            <div className="login-container">
                <div className="login-content">
                    <div className="input-item">
                        <p className="input-label">用户名：</p>
                        <TextField 
                            className="input-text"
                            underlineFocusStyle={styles.underlineFocusStyle}
                            value={this.state.userName} 
                            errorText={this.state.nameErr}
                            fullWidth={true}
                            onChange={this.changeName.bind(this)}
                            // onBlur={this.nameblur.bind(this)}
                        />
                        {/* <Input label="用户名：" change={this.changeName} inputErr={this.state.nameErr} value={this.state.userName} /> */}
                    </div>
                
                    <div className="input-item">
                        <p className="input-label">密码：</p>
                        <TextField 
                            className="input-text"
                            type="password"
                            underlineFocusStyle={styles.underlineFocusStyle}
                            value={this.state.password} 
                            errorText={this.state.passwordErr}
                            fullWidth={true}
                            onChange={this.changePassword.bind(this)}
                            // onBlur={this.passwordblur.bind(this)}
                        />
                        {/* <Input type="password" label="密码：" change={this.changePassword} inputErr={this.state.passwordErr} value={this.state.password} /> */}
                    </div>
                
                    <div className="form-btn">
                        <FlatButton label="登陆" onClick={this.login} className="confirm-btn" />
                        {/* <Button variant="contained" className="btn" color="primary" onClick={this.login}>登陆</Button> */}
                    </div>
        
                </div>
                <LoadingComp />
                <Snack />
            </div>
        )
    }
}

function mapPropsToState(state) {
    return state
}

export default  connect(mapPropsToState)(Login)