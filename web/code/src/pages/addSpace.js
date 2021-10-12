import React, { version } from 'react'
import { TextField, SelectField, MenuItem } from 'material-ui';
import SearchBox from '../components/search/index'
import BaseDialog from '../components/dialog/baseDialog'
import store from '../store/index'
import * as TYPE from '../store/actions'
import axios from 'axios';
import { browserHistory } from 'react-router'
// const appInfo = `{"app_id":"app-jlyN52Rm","instance_id":"ais-891ycdvw","version_id":"v0.0.6","dependences":[{"app_id":"app-Gd9sfVSa","instance_id":"ais-lnsgoa01","version_id":"v0.0.2","dependences":[{"app_id":"app-b3Dwlhzc","instance_id":"ais-nruebpei","version_id":"v0.0.2","dependences":[{"app_id":"app-jlyN52Rm","version_id":"v0.0.6"}]}]}]}`
const appInfoData =`{
    "app_id": "app-jlyN52Rm",
    "app_name": "app_demo_for_asb_1",
    "app_versions": [
        {"selected": false, "version": "external"},
        {"selected": false, "version": "v0.0.1"},
        {"selected": false, "version": "v0.0.2"},
        {"selected": false, "version": "v0.0.3"},
        {"selected": false, "version": "v0.0.4"},
        {"selected": false, "version": "v0.0.5"},
        {"selected": true, "version": "v0.0.6"}
    ],
    "conf":{},
    "external": 0,
    "instance_id": "ais-891ycdvw",
    "selectWillInstance": null,
    "version_id": "v0.0.6",
    "willInstances": [],
    "dependences": [
        {
            "app_id": "app-Gd9sfVSa",
            "app_name": "app_demo_for_asb_2",
            "app_versions": [
                {"selected": false, "version": "external"},
                {"selected": false, "version": "v0.0.1"},
                {"selected": true, "version": "v0.0.2"}
            ],
            "conf":{},
            "external": 0,
            "instance_id": "ais-nruebpei",
            "selectWillInstance": null,
            "version_id": "v0.0.2",
            "willInstances": [],
            "dependences": [
                {
                    "app_id": "app-b3Dwlhzc",
                    "app_name": "app_demo_for_asb_3",
                    "app_versions":  [
                        {"selected": false, "version": "external"},
                        {"selected": false, "version": "v0.0.1"},
                        {"selected": true, "version": "v0.0.2"}
                    ],
                    "conf": {},
                    "external": 0,
                    "instance_id": "ais-nruebpei",
                    "selectWillInstance": null,
                    "version_id": "v0.0.2",
                    "willInstances": [],
                    "dependences": [
                       {
                            "app_id": "app-jlyN52Rm",
                            "app_name": "app_demo_for_asb_1",
                            "app_versions": [
                                {"selected": false, "version": "external"},
                                {"selected": true, "version": "willHadInstance"},
                                {"selected": false, "version": "v0.0.1"},
                                {"selected": false, "version": "v0.0.2"},
                                {"selected": false, "version": "v0.0.3"},
                                {"selected": false, "version": "v0.0.4"},
                                {"selected": false, "version": "v0.0.5"},
                                {"selected": false, "version": "v0.0.6"}
                            ],
                            "conf": {},
                            "dependences": [],
                            "external": 0,
                            "instance_id": "ais-oy1lnpv1",
                            "selectWillInstance": {"app_id": "app-jlyN52Rm", "version_id": "v0.0.6", "instance_id": "ais-891ycdvw"},
                            "version_id": "willHadInstance",
                            "willInstances": [
                               {
                                    "app_id": "app-jlyN52Rm",
                                    "app_name": "app_demo_for_asb_1",
                                    "instance_id": "ais-891ycdvw",
                                    "version_id": "v0.0.6"
                               }
                            ]
                       }
                    ]
                }
            ]
        }
    ]
}`


const styles = {
    errorStyle: {
      color: '#EC5858',
    },
    underlineFocusStyle: {
      borderColor: '#3986FF',
    },
};

const appListData = [
    {name: '111', id: '1', color: 'red'},
    {name: '222', id: '2', color: 'blue'},
    {name: '333', id: '3', color: 'green'},

    {name: '111', id: '1', color: 'red'},
    {name: '222', id: '2', color: 'blue'},
    {name: '333', id: '3', color: 'green'},

    {name: '111', id: '1', color: 'red'},
    {name: '222', id: '2', color: 'blue'},
    {name: '333', id: '3', color: 'green'},

    {name: '111', id: '1', color: 'red'},
    {name: '222', id: '2', color: 'blue'},
    {name: '333', id: '3', color: 'green'},

    {name: '111', id: '1', color: 'red'},
    {name: '222', id: '2', color: 'blue'},
    {name: '333', id: '3', color: 'green'},
]

export default class AddSpace extends React.Component {
    constructor(props) {
        super(props)
        this.state ={
            name: '',
            errorName: '',
            description: '',
            errorDesc: '',
            cluster: '1',
            clusterList: [
                {name: '111', id: '1'},
                {name: '222', id: '2'},
                {name: '333', id: '3'},
            ], 
            addAppList: [],
            originAppList: [],
            appList: [], //appListData,//[],
            currentApp: null,
            current_page: 0,
            pageSize: 20,
            total: 0
            // type: this.props.params.type,
            // id: this.props.params.id

        }
    }

    componentDidMount() {
        let _this = this
       window.addEventListener('message', (e) => {
           console.log('message===',e)
            if(e.origin === window.location.origin ) {
                if(e.data) {
                    let data = JSON.parse(e.data)
                    if(data.action == 'refreshAppList' && data.appData && data.appData.length) {
                        let appInfo = data.appData[0]
                        let newAppList = _this.state.addAppList.slice()
                        if(data.isEdit) {
                            // 修改内容
                            newAppList.forEach((item, index) => {
                                if(item.app_id == appInfo.app_id){
                                    newAppList.splice(index, 1, appInfo)
                                }
                            })
                        }else {
                            // 添加内容
                            newAppList.push(appInfo)
                            // sessionStorage.setItem('spaceAppList', JSON.stringify(newAppList))
                           
                        }

                        _this.setState({
                            addAppList: newAppList
                        })
                       
                    }
                }
            }
        })

        
        // console.log('=-----',JSON.parse(appInfoData))
        // //测试用的部分
        // this.state.addAppList.push(JSON.parse(appInfoData))
        // this.setState({
        //     addAppList: this.state.addAppList
        // })


        this.getClusterData()
        
    }

    getClusterData() {
        let that = this

        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })
 
        axios({
            method: 'GET',
            url: `${window._BASEPATH}/api/getAllCluster`
        }).then((res) => {
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
            if(res.data.list && res.data.list.length) {
                that.setState({
                    clusterList: res.data.list || [],
                    cluster: res.data.list[0].id || ''
                })
            }else {
                that.setState({
                    clusterList: [],
                    cluster: ''
                })
            }
          
           
        }).catch((err) => {
            console.log(err)
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
        })
    }


    filterApp(name) {
        let reg = new RegExp(name, 'ig')
        let newData = this.state.originAppList.filter((item) => reg.test(item.name))
        this.setState({
            appList: newData
        })
    }


    nameBlur() {
        // 如果需要校验，可以在这里写
        this.verifyName()
        
    }
    descBlur() {
        // 如果需要校验，可以在这里写
    }

    verifyName() {
        if(this.state.name == ''){
            this.setState({
                errorName: '请输入名称'
            })
            return false
        }
        if(this.state.name.length > 30 ) {
            this.setState({
                errorName: '名字长度不能超过30个字符'
            })
            return false
        }
        return true
    }

    getAppList() {
        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })
        let data = {
            offset: this.state.current_page * this.state.pageSize,
            limit: this.state.pageSize,
            // search: this.state.search
        }
        axios({
            method: 'GET',
            url: `${window._BASEPATH}/api/getAllApp`,
            params: data
        }).then((res) => {
            let curAppList = this.state.originAppList.slice()
            let newList = curAppList.concat(res.data.list || [])
           
            this.setState({
                total: res.data.total ,
                appList: newList,
                originAppList: newList//res.data.list || []
            })


            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })


        }).catch((err) => {
            console.log(err)
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
        })
    }

    addApp() {
        this.setState({
            showDialog: true
        })
        
        this.getAppList()
    }

    appListScroll() {
        if(this.state.total != this.state.originAppList.length && (event.target.scrollHeight == (event.target.scrollTop + event.target.offsetHeight))){
            this.setState({
                current_page: ++this.state.current_page
            },() => {
                this.getAppList()
            })
        }
    }   

    closeDialog() {
        this.setState({
            showDialog: false,
            offset: 0,
            appList: [],
            originAppList: []
        })
    }

    search(value) {
        this.filterApp(value)
    }

    goAppDetail(curApp) {
        this.setState({
            currentApp:  curApp,
            showDialog: false,
            offset: 0,
            appList: [],
            originAppList: []
        })
        window.open(`/appdetail/${curApp.id}/0`, "_blank")
    }

    

    editApp(index, app_id){
        sessionStorage.setItem("appInfo", JSON.stringify(this.state.addAppList[index]))
        window.open(`/appdetail/${app_id}/1`, "_blank")
    }

    deleteApp(index){
        let newList = this.state.addAppList.slice()
        newList.splice(index,1)
        this.setState({
            addAppList: newList
        })
    }

    generateInstanceId(randomFlag, min, max){
        var str = "",
            range = min,
            arr = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'];
        
        // 随机产生
        if(randomFlag){
            range = Math.round(Math.random() * (max-min)) + min;
        }
        for(var i=0; i<range; i++){
            let pos = Math.round(Math.random() * (arr.length-1));
            str += arr[pos];
        }
        return 'ais-'+str;
    }

    submit(){
        //把里面的数据选择的数据拿出来
        if(!this.verifyName() ){
            return
        }
        if(this.state.cluster == '') {
            store.dispatch({
                type: TYPE.SHOW_SNACKBAR,
                val: {open: true, message: '请选择集群'}
            })
            return
        }

        let app_info = []
        const getData = (apps) => {
            if(apps && apps.length) {
                let depApp = []
                apps.forEach((item) => {
                    let tmp = {}
                    tmp['app_id'] = item.app_id
                    tmp['instance_id'] = item.instance_id || this.generateInstanceId()//'ais-asdf1234'//new Date().getTime().toString()
                    let selectedVersion = item.app_versions.find((val) => val.selected)
                    if(selectedVersion && selectedVersion.version == 'external') {
                        tmp['external'] = 1
                        tmp['conf'] = item.conf
                    }else if(selectedVersion && selectedVersion.version == 'willHadInstance') {
                        tmp['app_id'] = item.selectWillInstance.app_id
                        tmp['version_id'] = item.selectWillInstance.version_id
                        tmp['instance_id'] = item.selectWillInstance.instance_id
                    }else if(selectedVersion){
                        tmp['version_id'] = selectedVersion.version
                    }else {
                        tmp['version_id'] = ''
                    }
                    tmp['dependences'] = []
                    if(item.dependences) {
                        tmp['dependences'] = getData(item.dependences)
                    }

                    depApp.push(tmp)
                   
                })
                return depApp
            }
        }

        this.state.addAppList.forEach((item, index) => {
            let data = getData([item])
            if(data.length) {
                app_info.push(data[0])
            }
        })

        if(app_info.length == 0){
            store.dispatch({
                type: TYPE.SHOW_SNACKBAR,
                val: {open: true, message: '请添加应用'}
            })
            return
        }

        let data = {
            name: this.state.name,
            description: this.state.description,
            cluster_id: this.state.cluster,
            app_info: JSON.stringify(app_info)
        }

        console.log('addSpaceData===',data)
    
        store.dispatch({
            type: TYPE.SHOW_LOADING,
            val: true
        })

        axios({
            method: 'POST',
            url: `${window._BASEPATH}/api/addSpace`,
            data: data
        }).then((res) => {
            if(res.data.id) {
                store.dispatch({
                    type: TYPE.SHOW_SNACKBAR,
                    val: {open: true, message: '添加成功'}
                })
                store.dispatch({
                    type: TYPE.SHOW_LOADING,
                    val: false
                })
                browserHistory.push('/space')
            }else {
                store.dispatch({
                    type: TYPE.SHOW_SNACKBAR,
                    val: {open: true, message: res.data.msg || '添加失败'}
                })
            }
        }).catch((err) => {
            store.dispatch({
                type: TYPE.SHOW_SNACKBAR,
                val: {open: true, message: '请求错误'}
            })
            store.dispatch({
                type: TYPE.SHOW_LOADING,
                val: false
            })
        })



    }



    render() {
        return (
           <div className="add-space-content">
                 <div className="space-title">
                    <span className="gray-title" onClick={() => {browserHistory.push("/space")}}>空间管理</span>
                    <span className="second-title" ><i className="iconfont icon_arrow-right"></i>添加</span>
                </div>

                <div className="space-input" style={{marginTop: '20px'}}> 
                    <label>名称</label>
                    <TextField
                        value={this.state.name}
                        errorText={this.state.errorName}
                        underlineFocusStyle={styles.underlineFocusStyle}
                        fullWidth={true}
                        onChange={(e) => {
                            this.setState({
                                name: e.target.value,
                                errorName: ''
                            })
                        }}
                        onBlur={this.nameBlur.bind(this)}
                    />
                </div>
                <div className="space-input">
                    <label>描述</label>
                    <TextField
                        value={this.state.description}
                        errorText={this.state.errorDesc}
                        underlineFocusStyle={styles.underlineFocusStyle}
                        fullWidth={true}
                        onChange={(e) => {
                            this.setState({
                                description: e.target.value
                            })
                        }}
                        onBlur={this.descBlur.bind(this)}
                    />
                </div>
                <div className="space-input">
                    <label>集群</label>
                    <SelectField
                        value={this.state.cluster}
                        underlineFocusStyle={styles.underlineFocusStyle}
                        onChange={(event,index,value) => {
                            this.setState({
                                cluster: value
                            })
                        }}
                        fullWidth={true}
                    >
                        {
                            this.state.clusterList.map((item, index) => {
                                return(
                                    <MenuItem innerDivStyle={{color: '#4a4a4a'}}  key={item.id} value={item.id} primaryText={item.name} />
                                )
                            })
                        }
                    </SelectField>
                </div>

                <div className="space-app-list">
                    <p className="list-title">应用</p>
                    <ul>
                        {
                            this.state.addAppList.map((item,index) => {
                                return (
                                    <li key={index}>
                                        <div className="list-item-name">
                                            <div className="item-name-icon" style={{backgroundColor: item.color||'#54CACB'}}><i  className="iconfont icon_grey600"></i></div>
                                            {item.app_name}
                                        </div>
                                        <div className="list-item-btn">
                                            <button onClick={this.editApp.bind(this, index, item.app_id)}>编辑</button>
                                            <button onClick={this.deleteApp.bind(this, index)}>删除</button>
                                        </div>

                                    </li>
                                )
                            })
                        }
                         <li>
                            <div className="list-item-add" onClick={this.addApp.bind(this)}>
                                <i className="iconfont icon_baseline_add_blackgarden"></i>
                                <button >添加应用</button>
                            </div>

                        </li>
                    </ul>
                </div>

                <div className="submit-space">
                    <button onClick={this.submit.bind(this)}>添加</button>
                </div>




                <BaseDialog
                    open={this.state.showDialog}
                    title="添加应用"
                    showCloseBtn={true}
                    close={this.closeDialog.bind(this)}
                    contentClass="middle-dialog"
                >
                   <div className="app-dialog-container">
                        <div className="search-apps">
                            <SearchBox placeholder="搜索" changeCallback={this.search.bind(this)} />
                        </div>
                        <div className="dialog-app-list"  onScroll={this.appListScroll.bind(this)}>
                            <ul>
                                {
                                    this.state.appList.map((item,index) => {
                                    // appListData.map((item,index) => {
                                        
                                        return (
                                            <li key={index}>
                                                <div className="list-item-name">
                                                    <div className="item-name-icon" style={{backgroundColor: item.color||'#54CACB'}}><i  className="iconfont icon_grey600"></i></div>
                                                    {item.name}
                                                </div>
                                                <div className="list-item-btn">
                                                    <button onClick={this.goAppDetail.bind(this, item)}><i className="iconfont icon_dir-right"></i></button>
                                                </div>
                                            </li>
                                        )
                                    })
                                }
                            </ul>
                        </div>
                   </div>
                </BaseDialog>


           </div>
        )
    }
}