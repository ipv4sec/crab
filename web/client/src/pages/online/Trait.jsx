import React, { useState, useEffect, useRef } from 'react'
import { connect } from 'react-redux'
import Button from '@material-ui/core/Button'
import '../../style/sass/online.scss'
import axios from 'axios'
import store from '../../store/store'
import * as TYPE from '../../store/actions'

const defaultMetadata = `apiVersion: aam.globalsphare.com/v1alpha1
kind: Trait
metadata:
  name: expose
spec:
  parameter: |
    k1: *"v1" | string`


const Trait = (props) => {
    const [trait, setTrait] = useState(defaultMetadata)

    const changeTrait = (e) => {
        setTrait(e.target.value)
    }

    const checkRule = () => {
        if(trait.trim() === '') {
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
            url: '/api/online/createtrait',
            data: {value: trait},
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
            <div className="page-title">创建Trait</div>
            <section className="trait-content">
                <div className="trait-textarea">
                    <textarea className="textarea-input" value={trait} onChange={changeTrait}></textarea>
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

export default connect(mapStateToProps)(Trait)