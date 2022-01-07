// 防抖函数
function fangdou() {
    let timer = null
    return function(fn, args) {
        clearTimeout(timer)
        timer = setTimeout(() => {
            fn(args)
            timer = null
        }, 300)
    }
}