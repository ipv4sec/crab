const axios = require('axios').default
const config = require('../config/config.json')

axios.defaults.baseURL = config.domain

function get(url, params, header, callback){
    axios({
        method: 'GET',
        url: url,
        headers: header,
        params: params
    }).then((data) => {
        callback(data)
    }).catch((err) => {
        console.log('--- request error ---')
        console.log(err)
        callback({ code: 400, result: '请求错误'})
    })
}

function post(url, data, header, callback) {
    let newHeader = Object.assign({}, header, {"Content-Type": "application/json"})
    axios({
        method: 'POST',
        url: url,
        headers: newHeader,
        data: data
    }).then((res) => {
        callback(res)
    }).catch((error) => {
        console.log('--- request error ---')
        console.log(error)
        callback({ code: 400, result: '请求错误'})
    })
}


function postForm(url, data, header, callback) {
    let newHeader = Object.assign({}, header, {"Content-Type": "multipart/form-data"})
    axios({
        method: 'POST',
        url: url,
        headers: newHeader,
        data: data
    }).then((res) => {
        callback(res)
    }).catch((error) => {
        console.log('--- request error ---')
        console.log(error.message)
        callback({ code: 400, result: '请求错误'})
    })
}


function put(url, data,header, callback) {
    // console.log(url, data)
    let newHeader = Object.assign({}, header, {"Content-Type": "application/json"})
    axios({
        method: 'PUT',
        url: url,
        headers: newHeader,
        data: data
    }).then((res) => {
        callback(res)
    }).catch((error) => {
        console.log('--- request error ---')
        console.log(error)
        callback({ code: 400, result: '请求错误'})
    })
}

function del(url, data, header, callback) { 
    let newHeader = Object.assign({}, header, {"Content-Type": "application/json"})
    axios({
        method: 'DELETE',
        url: url,
        headers: newHeader,
        data: data
    }).then((res) => {
        callback(res)
    }).catch((error) => {
        console.log('--- request error ---')
        console.log(error)
        callback({ code: 400, result: '请求错误'})
    })
}


module.exports = {
    get,
    post,
    put, 
    del,
    postForm
}

