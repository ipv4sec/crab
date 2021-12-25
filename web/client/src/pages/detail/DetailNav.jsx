import React, { useState, useEffect } from 'react'

const DetailNav = (props) => {
    const [navList, setNavList] = useState([])

    const dealData = (data) => {
        let newData = []
        if(!(Object.keys(data).length)) { return newData}

        Object.keys(data).forEach((key) => {
            const tmp = {
                name: key,
                children: []
            }
            data[key].forEach(el => {
                tmp.children.push(el.metadata)
            })
            newData.push(tmp)
        })

        return newData
    }

    useEffect(() => {
        // setNavList(dealData(props.data))
        setNavList(props.data)
    }, [props.data])

    const changeNav = (e) => {
        const type = e.currentTarget.dataset.type
        const idx = e.currentTarget.dataset.idx
        const curNav = navList[idx]
        if(type === 'parent') {
            props.change({index: idx, parent: curNav, child: ''})
        }else if(type === 'child') {
            const cidx = e.currentTarget.dataset.cidx
            props.change({index: idx, cIndex: cidx, parent: curNav, child: curNav.navList[cidx]})
        }
    }

    return (
        <div className="detail-nav">
            {
                navList.map((item, index) => {

                    return (
                        <div className="dnav-item" key={item.id}>
                            <div className={`dnav-title ${props.curNav === item.id ? 'blueBorder' : ''}`} data-idx={index} data-type="parent" onClick={changeNav}>{item.name}</div>
                            {/* <ul className="dnav-list">
                                {
                                    item.navList.map((el, idx) => {
                                        return (
                                            <li className={`dnav-li ${props.curNav === el.id ? 'blueBorder' : ''}`} key={el.id} data-idx={index} data-cidx={idx} data-type="child" onClick={changeNav}>{el.name}</li>             
                                        )
                                    })
                                }
                            </ul> */}
                        </div>
                    )
                })
            }
           
        </div>
    )
}

export default DetailNav