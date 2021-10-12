/**
 * Created by tianyaming on 2019/2/26.
 */
import React from 'react';

export default class Pagination extends React.Component{
    constructor(props){
        super(props);

        this.state = {
            total: this.props.total || 0,       //分页的总页数
            totalRes: this.props.totalRes || 0, //用来分页的数据总数
            current: this.props.current || 0,   //当前的分页页数
            showPage: 10,                       //显示10个按钮，其余用...代替
        };
    }

    componentWillReceiveProps(nextProps){
        if(JSON.stringify(nextProps) != JSON.stringify(this.props)){
            this.setState({
                total: nextProps.total || 0,       //分页的总页数
                totalRes: nextProps.totalRes || 0, //用来分页的数据总数
                current: nextProps.current || 0,   //当前的分页页数
            });
        }
    }

    handleChange(number){
        this.props.handleChange && this.props.handleChange(number);
    }

    renderPage(){
        let total = this.state.total,
            showPage = this.state.showPage,
            current = this.state.current;

        let page = [];

        if(total <= showPage){
            for(let i=0;i<total;i++){
                let cname = current === i ? 'active' : '';

                page.push(
                    <li key={i+1} className={`page-number ${cname}`} onClick={ this.handleChange.bind(this, i) }>
                        <a>{i+1}</a>
                    </li>
                );
            }

        }else{
            let mid = showPage/2,
                begin = 0;

            if(current <= mid){
                begin = 0;
            }else if(current >= total - mid + 1){
                begin = total - showPage;
            }else{
                begin = current - mid;
            }

            for(let i=begin;i<begin+showPage;i++){
                let cname = current === i ? 'active' : '';

                page.push(
                    <li key={i+1} className={`page-number ${cname}`} onClick={ this.handleChange.bind(this, i) }>
                        <a>{i+1}</a>
                    </li>
                );
            }
        }

        return page;
    }

    /**
     * 上一页
     * */
    prev(){
        let current = this.state.current;
        if(current == 0){
            return;
        }else{
            this.handleChange(current - 1);
        }
    }
    /**
     * 下一页
     * */
    next(){
        let current = this.state.current,
            total = this.state.total;

        if(current == total - 1){
            return;
        }else{
            this.handleChange(current + 1);
        }
    }

    render(){
        let className = this.props.className || '';

        let prevDisabled = this.state.current === 0 ? true : false;
        let nextDisabled = this.state.current === this.state.total - 1 ? true : false;

        let allStyle = this.state.total > 1 ? {marginLeft:15} : {justifyContent:'center',flexGrow: 1};

        return(
            <div className="page-pagination">
                <div className={`pagination ${className}`}>
                    {
                        this.state.total > 1 ?
                            <ul>
                                {
                                    this.state.current == 0 ?
                                        null
                                        :
                                        <li className={`page-text page-prev ${prevDisabled ? 'disabled' : ''}`}
                                            onClick={ this.prev.bind(this) }>
                                            <a> 上一页 </a>
                                        </li>
                                }
                                { this.renderPage() }
                                {
                                    this.state.current == this.state.total - 1 ?
                                        null
                                        :
                                        <li className={`page-text page-next ${nextDisabled ? 'disabled' : ''}`}
                                            onClick={ this.next.bind(this) }>
                                            <a> 下一页 </a>
                                        </li>
                                }
                            </ul>
                        :
                            null
                    }


                    {
                        (this.state.totalRes === 0 || this.state.totalRes) ?
                            <div className="all" style={ allStyle }>
                                {this.state.totalRes > 4*10000 ? '约' : '共'} {this.state.totalRes} 条结果
                            </div>
                            :
                            null
                    }
                </div>
            </div>
        );
    }
}