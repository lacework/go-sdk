package metric

import (
	"context"
	"os"

	"[[.Component]]/internal/logger"
	"[[.Component]]/internal/version"

	cdk "github.com/lacework/go-sdk/cli/cdk/go/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var cdkClient cdk.CoreClient

func init() {
	// Set up a connection to the CDK server
	logger.Log.Infow("connecting to gRPC server", "address", os.Getenv("LW_CDK_TARGET"))
	conn, err := grpc.Dial(os.Getenv("LW_CDK_TARGET"),
		// we allow insecure connections since we are connecting to 'localhost'
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Warn("Cannot initialize CDK client", "error", err.Error())
	} else {
		cdkClient = cdk.NewCoreClient(conn)
	}
}

func SendMetricData(feature string, data map[string]string) {
	if cdkClient == nil {
		logger.Log.Warn("unable to send telemetry",
			"type", "data",
			"error", "client not initialized",
		)
		return
	}

	// add preflight version
	data["version"] = version.Version

	go func() {
		_, err := cdkClient.Honeyvent(context.Background(), &cdk.HoneyventRequest{
			Feature: feature, FeatureData: data,
		})
		if err != nil {
			logger.Log.Warn("unable to send telemetry",
				"type", "data", "error", err.Error(),
			)
		}
	}()
}

func SendMetricError(e error) {
	if cdkClient == nil {
		logger.Log.Warn("unable to send telemetry",
			"type", "error",
			"error", "client not initialized",
		)
		return
	}

	go func() {
		_, err := cdkClient.Honeyvent(context.Background(), &cdk.HoneyventRequest{
			Error: e.Error(),
		})
		if err != nil {
			logger.Log.Warn("unable to send telemetry",
				"type", "error", "error", err.Error(),
			)
		}
	}()
}
