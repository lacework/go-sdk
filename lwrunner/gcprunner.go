//
// Author:: Nicholas Schmeller (<nick.schmeller@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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

package lwrunner

import (
	"context"
	"fmt"
	"time"

	oslogin "cloud.google.com/go/oslogin/apiv1"
	osloginpb "cloud.google.com/go/oslogin/apiv1/osloginpb"
	"golang.org/x/crypto/ssh"
	osloginpb_common "google.golang.org/genproto/googleapis/cloud/oslogin/common"
)

type GCPRunner struct {
	Runner           Runner
	ParentUsername   string
	ProjectID        string
	AvailabilityZone string
	InstanceID       string
}

func NewGCPRunner(host, parentUsername, projectID, availabilityZone, instanceID string, callback ssh.HostKeyCallback) (*GCPRunner, error) {
	defaultCallback, err := DefaultKnownHosts()
	if err == nil && callback == nil {
		callback = defaultCallback
	}

	runner := New("", host, callback) // populate username during `SendAndUseIdentityFile()`

	return &GCPRunner{
		*runner,
		parentUsername,
		projectID,
		availabilityZone,
		instanceID,
	}, nil
}

func (run GCPRunner) SendAndUseIdentityFile() error {
	pubBytes, privBytes, err := GetKeyBytes()
	if err != nil {
		return err
	}

	err = run.SendPublicKey(pubBytes)
	if err != nil {
		return err
	}

	signer, err := ssh.ParsePrivateKey(privBytes)
	if err != nil {
		return err
	}
	run.Runner.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}

	return nil
}

// SendPublicKey is a helper function to send a public key to a GCP account
// for OSLogin authentication. The account must have the "Compute OS Login" IAM role
// and "Service Account User" authorization for the GCE default service account.
// When the SSH key is sent, it will persist in the GCP account for 10min.
func (run GCPRunner) SendPublicKey(pubBytes []byte) error {
	ctx := context.Background()
	c, err := oslogin.NewClient(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	key := &osloginpb_common.SshPublicKey{
		Key:                string(pubBytes),
		ExpirationTimeUsec: time.Now().UnixMicro() + (10 * time.Minute.Microseconds()), // expiration time is 10min from now
	}

	req := &osloginpb.ImportSshPublicKeyRequest{
		Parent:       "users/oslogin-account@lw-agent-team.iam.gserviceaccount.com",
		SshPublicKey: key,
		ProjectId:    run.ProjectID,
	}
	resp, err := c.ImportSshPublicKey(ctx, req)
	if err != nil {
		return err
	}
	fmt.Printf("ImportSshPublicKey resp: %v\n", *resp)

	// Get login info from the OSLogin profile for our SSH login
	posixAccounts := resp.LoginProfile.GetPosixAccounts()
	for _, account := range posixAccounts {
		if account.Primary { // there will only be one Primary account per profile
			run.Runner.User = account.Username
		}
	}

	return nil
}
