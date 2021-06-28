# apigee-ldap-callout

This repo shows a sample of how to use the [External Callout policy](https://cloud.google.com/apigee/docs/api-platform/reference/policies/external-callout-policy) with [LDAP](https://github.com/go-ldap/ldap) for Apigee X/hybrid.

## Description

The external callout service implements two actions:

1. Authenticate username and password

Accepts a username/password combination, tests the credentials against the configured LDAP and a flow variable `ldap-result` with the bool `true` or `false`. The message template is:

```json
{"ldapAuthRequest": {"username":"billy","password":"test"}}
```

Optional parameters include: `scope` (oneOf object, onelevel or subtree), `baseDN` (default = dc=example,dc=org) and `objectClass` (default = inetOrgPerson).

2. Search entity

Accepts a query and attributes, searches against the configured LDAP and returns the entities as well as a flow variable `ldap-search-result` with the count of entities. The message template is:

```json
{"ldapSearchRequest": {"query":"uid=billy","attributes":["dn"]}}
```

Optional parameters include `scope` (oneOf object, onelevel or subtree) and `baseDN` (default = dc=example,dc=org).

## Testing Locally

OpenLDAP can be started locally on docker. 

```sh
docker build -t openldap -f Dockerfile.ldap .
docker run --rm -p 389:389 -p 636:636 openldap
```

Start the External callout service

```sh
go run cmd/server/main.go --host=ldap://localhost:389 --binduser="cn=readonly,dc=example,dc=org" --bindpasswd=passwr0rd!
```

Start a REST External callout client

```sh
go run cmd/client/main.go
```

Test the configuration (Successful Message)

```sh
curl localhost:8080/processmessage -H "Content-Type: application/json" -d '{"ldapAuthRequest": {"username":"billy","password":"test"}}'

{"response":{},"additionalFlowVariables":{"ldap-result":{"bool":true}}}
```

Test the configuration (Invalid Message)

```sh
curl localhost:8080/processmessage -H "Content-Type: application/json" -d '{"ldapAuthRequest": {"username":"billy","password":"test123"}}'

{"response":{"content":"{\"status_code\":500,\"message\":\"LDAP Result Code 49 \\\"Invalid Credentials\\\": \"}"},"additionalFlowVariables":{"ldap-result":{"bool":false}}}
```

Test Search (Successful Message)

```sh
curl localhost:8080/processmessage -H "Content-Type: application/json" -d '{"ldapSearchRequest": {"query":"uid=billy","attributes":["dn"]}}'

{"response":{"content":"{\"Entries\":[{\"DN\":\"uid=billy,dc=example,dc=org\",\"Attributes\":null}],\"Referrals\":[],\"Controls\":[]}"},"additionalFlowVariables":{"ldap-search-count":{"int32":1}}}
```

Test Search (Unsuccessful Search)

```sh
curl localhost:8080/processmessage -H "Content-Type: application/json" -d '{"ldapSearchRequest": {"query":"uid=bill","attributes":["dn"]}}'
{"response":{"content":"{\"Entries\":[],\"Referrals\":[],\"Controls\":[]}"},"additionalFlowVariables":{"ldap-search-count":{"int32":0}}}
```

## Setup

1. Build and deploy the External callout service

```sh
skaffold run -p gcb --default-repo=gcr.io/xxx -f skaffold-apigee-ldap.yaml
```

2. Deploy the API Proxy included [here](./apiproxy)

3. Create the target server

```sh
export ORG=XXXX
export ENV=XXXX
export TOKEN=$(gcloud auth print-access-token)
export APIGEE_LDAP_HOSTNAME=XXXXX

curl "https://apigee.googleapis.com/v1/organizations/$ORG/environments/$ENV/targetservers" -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" --data-raw '{
  "name": "grpcserver",
  "host": $APIGEE_LDAP_HOSTNAME,
  "port": 50051,
  "isEnabled": true,
  "protocol": "GRPC"
}'
```

4. To test the setup (invalid creds),

```sh
curl https://api.acme.com/ldap -u test:test -v

...
< HTTP/2 403
< content-length: 115
< date: Thu, 24 Jun 2021 19:28:31 GMT
< server: istio-envoy
<

{
    "error":"authentication_error",
    "error_description": "username or password did not match"
}
```

5. To test the setup (valid creds),

```sh
curl https://api.acme.com/ldap -u billy:test -v

...
< via: 1.1 google
< alt-svc: clear
< server: istio-envoy
<

Hello, Guest!
```

___

## Support

This is not an officially supported Google product