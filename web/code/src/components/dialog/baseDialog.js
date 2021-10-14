import React from 'react';
import { Dialog, FlatButton } from 'material-ui'

/**
 * 单独一个dialog框，不包含内部内容
 * @param open 打开/关闭
 * @param title 标题
 * @function close 关闭后触发的方法
 * @function confirm 点击确认后触发的方法
 * @param contentClass 弹框主体部分样式
 * @param closeText
 * @param submitText
 */

export default class BaseDialog extends React.Component {

    constructor(props) {
        super(props);
    }

    handleClose() {
        this.props.close()
    }

    render() {
        let actions = [];
        if(!this.props.showCloseBtn) {
            if(this.props.close){
                actions = [
                    <FlatButton key={0} label={this.props.closeText || "取消"} onClick={() => {this.props.close()}} className="close-btn" />,
                    <FlatButton key={1} label={this.props.submitText || "确认"} onClick={() => {this.props.confirm()}} className="confirm-btn"/>,
                ];
            }else {
                actions = [
                    // <FlatButton key={0} label={this.props.closeText || "取消"} onClick={() => {this.props.close()}} className="close-btn" />,
                    <FlatButton key={1} label={this.props.submitText || "确认"} onClick={() => {this.props.confirm()}} className="confirm-btn"/>,
                ];
            }
           
        }
        // console.log('=======', this.props.open, this.props.showCloseBtn, this.props.title)

        return (
            <Dialog 
                className="base-dialog-container"
                titleClassName="dialog-title"
                title={this.props.title || "提示"}
                actions={actions}
                modal={false}
                open={this.props.open}
                // onRequestClose={this.handleClose.bind(this)}
                contentClassName={`base-dialog-content ${this.props.contentClass || ''}`}
            >
                {this.props.showCloseBtn ? (
                    <button className="base-dialog-close" onClick={() => {this.props.close()}}><i className="iconfont icon_navigation_close"></i></button>
                ) : null}
                {this.props.children}
            </Dialog>
        )
    }
}