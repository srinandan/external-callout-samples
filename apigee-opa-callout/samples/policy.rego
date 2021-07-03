package httpapi.authz

import input

default allow = false

developers = {
    "apps@samples.com": {
        "/opa/items": ["GET", "POST"]
    },
    "test@samples.com": {
        "/foo": ["GET"]
    }
}

allow {
    developers[input.developer_email][input.path][_] = input.method
}