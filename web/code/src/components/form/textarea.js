import React from 'react';
import {TextField} from 'material-ui'

/**
 * 包含label的input框
 * @param label input的label
 */
export default class TextAreaComp extends React.Component {
    constructor(props) {
        super(props) 
        this.state = {
            value: '',
            errorText: ''
        }
        this.setData = this.setData.bind(this)
    }


    setData(data) {
        this.setState({
            value: data.value || this.state.value,
            errorText: data.errorText || ''
        })
    }

    change(e) {
        this.setState({
            value: e.target.value,
            errorText: ''
        });
    }

    blur() {
        this.props.blur(this.state.value)
    }


    render(){
        return (
            <div className="base-input-container">
                <p className="input-label">{this.props.label}</p>
                <textarea
                    className="textarea-content"
                    value={this.state.value}
                    onChange={this.change.bind(this)}
                    onBlur={this.blur.bind(this)}
                    style={{border: this.state.errorText == '' ? "1px solid #E6E6E6" : "2px solid red"}}
                ></textarea>
            </div>
        )
    }
}

