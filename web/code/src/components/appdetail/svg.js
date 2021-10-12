import React from 'react'

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

export default class SVGComp extends React.Component{
    constructor(props){
        super(props)
        this.state = {
            appList:[],// this.props.allApps,
            // appList: [
            //     {
            //         name: 'app1', 
            //         id: '1', 
            //         value: '1',
            //         versions: [{version: 'v0.0.1', value: '0.0.1',selected: true, config: true},{version: 'v0.0.2', value: '0.0.2',selected: false, config: false}],
            //         parent: [
            //             {name: 'app2', id: '2', value: '2', versions: [{version: 'v0.0.1', value: '0.0.1',selected: true, config: true},{version: 'v0.0.2', value: '0.0.2',selected: false, config: false}],},
            //             {
            //                 name: 'app3', 
            //                 id: '3',
            //                 value: '3',
            //                 versions: [{version: 'v0.0.1', value: '0.0.1',selected: true, config: true},{version: 'v0.0.2', value: '0.0.2',selected: false, config: false}],
            //                 parent: [
            //                     {name: 'app4', id: '4', value: '4', versions: [{version: 'v0.0.1', value: '0.0.1',selected: true, config: true},{version: 'v0.0.2', value: '0.0.2',selected: false, config: false}],},
            //                     {name: 'app5', id: '5', value: '5',versions:[{version: 'v0.0.1', value: '0.0.1',selected: true, config: true},{version: 'v0.0.2', value: '0.0.2',selected: false, config: false}],},
            //                     {
            //                         name: 'app6', id: '6', value: '6',
            //                         versions: [{version: 'v0.0.1', value: '0.0.1', selected: true, config: true},{version: 'v0.0.2', value: '0.0.2',selected: false, config: false}],
            //                         parent: [
            //                             {name: 'app7', id: '7', value: '7', versions: [{version: 'v0.0.1', value: '0.0.1',selected: true, config: true},{version: 'v0.0.2', value: '0.0.2',selected: false, config: false}],},
            //                             {name: 'app8', id: '8', value: '8', versions: [{version: 'v0.0.1', value: '0.0.1',selected: true, config: true},{version: 'v0.0.2', value: '0.0.2',selected: false, config: false}],},
            //                             {name: 'app8', id: '8', value: '8', versions: [{version: 'v0.0.1', value: '0.0.1',selected: true, config: true},{version: 'v0.0.2', value: '0.0.2',selected: false, config: false}],},
            //                         ]
            //                     }
            //                 ]
            //             }
            //         ]
            //     }
            // ]
        }
        this.nameSpace = "http://www.w3.org/2000/svg"
        this.rectWidth = 140;
        this.rectHeight = 64;
        this.imgWidth = 24
        this.imgHeight = 24; 
        this.svgWidth = 0
        this.stepY = 50;
        this.curX = 0;
        this.curY = 0
        this.downTransX = 0
        this.downTransY = 0
        this.setData = this.setData.bind(this)
    }

    componentWillMount(){
        // console.log('---svg component will---')
    }

    componentDidMount() {
        // console.log('---svg component ---')
      
    }

    setData(data) {
        this.setState({
            appList: data
        },() => {
            this.renderSvg()
        })
    }
    
    svgload(){
        this.svgDom = document.getElementById("svgId").getSVGDocument()
        this.svgRoot = this.svgDom.getElementById("imageSvg")
        this.svgRoot.setAttribute("width", this.svgConRef.offsetWidth)
        this.svgRoot.setAttribute("height", this.svgConRef.offsetHeight)
        this.svgWidth = this.svgConRef.offsetWidth
        let _this = this

        this.props.svgLoaded(true)


        /* ----- 给svg绑定拖动事件  开始---- */ 
        this.svgDom.addEventListener('mousedown', function(e){
            _this.canMove = true
           _this.downX = e.pageX
           _this.downY = e.pageY 
        })
        this.svgDom.addEventListener('mousemove', function(e){
            if( _this.canMove) {
                _this.svgDom.querySelectorAll('g').forEach((item) => {
                    item.setAttribute("transform", "translate("+(_this.downTransX + (e.pageX - _this.downX))+","+(_this.downTransY + (e.pageY - _this.downY))+")")
                })
            }
           
        })
        this.svgDom.addEventListener('mouseup', function(e){
            _this.canMove = false
            _this.downTransX = _this.downTransX + (e.pageX - _this.downX)
            _this.downTransY = _this.downTransY + (e.pageY - _this.downY)
        })

         /* ----- 给svg绑定拖动事件  结束---- */ 

    }

    // 渲染svg结构
    renderSvg() {

        // 先清空
        this.curX = 0;
        this.curY = 0
        this.downTransX = 0
        this.downTransY = 0
        this.svgRoot.innerHTML = "";

        let level = 1
        const renderApp = (list, prevNodeX) => {
            if(list && list.length) {
                list.forEach((item, index) => {
                    let selectedVersion = []
                    if(item.app_versions) {
                        selectedVersion = item.app_versions.filter((el) => el.selected)
                    }
                    if(level == 1) {
                        this.curX =  this.svgWidth / 2 - this.rectWidth / 2    
                    }else {
                        this.curX = (prevNodeX+this.rectWidth/2 - (this.rectWidth + 30) * list.length / 2) + index * (this.rectWidth + 30) + 15

                        this.createLine(prevNodeX + this.rectWidth/2, this.curY - 50, this.curX+ this.rectWidth/2, this.curY)
                    }
                    

                    this.createOneRect(this.curX, this.curY, item.app_name, selectedVersion&&selectedVersion.length ? selectedVersion[0].version : '')
                    if(item.dependences && item.dependences.length){
                        this.curY += (50 + this.rectHeight)
                        level += 1
                        renderApp(item.dependences, this.curX)
                    }
                })
            }
        }

        renderApp(this.state.appList, this.curX)
    }


    createOneRect(x, y, appName, version) {
        let app_name = appName
        if(appName && appName.length > 6) {
            app_name = appName.substring(0,6) + '...'
        }
        let g = this.svgDom.createElementNS(this.nameSpace, 'g')
        let rect = this.createRect(x, y, this.rectWidth, this.rectHeight, 10, 10, "#54CACB")
        let image = this.createImage(x+20, y+this.rectHeight/2-this.imgHeight/2, this.imgWidth, this.imgHeight, "/images/app-icon.svg")
        let textTitle = this.createText(x+30+this.imgWidth, y+this.rectHeight/2-4, app_name , "white", "14px" )
        let textVersion = this.createText(x+30+this.imgWidth+4, y+this.rectHeight-16, version == 'external' ? '外部配置' : (version === 'willHadInstance' ? '即有实例' : version), "white", "11px" )
        g.appendChild(rect)
        g.appendChild(image)
        g.appendChild(textTitle)
        g.appendChild(textVersion)
        this.svgRoot.appendChild(g)
    }

    createRect( x, y, width, height,  rx, ry, color) {
       
        let rect = this.svgDom.createElementNS(this.nameSpace, 'rect')
        rect.setAttribute("x", x)
        rect.setAttribute("y", y)
        rect.setAttribute("width", width)
        rect.setAttribute("height", height)
        rect.setAttribute("rx", rx)
        rect.setAttribute("ry", ry)
        rect.setAttribute("fill", color || '#54CACB')

        return rect
      
        // this.svgRoot.appendChild(g)
    }

    createLine(x1, y1, x2, y2) {
        let g = this.svgDom.createElementNS(this.nameSpace, 'g')
        let line = this.svgDom.createElementNS(this.nameSpace, 'line')
        line.setAttribute("x1", x1)
        line.setAttribute("y1", y1)
        line.setAttribute("x2", x2)
        line.setAttribute("y2", y2)
        line.setAttribute("stroke", "#B4B4B4")
        line.setAttribute("stroke-width", "2")
      
        g.appendChild(line)
        this.svgRoot.appendChild(g)
    }

    createText(x, y, content, color, size) {
        let text = this.svgDom.createElementNS(this.nameSpace, 'text')
        text.setAttribute("x", x)
        text.setAttribute("y", y)
        text.setAttribute("fill", color)
        text.setAttribute("text-anchor", "start")
        text.setAttribute("font-size", size)
        text.textContent = content
        // this.svgRoot.appendChild(text)
        return text
    }

    createImage( x, y, width, height, url) {
        let image = this.svgDom.createElementNS(this.nameSpace, 'image')
        image.setAttribute("x", x)
        image.setAttribute("y", y)
        image.setAttribute("width", width)
        image.setAttribute("height", height)
        image.href.baseVal=url
        // this.svgRoot.appendChild(image)
        return image
    }


    render(){
        return (
            <div className="svg-container" ref={(ele) => {this.svgConRef = ele}}>
                <object id="svgId" data={`/images/appsvg.svg`} type="image/svg+xml"  onLoad={this.svgload.bind(this)}/>
            </div>
        )
    }
}