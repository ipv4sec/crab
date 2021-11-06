import React, { Suspense, lazy } from 'react'
import { BrowserRouter, Switch, Route } from 'react-router-dom'
import CircularProgress from '@material-ui/core/CircularProgress'

const Home = lazy(() => import('../pages/home/Home'))
const Login = lazy(() => import('../pages/login/Login'))
const Portal = lazy(() => import('../pages/portal/Portal'))

const loading = () => (
    <div className="error-loading">
        <CircularProgress size={80}/>
    </div>
)

const RouterDOM = () => (
    <BrowserRouter>
        <Suspense fallback={loading()}>
            <Switch>
                <Route path="/login" component={Login} />
                <Route path="/home" component={Home} />
                <Route path="/" component={Portal} />
            </Switch>
        </Suspense>
    </BrowserRouter>
)


export default RouterDOM