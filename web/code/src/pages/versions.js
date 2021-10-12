import React from 'react'
import BaseDialog from '../components/dialog/baseDialog'
import store from '../store/index'
import * as TYPE from '../store/actions'
import axios from 'axios'
import JsBeautify from 'js-beautify'
export default class Version extends React.Component {
    constructor(props){
        super(props)
        this.state = {
            list: [
                // 'v0.0.1', 'v0.0.2'
                // {version: 'v0.0.1', id: 'dsfdsf'},
                // {version: 'v0.0.1', id: 'dsfdsf'},
                // {version: 'v0.0.1', id: 'dsfdsf'},
                // {version: 'v0.0.1', id: 'dsfdsf'},
                // {version: 'v0.0.1', id: 'dsfdsf'},
                // {version: 'v0.0.1', id: 'dsfdsf'},
                // {version: 'v0.0.1', id: 'dsfdsf'},
                // {version: 'v0.0.1', id: 'dsfdsf'},
                // {version: 'v0.0.1', id: 'dsfdsf'},
                // {version: 'v0.0.1', id: 'dsfdsf'},
                // {version: 'v0.0.1', id: 'dsfdsf'},
                // {version: 'v0.0.1', id: 'dsfdsf'},
                // {version: 'v0.0.1', id: 'dsfdsf'},
                // {version: 'v0.0.1', id: 'dsfdsf'},
            ],
            open:false,
            code: ''
        }
    }

    componentDidMount() {
        this.getVersion()
    }

    getVersion() {
        let that = this

        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })

        axios({
            method: 'GET',
            url: `${window._BASEPATH}/api/getAppVersions`,
            params: {appId: this.props.data.id}
        }).then((res) => {
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
            that.setState({
                list: res.data.list || []
            })
        }).catch((err) => {
            console.log(err)
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
        })
    }
    click(data) {
        let that = this
        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })

        axios({
            method: 'GET',
            url: `${window._BASEPATH}/api/getVersionDetail`,
            params: {appId: this.props.data.id, version: data}
        }).then((res) => {
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
            that.setState({
                open: true,
                code: res.data ? JsBeautify.js(JSON.stringify(res.data)) : ''
            })
        }).catch((err) => {
            console.log(err)
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
        })

    }

    close(){
        this.setState({
            open: false
        })
    }

    render() {
        return (
            <div className="version-container">
                <ul>
                    {this.state.list.map((item, index) => {
                        return (
                            <li key={index}>
                                <span>{item}</span> 
                                <button onClick={this.click.bind(this, item)}>描述文件</button>
                            </li>
                        )
                    })}
                </ul>

                <BaseDialog
                    open={this.state.open}
                    title="版本描述文件"
                    showCloseBtn={true}
                    close={this.close.bind(this)}
                    contentClass="middle-dialog"
                >
                   <div className="show-code">
                       <pre>
                            {this.state.code}
                       </pre>
                   </div>
                </BaseDialog>

            </div>
        )
    }

}