//
// Copyright:: Copyright 2023, Lacework Inc.
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

package logger

import (
	"os"

	"github.com/lacework/go-sdk/lwlogger"
)

// Log allows this component to leverage our Go-SDK lwlogger
//
// Example:
//
//	import "preflight/internal/logger"
//	logger.Log.Info("an informational message")
//	logger.Log.Debug("a debug message")
//	logger.Log.Infow("info message with variables", "foo", "bar")
var Log = lwlogger.New(os.Getenv("LW_LOG")).Sugar()
