import * as TYPE from '../../store/actions'

let initState = {
    userName: ''
}

const userReducer = (state = initState, action) => {
    switch(action.type) {
        case TYPE.USER_NAME:
            return Object.assign({}, state, {userName: action.val})
        default:
            return state
    }
}

export default userReducer