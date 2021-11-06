import { combineReducers } from 'redux'
import homeReducer from './homeReducer'
import commonReducer from './commonReducer'
import userReducer from './userReducer'

const reducers = combineReducers({
    common: commonReducer,
    home: homeReducer,
    user: userReducer
})

export default reducers
