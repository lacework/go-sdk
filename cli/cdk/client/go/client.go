package cdk

import (
	"context"
	"encoding/json"
	"os"
	"time"

	cdk "github.com/lacework/go-sdk/cli/cdk/go/proto/v1"
	"github.com/lacework/go-sdk/lwlogger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ComponentCDKClient struct {
	logger           ComponentCDKLogger
	coreClient       cdk.CoreClient
	conn             *grpc.ClientConn
	componentVersion string
}

type ComponentCDKLogger interface {
	Infow(msg string, keysAndValues ...interface{})
	Debugf(template string, args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
}

type CDKClientOption func(c *ComponentCDKClient)

func WithLogger(logger *zap.SugaredLogger) CDKClientOption {
	return func(c *ComponentCDKClient) {
		c.logger = logger
	}
}

// NewCDKClient creates a new component CDK client
//
// This client provides opinionated access to the services offerred from gRPC in the CDK (caching, metric data, etc)
//
// Note, ensure you are closing the gRPC connection when your component ends using the `Close()` method
// on the ComponentCDKClient
func NewCDKClient(componentVersion string, opts ...CDKClientOption) (*ComponentCDKClient, error) {
	// set default logger
	defaultLogger := lwlogger.New(os.Getenv("LW_LOG")).Sugar()
	client := &ComponentCDKClient{logger: defaultLogger, componentVersion: componentVersion}

	for _, o := range opts {
		o(client)
	}

	client.logger.Infow("connecting to gRPC server", "address", os.Getenv("LW_CDK_TARGET"))
	conn, err := grpc.Dial(os.Getenv("LW_CDK_TARGET"),
		// we allow insecure connections since we are connecting to 'localhost'
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrap(err, "cannot initilize cdk client")
	}

	client.conn = conn
	client.coreClient = cdk.NewCoreClient(conn)
	return client, nil
}

// Close terminates the gRPC connection to the CDK service
func (c *ComponentCDKClient) Close() error {
	return c.conn.Close()
}

// SetLogger enables overwriting the built-in logger to a custom logger that satifies the ComponentCDKLogger interface
func (c *ComponentCDKClient) SetLogger(logger ComponentCDKLogger) {
	c.logger = logger
}

type CDKCacheMissError struct {
	Err error
}

func (c *CDKCacheMissError) Error() string {
	return c.Err.Error()
}

// ReadCacheAsset fetch key from cache
//
// when an error is returned, if the reason for the error is a cache miss it will be of type
// CDKCacheMissError which should be handled/treated as non-fatal
//
// Response data is in []byte format and will need to be unmarshalled to the correct data type
func (c *ComponentCDKClient) ReadCacheAsset(key string) ([]byte, error) {
	response, err := c.coreClient.ReadCache(context.Background(), &cdk.ReadCacheRequest{
		Key: key,
	})

	if err != nil {
		c.logger.Debugf("error reading cache; %s", err.Error())
		return nil, err
	}

	if response.Hit {
		c.logger.Debug("cache hit",
			"type", "data",
		)
		return response.Data, nil
	}

	c.logger.Debug("cache miss",
		"type", "data",
	)
	return nil, &CDKCacheMissError{Err: errors.New("cache miss")}
}

// WriteCacheAsset persists data to the Lacework CLI on-disk cache via the CDK service
//
// Note, data written to the cache is marshalled to JSON first.
//
// If there is an error writing to cache the error is logged but the return will be nil. Errors writing to
// should never be fatal and stop a component. However, if the data supplied cannot be marshalled into JSON
// an actual error will be returned.
func (c *ComponentCDKClient) WriteCacheAsset(key string, expires time.Time, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.Wrapf(err, "failed to convert data to be cached")
	}

	res, err := c.coreClient.WriteCache(context.Background(), &cdk.WriteCacheRequest{
		Key:     key,
		Expires: timestamppb.New(expires),
		Data:    jsonData,
	})

	if res != nil && res.Error {
		c.logger.Debugf("error writing to cache; %s", res.Message)
	}

	if err != nil {
		c.logger.Debugf("error writing to cache; %s", err.Error())
	}

	return nil
}

// MetricData is used when sending data to Honeycomb
//
// For convience, use the `Metric()` method on the ComponentCDKClient instead of this struct directly
type MetricData struct {
	// Feature name in Honeycomb
	Feature string

	// Feature data in Honeycomb
	FeatureData map[string]string

	// Duration for this span (each MetricData is a unique span)
	Duration int64

	client *ComponentCDKClient
}

// WithDuration attaches a duration to the given span that will be created in Honeycomb (which is optional)
func (m *MetricData) WithDuration(duration int64) *MetricData {
	m.Duration = duration
	return m
}

// Send Writes the MetricData to Honeycomb
func (m *MetricData) Send() error {
	return m.client.sendMetricData(m)
}

// Metric is used to create a new MetricData struct that can be augmented and ultimately sent
//
//	 now := time.Now()
//	 c, _ := NewCDKClient("0.0.1")
//	 _ = c.Metric("example", map[string]string{"example": "data"}).Send()
//	 _ = c.Metric("example2", map[string]string{"example2": "data"}).
//		       WithDuration(time.Since(now).Milliseconds()).
//		       Send()
func (c *ComponentCDKClient) Metric(feature string, featureData map[string]string) *MetricData {
	return &MetricData{Feature: feature, FeatureData: featureData, client: c}
}

func (c *ComponentCDKClient) sendMetricData(data *MetricData) error {
	data.FeatureData["version"] = c.componentVersion

	request := &cdk.HoneyventRequest{
		Feature: data.Feature, FeatureData: data.FeatureData,
	}

	if data.Duration != 0 {
		request.DurationMs = data.Duration
	}

	_, err := c.coreClient.Honeyvent(context.Background(), request)
	if err != nil {
		c.logger.Warn("unable to send telemetry",
			"type", "data", "error", err.Error(),
		)
	}

	return nil
}

// MetricError is used to write an error to Honeycomb
func (c *ComponentCDKClient) MetricError(e error) {
	_, err := c.coreClient.Honeyvent(context.Background(), &cdk.HoneyventRequest{
		Error: e.Error(),
	})
	if err != nil {
		c.logger.Warn("unable to send telemetry",
			"type", "error", "error", err.Error(),
		)
	}
}
