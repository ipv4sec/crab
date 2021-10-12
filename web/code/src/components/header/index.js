import React from 'react';
import {browserHistory} from 'react-router';
import IconButton from 'material-ui/IconButton';
import AccountBox from 'material-ui/svg-icons/action/account-box';
import store from '../../store/index'

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
        $.ajax({
            type: 'get',
            url: `/global/user`,
            dataType: 'json',
            success: (msg) => {
                if (msg.code === 200) {
                    this.setState({
                        user: msg.data.name,
                    });
                    store.dispatch({
                        type: TYPE.SAVE_USER,
                        val: msg.data
                    })
                } else {
                    events.customEvent.emit(events.SHOW_SNACK, msg.msg);
                }
            },
            error: (msg) => {
                console.error(msg);
            }
        });
    }

    logout() {
        location.href = `/auth/logout`;
    }
 
    render() {
        return (
            <div className="topbar-wrap clearfix">
                <div className="logo"
                     onClick={() => {
                         events.customEvent.emit(events.CHANGE_TAB)
                         browserHistory.push(`${window._BASEPATH}/`);
                     }}
                >
                    <img src={`${window._STATICPATH}/images/logo.png`}/>
                </div>
 
                <div className="option">
                    <div className="btn">
                        <IconButton className={'icon'} 
                            tooltip="账户信息"
                        onClick={() => {
                            this.setState({
                                userbox: !this.state.userbox,
                            })
                        }}>
                            <AccountBox color="#878787"/>
                        </IconButton>
                        {
                            this.state.userbox ?
                                <span>
                                    <div className="mask" onClick={() => {
                                        this.setState({
                                            userbox: false,
                                        })
                                    }}></div>
                                    <div className="userbox">
                                        <div className="name">你好,{this.state.user}</div>
                                        <div className="logout" onClick={this.logout.bind(this)}>
                                            <i className="iconfont"></i>
                                            <span>退出</span>
                                        </div>
                                    </div>
                                </span> : null
                        }
                    </div>
                </div>
            </div>
        )
    }
}