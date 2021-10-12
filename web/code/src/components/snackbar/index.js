import React from 'react';
import Snackbar from 'material-ui/Snackbar';
import store from '../../store'
import * as TYPE from '../../store/actions'

export default class Snack extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            open: false,
            message: '',
            duration: 2000
        }
    }

    snackbarClose() {
        store.dispatch({
            type: TYPE.SHOW_SNACKBAR,
            val: {
                open: false,
                message: '',
                duration: 2000
            }
        })
    } 

    render() {
        return (
            <Snackbar className="snackbar"
                style={{textAlign: 'center'}}
                open={store.getState().common.snackbar.open || false}
                message={store.getState().common.snackbar.message || ''}
                autoHideDuration={store.getState().common.snackbar.duration || 2000}
                autoHideDuration={2000}
                onRequestClose={this.snackbarClose.bind(this)}
            />
        );
    }
}