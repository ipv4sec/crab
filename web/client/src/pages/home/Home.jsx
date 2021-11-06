import React, { useState , lazy, Suspense} from 'react'
import '../../style/sass/home.scss'
import LeftNav from '../../components/LeftNav'
import { Route, Switch, useHistory } from 'react-router-dom'

import Loading from '../../components/Loading'
import SnackbarCmp from '../../components/Snackbar'

import Manager from '../manager/Manager'
import WorkLoad from '../workload/WorkLoad'
import Domain from '../domain/Domain'
import Reset from '../reset/Reset'

const Home = (props) => {

    let history = useHistory()

    function changeNav(data) {
        history.push(data)
    }

    return (
        <div className="home-container">
            <div className="content-left">
                <LeftNav change={changeNav}/>
            </div>
            <div className="content-right">
                <Switch> 
                    <Route path="/home/workload" component={WorkLoad} />
                    <Route path="/home/domain" component={Domain} />
                    <Route path="/home/reset" component={Reset} />
                    <Route path="/home" component={Manager} />
                </Switch>
            </div>

            <Loading />
            <SnackbarCmp />
        </div>  
    )
}

export default Home 

