import React, { useState, useEffect } from 'react'

import Dialog from '@material-ui/core/Dialog'
import DialogTitle from '@material-ui/core/DialogTitle'
import DialogContent from '@material-ui/core/DialogContent'
import DialogActions from '@material-ui/core/DialogActions'
import Button from '@material-ui/core/Button'

import '../../style/sass/workload.scss'

const DetailEdit = (props) => {
    const [editValue, setEditValue] = useState('')

    useEffect(() => {
        if(props.open) {
            setEditValue(props.data.value || '')
        }
    }, [props.open])

    
    const changeEditValue = (e) => {
        setEditValue(e.target.value)
    }

    const closeDialog = () => {
        props.close()
    }

    const confirmDialog = () => {
        props.confirm(editValue)
    }

    return (
        <Dialog
            open={props.open}
            close={closeDialog}
            aria-labelledby="upload-file-title"
        >
            <DialogTitle>{props.title}</DialogTitle>
            <DialogContent>
                <section className="detail-dialog">
                    <div className="detail-edit">
                        <textarea className="textarea-input" value={editValue} onChange={changeEditValue}></textarea>
                    </div>
                </section>
            </DialogContent>
            <DialogActions>
                <Button className="common-btn" color="primary" onClick={closeDialog}>关闭</Button>
                <Button className="common-btn" color="primary" onClick={confirmDialog}>保存</Button>
            </DialogActions>
        </Dialog>
    )

}

export default DetailEdit