import React from 'react'
import echarts from 'echarts'

const appInfoData =`[{
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
                    "instance_id": "ais-nruebpei2",
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
}]`


export default class Graph extends React.Component{
    constructor(props){
        super(props)
        this.state = {
            appList: []
        }
        this.rectWidth = 140;
        this.rectHeight = 64;
        this.imgWidth = 24
        this.imgHeight = 24; 
        this.svgWidth = 0
       
        this.curX = 0;
        this.curY = 0


    }

    componentDidMount(){
        // this.renderGraph()
        this.props.svgLoaded(true)
    }

    setData(data) {
        this.setState({
            appList: data
        },() => {
            this.renderGraph()
        })
    }

    renderGraph(){
        let echart = echarts.init(this.echartRef)

        echart.clear();

        let option = {
            // title: {
            //     text: 'Graph 简单示例'
            // },
            tooltip: {
                formatter: '{c}'
            },
           
            animationDurationUpdate: 1500,
            animationEasingUpdate: 'quinticInOut',
            series: [
                {
                    type: 'graph',
                    layout: 'none',
                    symbolSize: [140,60],
                    symbol: 'rect',
                    roam: true,
                    nodeScaleRatio: 0,
                    hoverAnimation: false,
                    animation: false,
                    label: {
                        show: true
                    },
                    // edgeSymbol: ['none', 'arrow'],
                    edgeSymbol: ['none', 'none'],
                    edgeSymbolSize: 10,
                    // edgeLabel: {
                    //     fontSize: 20
                    // },
                    data: [],
                    links:[],
                   
                    // data: [{
                    //     name: '节点1',
                    //     x: 300,
                    //     y: 300,
                    //     label: {
                    //         formatter: [
                    //             '{showImage|}{title|应用名称\nv0.0.1}',
                    //         ].join('\n'),
                    //         rich: {
                    //             showImage: {
                    //                 height: 30,
                    //                 align: 'left',
                    //                 backgroundColor: {
                    //                     image: ''
                    //                 }
                    //             },
                    //             title: {
                    //                 color: 'white',
                    //                 fontSize: 14,
                    //                 align: 'center'
                    //             },
                    //             subTitle: {
                    //                 color: 'white',
                    //                 fontSize: 12,
                    //                 align: 'center'
                    //             }
                    //         }
                    //     }
                    // }, {
                    //     name: '节点2',
                    //     x: 800,
                    //     y: 300
                    // }, {
                    //     name: '节点3',
                    //     x: 550,
                    //     y: 100
                    // }, {
                    //     name: '节点4',
                    //     x: 550,
                    //     y: 500
                    // }],
                    // links: [{
                    //     source: 0,
                    //     target: 1,
                    //     symbolSize: [5, 20],
                    //     label: {
                    //         show: true
                    //     },
                    //     lineStyle: {
                    //         width: 5,
                    //         curveness: 0.2
                    //     }
                    // }, {
                    //     source: '节点2',
                    //     target: '节点1',
                    //     label: {
                    //         show: true
                    //     },
                    //     lineStyle: {
                    //         curveness: 0.2
                    //     }
                    // }, {
                    //     source: '节点1',
                    //     target: '节点3'
                    // }, {
                    //     source: '节点2',
                    //     target: '节点3'
                    // }, {
                    //     source: '节点2',
                    //     target: '节点4'
                    // }, {
                    //     source: '节点1',
                    //     target: '节点4'
                    // }],
                    lineStyle: {
                        opacity: 0.9,
                        width: 2,
                        curveness: 0
                    }
                }
            ]
        };

       
        let data = this.renderSvg()
        // console.log('==graphdata==', data.graphData)

        let datas = []
        data.graphData.forEach((item) => {
            datas.push({
                name: item.instance_id,
                x: item.x,
                y: item.y,
                value: item.app_name,
                tooltip: {
                    formatter: '{c}'
                },
                itemStyle: {
                    color: '#54CACB',
                    shadowColor: 'black',
                    shadowBlur: 4
                },
                label: {
                    position: [14, 8],
                    formatter: [
                        '{showImage|}{title|\t\t\t'+(item.app_name.length > 7 ? (item.app_name.substring(0,7) + '...') : item.app_name )+'}\n{subTitle|\t\t\t'+item.version_id+'}',
                        
                    ].join('\n'),
                    rich: {
                        showImage: {
                            height: 20,
                            width: 20,
                            align: 'left',
                            lineHeight: 30,
                            verticalAlign: 'bottom',
                            backgroundColor: {
                                image: window._BASEPATH+'/images/grey600.png'
                            }
                        },
                        title: {
                            color: 'white',
                            fontSize: 14,
                            align: 'left',
                            verticalAlign: 'top',
                        },
                        subTitle: {
                            color: 'white',
                            fontSize: 12,
                            align: 'center'
                        }
                    }
                }
            })
        })

        

        option.series[0].data = datas
        option.series[0].links = data.linksData

        echart.setOption(option)

    }

     // 渲染svg结构
     renderSvg() {

        let graphData = []
        let linksData = []

        // 先清空
        this.curX = 0;
        this.curY = 0
        this.svgWidth = this.svgConRef.offsetWidth

        let level = 1
        const renderApp = (list, prevNodeX) => {
            if(list && list.length) {
                list.forEach((item, index) => {
                    // let selectedVersion = null
                    // if(item.app_versions) {
                    //     selectedVersion = item.app_versions.find((el) => el.selected)
                    // }
                    if(level == 1) {
                        this.curX =  this.svgWidth / 2 - this.rectWidth / 2   
                        graphData.push({
                            instance_id: item.instance_id,
                            app_name: item.app_name,
                            version_id: item.version_id == 'external' ? '外部配置' : item.version_id,
                            x: this.curX,
                            y: this.curY,
                        }) 
                    }else {
                        this.curX = (prevNodeX+this.rectWidth/2 - (this.rectWidth + 30) * list.length / 2) + index * (this.rectWidth + 30) + 15
                        // console.log()
                        if(item.version_id !== 'willHadInstance'){
                            graphData.push({
                                instance_id: item.instance_id,
                                app_name: item.app_name,
                                version_id: item.version_id == 'external' ? '外部配置' : item.version_id,
                                x: this.curX,
                                y: this.curY,
                            }) 
                        }
                       
                    }
                    

                    if(item.dependences && item.dependences.length){
                        this.curY += (50 + this.rectHeight)
                        level += 1

                        item.dependences.forEach((dpd) => {
                            if(dpd.version_id == 'willHadInstance') {
                                linksData.push({
                                    source: item.instance_id,
                                    target: dpd.selectWillInstance.instance_id,
                                    lineStyle: {
                                        curveness: 0.6
                                    }

                                })
                            }else {
                                linksData.push({
                                    source: item.instance_id,
                                    target: dpd.instance_id
                                })
                            }
                           
                        })

                        renderApp(item.dependences, this.curX)
                    }
                })
            }
        }

      
        // renderApp(JSON.parse(appInfoData), this.curX)
        renderApp(this.state.appList, this.curX)
        // console.log('graphData',graphData)
        // console.log('linksData',linksData)

        return {
            graphData,
            linksData
        }

    }


    render() {
        return (
            <div className="svg-container" ref={(ele) => {this.svgConRef = ele}}>
                <div style={{width: '100%', height: '100%'}} ref={(ele) => {this.echartRef = ele}}></div>
            </div>
        )
    }
}