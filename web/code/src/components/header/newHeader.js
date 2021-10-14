import React from 'react';
import {browserHistory} from 'react-router';
import IconButton from 'material-ui/IconButton';
import AccountBox from 'material-ui/svg-icons/action/account-box';
import store from '../../store/index'
import axios from 'axios';
import * as TYPE from '../../store/actions'

export default class Header extends React.Component {
    constructor() {
        super();
        this.state = {
            user: '',
        };
    }

    componentDidMount() {
        // this.getUserInfo();
    }

    componentWillUnmount() {

    }

    

    getUserInfo() {
        axios({
            method: "GET",
            url: `/global/global/user`
        }).then((res) => {
            // console.log(res)
            if(res.data.code == 200) {
                store.dispatch({
                    type: TYPE.SET_USERINFO,
                    val: res.data.data
                })
            }

        }).catch((err) => {
            console.log(err)
        })
    }

    logout() {
        // location.href = `/auth/logout`;
        window.sessionStorage.removeItem('user')
        window.location.replace('/login')
    }
 
    render() {
        return (
            <div className="header-container">
                <div className="header-item">
                    <i className="iconfont icon_user-account-box"></i>
                    <span>{store.getState().common.userInfo.name || '用户'}</span>
                </div>
                <div className="header-item" onClick={this.logout.bind(this)}>
                    <i className="iconfont icon_exit-to-app-button"></i>
                    <span onClick={this.logout.bind(this)}>登出</span>
                </div>
            </div>
        )
    }
}