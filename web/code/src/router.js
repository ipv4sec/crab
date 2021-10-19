import React from 'react';
import {Router, Route, IndexRoute, browserHistory} from 'react-router';
import Base from './pages/base';
import axios from 'axios'

axios.defaults.headers = {
    'X-Requested-With': 'XMLHttpRequest',
    'Content-Type': 'application/json; charset=UTF-8',
    'Cache-Control': 'no-cache'
  }


const AppPageComp = (location, callback) => {
    require.ensure([], require => {
        callback(null, require('./containers/appContainer').AppPageComp)
    }, 'appContainer')
};

import Login from './pages/login'



const RouterComp = () => {
    return (
        <Router history={browserHistory}>
            <Route path="/login" component={Login} />
            <Route path={`/`} component={Base}>
                <IndexRoute getComponent={AppPageComp}/>
                <Route path={`/app`} getComponent={AppPageComp}></Route>
            </Route>
        </Router>
    )
}

export default RouterComp