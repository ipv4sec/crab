// var config = require('../../config/config.json');
 
exports.getHeader = function (req) {

    let requestId = req.get('X-Request-Id');
    // let uid = req.get('C-Uid');
    // let name = req.get('C-Name');
    // let firstname = req.get('C-Firstname');
    // let lastname = req.get('C-Lastname');
    // let avatar = req.get('C-Avatar');
    // let email = req.get('C-Email');
    // let domain = req.get('C-Domain');
    // let csId = req.get('C-Cs-Id');
    // let instanceId = req.get('C-Instance-Id');
    // let groups = req.get('C-Groups');
    // let roles = req.get('C-Roles');


    let headers = {
        'Content-Type': 'application/json',
        'X-Request-Id': requestId,

        // 'C-Uid': uid,
        // 'C-Name': name,
        // 'C-Firstname': firstname,
        // 'C-Lastname': lastname,
        // 'C-Avatar': avatar,
        // 'C-Email': email,
        // 'C-Domain': domain,
        // 'C-Cs-Id': csId,
        // 'C-Instance-Id': instanceId,
        // 'C-Groups': groups,
        // 'C-Roles': roles,
    };

   

    // if (config.appConfig.env === 'dev') {
    //     let headers_test = config.appConfig.dev_header;
    //     return headers_test;
    // } else {
    //     return headers;
    // }

    return headers;

};

