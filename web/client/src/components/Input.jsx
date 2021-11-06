import React from 'react'
import '../style/sass/components.scss'

const Input = (props) => {

    const change = () => {
        if(props.change) {
            props.change(event.target.value)
        }
    }

    return (
        <div className="input-cmp-container">
            <div className="input-cmp-content">
                {
                    props.label ? (
                        <div className="input-cmp-label"><label className="input-label">{props.label}</label></div>
                    ) : null
                }
                <div className="input-context">
                    <input 
                        type={props.type || 'text'} 
                        className={`${props.inputErr ? 'input-border-hl' : ''} input-cmp-input`} 
                        value={props.value}
                        onChange={change}
                        onBlur={blur}
                        placeholder={props.placeholder || '请输入'}
                    />
                    {
                        props.inputErr ? (
                            <div className="input-cmp-error"><p>{props.inputErr}</p></div>
                        ) : null
                    }
                </div> 
                
            </div>
           
           
        </div>  
    )
}

export default Input