import React from 'react';
import {TextField} from 'material-ui'

/**
 * 包含label的input框
 * @param label input的label
 */

const styles = {
    errorStyle: {
      color: '#EC5858',
    },
    underlineFocusStyle: {
      borderColor: '#3986FF',
    },
  };
export default class InputComp extends React.Component {
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
                <TextField 
                    underlineFocusStyle={styles.underlineFocusStyle}
                    value={this.state.value} 
                    errorText={this.state.errorText}
                    fullWidth={true}
                    onChange={this.change.bind(this)}
                    onBlur={this.blur.bind(this)}
                />
            </div>
        )
    }
}

