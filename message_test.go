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

var (
	lpstrings []string = []string{
		"this is a test",
		"hope it succeeds",
		"but just in case",
		"send me your millions",
		"",
	}

	lpstringBytes []byte = []byte{
		0x0, 0xe, 't', 'h', 'i', 's', ' ', 'i', 's', ' ', 'a', ' ', 't', 'e', 's', 't',
		0x0, 0x10, 'h', 'o', 'p', 'e', ' ', 'i', 't', ' ', 's', 'u', 'c', 'c', 'e', 'e', 'd', 's',
		0x0, 0x10, 'b', 'u', 't', ' ', 'j', 'u', 's', 't', ' ', 'i', 'n', ' ', 'c', 'a', 's', 'e',
		0x0, 0x15, 's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'y', 'o', 'u', 'r', ' ', 'm', 'i', 'l', 'l', 'i', 'o', 'n', 's',
		0x0, 0x0,
	}
)

func TestReadLPBytes(t *testing.T) {
	total := 0

	for _, str := range lpstrings {
		b, n, err := readLPBytes(lpstringBytes[total:])

		require.NoError(t, err)
		require.Equal(t, str, string(b))
		require.Equal(t, len(str)+2, n)

		total += n
	}
}

func TestWriteLPBytes(t *testing.T) {
	total := 0
	buf := make([]byte, 1000)

	for _, str := range lpstrings {
		n, err := writeLPBytes(buf[total:], []byte(str))

		require.NoError(t, err)
		require.Equal(t, 2+len(str), n)

		total += n
	}

	require.Equal(t, lpstringBytes, buf[:total])
}

func TestMessageTypes(t *testing.T) {
	if CONNECT != 1 ||
		CONNACK != 2 ||
		PUBLISH != 3 ||
		PUBACK != 4 ||
		PUBREC != 5 ||
		PUBREL != 6 ||
		PUBCOMP != 7 ||
		SUBSCRIBE != 8 ||
		SUBACK != 9 ||
		UNSUBSCRIBE != 10 ||
		UNSUBACK != 11 ||
		PINGREQ != 12 ||
		PINGRESP != 13 ||
		DISCONNECT != 14 {

		t.Errorf("Message types have invalid code")
	}
}

func TestQosCodes(t *testing.T) {
	if QosAtMostOnce != 0 || QosAtLeastOnce != 1 || QosExactlyOnce != 2 {
		t.Errorf("QOS codes invalid")
	}
}

func TestConnackReturnCodes(t *testing.T) {
	require.Equal(t, ErrInvalidProtocolVersion.Error(), ConnackCode(1).Error(), "Incorrect ConnackCode error value.")

	require.Equal(t, ErrIdentifierRejected.Error(), ConnackCode(2).Error(), "Incorrect ConnackCode error value.")

	require.Equal(t, ErrServerUnavailable.Error(), ConnackCode(3).Error(), "Incorrect ConnackCode error value.")

	require.Equal(t, ErrBadUsernameOrPassword.Error(), ConnackCode(4).Error(), "Incorrect ConnackCode error value.")

	require.Equal(t, ErrNotAuthorized.Error(), ConnackCode(5).Error(), "Incorrect ConnackCode error value.")
}

func TestFixedHeaderFlags(t *testing.T) {
	type detail struct {
		name  string
		flags byte
	}

	details := map[MessageType]detail{
		RESERVED:    detail{"RESERVED", 0},
		CONNECT:     detail{"CONNECT", 0},
		CONNACK:     detail{"CONNACK", 0},
		PUBLISH:     detail{"PUBLISH", 0},
		PUBACK:      detail{"PUBACK", 0},
		PUBREC:      detail{"PUBREC", 0},
		PUBREL:      detail{"PUBREL", 2},
		PUBCOMP:     detail{"PUBCOMP", 0},
		SUBSCRIBE:   detail{"SUBSCRIBE", 2},
		SUBACK:      detail{"SUBACK", 0},
		UNSUBSCRIBE: detail{"UNSUBSCRIBE", 2},
		UNSUBACK:    detail{"UNSUBACK", 0},
		PINGREQ:     detail{"PINGREQ", 0},
		PINGRESP:    detail{"PINGRESP", 0},
		DISCONNECT:  detail{"DISCONNECT", 0},
		RESERVED2:   detail{"RESERVED2", 0},
	}

	for m, d := range details {
		if m.Name() != d.name {
			t.Errorf("Name mismatch. Expecting %s, got %s.", d.name, m.Name())
		}

		if m.DefaultFlags() != d.flags {
			t.Errorf("Flag mismatch for %s. Expecting %d, got %d.", m.Name(), d.flags, m.DefaultFlags())
		}
	}
}

func TestSupportedVersions(t *testing.T) {
	for k, v := range SupportedVersions {
		if k == 0x03 && v != "MQIsdp" {
			t.Errorf("Protocol version and name mismatch. Expect %s, got %s.", "MQIsdp", v)
		}
	}
}
