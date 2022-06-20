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
    "current":   "http://backend-current:8091",
    "candidate": "http://backend-candidate:8091",
}

// implement /getRecommendation endpoint
// calls backend service /recommend endpoint
app.get('/getRecommendation', (req, res) => {
	// Get user (session) identifier, for example by inspection of header X-User
    const user = req.header('X-User')

    // Get endpoint of backend endpoint "/recommend"
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
 
        // call backend service using url
        http.get(route + '/recommend', (resp) => {
            let str = '';
            resp.on('data', function(chunk) {
                str += chunk;
            });
            resp.on('end', function () {
                // write response to query
                res.send(`Recommendation: ${str}`);
            });
        }).on("error", (err) => {
            console.log("Error: " + err.message)
        });
    });
});

// implement /buy endpoint
// writes value for sample_metric which may have spanned several calls to /getRecommendation
app.get('/buy', (req, res) => {
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

app.listen(8090, '0.0.0.0');
