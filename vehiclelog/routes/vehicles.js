var express = require('express');
var router = express.Router();

/* GET vehicles listing. */
router.get('/', function(req, res, next) {
  var query = req.query;
  var promise = query.queryChaincode('queryAllVehicles', ['']);
  promise.then((result) => {
    res.writeHead(200, {'Content-Type': 'application/json'});
    if (result) {
      res.write(result);
    }
    res.end();
  });
});

/*
 * GET mileage of one vehicle.
 */
router.get('/mileage/:id', function(req, res, next) {
  var query = req.query;
  var id = req.params.id;
  var promise = query.queryChaincode('queryMileage', [id]);
  promise.then((result) => {
    res.writeHead(200, {'Content-Type': 'application/json'});
    if (result) {
      res.write(result);
    }
    res.end();
  });
});

/*
 * POST to adduser.
 */
router.post('/addvehicle', function(req, res) {
  var invoke = req.invoke;
  var vehicle = req.body;
  var promise = invoke.invokeChaincode('createVehicle', [vehicle.id, vehicle.type, vehicle.manufactor, vehicle.model, vehicle.registration]);
  promise.then((result) => {
    res.send(
      (result.startsWith("Failed")) ? { msg: result } : { msg: '' }
   );
  });
});

/*
 * POST to adduser.
 */
router.post('/addmileage/:id', function(req, res) {
  var invoke = req.invoke;
  var id = req.params.id;
  var mileage = req.body;
  var promise = invoke.invokeChaincode('createMileage', [id, mileage.date, mileage.mileage, mileage.comment]);
  promise.then((result) => {
    res.send(
      (result.startsWith("Failed")) ? { msg: result } : { msg: '' }
   );
  });
});

/*
 * GET all mileages of one vehicle.
 */
router.get('/mileageLog/:id', function(req, res, next) {
  var query = req.query;
  var id = req.params.id;
  var promise = query.queryChaincode('queryMileageHistory', [id]);
  promise.then((result) => {
    res.writeHead(200, {'Content-Type': 'application/json'});
    if (result) {
      res.write(result);
    }
    res.end();
  });
});


module.exports = router;
