// Code generated by smithy-go-codegen DO NOT EDIT.

package ec2instanceconnect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/restjson"
	"github.com/aws/aws-sdk-go-v2/service/ec2instanceconnect/types"
	smithy "github.com/aws/smithy-go"
	smithyio "github.com/aws/smithy-go/io"
	"github.com/aws/smithy-go/middleware"
	"github.com/aws/smithy-go/ptr"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"io"
	"strings"
)

type awsAwsjson11_deserializeOpSendSerialConsoleSSHPublicKey struct {
}

func (*awsAwsjson11_deserializeOpSendSerialConsoleSSHPublicKey) ID() string {
	return "OperationDeserializer"
}

func (m *awsAwsjson11_deserializeOpSendSerialConsoleSSHPublicKey) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	out, metadata, err = next.HandleDeserialize(ctx, in)
	if err != nil {
		return out, metadata, err
	}

	response, ok := out.RawResponse.(*smithyhttp.Response)
	if !ok {
		return out, metadata, &smithy.DeserializationError{Err: fmt.Errorf("unknown transport type %T", out.RawResponse)}
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return out, metadata, awsAwsjson11_deserializeOpErrorSendSerialConsoleSSHPublicKey(response, &metadata)
	}
	output := &SendSerialConsoleSSHPublicKeyOutput{}
	out.Result = output

	var buff [1024]byte
	ringBuffer := smithyio.NewRingBuffer(buff[:])

	body := io.TeeReader(response.Body, ringBuffer)
	decoder := json.NewDecoder(body)
	decoder.UseNumber()
	var shape interface{}
	if err := decoder.Decode(&shape); err != nil && err != io.EOF {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return out, metadata, err
	}

	err = awsAwsjson11_deserializeOpDocumentSendSerialConsoleSSHPublicKeyOutput(&output, shape)
	if err != nil {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return out, metadata, err
	}

	return out, metadata, err
}

func awsAwsjson11_deserializeOpErrorSendSerialConsoleSSHPublicKey(response *smithyhttp.Response, metadata *middleware.Metadata) error {
	var errorBuffer bytes.Buffer
	if _, err := io.Copy(&errorBuffer, response.Body); err != nil {
		return &smithy.DeserializationError{Err: fmt.Errorf("failed to copy error response body, %w", err)}
	}
	errorBody := bytes.NewReader(errorBuffer.Bytes())

	errorCode := "UnknownError"
	errorMessage := errorCode

	code := response.Header.Get("X-Amzn-ErrorType")
	if len(code) != 0 {
		errorCode = restjson.SanitizeErrorCode(code)
	}

	var buff [1024]byte
	ringBuffer := smithyio.NewRingBuffer(buff[:])

	body := io.TeeReader(errorBody, ringBuffer)
	decoder := json.NewDecoder(body)
	decoder.UseNumber()
	code, message, err := restjson.GetErrorInfo(decoder)
	if err != nil {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	errorBody.Seek(0, io.SeekStart)
	if len(code) != 0 {
		errorCode = restjson.SanitizeErrorCode(code)
	}
	if len(message) != 0 {
		errorMessage = message
	}

	switch {
	case strings.EqualFold("AuthException", errorCode):
		return awsAwsjson11_deserializeErrorAuthException(response, errorBody)

	case strings.EqualFold("EC2InstanceNotFoundException", errorCode):
		return awsAwsjson11_deserializeErrorEC2InstanceNotFoundException(response, errorBody)

	case strings.EqualFold("EC2InstanceStateInvalidException", errorCode):
		return awsAwsjson11_deserializeErrorEC2InstanceStateInvalidException(response, errorBody)

	case strings.EqualFold("EC2InstanceTypeInvalidException", errorCode):
		return awsAwsjson11_deserializeErrorEC2InstanceTypeInvalidException(response, errorBody)

	case strings.EqualFold("EC2InstanceUnavailableException", errorCode):
		return awsAwsjson11_deserializeErrorEC2InstanceUnavailableException(response, errorBody)

	case strings.EqualFold("InvalidArgsException", errorCode):
		return awsAwsjson11_deserializeErrorInvalidArgsException(response, errorBody)

	case strings.EqualFold("SerialConsoleAccessDisabledException", errorCode):
		return awsAwsjson11_deserializeErrorSerialConsoleAccessDisabledException(response, errorBody)

	case strings.EqualFold("SerialConsoleSessionLimitExceededException", errorCode):
		return awsAwsjson11_deserializeErrorSerialConsoleSessionLimitExceededException(response, errorBody)

	case strings.EqualFold("SerialConsoleSessionUnavailableException", errorCode):
		return awsAwsjson11_deserializeErrorSerialConsoleSessionUnavailableException(response, errorBody)

	case strings.EqualFold("ServiceException", errorCode):
		return awsAwsjson11_deserializeErrorServiceException(response, errorBody)

	case strings.EqualFold("ThrottlingException", errorCode):
		return awsAwsjson11_deserializeErrorThrottlingException(response, errorBody)

	default:
		genericError := &smithy.GenericAPIError{
			Code:    errorCode,
			Message: errorMessage,
		}
		return genericError

	}
}

type awsAwsjson11_deserializeOpSendSSHPublicKey struct {
}

func (*awsAwsjson11_deserializeOpSendSSHPublicKey) ID() string {
	return "OperationDeserializer"
}

func (m *awsAwsjson11_deserializeOpSendSSHPublicKey) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	out, metadata, err = next.HandleDeserialize(ctx, in)
	if err != nil {
		return out, metadata, err
	}

	response, ok := out.RawResponse.(*smithyhttp.Response)
	if !ok {
		return out, metadata, &smithy.DeserializationError{Err: fmt.Errorf("unknown transport type %T", out.RawResponse)}
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return out, metadata, awsAwsjson11_deserializeOpErrorSendSSHPublicKey(response, &metadata)
	}
	output := &SendSSHPublicKeyOutput{}
	out.Result = output

	var buff [1024]byte
	ringBuffer := smithyio.NewRingBuffer(buff[:])

	body := io.TeeReader(response.Body, ringBuffer)
	decoder := json.NewDecoder(body)
	decoder.UseNumber()
	var shape interface{}
	if err := decoder.Decode(&shape); err != nil && err != io.EOF {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return out, metadata, err
	}

	err = awsAwsjson11_deserializeOpDocumentSendSSHPublicKeyOutput(&output, shape)
	if err != nil {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return out, metadata, err
	}

	return out, metadata, err
}

func awsAwsjson11_deserializeOpErrorSendSSHPublicKey(response *smithyhttp.Response, metadata *middleware.Metadata) error {
	var errorBuffer bytes.Buffer
	if _, err := io.Copy(&errorBuffer, response.Body); err != nil {
		return &smithy.DeserializationError{Err: fmt.Errorf("failed to copy error response body, %w", err)}
	}
	errorBody := bytes.NewReader(errorBuffer.Bytes())

	errorCode := "UnknownError"
	errorMessage := errorCode

	code := response.Header.Get("X-Amzn-ErrorType")
	if len(code) != 0 {
		errorCode = restjson.SanitizeErrorCode(code)
	}

	var buff [1024]byte
	ringBuffer := smithyio.NewRingBuffer(buff[:])

	body := io.TeeReader(errorBody, ringBuffer)
	decoder := json.NewDecoder(body)
	decoder.UseNumber()
	code, message, err := restjson.GetErrorInfo(decoder)
	if err != nil {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	errorBody.Seek(0, io.SeekStart)
	if len(code) != 0 {
		errorCode = restjson.SanitizeErrorCode(code)
	}
	if len(message) != 0 {
		errorMessage = message
	}

	switch {
	case strings.EqualFold("AuthException", errorCode):
		return awsAwsjson11_deserializeErrorAuthException(response, errorBody)

	case strings.EqualFold("EC2InstanceNotFoundException", errorCode):
		return awsAwsjson11_deserializeErrorEC2InstanceNotFoundException(response, errorBody)

	case strings.EqualFold("EC2InstanceStateInvalidException", errorCode):
		return awsAwsjson11_deserializeErrorEC2InstanceStateInvalidException(response, errorBody)

	case strings.EqualFold("EC2InstanceUnavailableException", errorCode):
		return awsAwsjson11_deserializeErrorEC2InstanceUnavailableException(response, errorBody)

	case strings.EqualFold("InvalidArgsException", errorCode):
		return awsAwsjson11_deserializeErrorInvalidArgsException(response, errorBody)

	case strings.EqualFold("ServiceException", errorCode):
		return awsAwsjson11_deserializeErrorServiceException(response, errorBody)

	case strings.EqualFold("ThrottlingException", errorCode):
		return awsAwsjson11_deserializeErrorThrottlingException(response, errorBody)

	default:
		genericError := &smithy.GenericAPIError{
			Code:    errorCode,
			Message: errorMessage,
		}
		return genericError

	}
}

func awsAwsjson11_deserializeErrorAuthException(response *smithyhttp.Response, errorBody *bytes.Reader) error {
	var buff [1024]byte
	ringBuffer := smithyio.NewRingBuffer(buff[:])

	body := io.TeeReader(errorBody, ringBuffer)
	decoder := json.NewDecoder(body)
	decoder.UseNumber()
	var shape interface{}
	if err := decoder.Decode(&shape); err != nil && err != io.EOF {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	output := &types.AuthException{}
	err := awsAwsjson11_deserializeDocumentAuthException(&output, shape)

	if err != nil {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	errorBody.Seek(0, io.SeekStart)
	return output
}

func awsAwsjson11_deserializeErrorEC2InstanceNotFoundException(response *smithyhttp.Response, errorBody *bytes.Reader) error {
	var buff [1024]byte
	ringBuffer := smithyio.NewRingBuffer(buff[:])

	body := io.TeeReader(errorBody, ringBuffer)
	decoder := json.NewDecoder(body)
	decoder.UseNumber()
	var shape interface{}
	if err := decoder.Decode(&shape); err != nil && err != io.EOF {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	output := &types.EC2InstanceNotFoundException{}
	err := awsAwsjson11_deserializeDocumentEC2InstanceNotFoundException(&output, shape)

	if err != nil {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	errorBody.Seek(0, io.SeekStart)
	return output
}

func awsAwsjson11_deserializeErrorEC2InstanceStateInvalidException(response *smithyhttp.Response, errorBody *bytes.Reader) error {
	var buff [1024]byte
	ringBuffer := smithyio.NewRingBuffer(buff[:])

	body := io.TeeReader(errorBody, ringBuffer)
	decoder := json.NewDecoder(body)
	decoder.UseNumber()
	var shape interface{}
	if err := decoder.Decode(&shape); err != nil && err != io.EOF {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	output := &types.EC2InstanceStateInvalidException{}
	err := awsAwsjson11_deserializeDocumentEC2InstanceStateInvalidException(&output, shape)

	if err != nil {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	errorBody.Seek(0, io.SeekStart)
	return output
}

func awsAwsjson11_deserializeErrorEC2InstanceTypeInvalidException(response *smithyhttp.Response, errorBody *bytes.Reader) error {
	var buff [1024]byte
	ringBuffer := smithyio.NewRingBuffer(buff[:])

	body := io.TeeReader(errorBody, ringBuffer)
	decoder := json.NewDecoder(body)
	decoder.UseNumber()
	var shape interface{}
	if err := decoder.Decode(&shape); err != nil && err != io.EOF {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	output := &types.EC2InstanceTypeInvalidException{}
	err := awsAwsjson11_deserializeDocumentEC2InstanceTypeInvalidException(&output, shape)

	if err != nil {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	errorBody.Seek(0, io.SeekStart)
	return output
}

func awsAwsjson11_deserializeErrorEC2InstanceUnavailableException(response *smithyhttp.Response, errorBody *bytes.Reader) error {
	var buff [1024]byte
	ringBuffer := smithyio.NewRingBuffer(buff[:])

	body := io.TeeReader(errorBody, ringBuffer)
	decoder := json.NewDecoder(body)
	decoder.UseNumber()
	var shape interface{}
	if err := decoder.Decode(&shape); err != nil && err != io.EOF {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	output := &types.EC2InstanceUnavailableException{}
	err := awsAwsjson11_deserializeDocumentEC2InstanceUnavailableException(&output, shape)

	if err != nil {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	errorBody.Seek(0, io.SeekStart)
	return output
}

func awsAwsjson11_deserializeErrorInvalidArgsException(response *smithyhttp.Response, errorBody *bytes.Reader) error {
	var buff [1024]byte
	ringBuffer := smithyio.NewRingBuffer(buff[:])

	body := io.TeeReader(errorBody, ringBuffer)
	decoder := json.NewDecoder(body)
	decoder.UseNumber()
	var shape interface{}
	if err := decoder.Decode(&shape); err != nil && err != io.EOF {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	output := &types.InvalidArgsException{}
	err := awsAwsjson11_deserializeDocumentInvalidArgsException(&output, shape)

	if err != nil {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	errorBody.Seek(0, io.SeekStart)
	return output
}

func awsAwsjson11_deserializeErrorSerialConsoleAccessDisabledException(response *smithyhttp.Response, errorBody *bytes.Reader) error {
	var buff [1024]byte
	ringBuffer := smithyio.NewRingBuffer(buff[:])

	body := io.TeeReader(errorBody, ringBuffer)
	decoder := json.NewDecoder(body)
	decoder.UseNumber()
	var shape interface{}
	if err := decoder.Decode(&shape); err != nil && err != io.EOF {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	output := &types.SerialConsoleAccessDisabledException{}
	err := awsAwsjson11_deserializeDocumentSerialConsoleAccessDisabledException(&output, shape)

	if err != nil {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	errorBody.Seek(0, io.SeekStart)
	return output
}

func awsAwsjson11_deserializeErrorSerialConsoleSessionLimitExceededException(response *smithyhttp.Response, errorBody *bytes.Reader) error {
	var buff [1024]byte
	ringBuffer := smithyio.NewRingBuffer(buff[:])

	body := io.TeeReader(errorBody, ringBuffer)
	decoder := json.NewDecoder(body)
	decoder.UseNumber()
	var shape interface{}
	if err := decoder.Decode(&shape); err != nil && err != io.EOF {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	output := &types.SerialConsoleSessionLimitExceededException{}
	err := awsAwsjson11_deserializeDocumentSerialConsoleSessionLimitExceededException(&output, shape)

	if err != nil {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	errorBody.Seek(0, io.SeekStart)
	return output
}

func awsAwsjson11_deserializeErrorSerialConsoleSessionUnavailableException(response *smithyhttp.Response, errorBody *bytes.Reader) error {
	var buff [1024]byte
	ringBuffer := smithyio.NewRingBuffer(buff[:])

	body := io.TeeReader(errorBody, ringBuffer)
	decoder := json.NewDecoder(body)
	decoder.UseNumber()
	var shape interface{}
	if err := decoder.Decode(&shape); err != nil && err != io.EOF {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	output := &types.SerialConsoleSessionUnavailableException{}
	err := awsAwsjson11_deserializeDocumentSerialConsoleSessionUnavailableException(&output, shape)

	if err != nil {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	errorBody.Seek(0, io.SeekStart)
	return output
}

func awsAwsjson11_deserializeErrorServiceException(response *smithyhttp.Response, errorBody *bytes.Reader) error {
	var buff [1024]byte
	ringBuffer := smithyio.NewRingBuffer(buff[:])

	body := io.TeeReader(errorBody, ringBuffer)
	decoder := json.NewDecoder(body)
	decoder.UseNumber()
	var shape interface{}
	if err := decoder.Decode(&shape); err != nil && err != io.EOF {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	output := &types.ServiceException{}
	err := awsAwsjson11_deserializeDocumentServiceException(&output, shape)

	if err != nil {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	errorBody.Seek(0, io.SeekStart)
	return output
}

func awsAwsjson11_deserializeErrorThrottlingException(response *smithyhttp.Response, errorBody *bytes.Reader) error {
	var buff [1024]byte
	ringBuffer := smithyio.NewRingBuffer(buff[:])

	body := io.TeeReader(errorBody, ringBuffer)
	decoder := json.NewDecoder(body)
	decoder.UseNumber()
	var shape interface{}
	if err := decoder.Decode(&shape); err != nil && err != io.EOF {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	output := &types.ThrottlingException{}
	err := awsAwsjson11_deserializeDocumentThrottlingException(&output, shape)

	if err != nil {
		var snapshot bytes.Buffer
		io.Copy(&snapshot, ringBuffer)
		err = &smithy.DeserializationError{
			Err:      fmt.Errorf("failed to decode response body, %w", err),
			Snapshot: snapshot.Bytes(),
		}
		return err
	}

	errorBody.Seek(0, io.SeekStart)
	return output
}

func awsAwsjson11_deserializeDocumentAuthException(v **types.AuthException, value interface{}) error {
	if v == nil {
		return fmt.Errorf("unexpected nil of type %T", v)
	}
	if value == nil {
		return nil
	}

	shape, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected JSON type %v", value)
	}

	var sv *types.AuthException
	if *v == nil {
		sv = &types.AuthException{}
	} else {
		sv = *v
	}

	for key, value := range shape {
		switch key {
		case "Message":
			if value != nil {
				jtv, ok := value.(string)
				if !ok {
					return fmt.Errorf("expected String to be of type string, got %T instead", value)
				}
				sv.Message = ptr.String(jtv)
			}

		default:
			_, _ = key, value

		}
	}
	*v = sv
	return nil
}

func awsAwsjson11_deserializeDocumentEC2InstanceNotFoundException(v **types.EC2InstanceNotFoundException, value interface{}) error {
	if v == nil {
		return fmt.Errorf("unexpected nil of type %T", v)
	}
	if value == nil {
		return nil
	}

	shape, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected JSON type %v", value)
	}

	var sv *types.EC2InstanceNotFoundException
	if *v == nil {
		sv = &types.EC2InstanceNotFoundException{}
	} else {
		sv = *v
	}

	for key, value := range shape {
		switch key {
		case "Message":
			if value != nil {
				jtv, ok := value.(string)
				if !ok {
					return fmt.Errorf("expected String to be of type string, got %T instead", value)
				}
				sv.Message = ptr.String(jtv)
			}

		default:
			_, _ = key, value

		}
	}
	*v = sv
	return nil
}

func awsAwsjson11_deserializeDocumentEC2InstanceStateInvalidException(v **types.EC2InstanceStateInvalidException, value interface{}) error {
	if v == nil {
		return fmt.Errorf("unexpected nil of type %T", v)
	}
	if value == nil {
		return nil
	}

	shape, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected JSON type %v", value)
	}

	var sv *types.EC2InstanceStateInvalidException
	if *v == nil {
		sv = &types.EC2InstanceStateInvalidException{}
	} else {
		sv = *v
	}

	for key, value := range shape {
		switch key {
		case "Message":
			if value != nil {
				jtv, ok := value.(string)
				if !ok {
					return fmt.Errorf("expected String to be of type string, got %T instead", value)
				}
				sv.Message = ptr.String(jtv)
			}

		default:
			_, _ = key, value

		}
	}
	*v = sv
	return nil
}

func awsAwsjson11_deserializeDocumentEC2InstanceTypeInvalidException(v **types.EC2InstanceTypeInvalidException, value interface{}) error {
	if v == nil {
		return fmt.Errorf("unexpected nil of type %T", v)
	}
	if value == nil {
		return nil
	}

	shape, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected JSON type %v", value)
	}

	var sv *types.EC2InstanceTypeInvalidException
	if *v == nil {
		sv = &types.EC2InstanceTypeInvalidException{}
	} else {
		sv = *v
	}

	for key, value := range shape {
		switch key {
		case "Message":
			if value != nil {
				jtv, ok := value.(string)
				if !ok {
					return fmt.Errorf("expected String to be of type string, got %T instead", value)
				}
				sv.Message = ptr.String(jtv)
			}

		default:
			_, _ = key, value

		}
	}
	*v = sv
	return nil
}

func awsAwsjson11_deserializeDocumentEC2InstanceUnavailableException(v **types.EC2InstanceUnavailableException, value interface{}) error {
	if v == nil {
		return fmt.Errorf("unexpected nil of type %T", v)
	}
	if value == nil {
		return nil
	}

	shape, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected JSON type %v", value)
	}

	var sv *types.EC2InstanceUnavailableException
	if *v == nil {
		sv = &types.EC2InstanceUnavailableException{}
	} else {
		sv = *v
	}

	for key, value := range shape {
		switch key {
		case "Message":
			if value != nil {
				jtv, ok := value.(string)
				if !ok {
					return fmt.Errorf("expected String to be of type string, got %T instead", value)
				}
				sv.Message = ptr.String(jtv)
			}

		default:
			_, _ = key, value

		}
	}
	*v = sv
	return nil
}

func awsAwsjson11_deserializeDocumentInvalidArgsException(v **types.InvalidArgsException, value interface{}) error {
	if v == nil {
		return fmt.Errorf("unexpected nil of type %T", v)
	}
	if value == nil {
		return nil
	}

	shape, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected JSON type %v", value)
	}

	var sv *types.InvalidArgsException
	if *v == nil {
		sv = &types.InvalidArgsException{}
	} else {
		sv = *v
	}

	for key, value := range shape {
		switch key {
		case "Message":
			if value != nil {
				jtv, ok := value.(string)
				if !ok {
					return fmt.Errorf("expected String to be of type string, got %T instead", value)
				}
				sv.Message = ptr.String(jtv)
			}

		default:
			_, _ = key, value

		}
	}
	*v = sv
	return nil
}

func awsAwsjson11_deserializeDocumentSerialConsoleAccessDisabledException(v **types.SerialConsoleAccessDisabledException, value interface{}) error {
	if v == nil {
		return fmt.Errorf("unexpected nil of type %T", v)
	}
	if value == nil {
		return nil
	}

	shape, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected JSON type %v", value)
	}

	var sv *types.SerialConsoleAccessDisabledException
	if *v == nil {
		sv = &types.SerialConsoleAccessDisabledException{}
	} else {
		sv = *v
	}

	for key, value := range shape {
		switch key {
		case "Message":
			if value != nil {
				jtv, ok := value.(string)
				if !ok {
					return fmt.Errorf("expected String to be of type string, got %T instead", value)
				}
				sv.Message = ptr.String(jtv)
			}

		default:
			_, _ = key, value

		}
	}
	*v = sv
	return nil
}

func awsAwsjson11_deserializeDocumentSerialConsoleSessionLimitExceededException(v **types.SerialConsoleSessionLimitExceededException, value interface{}) error {
	if v == nil {
		return fmt.Errorf("unexpected nil of type %T", v)
	}
	if value == nil {
		return nil
	}

	shape, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected JSON type %v", value)
	}

	var sv *types.SerialConsoleSessionLimitExceededException
	if *v == nil {
		sv = &types.SerialConsoleSessionLimitExceededException{}
	} else {
		sv = *v
	}

	for key, value := range shape {
		switch key {
		case "Message":
			if value != nil {
				jtv, ok := value.(string)
				if !ok {
					return fmt.Errorf("expected String to be of type string, got %T instead", value)
				}
				sv.Message = ptr.String(jtv)
			}

		default:
			_, _ = key, value

		}
	}
	*v = sv
	return nil
}

func awsAwsjson11_deserializeDocumentSerialConsoleSessionUnavailableException(v **types.SerialConsoleSessionUnavailableException, value interface{}) error {
	if v == nil {
		return fmt.Errorf("unexpected nil of type %T", v)
	}
	if value == nil {
		return nil
	}

	shape, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected JSON type %v", value)
	}

	var sv *types.SerialConsoleSessionUnavailableException
	if *v == nil {
		sv = &types.SerialConsoleSessionUnavailableException{}
	} else {
		sv = *v
	}

	for key, value := range shape {
		switch key {
		case "Message":
			if value != nil {
				jtv, ok := value.(string)
				if !ok {
					return fmt.Errorf("expected String to be of type string, got %T instead", value)
				}
				sv.Message = ptr.String(jtv)
			}

		default:
			_, _ = key, value

		}
	}
	*v = sv
	return nil
}

func awsAwsjson11_deserializeDocumentServiceException(v **types.ServiceException, value interface{}) error {
	if v == nil {
		return fmt.Errorf("unexpected nil of type %T", v)
	}
	if value == nil {
		return nil
	}

	shape, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected JSON type %v", value)
	}

	var sv *types.ServiceException
	if *v == nil {
		sv = &types.ServiceException{}
	} else {
		sv = *v
	}

	for key, value := range shape {
		switch key {
		case "Message":
			if value != nil {
				jtv, ok := value.(string)
				if !ok {
					return fmt.Errorf("expected String to be of type string, got %T instead", value)
				}
				sv.Message = ptr.String(jtv)
			}

		default:
			_, _ = key, value

		}
	}
	*v = sv
	return nil
}

func awsAwsjson11_deserializeDocumentThrottlingException(v **types.ThrottlingException, value interface{}) error {
	if v == nil {
		return fmt.Errorf("unexpected nil of type %T", v)
	}
	if value == nil {
		return nil
	}

	shape, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected JSON type %v", value)
	}

	var sv *types.ThrottlingException
	if *v == nil {
		sv = &types.ThrottlingException{}
	} else {
		sv = *v
	}

	for key, value := range shape {
		switch key {
		case "Message":
			if value != nil {
				jtv, ok := value.(string)
				if !ok {
					return fmt.Errorf("expected String to be of type string, got %T instead", value)
				}
				sv.Message = ptr.String(jtv)
			}

		default:
			_, _ = key, value

		}
	}
	*v = sv
	return nil
}

func awsAwsjson11_deserializeOpDocumentSendSerialConsoleSSHPublicKeyOutput(v **SendSerialConsoleSSHPublicKeyOutput, value interface{}) error {
	if v == nil {
		return fmt.Errorf("unexpected nil of type %T", v)
	}
	if value == nil {
		return nil
	}

	shape, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected JSON type %v", value)
	}

	var sv *SendSerialConsoleSSHPublicKeyOutput
	if *v == nil {
		sv = &SendSerialConsoleSSHPublicKeyOutput{}
	} else {
		sv = *v
	}

	for key, value := range shape {
		switch key {
		case "RequestId":
			if value != nil {
				jtv, ok := value.(string)
				if !ok {
					return fmt.Errorf("expected RequestId to be of type string, got %T instead", value)
				}
				sv.RequestId = ptr.String(jtv)
			}

		case "Success":
			if value != nil {
				jtv, ok := value.(bool)
				if !ok {
					return fmt.Errorf("expected Success to be of type *bool, got %T instead", value)
				}
				sv.Success = jtv
			}

		default:
			_, _ = key, value

		}
	}
	*v = sv
	return nil
}

func awsAwsjson11_deserializeOpDocumentSendSSHPublicKeyOutput(v **SendSSHPublicKeyOutput, value interface{}) error {
	if v == nil {
		return fmt.Errorf("unexpected nil of type %T", v)
	}
	if value == nil {
		return nil
	}

	shape, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected JSON type %v", value)
	}

	var sv *SendSSHPublicKeyOutput
	if *v == nil {
		sv = &SendSSHPublicKeyOutput{}
	} else {
		sv = *v
	}

	for key, value := range shape {
		switch key {
		case "RequestId":
			if value != nil {
				jtv, ok := value.(string)
				if !ok {
					return fmt.Errorf("expected RequestId to be of type string, got %T instead", value)
				}
				sv.RequestId = ptr.String(jtv)
			}

		case "Success":
			if value != nil {
				jtv, ok := value.(bool)
				if !ok {
					return fmt.Errorf("expected Success to be of type *bool, got %T instead", value)
				}
				sv.Success = jtv
			}

		default:
			_, _ = key, value

		}
	}
	*v = sv
	return nil
}