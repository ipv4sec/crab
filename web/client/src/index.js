import React from 'react'
import { render } from 'react-dom'
import { Provider } from 'react-redux'
import store from './store/store'
import APP from './APP'

const renderApp = () => render( 
    <Provider store={store}>
        <APP />
    </Provider>, 
    document.getElementById('root'))

// 热更新
if(process.env.NODE_ENV !== 'production' && module.hot) {
    module.hot.accept('./APP', renderApp)
}

renderApp()
