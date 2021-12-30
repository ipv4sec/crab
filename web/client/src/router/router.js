import React, { Suspense, lazy } from 'react'
import { BrowserRouter, Switch, Route } from 'react-router-dom'
import SuspenseLoading from '../components/SuspenseLoading'

const Home = lazy(() => import('../pages/home/Home'))
const Login = lazy(() => import('../pages/login/Login'))
const Portal = lazy(() => import('../pages/portal/Portal'))
const Detail = lazy(() => import('../pages/detail/Detail'))
const Trait = lazy(() => import('../pages/online/Trait'))
const WorkloadType = lazy(() => import('../pages/online/WorkloadType'))
const WorkloadVendor = lazy(() => import('../pages/online/WorkloadVendor'))
const CreateApp = lazy(() => import('../pages/online/CreateApp'))


const RouterDOM = () => (
    <BrowserRouter>
        <Suspense fallback={SuspenseLoading()}>
            <Switch>
                <Route path="/login" component={Login} />
                <Route path="/home" component={Home} />
                <Route path="/detail/:id/:name" component={Detail} />
                <Route path="/trait" component={Trait} />
                <Route path="/workloadtype" component={WorkloadType} />
                <Route path="/workloadvendor" component={WorkloadVendor} />
                <Route path="/createapp" component={CreateApp} />
                <Route path="/" component={Portal} />
            </Switch>
        </Suspense>
    </BrowserRouter>
)


export default RouterDOM