import React, { useState , lazy, Suspense} from 'react'
import '../../style/sass/home.scss'
import LeftNav from '../../components/LeftNav'
import { Route, Switch, useHistory } from 'react-router-dom'
import SuspenseLoading from '../../components/SuspenseLoading'
import Loading from '../../components/Loading'
import SnackbarCmp from '../../components/Snackbar'
import axios from 'axios'

// import Manager from '../manager/Manager'
// import WorkLoad from '../workload/WorkLoad'
// import Domain from '../domain/Domain'
// import Reset from '../reset/Reset'

const Manager = lazy(() => import('../manager/Manager'))
const WorkLoad = lazy(() => import('../workload/WorkLoad'))
const Domain = lazy(() => import('../domain/Domain'))
const Reset = lazy(() => import('../reset/Reset'))
const Online = lazy(() => import('../online/Online'))

const Home = (props) => {

    let history = useHistory()

    function changeNav(data) {
        history.push(data)
    }

    const logout = () => {
        // axios({
        //     method: 'GET',
        //     url: '/api/user/logout'
        // }).then(res => {
            
        // }).catch(err => {
        //     console.log('退出登陆失败')
        //     console.log(err)
        // })

        sessionStorage.removeItem('user')
        window.location.replace('/login')
    }

    return (
        <div className="home-container">
            {/* <div className="user-info">
                <p className="userinfo-name">{window.sessionStorage.getItem('user') || ''}</p>
                <button onClick={logout} className="userinfo-logout">退出登陆</button>
            </div> */}
            <div className="content-left">
                <LeftNav change={changeNav}/>
            </div>
            <div className="content-right">
                <Suspense fallback={SuspenseLoading()}>
                    <Switch> 
                        <Route path="/home/online" component={Online} />
                        <Route path="/home/workload" component={WorkLoad} />
                        <Route path="/home/domain" component={Domain} />
                        <Route path="/home/reset" component={Reset} />
                        <Route path="/home" component={Manager} />
                    </Switch>
                </Suspense>
            </div>

            <Loading />
            <SnackbarCmp />
        </div>  
    )
}

export default Home 

