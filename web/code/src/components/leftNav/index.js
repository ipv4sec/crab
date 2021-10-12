import React, {Component} from 'react';
import {browserHistory} from 'react-router';
import store from '../../store/index';

import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import {List, ListItem, makeSelectable} from 'material-ui/List';
import FontIcon from 'material-ui/FontIcon';

export default class NavList extends React.Component {
    constructor() {
        super();
        this.state = {
            curNav: '',
            navList: [
                // {label: '空间管理', value: 'space', icon: 'icon_grey600beifen'},
                {label: '应用管理', value: 'app', icon: 'icon_grey600'},
                // {label: '集群管理', value: 'cluster', icon: 'icon_grey6002'},
            ],
        };
        this.changeNav = this.changeNav.bind(this)
    }
 
    componentDidMount() {
        if(sessionStorage.getItem('curNav')) {
            this.setState({
                curNav: sessionStorage.getItem('curNav')
            })
        }else {
            this.setState({
                curNav: 'app'
            })
        }
    }

    componentWillUnmount() {

    }

    changeNav(value){
        sessionStorage.setItem("curNav", value)
        this.setState({
            curNav: value
        })
        browserHistory.push("/"+value)
    }

    render() {
        return ( 
            <div className="left-nav">
                <div className="logo">
                    {/* <img src={`${window._STATICPATH}/images/logo.png`}/> */}
                    <p className="logo-title">islandConsole</p>
                </div>
               
                <div className="nav-list">
                    {/* <List> */}
                        {
                            this.state.navList.map((item, index) => {
                                return (
                                    <div className="list-item" >
                                        <div className={`item-content ${this.state.curNav == item.value ? "blueBorder" :""}`} onClick={() => {this.changeNav(item.value)}}>
                                            <i className={`iconfont ${item.icon}`}></i>
                                            <span >{item.label}</span> 
                                        </div>
                                    </div>
                                )
                            } )
                        }
                </div>
            </div>
        )
    }
}
