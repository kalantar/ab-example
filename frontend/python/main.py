import random

from http import HTTPStatus
import requests

import grpc
import abn_pb2
import abn_pb2_grpc

from flask import Flask, request

# map of track to route to backend service
trackToRoute = {
    "current": "http://backend-current:8090",
	"candidate": "http://backend-candidate:8090"
}

app = Flask(__name__)

# implement /hello endpoint
# calls backend service /world endpoint
@app.route('/hello')
def hello():
    # Get user (session) identifier, for example, by inspection of header X-User
    if not ('X-User' in request.headers):
        return "header X-User missing", HTTPStatus.INTERNAL_SERVER_ERROR
    user = request.headers['X-User']

    # Get endpoint of backend endpoint "/world"
    # In this example, the backend endpoint depends on the version (track) of the backend service
    # the user is assigned by the Iter8 SDK Lookup() method

    # establish connection to ABn service
    with grpc.insecure_channel("abn:50051") as channel:
        stub = abn_pb2_grpc.ABNStub(channel)

        # call ABn service API Lookup() to get an assigned track for the user
        s = stub.Lookup( \
            abn_pb2.Application(name="default/backend", \
            user=user) \
        )

        # lookkup route using track
        if not (s.track in trackToRoute):
            return "unknwon track returned: {0}".format(s.track), HTTPStatus.INTERNAL_SERVER_ERROR
        route = trackToRoute[s.track]

        # call backend service using url
        try:
            r = requests.get(url=route + "/world", allow_redirects=True)
            r.raise_for_status()
            world = r.text
        except Exception as e:
            return "call to backend endpoint /world failed: {0}".format(e), HTTPStatus.INTERNAL_SERVER_ERROR

        return "hello world {0}".format(world)
    
# implement /goodbye endpoint
# writes value for sample_metrc which may have spanned several calls to /hello
@app.route('/goodbye')
def goodbye():
    # Get user (session) identifier, for example, by inspection of header X-User
    if not ('X-User' in request.headers):
        return "header X-User missing", HTTPStatus.INTERNAL_SERVER_ERROR
    user = request.headers['X-User']

	# export metric to metrics database
	# this is best effort; we ignore any failure

    # establish connection to ABn service
    with grpc.insecure_channel("abn:50051") as channel:
        stub = abn_pb2_grpc.ABNStub(channel)

        # export metric to metrics database
        # this is best effort; we ignore any failure
        stub.WriteMetric( \
            abn_pb2.MetricValue(name="sample_metric", \
            value=str(random.randint(0,100)), \
            application="default/backend", \
            user=user) \
        )

    return ""