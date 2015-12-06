// Copyright (c) 2014 The gomqtt Authors. All rights reserved.
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

package message

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubscribeInterface(t *testing.T) {
	msg := NewSubscribeMessage()
	msg.Subscriptions = []Subscription{
		{Topic: []byte("hello"), QoS: QosAtMostOnce},
	}

	require.Equal(t, msg.Type(), SUBSCRIBE)
	require.NotNil(t, msg.String())
}

func TestSubscribeMessageDecode(t *testing.T) {
	msgBytes := []byte{
		byte(SUBSCRIBE<<4) | 2,
		36,
		0, // packet ID MSB
		7, // packet ID LSB
		0, // topic name MSB
		7, // topic name LSB
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0, // QoS
		0, // topic name MSB
		8, // topic name LSB
		'/', 'a', '/', 'b', '/', '#', '/', 'c',
		1,  // QoS
		0,  // topic name MSB
		10, // topic name LSB
		'/', 'a', '/', 'b', '/', '#', '/', 'c', 'd', 'd',
		2, // QoS
	}

	msg := NewSubscribeMessage()
	n, err := msg.Decode(msgBytes)

	require.NoError(t, err)
	require.Equal(t, len(msgBytes), n)
	require.Equal(t, 3, len(msg.Subscriptions))
	require.Equal(t, []byte("surgemq"), msg.Subscriptions[0].Topic)
	require.Equal(t, 0, int(msg.Subscriptions[0].QoS))
	require.Equal(t, []byte("/a/b/#/c"), msg.Subscriptions[1].Topic)
	require.Equal(t, 1, int(msg.Subscriptions[1].QoS))
	require.Equal(t, []byte("/a/b/#/cdd"), msg.Subscriptions[2].Topic)
	require.Equal(t, 2, int(msg.Subscriptions[2].QoS))
}

func TestSubscribeMessageDecodeError1(t *testing.T) {
	msgBytes := []byte{
		byte(SUBSCRIBE<<4) | 2,
		9, // <- too much
	}

	msg := NewSubscribeMessage()
	_, err := msg.Decode(msgBytes)

	require.Error(t, err)
}

func TestSubscribeMessageDecodeError2(t *testing.T) {
	msgBytes := []byte{
		byte(SUBSCRIBE<<4) | 2,
		0,
		// <- missing packet id
	}

	msg := NewSubscribeMessage()
	_, err := msg.Decode(msgBytes)

	require.Error(t, err)
}

func TestSubscribeMessageDecodeError3(t *testing.T) {
	msgBytes := []byte{
		byte(SUBSCRIBE<<4) | 2,
		2,
		0, // packet ID MSB
		7, // packet ID LSB
		// <- missing subscription
	}

	msg := NewSubscribeMessage()
	_, err := msg.Decode(msgBytes)

	require.Error(t, err)
}

func TestSubscribeMessageDecodeError4(t *testing.T) {
	msgBytes := []byte{
		byte(SUBSCRIBE<<4) | 2,
		5,
		0, // packet ID MSB
		7, // packet ID LSB
		0, // topic name MSB
		2, // topic name LSB <- wrong size
		's',
	}

	msg := NewSubscribeMessage()
	_, err := msg.Decode(msgBytes)

	require.Error(t, err)
}

func TestSubscribeMessageDecodeError5(t *testing.T) {
	msgBytes := []byte{
		byte(SUBSCRIBE<<4) | 2,
		5,
		0, // packet ID MSB
		7, // packet ID LSB
		0, // topic name MSB
		1, // topic name LSB
		's',
		// <- missing qos
	}

	msg := NewSubscribeMessage()
	_, err := msg.Decode(msgBytes)

	require.Error(t, err)
}

func TestSubscribeMessageEncode(t *testing.T) {
	msgBytes := []byte{
		byte(SUBSCRIBE<<4) | 2,
		36,
		0, // packet ID MSB
		7, // packet ID LSB
		0, // topic name MSB
		7, // topic name LSB
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0, // QoS
		0, // topic name MSB
		8, // topic name LSB
		'/', 'a', '/', 'b', '/', '#', '/', 'c',
		1,  // QoS
		0,  // topic name MSB
		10, // topic name LSB
		'/', 'a', '/', 'b', '/', '#', '/', 'c', 'd', 'd',
		2, // QoS
	}

	msg := NewSubscribeMessage()
	msg.PacketId = 7
	msg.Subscriptions = []Subscription{
		{[]byte("surgemq"), 0},
		{[]byte("/a/b/#/c"), 1},
		{[]byte("/a/b/#/cdd"), 2},
	}

	dst := make([]byte, msg.Len())
	n, err := msg.Encode(dst)

	require.NoError(t, err)
	require.Equal(t, len(msgBytes), n)
	require.Equal(t, msgBytes, dst)
}

func TestSubscribeMessageEncodeError1(t *testing.T) {
	msg := NewSubscribeMessage()

	dst := make([]byte, 1) // <- too small
	_, err := msg.Encode(dst)

	require.Error(t, err)
}

func TestSubscribeMessageEncodeError2(t *testing.T) {
	msg := NewSubscribeMessage()
	msg.Subscriptions = []Subscription{
		{make([]byte, 65536), 0}, // too big
	}

	dst := make([]byte, msg.Len())
	_, err := msg.Encode(dst)

	require.Error(t, err)
}

func TestSubscribeEqualDecodeEncode(t *testing.T) {
	msgBytes := []byte{
		byte(SUBSCRIBE<<4) | 2,
		36,
		0, // packet ID MSB
		7, // packet ID LSB
		0, // topic name MSB
		7, // topic name LSB
		's', 'u', 'r', 'g', 'e', 'm', 'q',
		0, // QoS
		0, // topic name MSB
		8, // topic name LSB
		'/', 'a', '/', 'b', '/', '#', '/', 'c',
		1,  // QoS
		0,  // topic name MSB
		10, // topic name LSB
		'/', 'a', '/', 'b', '/', '#', '/', 'c', 'd', 'd',
		2, // QoS
	}

	msg := NewSubscribeMessage()
	n, err := msg.Decode(msgBytes)

	require.NoError(t, err)
	require.Equal(t, len(msgBytes), n)

	dst := make([]byte, msg.Len())
	n2, err := msg.Encode(dst)

	require.NoError(t, err)
	require.Equal(t, len(msgBytes), n2)
	require.Equal(t, msgBytes, dst[:n2])

	n3, err := msg.Decode(dst)

	require.NoError(t, err)
	require.Equal(t, len(msgBytes), n3)
}

func BenchmarkSubscribeEncode(b *testing.B) {
	msg := NewSubscribeMessage()
	msg.PacketId = 7
	msg.Subscriptions = []Subscription{
		{[]byte("t"), 0},
	}

	buf := make([]byte, msg.Len())

	for i := 0; i < b.N; i++ {
		_, err := msg.Encode(buf)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkSubscribeDecode(b *testing.B) {
	msgBytes := []byte{
		byte(SUBSCRIBE<<4) | 2,
		6,
		0, // packet ID MSB
		1, // packet ID LSB
		0, // topic name MSB
		1, // topic name LSB
		't',
		0, // QoS
	}

	msg := NewSubscribeMessage()

	for i := 0; i < b.N; i++ {
		_, err := msg.Decode(msgBytes)
		if err != nil {
			panic(err)
		}
	}
}
