/*
 Copyright 2020 The Kubernetes Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYAMLValidateTokenConfig(t *testing.T) {
	testCases := []struct {
		name               string
		cfg                *NsxtYAML
		expectedErrMessage string
	}{
		{
			name:               "empty file",
			cfg:                &NsxtYAML{},
			expectedErrMessage: "user or vmc access token or client cert file or secret for user must be set",
		},
		{
			name: "empty auth host",
			cfg: &NsxtYAML{
				VMCAccessToken: "token",
			},
			expectedErrMessage: "vmc auth host is required if auth token is provided",
		},
		{
			name: "empty auth token",
			cfg: &NsxtYAML{
				VMCAuthHost: "host",
			},
			expectedErrMessage: "vmc auth token is required if auth host is provided",
		},
		{
			name: "empty host",
			cfg: &NsxtYAML{
				VMCAccessToken: "token",
				VMCAuthHost:    "host",
			},
			expectedErrMessage: "host is empty",
		},
		{
			name: "valid config",
			cfg: &NsxtYAML{
				VMCAccessToken: "token",
				VMCAuthHost:    "host",
				Host:           "server",
			},
		},
	}

	for _, testCase := range testCases {
		err := testCase.cfg.validateConfig()
		if err != nil {
			assert.EqualError(t, err, testCase.expectedErrMessage)
		} else {
			assert.Equal(t, "", testCase.expectedErrMessage)
		}
	}
}

func TestYAMLValidateUserConfig(t *testing.T) {
	cfg := &NsxtYAML{
		User: "admin",
	}
	err := cfg.validateConfig()
	assert.EqualError(t, err, "password is empty")

	cfg.Password = "secret"
	err = cfg.validateConfig()
	assert.EqualError(t, err, "host is empty")

	cfg.Host = "server"
	err = cfg.validateConfig()
	assert.Nil(t, err)
}

func TestYAMLValidateCertConfig(t *testing.T) {
	testCases := []struct {
		name               string
		cfg                *NsxtYAML
		expectedErrMessage string
	}{
		{
			name: "empty client cert file",
			cfg: &NsxtYAML{
				ClientAuthKeyFile: "client-key",
			},
			expectedErrMessage: "client cert file is required if client key file is provided",
		},
		{
			name: "empty client key file",
			cfg: &NsxtYAML{
				ClientAuthCertFile: "client-cert",
			},
			expectedErrMessage: "client key file is required if client cert file is provided",
		},
		{
			name: "empty host",
			cfg: &NsxtYAML{
				ClientAuthKeyFile:  "client-key",
				ClientAuthCertFile: "client-cert",
			},
			expectedErrMessage: "host is empty",
		},
		{
			name: "valid config",
			cfg: &NsxtYAML{
				ClientAuthKeyFile:  "client-key",
				ClientAuthCertFile: "client-cert",
				Host:               "server",
			},
			expectedErrMessage: "",
		},
	}

	for _, testCase := range testCases {
		err := testCase.cfg.validateConfig()
		if err != nil {
			assert.EqualError(t, err, testCase.expectedErrMessage)
		} else {
			assert.Equal(t, "", testCase.expectedErrMessage)
		}
	}
}

func TestYAMLValidateSecretConfig(t *testing.T) {
	testCases := []struct {
		name               string
		cfg                *NsxtYAML
		expectedErrMessage string
	}{
		{
			name: "empty secret namespace",
			cfg: &NsxtYAML{
				SecretName: "secret-name",
			},
			expectedErrMessage: "secret namespace is required if secret name is provided",
		},
		{
			name: "empty secret name",
			cfg: &NsxtYAML{
				SecretNamespace: "secret-ns",
			},
			expectedErrMessage: "secret name is required if secret namespace is provided",
		},
		{
			name: "empty host",
			cfg: &NsxtYAML{
				SecretName:      "secret-name",
				SecretNamespace: "secret-ns",
			},
			expectedErrMessage: "host is empty",
		},
		{
			name: "valid config",
			cfg: &NsxtYAML{
				SecretName:      "secret-name",
				SecretNamespace: "secret-ns",
				Host:            "server",
			},
			expectedErrMessage: "",
		},
	}

	for _, testCase := range testCases {
		err := testCase.cfg.validateConfig()
		if err != nil {
			assert.EqualError(t, err, testCase.expectedErrMessage)
		} else {
			assert.Equal(t, "", testCase.expectedErrMessage)
		}
	}
}

func TestReadYAMLConfig(t *testing.T) {
	contents := `
nsxt:
  user: admin
  password: secret
  host: nsxt-server
  insecureFlag: false
  remoteAuth: true
  vmcAccessToken: vmc-token
  vmcAuthHost: vmc-host
  clientAuthCertFile: client-cert-file
  clientAuthKeyFile: client-key-file
  caFile: ca-file
  secretName: secret-name
  secretNamespace: secret-ns
`
	config, err := ReadRawConfigYAML([]byte(contents))
	if err != nil {
		t.Error(err)
		return
	}

	assertEquals := func(name, left, right string) {
		if left != right {
			t.Errorf("%s %s != %s", name, left, right)
		}
	}
	assertEquals("NSXT.user", config.NSXT.User, "admin")
	assertEquals("NSXT.password", config.NSXT.Password, "secret")
	assertEquals("NSXT.host", config.NSXT.Host, "nsxt-server")
	assert.Equal(t, false, config.NSXT.InsecureFlag)
	assert.Equal(t, true, config.NSXT.RemoteAuth)
	assertEquals("NSXT.vmcAccessToken", config.NSXT.VMCAccessToken, "vmc-token")
	assertEquals("NSXT.vmcAuthHost", config.NSXT.VMCAuthHost, "vmc-host")
	assertEquals("NSXT.clientAuthCertFile", config.NSXT.ClientAuthCertFile, "client-cert-file")
	assertEquals("NSXT.clientAuthKeyFile", config.NSXT.ClientAuthKeyFile, "client-key-file")
	assertEquals("NSXT.caFile", config.NSXT.CAFile, "ca-file")
	assertEquals("NSXT.secretName", config.NSXT.SecretName, "secret-name")
	assertEquals("NSXT.secretNamespace", config.NSXT.SecretNamespace, "secret-ns")
}
