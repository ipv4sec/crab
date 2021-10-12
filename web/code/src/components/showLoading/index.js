import React from 'react'
import CircularProgress from 'material-ui/CircularProgress'
import store from '../../store/index'

export default class LoadingComp extends React.Component{
    constructor(props) {
        super(props)
    }

    render() {
        return(
            <React.Fragment>
                {
                    store.getState().common.showLoading ? (
                        <CircularProgress className="circularStyle" color="#2D6DFE" size={40} thickness={3}/>
                    ) : null
                }
            </React.Fragment>
           
        )
    }
}