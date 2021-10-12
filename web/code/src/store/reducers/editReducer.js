import * as TYPE from '../actions';

const initState = {
    userInfo: {}, // 用户信息
    permission: {}, // 权限信息
    compId: '',// 当前的自定义组件id
    body: '',
    style: '',
    config: '',
    configData: {},
    snackbar: {},
    showPreview: false,
    showLoading: false,

    // 没用了
    snackMessage: '',
    snackDuration: 2000
};

function edit(state = initState, action) {
    switch (action.type) {
        case TYPE.SET_USERINFO:
            return Object.assign({}, state, {userInfo: action.val});
        case TYPE.SET_PERMISSION:
            return Object.assign({}, state, {permission: action.val});
        case TYPE.SET_COMPID:
            return Object.assign({}, state, {compId: action.val});
        case TYPE.SET_BODY:
            return Object.assign({}, state, {body: action.val});
        case TYPE.SET_STYLE:
            return Object.assign({}, state, {style: action.val});
        case TYPE.SET_CONFIG:
            return Object.assign({}, state, {config: action.val});
        case TYPE.SET_CONFIG_DATA:
            return Object.assign({}, state, {configData: action.val});
        case TYPE.SHOW_SNACKBAR:
            return Object.assign({}, state, {snackbar: action.val});
        case TYPE.SNACK_MESSAGE:
            return Object.assign({}, state, {snackMessage: action.val});
        case TYPE.SNACK_DURATION:
            return Object.assign({}, state, {snackDuration: action.val});
        case TYPE.SHOW_COMP_PREVIEW:
            return Object.assign({}, state, {showPreview: action.val});
        case TYPE.SHOW_LOADING:
            return Object.assign({}, state, {showLoading: action.val});
            
        default:
            return state;
    }
}

export default edit;