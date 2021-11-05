import { combineReducers } from 'redux'
import homeReducer from './homeReducer'
import commonReducer from './commonReducer'

const reducers = combineReducers({
    common: commonReducer,
    home: homeReducer
})

export default reducers
