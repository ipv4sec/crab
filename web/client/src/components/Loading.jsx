import React from 'react'
import store from '../store/store'
import { connect } from 'react-redux'
import CircularProgress from '@material-ui/core/CircularProgress'

const Loading = (props) => {
    let loading = store.getState().common.loading
    return (
        <div className={` circular-progress ${loading ? "show-progress" : ""}`}  >
            <CircularProgress size={60} color="primary" />
        </div>
        
    )
}
function mapPropsToState(state) {
    return state
}
export default connect(mapPropsToState)(Loading)