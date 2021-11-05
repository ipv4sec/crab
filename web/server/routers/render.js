const express = require('express')
const router = express.Router()
router.use((req, res, next) => {
    console.log('['+new Date() + '] ' + req.originalUrl)
    next()
})

router.get('/*', (req, res) => {
    res.render('index')
})

module.exports = router


