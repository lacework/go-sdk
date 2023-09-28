//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
// License:: Apache License, Version 2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// A Lacework component package to help facilitate the loading and execution of components
package lwcomponent

import (
	"testing"

	"aead.dev/minisign"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type verifySignatureTest struct {
	Name      string
	Message   []byte
	Signature []byte
	Error     error
}

var verifySignatureTests = []verifySignatureTest{
	{
		Name:  "UnableToParseSignature",
		Error: errors.New("unable to parse signature"),
	},
	{
		Name: "InvalidTrustedComment",
		Signature: []byte(`untrusted comment: 
RWQnyHnTKz1RiRl+vm7dNpp7Jrt+qLANsu9m+TaHSgVyJO9Ldeo5ZPw/GQJu1fJZBuCrsMVW5KbuB6dEPzQvi8ct1zJfMkgVsgs=
trusted comment: RWQnyHnTKz1RidqMl45Ou5XbJCroV7zohE1jpbigTUGXyQY94AF5uo4v.1.dW50cnVzdGVkIGNvbW1lbnQ6IApSV1RjbVljdjlQMHlNcExyVnA3V3lIdmx2RU02VlhJbFhzUGE0QVZVL1MraHlwUCtQdHBFWEVvck5MS3ZReTJSV0QxTmNKNzNlU1NMYUE2TnpLNnBvczlpb2RhbENMczhOZzg9CnRydXN0ZWQgY29tbWVudDogCm1iMmpJRVZBUTNPZlIvL0pRcjUxVzdQR3FrWDhGbWM1YUFKcCt0Q3JXOHZmWmZ6YnVMM2ZvejdtQmV6MHpQdlBldmQxU0FrUklmSW1EMkdjSkVlUEJ3PT0=
cJwh7XDctGR8dpc1NEJnPGU5IXNvMNXa6hdMA0BBZvvnsKQHeVVEkMc7zo72iFsHEKRxe+AmXkryNgVi7Gg0DQ==`),
		Error: errors.New("invalid signature trusted comment"),
	},
	{
		Name: "UnableToParseRootSignature",
		Signature: []byte(`untrusted comment: 
RWQnyHnTKz1RiRl+vm7dNpp7Jrt+qLANsu9m+TaHSgVyJO9Ldeo5ZPw/GQJu1fJZBuCrsMVW5KbuB6dEPzQvi8ct1zJfMkgVsgs=
trusted comment: 1.RWQnyHnTz1RidqMl45Ou5XbJCroV7zohE1jpbigTUGXyQY94AF5uo4v.1.W50cnVzdGVkIGNvbW1lbnQ6IApSV1RjbVljdjlQMHlNcExyVnA3V3lIdmx2RU02VlhJbFhzUGE0QVZVL1MraHlwUCtQdHBFWEVvck5MS3ZReTJSV0QxTmNKNzNlU1NMYUE2TnpLNnBvczlpb2RhbENMczhOZzg9CnRydXN0ZWQgY29tbWVudDogCm1iMmpJRVZBUTNPZlIvL0pRcjUxVzdQR3FrWDhGbWM1YUFKcCt0Q3JXOHZmWmZ6YnVMM2ZvejdtQmV6MHpQdlBldmQxU0FrUklmSW1EMkdjSkVlUEJ3PT0=
cJwh7XDctGR8dpc1NEJnPGU5IXNvMNXa6hdMA0BBZvvnsKQHeVVEkMc7zo72iFsHEKRxe+AmXkryNgVi7Gg0DQ==`),
		Error: errors.New("unable to parse root signature from trusted comment: illegal base64 data at input byte 303"),
	},
	{
		Name: "InvalidRootSignatureOverSigningKey",
		Signature: []byte(`untrusted comment: 
RWQnyHnTKz1RiRl+vm7dNpp7Jrt+qLANsu9m+TaHSgVyJO9Ldeo5ZPw/GQJu1fJZBuCrsMVW5KbuB6dEPzQvi8ct1zJfMkgVsgs=
trusted comment: 1.RWQnyHnTz1RidqMl45Ou5XbJCroV7zohE1jpbigTUGXyQY94AF5uo4v.1.dW50cnVzdGVkIGNvbW1lbnQ6IApSV1RjbVljdjlQMHlNcExyVnA3V3lIdmx2RU02VlhJbFhzUGE0QVZVL1MraHlwUCtQdHBFWEVvck5MS3ZReTJSV0QxTmNKNzNlU1NMYUE2TnpLNnBvczlpb2RhbENMczhOZzg9CnRydXN0ZWQgY29tbWVudDogCm1iMmpJRVZBUTNPZlIvL0pRcjUxVzdQR3FrWDhGbWM1YUFKcCt0Q3JXOHZmWmZ6YnVMM2ZvejdtQmV6MHpQdlBldmQxU0FrUklmSW1EMkdjSkVlUEJ3PT0=
cJwh7XDctGR8dpc1NEJnPGU5IXNvMNXa6hdMA0BBZvvnsKQHeVVEkMc7zo72iFsHEKRxe+AmXkryNgVi7Gg0DQ==`),
		Error: errors.New("invalid root signature over signing key"),
	},
	{
		Name: "InvalidSignatureOverComponent",
		Signature: []byte(`untrusted comment: 
RWQnyHnTKz1RiRl+vm7dNpp7Jrt+qLANsu9m+TaHSgVyJO9Ldeo5ZPw/GQJu1fJZBuCrsMVW5KbuB6dEPzQvi8ct1zJfMkgVsgs=
trusted comment: 1.RWQnyHnTKz1RidqMl45Ou5XbJCroV7zohE1jpbigTUGXyQY94AF5uo4v.1.dW50cnVzdGVkIGNvbW1lbnQ6IApSV1RjbVljdjlQMHlNcExyVnA3V3lIdmx2RU02VlhJbFhzUGE0QVZVL1MraHlwUCtQdHBFWEVvck5MS3ZReTJSV0QxTmNKNzNlU1NMYUE2TnpLNnBvczlpb2RhbENMczhOZzg9CnRydXN0ZWQgY29tbWVudDogCm1iMmpJRVZBUTNPZlIvL0pRcjUxVzdQR3FrWDhGbWM1YUFKcCt0Q3JXOHZmWmZ6YnVMM2ZvejdtQmV6MHpQdlBldmQxU0FrUklmSW1EMkdjSkVlUEJ3PT0=
cJwh7XDctGR8dpc1NEJnPGU5IXNvMNXa6hdMA0BBZvvnsKQHeVVEkMc7zo72iFsHEKRxe+AmXkryNgVi7Gg0DQ==`),
		Error: errors.New("invalid signature over component"),
	},
	{
		Name:      "OK",
		Message:   helloWorld,
		Signature: helloWorldSigDecoded,
		Error:     nil,
	},
}

func TestVerifySignature(t *testing.T) {
	for _, vst := range verifySignatureTests {
		t.Run(vst.Name, func(t *testing.T) {
			rootPublicKey := minisign.PublicKey{}
			if err := rootPublicKey.UnmarshalText([]byte(publicKey)); err != nil {
				assert.FailNow(t, err.Error())
			}

			actualError := verifySignature(rootPublicKey, vst.Message, vst.Signature)
			if vst.Error != nil {
				assert.Equal(t, vst.Error.Error(), actualError.Error())
			} else {
				assert.Nil(t, actualError)
			}
		})
	}
}
