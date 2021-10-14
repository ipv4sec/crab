import React from 'react';

export default class HelpTip extends React.Component {
    constructor() {
        super();
    }

    render() {
        let content = this.props.content;
        let width = this.props.width || '180px';
        let left = this.props.left || '-24px';
        let top = this.props.top || '30px';
        let size = this.props.size || '16px';
        let icon = this.props.icon || 'icon_help';
        let iconColor = this.props.iconColor || '#9B9B9B';
        let classname = this.props.className || '';
        return (
            <div className={"global-help-tip " + classname}>
                <span className={"iconfont " + icon}
                      style={{color: iconColor, fontSize: size, lineHeight: size}}></span>
                <div className="box" style={{width: width}}>
                    <div className="help-text" style={{left: left, top: top}}
                         dangerouslySetInnerHTML={{__html: content}}>
                    </div>
                </div>
            </div>
        )
    }
}