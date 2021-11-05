import React from 'react'
import Dialog from '@material-ui/core/Dialog'
import DialogTitle from '@material-ui/core/DialogTitle'
import DialogContent from '@material-ui/core/DialogContent'
import DialogActions from '@material-ui/core/DialogActions'
import Button from '@material-ui/core/Button'
import '../style/sass/readlog.scss'

const ReadLog = (props) => {

    function closeDialog() {
        props.close()
    }

    return (
        <Dialog
            open={props.open}
            onClose={closeDialog}
            aria-labelledby="read-log-title"    
        >
           <DialogTitle id="read-log-title">{props.title}</DialogTitle>
           <DialogContent>
                <div className="log-list">
                    {
                        props.data.map((item, index) => {
                            return (
                                <div key={item.name} className="log-item">
                                    <p>{item.name}：</p>
                                    <p className="item-desc">{item.message}</p>
                                </div>
                            )
                        })
                    }
                </div>
           </DialogContent>
           <DialogActions>
                {/* <Button className="common-btn" color="primary" onClick={submitDialog}>确定</Button> */}
                <Button className="common-btn" color="primary" onClick={closeDialog}>取消</Button>
            </DialogActions>
        </Dialog>
    )
}

export default ReadLog