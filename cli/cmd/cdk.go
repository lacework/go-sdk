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
	"fmt"
	"net"
	"os"

	cdk "github.com/lacework/go-sdk/cli/cdk/go/proto/v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// default gRPC target if not specified via LW_CDK_TARGET
const defaultGrpcTarget string = "localhost:1123"

// envs are all the environment variables passed to CDK components
func (c *cliState) envs() []string {
	return []string{
		fmt.Sprintf("LW_ACCOUNT=%s", c.Account),
		fmt.Sprintf("LW_SUBACCOUNT=%s", c.Subaccount),
		fmt.Sprintf("LW_API_KEY=%s", c.KeyID),
		fmt.Sprintf("LW_API_SECRET=%s", c.Secret),
		fmt.Sprintf("LW_API_TOKEN=%s", c.Token),
		fmt.Sprintf("LW_ORGANIZATION=%v", c.OrgLevel),
		fmt.Sprintf("LW_NONINTERACTIVE=%v", c.nonInteractive),
		fmt.Sprintf("LW_NOCACHE=%v", c.noCache),
		fmt.Sprintf("LW_NOCOLOR=%s", os.Getenv("NO_COLOR")),
		fmt.Sprintf("LW_LOG=%s", c.LogLevel),
		fmt.Sprintf("LW_JSON=%v", c.jsonOutput),
		fmt.Sprintf("LW_CDK_TARGET=%s", c.GrpcTarget()),
	}
}

// GrpcTarget returns the gRPC target that the CDK architecture will use
// to allow components to communicate back to the Lacework CLI
func (c *cliState) GrpcTarget() string {
	if target := os.Getenv("LW_CDK_TARGET"); target != "" {
		return target
	}
	return defaultGrpcTarget
}

// Ping implements CDK.Ping
func (c *cliState) Ping(ctx context.Context, in *cdk.PingRequest) (*cdk.PongReply, error) {
	c.Log.Debugw("message", "from", "CDK/Ping", "component_name", in.GetComponentName())
	return &cdk.PongReply{Message: fmt.Sprintf("Pong %s", in.GetComponentName())}, nil
}

// Honeyvent implements CDK.Honeyvent
func (c *cliState) Honeyvent(ctx context.Context, in *cdk.HoneyventRequest) (*cdk.Reply, error) {
	c.Log.Debugw("message", "from", "CDK/Honeyvent", "feature", in.GetFeature())

	// Set event feature, if provided
	if f := in.GetFeature(); f != "" {
		c.Event.Feature = f
	}

	// Add feature fields
	for key, value := range in.GetFeatureData() {
		c.Event.AddFeatureField(key, value)
	}

	// Set any error, if any
	if err := in.GetError(); err != "" {
		c.Event.Error = err
	}

	// Set duration in millisecond, if provided
	if durationMs := in.GetDurationMs(); durationMs != 0 {
		c.Event.DurationMs = durationMs
	}

	// Send the Honeyvent
	c.SendHoneyvent()

	return &cdk.Reply{}, nil
}

// Serve will start the CDK gRPC server at the provided port and log level
func (c *cliState) Serve(target string) error {
	c.Stop() // make sure server is not running

	lis, err := net.Listen("tcp", target)
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}

	c.cdkServer = grpc.NewServer() // guardrails-disable-line
	cdk.RegisterCoreServer(c.cdkServer, c)

	c.Log.Infow("gRPC server started", "address", lis.Addr())
	if err := c.cdkServer.Serve(lis); err != nil {
		return errors.Wrap(err, "failed to serve")
	}

	return nil
}

// Stop will stop the CDK gRPC server gracefully. It stops the server from
// accepting new connections and RPCs and blocks until all the pending RPCs
// are finished.
func (c *cliState) Stop() {
	if c.cdkServer != nil {
		c.Log.Info("stopping gRPC server")
		c.cdkServer.GracefulStop()
	}
}
