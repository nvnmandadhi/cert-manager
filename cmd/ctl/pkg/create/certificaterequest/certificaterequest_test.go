/*
Copyright 2020 The Jetstack cert-manager contributors.

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

package certificaterequest

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := map[string]struct {
		inputFile    string
		inputArgs    []string
		keyFilename  string
		certFilename string
		fetchCert    bool

		expErr    bool
		expErrMsg string
	}{
		"CR name not passed as arg throws error": {
			inputFile: "example.yaml",
			inputArgs: []string{},
			expErr:    true,
			expErrMsg: "the name of the CertificateRequest to be created has to be provided as argument",
		},
		"More than one arg throws error": {
			inputFile: "example.yaml",
			inputArgs: []string{"hello", "World"},
			expErr:    true,
			expErrMsg: "only one argument can be passed in: the name of the CertificateRequest",
		},
		"not specifying path to yaml manifest throws error": {
			inputFile: "",
			inputArgs: []string{"hello"},
			expErr:    true,
			expErrMsg: "the path to a YAML manifest of a Certificate resource cannot be empty, please specify by using --from-certificate-file flag",
		},
		"key filename and cert filename are optional flags": {
			inputFile:    "example.yaml",
			inputArgs:    []string{"hello"},
			keyFilename:  "",
			certFilename: "",
			expErr:       false,
		},
		"identical key filename and cert filename throws error": {
			inputFile:    "example.yaml",
			inputArgs:    []string{"hello"},
			keyFilename:  "same",
			certFilename: "same",
			expErr:       true,
			expErrMsg:    "the file to store private key cannot be the same as the file to store certificate",
		},
		"cannot specify cert filename without fetch-certificate flag": {
			inputFile:    "example.yaml",
			inputArgs:    []string{"hello"},
			certFilename: "cert.crt",
			fetchCert:    false,
			expErr:       true,
			expErrMsg:    "cannot specify file to store certificate if not waiting for and fetching certificate, please set --fetch-certificate flag",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			opts := &Options{
				InputFilename: test.inputFile,
				KeyFilename:   test.keyFilename,
				CertFileName:  test.certFilename,
			}

			// Validating args and flags
			err := opts.Validate(test.inputArgs)
			if err != nil {
				if !test.expErr {
					t.Errorf("got unexpected error when validating args and flags: %v", err)
				} else if err.Error() != test.expErrMsg {
					t.Errorf("got unexpected error when validating args and flags, expected: %v; actual: %v", test.expErrMsg, err)
				}
				return
			} else {
				// got no error
				if test.expErr {
					t.Errorf("expected but got no error validating args and flags")
				}
			}
		})
	}
}
