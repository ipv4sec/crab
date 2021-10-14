import React from 'react';
import {Router, Route, IndexRoute, browserHistory} from 'react-router';
import Base from './pages/base';
import { SpacePageComp, AddSpaceComp, EditSpaceComp } from './containers/spaceContainer';
import { AppPageComp } from './containers/appContainer';
import { ClusterPageComp } from './containers/clusterContainer';
import {AppDetailComp} from './containers/appDetailContainer'
import axios from 'axios'

axios.defaults.headers = {
    'X-Requested-With': 'XMLHttpRequest',
    'Content-Type': 'application/json; charset=UTF-8',
    'Cache-Control': 'no-cache'
  }


const RouterComp = () => {
    return (
        <Router history={browserHistory}>
            <Route path={`/`} component={Base}>
                <IndexRoute component={SpacePageComp}/>
                <Route path={`/space`}>
                    <IndexRoute component={SpacePageComp}/>
                    <Route path="/space/add" component={AddSpaceComp} ></Route>
                    <Route path="/space/edit/:spaceId" component={EditSpaceComp} ></Route>
                </Route>
                
                <Route path={`/app`} component={AppPageComp}></Route>
                <Route path={`/cluster`} component={ClusterPageComp}></Route>
            </Route>
            <Route path="/appdetail/:appId/:isEdit" component={AppDetailComp}></Route>
        </Router>
    )
}

export default RouterComp