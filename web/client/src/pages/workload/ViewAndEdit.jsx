import React, { useState, useEffect } from 'react'

import Dialog from '@material-ui/core/Dialog'
import DialogTitle from '@material-ui/core/DialogTitle'
import DialogContent from '@material-ui/core/DialogContent'
import DialogActions from '@material-ui/core/DialogActions'
import Button from '@material-ui/core/Button'

import '../../style/sass/workload.scss'

const ViewAndEdit = (props) => {
    const [editValue, setEditValue] = useState('')
    const [previewData, setPreviewData] = useState('dsfdsfdfdsfdsf')

    useEffect(() => {
        if(props.open) {
            if(props.type == 'view') {
                setPreviewData(props.data.value || '')
            }else if(props.type === 'edit') {  
                setEditValue(props.data.value || '')
            }
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
                <section className="workload-dialog">
                    {
                        props.type === 'view' ? (
                            <div className="workdialog-view">
                                 <pre className="preview-pre">{previewData}</pre>
                            </div>
                        ) : (
                            props.type === 'edit' ? (
                                <div className="workdialog-edit">
                                    <textarea className="textarea-input" value={editValue} onChange={changeEditValue}></textarea>
                                </div>
                            ) : null
                        )
                    }
                </section>

            </DialogContent>
            <DialogActions>
                <Button className="common-btn" color="primary" onClick={closeDialog}>关闭</Button>

                {props.type === 'edit' ? (
                    <Button className="common-btn" color="primary" onClick={confirmDialog}>保存</Button>
                ) : null}
            </DialogActions>
        </Dialog>
    )

}

export default ViewAndEdit