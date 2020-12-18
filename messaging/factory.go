//
// Copyright (c) 2019 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package messaging

import (
	"fmt"
	"strings"

	"github.com/edgexfoundry/go-mod-messaging/internal/pkg/mqtt"
	"github.com/edgexfoundry/go-mod-messaging/internal/pkg/redis/streams"
	"github.com/edgexfoundry/go-mod-messaging/internal/pkg/zeromq"
	"github.com/edgexfoundry/go-mod-messaging/pkg/types"
)

const (
	// ZeroMQ messaging implementation
	ZeroMQ = "zero"

	// MQTT messaging implementation
	MQTT = "mqtt"

	// RedisStreams messaging implementation
	RedisStreams = "redisstreams"
)

var customTypes = make(map[string]func(config types.MessageBusConfig) (MessageClient, error))

// RegisterCustomType allows registering custom messagebus client types for use by NewMessageClient
func RegisterCustomType(msgType string, builder func(config types.MessageBusConfig) (MessageClient, error)) {
	lowerType := strings.ToLower(msgType)
	customTypes[lowerType] = builder
}

// NewMessageClient is a factory function to instantiate different message client depending on
// the "Type" from the configuration
func NewMessageClient(msgConfig types.MessageBusConfig) (MessageClient, error) {

	if msgConfig.PublishHost.IsHostInfoEmpty() && msgConfig.SubscribeHost.IsHostInfoEmpty() {
		return nil, fmt.Errorf("unable to create messageClient: host info not set")
	}

	switch lowerMsgType := strings.ToLower(msgConfig.Type); lowerMsgType {
	case ZeroMQ:
		return zeromq.NewZeroMqClient(msgConfig)
	case MQTT:
		return mqtt.NewMQTTClient(msgConfig)
	case RedisStreams:
		return streams.NewClient(msgConfig)
	default:
		for t, b := range customTypes {
			if lowerMsgType == t {
				return b(msgConfig)
			}
		}
		return nil, fmt.Errorf("unknown message type '%s' requested", msgConfig.Type)
	}
}
