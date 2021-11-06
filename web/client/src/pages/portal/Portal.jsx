import React, { useEffect } from 'react'

const Portal = () => {
    useEffect(() => {
        if(window.sessionStorage.getItem('user')) {
            window.location.replace('/home') 
            window.sessionStorage.setItem('curNav', '/home')
        }else {
            window.location.replace('/login') 
        }
    }, [])
    return (<div></div>)
}

export default Portal