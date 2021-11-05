import React, { useState, useEffect } from 'react'
import '../style/sass/components.scss'

const Input = (props) => {

    const [value, setValue ] = useState(props.value || '')

    useEffect(() => {
        console.log('---input set --', props.set)
        setValue(props.set)
    }, [props.set])

    const change = () => {
        let value = event.target.value
        setValue(value)
        if(props.change) {
            props.change(value)
        }
    }

    return (
        <div className="input-cmp-container">
            <div className="input-cmp-content">
                {
                    props.label ? (
                        <label className="input-cmp-label">{props.label}</label>
                    ) : null
                }
                <input 
                    type={props.type || 'text'} 
                    className={`${props.inputErr ? 'input-border-hl' : ''} input-cmp-input`} 
                    value={value}
                    onChange={change}
                    onBlur={blur}
                    placeholder={props.placeholder || '请输入'}
                />
            </div>
            {
                props.inputErr ? (
                    <div className="input-cmp-error"><p>{props.inputErr}</p></div>
                ) : null
            }
           
        </div>  
    )
}

export default Input