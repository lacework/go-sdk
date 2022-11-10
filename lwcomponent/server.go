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

package lwcomponent

import (
	"context"
	"fmt"
	"net"

	"github.com/lacework/go-sdk/lwcomponent/cdk"
	"github.com/lacework/go-sdk/lwlogger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Our CDK Server that implements all proto services defined at 'lwcomponent/cdk'
type server struct {
	Log *zap.SugaredLogger

	cdk.UnimplementedStatusServer
}

// Ping implements cdk.StatusServer
func (s *server) Ping(ctx context.Context, in *cdk.PingRequest) (*cdk.PingReply, error) {
	s.Log.Debugw("message", "from", "cdk.Status/Ping", "component_name", in.GetComponentName())
	return &cdk.PingReply{Message: fmt.Sprintf("Pong %s", in.GetComponentName())}, nil
}

// Serve will start the CDK gRPC Server at the provided port and log level
func Serve(port int, logLevel string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}

	var (
		grpcServer = grpc.NewServer()
		log        = lwlogger.New(logLevel).Sugar()
	)
	cdk.RegisterStatusServer(grpcServer, &server{Log: log})
	log.Infow("gRPC server started", "address", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		return errors.Wrap(err, "failed to serve")
	}

	return nil
}
