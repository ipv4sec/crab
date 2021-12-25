import * as TYPE from '../actions'


let initState = {
    instanceInfo: {},
    service: '', // 当前是哪个服务
    serviceDetail: '', // 服务详情名字
}

function detailReducer(state=initState, action) {

    switch(action.type) {
        case TYPE.SET_INSTANCE_INFO:
            return Object.assign({}, state, {instanceof: action.val});
        case TYPE.SET_SERVICE:
            return Object.assign({}, state, {service: action.val});
        case TYPE.SET_SERVICE_DETAIL:
            return Object.assign({}, state, {serviceDetail: action.val});
        default:
            return state
    }
}

export default detailReducer