
// // var PROTO_PATH = __dirname + '../../../../iter8-tools/iter8/abn/grpc/abn.proto';
// var PROTO_PATH = '../../../../iter8-tools/iter8/abn/grpc/abn.proto';
// var grpc = require('@grpc/grpc-js');
// var protoLoader = require('@grpc/proto-loader');
// // Suggested options for similarity to existing grpc.load behavior
// var packageDefinition = protoLoader.loadSync(
//     PROTO_PATH,
//     {keepCase: true,
//      longs: String,
//      enums: String,
//      defaults: true,
//      oneofs: true
//     });
// var protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
// // The protoDescriptor object has the full package hierarchy
// var abn_proto = protoDescriptor.abnservice
// console.log(`${abn_proto}`)

var messages = require('./abn_pb.js');
var services = require('./abn_grpc_pb.js');
var grpc = require('@grpc/grpc-js');
var http = require('http');

'use strict';
const express = require('express');
const { registerChannelzSubchannel } = require('@grpc/grpc-js/build/src/channelz.js');
const { application } = require('express');

const app  = express();

// map of track to route to backend service
const trackToRoute = {
    "current":   "http://backend-current:8090",
    "candidate": "http://backend-candidate:8090",
}

// implement /hello endpoint
// calls backend service /world endpoint
app.get('/hello', (req, res) => {
	// Get user (session) identifier, for example by inspection of header X-User
    const user = req.header('X-User')

    // Get endpoint of backend endpoint "/world"
	// In this example, the backend endpoint depends on the version (track) of the backend service
	// the user is assigned by the Iter8 SDK Lookup() method

    // establish connection to ABn service
    var client = new services.ABNClient('abn:50051', grpc.credentials.createInsecure());

    // call ABn service API Lookup() to get an assigned track for the user
    var application = new messages.Application();
    application.setName('default/backend');
    application.setUser(user);
    client.lookup(application, function(err, session) {
        // lookup route using track
        var route = trackToRoute[session.getTrack()];
        console.log(route)

        // call backend service using url
        http.get(route + '/world', (resp) => {
            let str = '';
            resp.on('data', function(chunk) {
                str += chunk;
            });
            resp.on('end', function () {
                // write response to query
                res.send(`Hello world ${str}`);
            });
        }).on("error", (err) => {
            console.log("Error: " + err.message)
        });
    });
});

// implement /goodbye endpoint
// writes value for sample_metric which may have spanned several calls to /hello
app.get('/goodbye', (req, res) => {
	// Get user (session) identifier, for example by inspection of header X-User
    const user = req.header('X-User')

	// export metric to metrics database
	// this is best effort; we ignore any failure

	// establish connection to ABn service
    var client = new services.ABNClient('abn:50051', grpc.credentials.createInsecure());

    // export metric
    var mv = new messages.MetricValue();
    mv.setName('sample_metric');
    mv.setValue(); 
    mv.setApplication('default/backend');
    mv.setUser(user);
    client.writeMetric(mv, function(err, session) {});
    res.sendStatus(200);
});

app.listen(8091, '0.0.0.0');
