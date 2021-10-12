import React from 'react';
import {Router, Route, IndexRoute, browserHistory} from 'react-router';
import Base from './pages/base';
// import { SpacePageComp, AddSpaceComp, EditSpaceComp } from './containers/spaceContainer';
// import { AppPageComp } from './containers/appContainer';
// import { ClusterPageComp } from './containers/clusterContainer';
// import {AppDetailComp} from './containers/appDetailContainer'
import axios from 'axios'

axios.defaults.headers = {
    'X-Requested-With': 'XMLHttpRequest',
    'Content-Type': 'application/json; charset=UTF-8',
    'Cache-Control': 'no-cache'
  }


// const SpacePageComp = (location, callback) => {
//     require.ensure([], require => {
//         callback(null, require('./containers/spaceContainer').SpacePageComp)
//     }, 'spaceContainer')
// };
// const AddSpaceComp = (location, callback) => {
//     require.ensure([], require => {
//         callback(null, require('./containers/spaceContainer').AddSpaceComp)
//     }, 'spaceContainer')
// };
// const EditSpaceComp = (location, callback) => {
//     require.ensure([], require => {
//         callback(null, require('./containers/spaceContainer').EditSpaceComp)
//     }, 'spaceContainer')
// };
const AppPageComp = (location, callback) => {
    require.ensure([], require => {
        callback(null, require('./containers/appContainer').AppPageComp)
    }, 'appContainer')
};
// const ClusterPageComp = (location, callback) => {
//     require.ensure([], require => {
//         callback(null, require('./containers/clusterContainer').ClusterPageComp)
//     }, 'clusterContainer')
// };

// const AppDetailComp = (location, callback) => {
//     require.ensure([], require => {
//         callback(null, require('./containers/appDetailContainer').AppDetailComp)
//     }, 'appDetailContainer')
// };



// const Login = (location, callback) => {
//     require.ensure([], require => {
//         callback(null, require('./pages/login'))
//     }, 'login')
// };


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