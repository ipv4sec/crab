import React from 'react'
import '../style/sass/components.scss'

export default class Input extends React.Component{
    constructor(props) {
        super(props)

        this.state = {
            value: ''
        }

        this.setValue = this.setValue.bind(this)
    }

    change(e) {
        this.setState({
            value: e.target.value
        })
        if(this.props.change) {
            this.props.change(e.target.value)
        }
    }

    setValue(value) {
        this.setState({
            value
        })
    }

    render() {
        return (
            <div className="input-cmp-container">
                <div className="input-cmp-content">
                    {
                        this.props.label ? (
                            <div className="input-cmp-label"><label className="input-label">{this.props.label}</label></div>
                        ) : null
                    }
                   
                    <div className="input-context">
                        {
                            this.props.icon ? (
                                <div className="input-cmp-icon"><span className={`iconfont ${this.props.icon}`}></span></div>
                            ) : null
                        }
                        <input 
                            type={this.props.type || 'text'} 
                            className={`${this.props.inputErr ? 'input-border-hl' : ''} input-cmp-input`} 
                            value={this.state.value}
                            onChange={this.change.bind(this)}
                            // onBlur={this.blur.bind(this)}
                            placeholder={this.props.placeholder || '请输入'}
                        />
                        {
                            this.props.inputErr ? (
                                <div className="input-cmp-error"><p>{this.props.inputErr}</p></div>
                            ) : null
                        }
                    </div> 
                    
                </div>
               
               
            </div>  
        )
    }
}


// const Input = (props) => {

//     const change = (e) => {
//         if(props.change) {
//             props.change(e.target.value)
//         }
//     }

//     return (
//         <div className="input-cmp-container">
//             <div className="input-cmp-content">
//                 {
//                     props.label ? (
//                         <div className="input-cmp-label"><label className="input-label">{props.label}</label></div>
//                     ) : null
//                 }
               
//                 <div className="input-context">
//                     {
//                         props.icon ? (
//                             <div className="input-cmp-icon"><span className={`iconfont ${props.icon}`}></span></div>
//                         ) : null
//                     }
//                     <input 
//                         type={props.type || 'text'} 
//                         className={`${props.inputErr ? 'input-border-hl' : ''} input-cmp-input`} 
//                         value={props.value}
//                         onChange={change}
//                         onBlur={blur}
//                         placeholder={props.placeholder || '请输入'}
//                     />
//                     {
//                         props.inputErr ? (
//                             <div className="input-cmp-error"><p>{props.inputErr}</p></div>
//                         ) : null
//                     }
//                 </div> 
                
//             </div>
           
           
//         </div>  
//     )
// }

// export default Input
