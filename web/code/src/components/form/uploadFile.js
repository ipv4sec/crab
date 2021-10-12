import React from 'react'

export default class UploadFile extends React.Component{
    constructor(props) {
        super(props)
        this.state = {
            btnText: '选择文件上传',
            fileName: '',
            errorText: ''
        }
        this.setData = this.setData.bind(this)
    }

    changeFile() {
        // console.log(event)
        // console.log(event.target.files)
        this.setState({
            fileName: event.target.files[0].name,
            btnText: '更改',
            errorText: ''
        })
        this.props.change(event.target.files[0])
    }

    setData(data) {
        this.setState({
            fileName: data.value || '',
            btnText: this.state.fileName == '' ? '选择文件上传' : '更改',
            errorText: data.errorText || ''
        })
    }


    render() {
        return (
            <div className="upload-file-container">
                {/* <p className="input-label">{this.props.label}</p> */}
               
                <div className="upload-content">
                    {
                        this.state.fileName ? (
                            <p>{this.state.fileName}</p>
                        ) : null
                    }
                    <div className="upload-btn">
                        <input type="file" onChange={this.changeFile.bind(this)}/>
                        <button>{this.state.btnText}</button>
                    </div> 
                   
                </div>
                {
                    this.state.errorText !== '' ? (
                        <p className="errorText-style">{this.state.errorText}</p>
                    ) : null
                }
            </div>
        )
    }
}
