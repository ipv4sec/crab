import React, { useState, useEffect, useRef } from 'react'
import '../style/sass/components.scss'

export default class AutoTextarea extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            value: ''
        }

        this.txaRef = React.createRef()
        this.setData = this.setData.bind(this)
        this.getData = this.getData.bind(this)
        this.lineHeight = 0
        this.padding = 0
    }

    componentDidMount() {
        this.lineHeight = parseInt(getComputedStyle(this.txaRef.current).lineHeight)
        this.padding = parseInt(getComputedStyle(this.txaRef.current).paddingTop) * 2
        this.txaRef.current.style.height = this.padding + this.lineHeight+ 'px'
    }

    setData(value) {
        this.setState({
            value
        }, () => {
            if(this.state.value.trim() === '') {
                this.txaRef.current.style.height = this.padding + this.lineHeight+ 'px'
            }else {
                this.txaRef.current.style.height = this.txaRef.current.scrollHeight+ 'px'
            }
        })
    }

    getData() {
        return this.state.value
    }

    changeValue(e){
        this.setState({
            value: e.target.value
        })
    }

    // 回车换行
    keyDown(e){
        if(e.keyCode === 13) {
            this.txaRef.current.style.height = this.txaRef.current.offsetHeight + this.lineHeight + 'px'
        }
    }

    // 删除减行
    keyUp(e) {
        if(e.keyCode === 8) {
            let curHeight = e.target.value.split('\n').length * this.lineHeight
            if(curHeight + this.padding <  this.txaRef.current.offsetHeight) {
                this.txaRef.current.style.height = this.txaRef.current.offsetHeight - this.lineHeight + 'px'
            }
        }
    }

    autoTxaClick() {
        this.txaRef.current.focus()
    }

    paste() {
        console.log('---paste')
        this.txaRef.current.style.height = this.txaRef.current.scrollHeight+ 'px'
    }

    render() {
        return (
            <div className={`auto-textarea ${this.props.class}`} onClick={this.autoTxaClick.bind(this)}>
                <textarea 
                    ref={this.txaRef}
                    className="auto-input"
                    placeholder="请输入..."
                    value={this.state.value}
                    onChange={this.changeValue.bind(this)}
                    onKeyDown={this.keyDown.bind(this)}
                    onKeyUp={this.keyUp.bind(this)}
                    onPaste={this.paste.bind(this)}
                ></textarea>
            </div>
        )
    }
}
