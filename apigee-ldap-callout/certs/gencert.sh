#ensure there is a blank index.txt file
#echo '01' > ca.srl


# CA
openssl genrsa -des3 -out ca.key
openssl req -x509 -new -days 750 -key ca.key -sha256 -out ca.crt

openssl genrsa -out cert.key 2048
openssl req -new -key cert.key -out cert.csr -config openssl_cert.cnf -extensions v3_req

#sign the cert
openssl ca -config ca.cnf -cert ca.crt -out cert.crt -extfile v3.ext -in cert.csr

#gen dhparam
#openssl dhparam -out dhparam.pem 2048

# print cert
# openssl x509 -in ca.crt -text -noout

#verify cert
# openssl verify -CAfile rootCA.crt -CApath /home/srinandans/workspace/tls -no-CAfile -no-CApath ca.crt 
