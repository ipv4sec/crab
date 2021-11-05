import * as TYPE from '../actions'

const initState = {
    loading: false,
    snackbar: ''//{open: false, message: ''}
}

function commonReducer (state = initState, action) {
    switch(action.type) {
        case TYPE.LOADING:
            return Object.assign({}, state, {loading: action.val})
        case TYPE.SNACKBAR:
            return Object.assign({}, state, {snackbar: action.val})
        default:
            return state;
    }
}

export default commonReducer