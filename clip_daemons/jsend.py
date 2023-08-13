import json

def New(data):
    resp = {}
    resp["status"] = "success"
    resp["data"] = data
    return json.dumps(resp)

def NewError(message: str):
    resp = {}
    resp["status"] = "error"
    resp["message"] = message
    return json.dumps(resp)

def NewFail(data):
    resp = {}
    resp["status"] = "fail"
    resp["data"] = data
    return json.dumps(resp)