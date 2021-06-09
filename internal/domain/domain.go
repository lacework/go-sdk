//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2021, Lacework Inc.
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

package domain

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Use this package to disseminate a domain URL into
// account, cluster and whether or not it is internal

type domain struct {
	Account  string
	Cluster  string
	Internal bool
}

// New returns domain information from the provided URL
//
// For instance, the following URL:
// ```
// d := domain.New("https://account.fra.lacework.net")
// ```
// Would be dessiminated into:
// * `account` as the account name
// * `fra` as the cluster name
func New(url string) (domain, error) {
	// for the full url https://ACCOUNT.lacework.net
	// remove the prefixes https:// or http://
	rx, err := regexp.Compile(`(http://|https://)`)
	if err == nil {
		url = rx.ReplaceAllString(url, "")
	}

	// for the full domain ACCOUNT[.CUSTER][.corp].lacework.net
	// subtract the account and cluster name, also detect if the
	// domain is internal (dev, qa, preprod, etc)
	rx, err = regexp.Compile(`\.lacework\.net.*`)
	if err == nil {
		domainSplit := rx.Split(url, -1)
		if len(domainSplit) > 1 {
			url = domainSplit[0]
		} else {
			return domain{}, errors.New("domain not supported")
		}
	}

	domainInfo := strings.Split(url, ".")
	switch len(domainInfo) {
	case 1:
		return domain{Account: domainInfo[0]}, nil
	case 2:
		return domain{
			Account: domainInfo[0],
			Cluster: domainInfo[1],
		}, nil
	case 3:
		if domainInfo[2] != "corp" {
			return domain{}, errors.New("unable to detect if domain is internal")
		}
		return domain{
			Account:  domainInfo[0],
			Cluster:  domainInfo[1],
			Internal: true,
		}, nil
	default:
		return domain{}, errors.New("unable to detect domain information")
	}
}

func (d *domain) String() string {
	if d.Internal {
		return fmt.Sprintf("%s.%s.corp", d.Account, d.Cluster)
	}

	if d.Cluster != "" {
		return fmt.Sprintf("%s.%s", d.Account, d.Cluster)
	}

	return d.Account
}
