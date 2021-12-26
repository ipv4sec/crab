import React, { useEffect } from 'react'
import { useHistory } from 'react-router-dom'

const Portal = () => {
    const history = useHistory()

    useEffect(() => {
        if(window.sessionStorage.getItem('user')) {
            // window.location.replace('/home') 
            history.replace('/home')
            window.sessionStorage.setItem('curNav', '/home')
        }else {
            history.replace('/login') 
        }
    }, [])
    return (<div></div>)
}

export default Portal