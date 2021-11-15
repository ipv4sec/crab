
import * as express from "express";
import {NextFunction, Request, Response} from "express-serve-static-core";

let app: express.Application = express();

app.use(function (req: Request, res: Response, next: NextFunction) {
    res.json({
        "headers": req.rawHeaders,
        "host": req.hostname,
    })
})
app.listen(3000);
