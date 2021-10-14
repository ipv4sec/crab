var express = require('express');
var config = require('../config/config.json');
var router = express.Router();
const path = require('path')
router.get('/*', (req, res) => {
  
    res.render('index', {
        title: 'islandConsole',
        endhost: config.webConfig.requestDomain
    })
})

module.exports = router;
