//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
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

package cmd

import (
	"context"
	"testing"
	"time"

	cdk "github.com/lacework/go-sdk/cli/cdk/go/proto/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestCDKServer(t *testing.T) {
	go cli.Serve(defaultGrpcTarget)
	defer cli.Stop()

	// wait for the gRPC server to come online
	time.Sleep(time.Millisecond * 200)

	conn, err := grpc.Dial(defaultGrpcTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if assert.Nil(t, err) {
		defer conn.Close()
	}

	var (
		cdkClient   = cdk.NewCoreClient(conn)
		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	)
	defer cancel()

	// Ping the CDK Server
	reply, err := cdkClient.Ping(ctx, &cdk.PingRequest{
		ComponentName: "cli-test",
	})
	if assert.Nil(t, err) {
		assert.Equalf(t, "Pong cli-test", reply.GetMessage(),
			"Expected a Ping -> Pong")
	} else {
		assert.Equal(t, "this should be empty", err.Error())
	}
}
