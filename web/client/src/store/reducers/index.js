import { combineReducers } from 'redux'
import homeReducer from './homeReducer'
import commonReducer from './commonReducer'
import userReducer from './userReducer'
import detailReducer from './detailReducer'

const reducers = combineReducers({
    common: commonReducer,
    home: homeReducer,
    user: userReducer,
    detail: detailReducer
})

export default reducers
