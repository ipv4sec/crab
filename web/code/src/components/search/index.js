import React from 'react';
import IconButton from 'material-ui/IconButton';

export default class SearchBox extends React.Component {
    constructor() {
        super();
        this.state = {
            value: '',
        }
        this.clearValue = this.clearValue.bind(this)
    }

    componentDidMount() {
        if (this.props.init) {
            this.setState({
                value: this.props.init,
            })
        }
    }

    clearValue() {
        this.setState({
            value: ''
        })
    }

    initData() {
        if (this.props.init) {
            this.setState({
                value: this.props.init,
            })
        }
    }

    handleKeyUp(e) {
        if (e.keyCode == 13) {
            if (this.props.callback) {
                this.props.callback(this.state.value);
            }
        }
    }

    render() {
        return (
            <div className="search-box-component clearfix">
                <input type="text" className="input-area" placeholder={this.props.placeholder || ''}
                    id={this.props.id} value={this.state.value}
                    onChange={(e) => {
                        this.setState({
                            value: e.target.value,
                        });
                        if(this.props.changeCallback) {
                            this.props.changeCallback(e.target.value)
                        }
                    }}
                    onKeyUp={this.handleKeyUp.bind(this)}/>
                <IconButton iconClassName="iconfont icon_search" className="search-btn"
                    onClick={() => {
                        if (this.props.callback) {
                            this.props.callback(this.state.value);
                        }
                    }}
                />
            </div>
        )
    }
}