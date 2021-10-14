import React from 'react'
import UploadFile from '../components/form/uploadFile'
import AppTree from '../components/appdetail/appTree'
import axios from 'axios'
import store from '../store'
import * as TYPE from '../store/actions'

export default class AddApp extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            appId: '',//this.props.params.appId,
            isEdit: true //this.props.params.isEdit == '1' ? true : false,
        }
      
        this.setData = this.setData.bind(this)
        this.getData = this.getData.bind(this)
    }

    componentDidMount() {
        // this.getAppInfo()
    }

    getAppInfo() {
        let appInfo = this.props.data
        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })
        axios({
            method: 'GET',
            url:  `${window._BASEPATH}/manager/uploadFile`,
            headers:  {'Content-Type': 'application/x-www-form-urlencoded'},
            data: formData
        }).then((res) => {
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
            if(res.data.status == 'success'){
                this.appTreeRef.setData(res.data)
            }
        }).catch((err) => {
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
            console.log(err)
        })

    }

    getData() {
        return this.appTreeRef.getData()
       
    }

    setData(data) {
        if(data && Object.keys(data).length) {
            this.name = data.name || '';
            this.git_url = data.git_url || '';
            this.description = data.description || '';
            this.nameRef.setData({value: this.name});
            this.gitRef.setData({value: this.git_url});
            this.descRef.setData({value: this.description});
        }
    }
    

    changeFile(value) {
        this.file = value
        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })
        // 如果需要校验，可以在这里写
        let formData = new FormData
        formData.append('file', value)
        axios({
            method: 'POST',
            url:  `${window._BASEPATH}/manager/addFile`,
            headers:  {'Content-Type': 'application/x-www-form-urlencoded'},
            data: formData
        }).then((res) => {
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
            if(res.status == 200){
                this.appTreeRef.setData(res.data)
            }
        }).catch((err) => {
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
            console.log(err)
        })
    }

     // 选择app版本后触发
     checkedAppList(data) {
        // 选择某个版本后会返回所有app树
        this.setState({
            allApps: data
        },() => {
            this.getedApps = true
            // 更新svg部分
            if(this.svgRef && this.svgLoaded){
                this.svgRef.setData(this.state.allApps)
            }
        })
    }


    render() {
        return (
           <div className="add-app-content">
                <UploadFile ref={(ele) => {this.fileRef = ele}} change={this.changeFile.bind(this)} type={this.props.type}/>
          
                <div className="add-app-tree">
                    <AppTree ref={(ele) => {this.appTreeRef = ele}} />
                </div>
           
           </div>
        )
    }
}