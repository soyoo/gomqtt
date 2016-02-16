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

package broker

import (
	"testing"

	"github.com/gomqtt/packet"
	"github.com/stretchr/testify/assert"
)

func abstractSessionSubscriptionTest(t *testing.T, session Session) {
	subscription := &packet.Subscription{
		Topic: "+",
		QOS:   1,
	}

	subs, err := session.AllSubscriptions()
	assert.Equal(t, 0, len(subs))

	sub, err := session.LookupSubscription("foo")
	assert.Nil(t, sub)
	assert.NoError(t, err)

	err = session.SaveSubscription(subscription)
	assert.NoError(t, err)

	sub, err = session.LookupSubscription("foo")
	assert.Equal(t, subscription, sub)
	assert.NoError(t, err)

	subs, err = session.AllSubscriptions()
	assert.Equal(t, 1, len(subs))

	err = session.DeleteSubscription("+")
	assert.NoError(t, err)

	sub, err = session.LookupSubscription("foo")
	assert.Nil(t, sub)
	assert.NoError(t, err)

	subs, err = session.AllSubscriptions()
	assert.Equal(t, 0, len(subs))
}

func abstractSessionWillTest(t *testing.T, session Session) {
	theWill := &packet.Message{"test", []byte("test"), 0, false}

	will, err := session.LookupWill()
	assert.Nil(t, will)
	assert.NoError(t, err)

	err = session.SaveWill(theWill)
	assert.NoError(t, err)

	will, err = session.LookupWill()
	assert.Equal(t, theWill, will)
	assert.NoError(t, err)

	err = session.ClearWill()
	assert.NoError(t, err)

	will, err = session.LookupWill()
	assert.Nil(t, will)
	assert.NoError(t, err)
}

func abstractBackendGetSessionTest(t *testing.T, backend Backend) {
	session1, err := backend.GetSession(nil, "foo")
	assert.NoError(t, err)
	assert.NotNil(t, session1)

	session2, err := backend.GetSession(nil, "foo")
	assert.NoError(t, err)
	assert.True(t, session1 == session2)

	session3, err := backend.GetSession(nil, "bar")
	assert.NoError(t, err)
	assert.False(t, session3 == session1)
	assert.False(t, session3 == session2)

	session4, err := backend.GetSession(nil, "")
	assert.NoError(t, err)
	assert.NotNil(t, session4)

	session5, err := backend.GetSession(nil, "")
	assert.NoError(t, err)
	assert.NotNil(t, session5)
	assert.True(t, session4 != session5)
}

func abstractBackendRetainedTest(t *testing.T, backend Backend) {
	msg1 := &packet.Message{
		Topic:   "foo",
		Payload: []byte("bar"),
		QOS:     1,
		Retain:  true,
	}

	msg2 := &packet.Message{
		Topic:   "foo/bar",
		Payload: []byte("bar"),
		QOS:     1,
		Retain:  true,
	}

	msg3 := &packet.Message{
		Topic:   "foo",
		Payload: []byte("bar"),
		QOS:     2,
		Retain:  true,
	}

	msg4 := &packet.Message{
		Topic:  "foo",
		QOS:    1,
		Retain: true,
	}

	// should be empty
	msgs, err := backend.Subscribe(nil, "foo")
	assert.NoError(t, err)
	assert.Empty(t, msgs)

	err = backend.Publish(nil, msg1)
	assert.NoError(t, err)

	// should have one
	msgs, err = backend.Subscribe(nil, "foo")
	assert.NoError(t, err)
	assert.Equal(t, msg1, msgs[0])

	err = backend.Publish(nil, msg2)
	assert.NoError(t, err)

	// should have two
	msgs, err = backend.Subscribe(nil, "#")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(msgs))

	err = backend.Publish(nil, msg3)
	assert.NoError(t, err)

	// should have another
	msgs, err = backend.Subscribe(nil, "foo")
	assert.NoError(t, err)
	assert.Equal(t, msg3, msgs[0])

	err = backend.Publish(nil, msg4)
	assert.NoError(t, err)

	// should have none
	msgs, err = backend.Subscribe(nil, "foo")
	assert.NoError(t, err)
	assert.Empty(t, msgs)
}

// store and look up subscriptions by client
// remove subscriptions by client
// store and look up subscriptions by topic
// QoS 0 subscriptions, restored but not matched
