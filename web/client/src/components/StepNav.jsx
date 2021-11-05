import React, { useState } from 'react'
import '../style/sass/components.scss'
import store from '../store/store'

const StepNav = (props) => {
    const steps = [
        {label: '1', value: 1},
        {label: '2', value: 2},
        {label: '3', value: 3},
        {label: '4', value: 4},
        {label: '5', value: 5},
        {label: '6', value: 6},
    ]

    function active(index) {
        const steps = store.getState().home.steps

        if(steps.indexOf(index+1) > -1) {
            return true
        }else {
            return false
        }
    }
    return (
        <div className="step-container">
            {
                steps.map((item, index) => {
                    if(index === 0) {
                        return (
                            <div key={item.value} className="step-item"><div className={`${active(index) ? 'active-step' : ''} item-label`}>{item.label}</div></div>
                        )
                    }else {
                        return (
                            <div key={item.value} className="step-item" >
                                <div className="item-arrow"></div>
                                <div className={`${active(index) ? 'active-step' : ''} item-label`}>{item.label}</div>
                            </div>
                        )
                    }
                })
            }
        </div>
    )
}

export default StepNav