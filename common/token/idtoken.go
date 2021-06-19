// Copyright 2020 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package token

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/lestrrat/go-jwx/jwa"
	"github.com/lestrrat/go-jwx/jwt"
	"github.com/srinandan/sample-apps/common"
	"golang.org/x/oauth2/google"

	oauthsvc "google.golang.org/api/oauth2/v2"
)

const audience = "https://www.googleapis.com/oauth2/v4/token"

// NewTokenSource creates a TokenSource that returns ID tokens with the audience
// provided and configured with the supplied options. The parameter audience may
// not be empty.
func NewTokenSource(ctx context.Context, targetAudience string) (string, error) {
	data, err := readServiceAccount()
	if err != nil {
		return "", err
	}
	return tokenSourceFromBytes(ctx, getHostname(targetAudience), data)
}

func readServiceAccount() (content []byte, err error) {
	content, err = ioutil.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		common.Error.Println("service account was not found")
		return nil, err
	}
	return content, nil
}

func getHostname(extCalloutServiceEndpoint string) string {
	if strings.Contains(extCalloutServiceEndpoint, ":") {
		names := strings.Split(extCalloutServiceEndpoint, ":")
		return names[0]
	} else {
		return extCalloutServiceEndpoint
	}
}

func doExchange(token string) (payload []byte, err error) {
	d := url.Values{}
	d.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	d.Add("assertion", token)

	client := &http.Client{}
	req, err := http.NewRequest("POST", audience, strings.NewReader(d.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func tokenSourceFromBytes(ctx context.Context, targetAudience string, data []byte) (identityToken string, err error) {

	conf, err := google.JWTConfigFromJSON(data, oauthsvc.UserinfoEmailScope)
	if err != nil {
		return identityToken, err
	}

	iat := time.Now()
	exp := iat.Add(time.Hour)

	token := jwt.New()
	_ = token.Set(jwt.AudienceKey, audience)
	_ = token.Set(jwt.IssuerKey, conf.Email)
	_ = token.Set(jwt.IssuedAtKey, iat.Unix())
	_ = token.Set(jwt.ExpirationKey, exp.Unix())
	_ = token.Set(`target_audience`, "https://"+targetAudience)

	// from https://github.com/golang/oauth2/blob/master/internal/oauth2.go#L23
	key := conf.PrivateKey
	block, _ := pem.Decode(key)
	if block != nil {
		key = block.Bytes
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		parsedKey, err = x509.ParsePKCS1PrivateKey(key)
		if err != nil {
			return identityToken, err
		}
	}
	parsed, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return identityToken, err
	}

	payload, err := token.Sign(jwa.RS256, parsed)
	if err != nil {
		return identityToken, err
	}

	body, err := doExchange(string(payload))
	if err != nil {
		return identityToken, err
	}

	m := make(map[string]string)
	if err = json.Unmarshal(body, &m); err != nil {
		return identityToken, err
	}

	if m["id_token"] != "" {
		return m["id_token"], nil
	} else {
		return identityToken, fmt.Errorf("unable to generate id_token %s\n", m)
	}
}
