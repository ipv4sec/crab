import React from 'react';

// import Header from '../components/header/index';
import Header from '../components/header/newHeader';
import LeftNav from '../components/leftNav/index';
import store from '../store';
import {connect} from 'react-redux'
// import * as events from '../libs/events';
// import Snack from '../components/snack';
// import UploadFile from '../components/upload/UploadFile';
// import Load from '../components/load/Load';
// import WinDel from '../components/win/WinDel';
// import TabSwitch from '../components/globalComponent/TabSwitch';

import LoadingComp from '../components/showLoading/index'
import Snack from '../components/snackbar/index'


class Base extends React.Component {
    constructor() {
        super();
    }

    componentWillMount() {
        if(!window.sessionStorage.getItem('user')) {
            window.location.replace('/login')
        }
    }

    componentDidMount() {
        this.init();

        // events.customEvent.on(events.REFRESH_ROOT_DOM, this.refreshRoot.bind(this));

    }

    componentWillUnmount() {
        // events.customEvent.removeAllListeners(events.REFRESH_ROOT_DOM);
    }

    init() {
        this.addInterceptor();
    }

    addInterceptor() {
        $(document).ajaxSuccess(function (event, xhr, settings) {
            var res = xhr.responseJSON;

            if (res) {
                if (res.code === 404) {
                    //退出登录
                    location.reload();
                } else if (res.code === 401) {
                    //未授权
                    // events.customEvent.emit(events.SHOW_SNACK, res.msg);
                } else {
                    //console.info(res.code);
                }
            }
        });
    }

    refreshRoot() {
        this.forceUpdate();
    }

    hideElement() {
        // events.customEvent.emit(events.HIDE_ELEMENT);
    }

    render() {
        return (
            <div className="wrap" onClick={this.hideElement.bind(this)}>
                <Header/>
                <LeftNav/>
                <div className="main">
                    {this.props.children}
                </div>

                <LoadingComp />
                <Snack />
            </div>
        )
    }
}

//将state.counter绑定到props的counter
function mapStateToProps(state) {
    return state;
}

export default connect(mapStateToProps)(Base);