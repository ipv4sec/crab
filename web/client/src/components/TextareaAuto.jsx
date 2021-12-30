import React, { useState , useEffect, useRef } from 'react'

const TextareaAuto = (props) => {
    const [value, changeValue] = useState('')
    const [height, changeHeight] = useState(1)
    const preRef = useRef(null)

    const textareaChange = (e) => {
        console.log(e)
        if(e.target.value === '\n') {
            changeHeight(height + 1)
        }
        changeValue(e.target.value)
    }


    return (
        <div className="textarea-auto">
            <textarea 
                className="textarea-common textarea-auto-edit" 
                height={height+'rem'}
                value={value}
                onChange={textareaChange}
            ></textarea>
            <div ref={preRef} className="textarea-common textarea-auto-view"></div>
        </div>
    )
}

export default TextareaAuto