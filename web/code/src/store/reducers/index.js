import {combineReducers} from 'redux';
// import edit from './editReducer';
import common from './commonReducer'

const rootReducer = combineReducers({
    common,
});

export default rootReducer;