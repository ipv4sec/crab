import React from 'react'
import AceEditor from 'react-ace'

import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-xcode'

import '../style/sass/components.scss'

export default class Editor extends React.Component{
    constructor(props) {
        super(props)
        this.state = {
            value: '',
            height: 19,
            cursorStyle: 'ace-cursor'
        }

        this.change = this.change.bind(this)
        this.focus = this.focus.bind(this)
        this.blur = this.blur.bind(this)
        this.lineHeight= 0
        this.propsHeight = 19
    }

    focus() {
        this.setState({
            cursorStyle: ''
        })
    }

    blur() {
        this.setState({
            cursorStyle: 'ace-cursor'
        })
    }

    change(val) {
        this.setState({
            value: val
        }, () => {
            this.resize()
        })
    }
    
    // 获取当前数据
    getData() {
        return this.state.value
    }

    // 设置数据
    setData(data) {
        this.setState({
            value: data
        }, () => {
            this.resize()
        })
    }

    // 外部设置高度
    setHeight(height) {
        this.propsHeight = height
        this.setState({
            height
        })
    }

    // 重新计算高度
    resize() {
        if(!(this.lineHeight)) {
            this.lineHeight = parseInt(getComputedStyle(document.querySelector('.ace_line')).lineHeight.replace('px', ''))
        }

        const len = this.state.value.split('\n').length

        if(len * this.lineHeight > this.propsHeight) {
            this.setState({
                height: len * this.lineHeight
            })
        }else {
            this.setState({
                height: this.propsHeight
            })
        }
    }

    aceLoaded(editor) {
        editor.renderer.setPadding(0)
    }

    render() {
        return  (
            <AceEditor 
                mode="yaml"
                theme="xcode"
                width="100%"
                height={this.state.height + 'px'}
                style={{margin: 0}}
                placeholder='请输入...'
                fontSize={14}
                tabSize={2}
                value={this.state.value}
                onLoad={this.aceLoaded}
                onChange={this.change}
                name={this.props.uniqueName}
                showGutter={false}
                showPrintMargin={false}
                highlightActiveLine={false}
                onFocus={this.focus}
                onBlur={this.blur}
                editorProps={{ $blockScrolling: true }}
            />
        )
    }
}
