import React from 'react'
import InputComp from '../components/form/input'
import UploadFile from '../components/form/uploadFile'
import axios from 'axios'

export default class AddCluster extends React.Component {
    constructor(props) {
        super(props)
        this.name = '';
        this.host = '';
        this.other = '';
        this.file = null
        this.setData = this.setData.bind(this)
        this.getData = this.getData.bind(this)
    }

    componentDidMount() {
        if(this.props.type == 'edit'){
            this.setData(this.props.data)
        }
        
    }

    getData() {
        if(this.verifyName() && this.verifyHost() && this.verifyFile() ) {
            return {
                name: this.name,
                host: this.host,
                other: this.other,
                file: this.file
            }
        }else {
            return false
        }

        
    }

    setData(data) {
        if(data && Object.keys(data).length) {
            this.name = data.name || '';
            this.host = data.url || '';
            this.other = data.description || '';
            this.file = data.file || '';
            this.nameRef.setData({value: this.name});
            this.hostRef.setData({value: this.host});
            this.otherRef.setData({value: this.other});
            this.fileRef.setData({value: this.file.toString()});
        }
     
    }

    nameBlur(value) {
        this.name = value
        // 如果需要校验，可以在这里写
        
    }
    hostBlur(value) {
        this.host = value
        // 如果需要校验，可以在这里写
    }
    otherBlur(value) {
        this.other = value
        // 如果需要校验，可以在这里写
    }
    changeFile(value) {
        this.file = value

        // 如果需要校验，可以在这里写
        // let formData = new FormData
        // formData.append('file', value)
        // axios({
        //     method: 'POST',
        //     url:  `${window._BASEPATH}/api/uploadFile`,
        //     headers:  {'Content-Type': 'application/x-www-form-urlencoded'},
        //     data: formData
        // }).then((res) => {
        //     console.log(res)
        //     if(res.data.status == 'success'){
        //         this.file = res.data.data
        //     }
           
        // }).catch((err) => {
        //     console.log(err)
        // })

        
        

    }

    verifyName() {
        if(this.name == ''){
            this.nameRef.setData({
                errorText: '请输入名称'
            })
            return false
        }
        if(this.name.length > 30 ) {
            this.nameRef.setData({
                errorText: '名字长度不能超过30个字符'
            })
            return false
        }
        return true
    }
    verifyHost() {
        if(this.host == ''){
            this.hostRef.setData({
                errorText: '请输入k8s-host'
            })
            return false
        }
        return true
    }
    verifyFile() {
        if(!this.file ){
            this.fileRef.setData({
                errorText: '请上传证书'
            })
            return false
        }
       
        return true
    }


    render() {
        return (
           <div className="add-cluster-content">
               <InputComp ref={(ele) => {this.nameRef = ele}} label="名称" blur={this.nameBlur.bind(this)}/>
               <InputComp ref={(ele) => {this.hostRef = ele}} label="k8s-host" blur={this.hostBlur.bind(this)}/>
               <UploadFile ref={(ele) => {this.fileRef = ele}} label="证书" change={this.changeFile.bind(this)} type={this.props.type}/>
               <InputComp ref={(ele) => {this.otherRef = ele}} label="其他配置" blur={this.otherBlur.bind(this)}/>
           </div>
        )
    }
}