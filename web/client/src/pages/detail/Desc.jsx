import React from 'react'

const Desc = (props) => {
    return (
        <section className='detail-desc'>
            <p className='detail-desc-name'>{props.name}</p>
            <pre>{props.data}</pre>
        </section>
    )
}

export default Desc