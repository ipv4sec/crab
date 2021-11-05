import * as TYPE from '../actions'


let initState = {
    title: 'xxxx',
    iframeSrc: ''
}

function homeReducer(state=initState, action) {

    switch(action.type) {
        case TYPE.TITLE:
            return Object.assign({}, state, {title: action.val});
        case TYPE.IFRAME_SRC:
            return Object.assign({}, state, {iframeSrc: action.val});
        default:
            return state
    }
}

export default homeReducer