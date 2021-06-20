package httpapi.authz

import input

default allow = false

developers = {
    "apps@sample.com": {
        "/opa/items": {"GET"},
    },
}

allow {
    developers[input.developer_email][input.path][_] = input.method
}