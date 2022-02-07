import React, { useState, useEffect } from 'react'

const Crumbs = (props) => {
    // const [crumbs, setCrumbs] = useState([])

    // useEffect(() => {
    //     setCrumbs(props.data)
    // }, [props.data])

    const changeCrumb = (e) => {
        const curNav = e.currentTarget.dataset.name
        const id = e.currentTarget.dataset.id
        props.change(curNav, id)
    }

    const goLog = () => {
        props.goLog()
    }

    return (
        <nav className="detail-crumbs">
            <ul>
                <li><button className="common-btn cursor-none">{props.name || ''}</button><span>&gt;</span></li>
                {
                    props.data.map((item, idx) => {
                        return (
                            <li 
                                key={item.id} 
                                data-name={item.name} 
                                data-id={item.id} 
                                onClick={changeCrumb}
                            >
                                <button className={`common-btn ${idx !== 0 ? 'cursor-none' : ''} ${idx === (props.data.length - 1) ? 'highlight-crumb' : ''}`}>
                                    {item.name}
                                </button>
                                {idx === (props.data.length - 1) ? '' : (<span>&gt;</span>) }
                            </li>
                        )
                    })
                }
            </ul>

            {/* {
                (props.data.length && props.data[props.data.length - 1].name === 'pod') ? (
                    <button className="common-btn log-btn" onClick={goLog}>查看日志</button>
                ) : null
            } */}
          
        </nav>
    )
}

export default Crumbs