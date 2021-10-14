import * as TYPE from '../actions';

const initState = {
    userInfo: {}, // 用户信息
    permission: {}, // 权限信息
    snackbar: {},
    showLoading: false,
};

function common(state = initState, action) {
    switch (action.type) {
        case TYPE.SET_USERINFO:
            return Object.assign({}, state, {userInfo: action.val});
        case TYPE.SET_PERMISSION:
            return Object.assign({}, state, {permission: action.val});
        case TYPE.SHOW_SNACKBAR:
            return Object.assign({}, state, {snackbar: action.val});
        case TYPE.SHOW_LOADING:
            return Object.assign({}, state, {showLoading: action.val});
        default:
            return state;
    }
}

export default common;