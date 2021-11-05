import React from 'react'

import Snackbar from '@material-ui/core/Snackbar'
import Slide from '@material-ui/core/Slide'
import store from '../store/store'
import * as TYPE from '../store/actions'
import { connect } from 'react-redux'

const TransitionUp = (props) => {
    return <Slide {...props} direction="up" />;
}

const SnackbarCmp = (props) => {
    function close() {
        store.dispatch({
            type: TYPE.SNACKBAR,
            val: ''
        })
    }
    
    return (
        <div className="snackbar-container">
            <Snackbar 
                open={Boolean(store.getState().common.snackbar)}
                anchorOrigin={{vertical: 'bottom', horizontal: 'center'}}
                autoHideDuration={500000}
                onClose={close}
                TransitionComponent={TransitionUp}
                message={store.getState().common.snackbar}
            />
        </div>
    )
}

function mapPropsToState(state) { return state }
export default connect(mapPropsToState)(SnackbarCmp)