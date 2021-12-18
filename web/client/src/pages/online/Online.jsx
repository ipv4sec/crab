import React, { lazy, Suspense } from 'react'
import {Route, Switch } from 'react-router-dom'
import SuspenseLoading from '../../components/SuspenseLoading'

const CreateApp = lazy(() => import('./CreateApp'))
const Trait = lazy(() => import('./Trait'))
const WorkloadType = lazy(() => import('./WorkloadType'))
const WorkloadVendor = lazy(() => import('./WorkloadVendor'))

const Online = (props) => {


    return (
        <section>
            <Suspense fallback={SuspenseLoading()}>
                <Switch>
                    <Route path="/home/online/create" component={CreateApp}></Route>
                    <Route path="/home/online/trait" component={Trait}></Route>
                    <Route path="/home/online/workloadtype" component={WorkloadType}></Route>
                    <Route path="/home/online/workloadvendor" component={WorkloadVendor}></Route>
                </Switch>
            </Suspense>
        </section>
    )
}

export default Online