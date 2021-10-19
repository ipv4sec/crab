import {connect} from 'react-redux'

import AppPage from '../pages/appPage';


//将state.counter绑定到props的counter
function mapStateToProps(state) {
    return state;
}

export const AppPageComp = connect(mapStateToProps)(AppPage);