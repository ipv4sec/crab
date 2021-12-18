import React, { useState, useEffect, useRef } from 'react'
import { connect } from 'react-redux'
import Button from '@material-ui/core/Button'
import '../../style/sass/online.scss'
import axios from 'axios'
import store from '../../store/store'
import * as TYPE from '../../store/actions'

const defaultMetadata = `apiVersion: aam.globalsphare.com/v1alpha1
kind: Application
metadata:
  name: example
  version: 0.0.1
  description: 样例应用
  keywords:
    - 样例应用
  author: example@example.com
  maintainers:
    - email: example@example.com
      name: example
      web: https://example.com
  repositories: ["https://github.com/example/example.git"]
  bugs: https://github.com/example/example/issues
  licenses:
    - type: LGPL
      url: https://license.spec.com`


const WorkloadType = (props) => {
    const [value, setValue] = useState(defaultMetadata)

    const changeValue = (e) => {
        setValue(e.target.value)
    }

    const checkRule = () => {
        if(value.trim() === '') {
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: '请输入trait内容'
            })
            return false
        }

        return true
    }

    const save = () => {
        if(!(checkRule())) { return }

        store.dispatch({
            type: TYPE.LOADING,
            val: true
        })
        axios({
            method: 'POST',
            url: '/api/online/createworkloadtype',
            data: {value},
            headers: { 'Content-Type': 'application/json'}
        }).then(res => {
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: res.data.result
            })
        }).catch(err => {
            console.log(err)
            store.dispatch({
                type: TYPE.SNACKBAR,
                val: '请求错误'
            })
        }).finally(() => {
            store.dispatch({
                type: TYPE.LOADING,
                val: false
            })
        })
    }


    return (
        <section className="page-container online-container">
            <div className="page-title">创建WorkloadType</div>
            <section className="trait-content">
                <div className="trait-textarea">
                    <textarea className="textarea-input" value={value} onChange={changeValue}></textarea>
                </div>
                <div className="online-btns">
                    <Button className="online-btn" variant="contained" color="primary" onClick={save}>保存</Button>
                </div>
            </section>
           
            
        </section>
    )
}

function mapStateToProps(state) {
    return state
}

export default connect(mapStateToProps)(WorkloadType)