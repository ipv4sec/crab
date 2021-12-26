import React from 'react'
import RouterDOM from './router/router'
import './style/sass/common.scss'
import { createTheme, ThemeProvider  } from '@material-ui/core/styles'
import axios from 'axios'
import { BrowserRouter }  from 'react-router-dom'

// 添加请求拦截器
axios.interceptors.request.use(function (config) {
    // 在发送请求之前做些什么
    // config.headers['Auth'] = window.sessionStorage.getItem('token') || ''
    // if(!window.sessionStorage.getItem('user')) {
    //     window.location.replace('/login')
    // }

    return config;
}, function (error) {
    // 对请求错误做些什么
    return Promise.reject(error);
});

// 添加响应拦截器
axios.interceptors.response.use(function (response) {
    // 对响应数据做点什么
    console.log('----response---')
    const browserHistory = new BrowserRouter()
    console.log('history==',browserHistory)

    if(response.data.code === 40404) {
        window.sessionStorage.setItem('user', '')
        window.location.replace('/login')
        return
        // browserHistory.history.replace('/login')
    }else {
        return response;
    }

  }, function (error) {
    // 对响应错误做点什么
    return Promise.reject(error);
});

const mytheme = createTheme({
    palette: {
        primary: {
            main: '#3986FF'
        },
        secondary: {
            main: '#EC5858'
        }
    }
})

const APP = () => (
    <ThemeProvider theme={mytheme} >
        <div className="root-container">
            <RouterDOM />
        </div>
    </ThemeProvider>
 
)

export default APP