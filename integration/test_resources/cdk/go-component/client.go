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

package main

import (
	"context"
	"os"
	"time"

	cdk "github.com/lacework/go-sdk/cli/cdk/go/proto/v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Connect to the CDK server
func Connect() (cdk.CoreClient, *grpc.ClientConn, error) {
	// Set up a connection to the CDK server
	log.Infow("connecting to gRPC server", "address", os.Getenv("LW_CDK_TARGET"))
	conn, err := grpc.Dial(os.Getenv("LW_CDK_TARGET"),
		// @afiune we do an insecure connection since we are
		// connecting to the server running on the same machine
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not connect")
	}

	var (
		cdkClient   = cdk.NewCoreClient(conn)
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	)
	defer cancel()

	// Ping the CDK Server
	reply, err := cdkClient.Ping(ctx, &cdk.PingRequest{
		ComponentName: os.Getenv("LW_COMPONENT_NAME"),
	})
	if err != nil {
		return cdkClient, conn, errors.Wrap(err, "could not ping")
	}
	log.Debugw("response", "from", "cdk.v1.Core/Ping", "message", reply.GetMessage())

	return cdkClient, conn, nil
}
