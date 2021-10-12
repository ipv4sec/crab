import React from 'react';
import ReactDOM from 'react-dom';
import {Router, Route, browserHistory} from 'react-router'
import MuiProvider from 'material-ui/styles/MuiThemeProvider';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import {Provider } from 'react-redux'
import store from './store'
import RouterComp from './router.js'

const muiTheme = getMuiTheme({
    palette: {
        primary: {
            main: '#4A90E2',
            light: '#4A90E2',
            dark: '#4A90E2'
        },
        secondary: {
            main: '#E93030',
            light: '#E93030',
            dark: '#E93030'
        },
        checkbox: {
            checkedColor: '#4A90E2',
        },
    },
    checkbox: {
        checkedColor: '#4A90E2',
    },
});

const App = () => ( 
    <Provider store = {store}>
        <MuiProvider muiTheme={muiTheme}>
            <RouterComp />
        </MuiProvider>
    </Provider>
)

ReactDOM.render(<App/>, document.getElementById('page'));