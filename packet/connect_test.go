package packet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectInterface(t *testing.T) {
	pkt := NewConnect()
	pkt.Will = &Message{
		Topic:   "w",
		Payload: []byte("m"),
		QOS:     QOSAtLeastOnce,
	}

	assert.Equal(t, pkt.Type(), CONNECT)
	assert.Equal(t, "<Connect ClientID=\"\" KeepAlive=0 Username=\"\" Password=\"\" CleanSession=true Will=<Message Topic=\"w\" QOS=1 Retain=false Payload=6d> Version=4>", pkt.String())
}

func TestConnectDecode1(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		58,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		4,   // Protocol Level
		206, // Connect Flags
		0,   // Keep Alive MSB
		10,  // Keep Alive LSB
		0,   // Client ID MSB
		6,   // Client ID LSB
		'g', 'o', 'm', 'q', 't', 't',
		0, // Will Topic MSB
		4, // Will Topic LSB
		'w', 'i', 'l', 'l',
		0,  // Will Message MSB
		12, // Will Message LSB
		's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e',
		0, // Username ID MSB
		6, // Username ID LSB
		'g', 'o', 'm', 'q', 't', 't',
		0,  // Password ID MSB
		10, // Password ID LSB
		'v', 'e', 'r', 'y', 's', 'e', 'c', 'r', 'e', 't',
	}

	pkt := NewConnect()
	n, err := pkt.Decode(pktBytes)

	assert.NoError(t, err)
	assert.Equal(t, len(pktBytes), n)
	assert.Equal(t, uint16(10), pkt.KeepAlive)
	assert.Equal(t, "gomqtt", pkt.ClientID)
	assert.Equal(t, "will", pkt.Will.Topic)
	assert.Equal(t, []byte("send me home"), pkt.Will.Payload)
	assert.Equal(t, "gomqtt", pkt.Username)
	assert.Equal(t, "verysecret", pkt.Password)
	assert.Equal(t, Version311, pkt.Version)
}

func TestConnectDecode2(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		58,
		0, // Protocol String MSB
		6, // Protocol String LSB
		'M', 'Q', 'I', 's', 'd', 'p',
		3,   // Protocol Level
		206, // Connect Flags
		0,   // Keep Alive MSB
		10,  // Keep Alive LSB
		0,   // Client ID MSB
		6,   // Client ID LSB
		'g', 'o', 'm', 'q', 't', 't',
		0, // Will Topic MSB
		4, // Will Topic LSB
		'w', 'i', 'l', 'l',
		0,  // Will Message MSB
		12, // Will Message LSB
		's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e',
		0, // Username ID MSB
		6, // Username ID LSB
		'g', 'o', 'm', 'q', 't', 't',
		0,  // Password ID MSB
		10, // Password ID LSB
		'v', 'e', 'r', 'y', 's', 'e', 'c', 'r', 'e', 't',
	}

	pkt := NewConnect()
	n, err := pkt.Decode(pktBytes)

	assert.NoError(t, err)
	assert.Equal(t, len(pktBytes), n)
	assert.Equal(t, uint16(10), pkt.KeepAlive)
	assert.Equal(t, "gomqtt", pkt.ClientID)
	assert.Equal(t, "will", pkt.Will.Topic)
	assert.Equal(t, []byte("send me home"), pkt.Will.Payload)
	assert.Equal(t, "gomqtt", pkt.Username)
	assert.Equal(t, "verysecret", pkt.Password)
	assert.Equal(t, Version31, pkt.Version)
}

func TestConnectDecodeError1(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		60, // < wrong remaining length
		0,  // Protocol String MSB
		5,  // Protocol String LSB
		'M', 'Q', 'T', 'T',
	}

	pkt := NewConnect()
	_, err := pkt.Decode(pktBytes)

	assert.Error(t, err)
}

func TestConnectDecodeError2(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		6,
		0, // Protocol String MSB
		5, // Protocol String LSB < invalid size
		'M', 'Q', 'T', 'T',
	}

	pkt := NewConnect()
	_, err := pkt.Decode(pktBytes)

	assert.Error(t, err)
}

func TestConnectDecodeError3(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		6,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		// Protocol Level < missing
	}

	pkt := NewConnect()
	_, err := pkt.Decode(pktBytes)

	assert.Error(t, err)
}

func TestConnectDecodeError4(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		8,
		0,                       // Protocol String MSB
		5,                       // Protocol String LSB
		'M', 'Q', 'T', 'T', 'X', // < wrong protocol string
		4,
	}

	pkt := NewConnect()
	_, err := pkt.Decode(pktBytes)

	assert.Error(t, err)
}

func TestConnectDecodeError5(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		7,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		5, // Protocol Level < wrong id
	}

	pkt := NewConnect()
	_, err := pkt.Decode(pktBytes)

	assert.Error(t, err)
}

func TestConnectDecodeError6(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		7,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		4, // Protocol Level
		// Connect Flags < missing
	}

	pkt := NewConnect()
	_, err := pkt.Decode(pktBytes)

	assert.Error(t, err)
}

func TestConnectDecodeError7(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		7,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		4, // Protocol Level
		1, // Connect Flags < reserved bit set to one
	}

	pkt := NewConnect()
	_, err := pkt.Decode(pktBytes)

	assert.Error(t, err)
}

func TestConnectDecodeError8(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		7,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		4,  // Protocol Level
		24, // Connect Flags < invalid qos
	}

	pkt := NewConnect()
	_, err := pkt.Decode(pktBytes)

	assert.Error(t, err)
}

func TestConnectDecodeError9(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		7,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		4, // Protocol Level
		8, // Connect Flags < will flag set to zero but others not
	}

	pkt := NewConnect()
	_, err := pkt.Decode(pktBytes)

	assert.Error(t, err)
}

func TestConnectDecodeError10(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		7,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		4,  // Protocol Level
		64, // Connect Flags < password flag set but not username
	}

	pkt := NewConnect()
	_, err := pkt.Decode(pktBytes)

	assert.Error(t, err)
}

func TestConnectDecodeError11(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		7,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		4, // Protocol Level
		0, // Connect Flags
		0, // Keep Alive MSB < missing
		// Keep Alive LSB < missing
	}

	pkt := NewConnect()
	_, err := pkt.Decode(pktBytes)

	assert.Error(t, err)
}

func TestConnectDecodeError12(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		7,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		4, // Protocol Level
		0, // Connect Flags
		0, // Keep Alive MSB
		1, // Keep Alive LSB
		0, // Client ID MSB
		2, // Client ID LSB < wrong size
		'x',
	}

	pkt := NewConnect()
	_, err := pkt.Decode(pktBytes)

	assert.Error(t, err)
}

func TestConnectDecodeError13(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		6,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		4, // Protocol Level
		0, // Connect Flags < clean session false
		0, // Keep Alive MSB
		1, // Keep Alive LSB
		0, // Client ID MSB
		0, // Client ID LSB
	}

	pkt := NewConnect()
	_, err := pkt.Decode(pktBytes)

	assert.Error(t, err)
}

func TestConnectDecodeError14(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		6,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		4, // Protocol Level
		6, // Connect Flags
		0, // Keep Alive MSB
		1, // Keep Alive LSB
		0, // Client ID MSB
		0, // Client ID LSB
		0, // Will Topic MSB
		1, // Will Topic LSB < wrong size
	}

	pkt := NewConnect()
	_, err := pkt.Decode(pktBytes)

	assert.Error(t, err)
}

func TestConnectDecodeError15(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		6,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		4, // Protocol Level
		6, // Connect Flags
		0, // Keep Alive MSB
		1, // Keep Alive LSB
		0, // Client ID MSB
		0, // Client ID LSB
		0, // Will Topic MSB
		0, // Will Topic LSB
		0, // Will Payload MSB
		1, // Will Payload LSB < wrong size
	}

	pkt := NewConnect()
	_, err := pkt.Decode(pktBytes)

	assert.Error(t, err)
}

func TestConnectDecodeError16(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		6,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		4,   // Protocol Level
		194, // Connect Flags
		0,   // Keep Alive MSB
		1,   // Keep Alive LSB
		0,   // Client ID MSB
		0,   // Client ID LSB
		0,   // Username MSB
		1,   // Username LSB < wrong size
	}

	pkt := NewConnect()
	_, err := pkt.Decode(pktBytes)

	assert.Error(t, err)
}

func TestConnectDecodeError17(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		6,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		4,   // Protocol Level
		194, // Connect Flags
		0,   // Keep Alive MSB
		1,   // Keep Alive LSB
		0,   // Client ID MSB
		0,   // Client ID LSB
		0,   // Username MSB
		0,   // Username LSB
		0,   // Password MSB
		1,   // Password LSB < wrong size
	}

	pkt := NewConnect()
	_, err := pkt.Decode(pktBytes)

	assert.Error(t, err)
}

func TestConnectEncode1(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		58,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		4,   // Protocol level 4
		204, // Connect Flags
		0,   // Keep Alive MSB
		10,  // Keep Alive LSB
		0,   // Client ID MSB
		6,   // Client ID LSB
		'g', 'o', 'm', 'q', 't', 't',
		0, // Will Topic MSB
		4, // Will Topic LSB
		'w', 'i', 'l', 'l',
		0,  // Will Message MSB
		12, // Will Message LSB
		's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e',
		0, // Username ID MSB
		6, // Username ID LSB
		'g', 'o', 'm', 'q', 't', 't',
		0,  // Password ID MSB
		10, // Password ID LSB
		'v', 'e', 'r', 'y', 's', 'e', 'c', 'r', 'e', 't',
	}

	pkt := NewConnect()
	pkt.Will = &Message{
		QOS:     QOSAtLeastOnce,
		Topic:   "will",
		Payload: []byte("send me home"),
	}
	pkt.CleanSession = false
	pkt.ClientID = "gomqtt"
	pkt.KeepAlive = 10
	pkt.Username = "gomqtt"
	pkt.Password = "verysecret"

	dst := make([]byte, pkt.Len())
	n, err := pkt.Encode(dst)

	assert.NoError(t, err)
	assert.Equal(t, len(pktBytes), n)
	assert.Equal(t, pktBytes, dst[:n])
}

func TestConnectEncode2(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		14,
		0, // Protocol String MSB
		6, // Protocol String LSB
		'M', 'Q', 'I', 's', 'd', 'p',
		3,  // Protocol level 4
		2,  // Connect Flags
		0,  // Keep Alive MSB
		10, // Keep Alive LSB
		0,  // Client ID MSB
		0,  // Client ID LSB
	}

	pkt := NewConnect()
	pkt.CleanSession = true
	pkt.KeepAlive = 10
	pkt.Version = Version31

	dst := make([]byte, pkt.Len())
	n, err := pkt.Encode(dst)

	assert.NoError(t, err)
	assert.Equal(t, len(pktBytes), n)
	assert.Equal(t, pktBytes, dst[:n])
}

func TestConnectEncode3(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		12,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		4,  // Protocol level 4
		2,  // Connect Flags
		0,  // Keep Alive MSB
		10, // Keep Alive LSB
		0,  // Client ID MSB
		0,  // Client ID LSB
	}

	pkt := NewConnect()
	pkt.CleanSession = true
	pkt.KeepAlive = 10
	pkt.Version = 0

	dst := make([]byte, pkt.Len())
	n, err := pkt.Encode(dst)

	assert.NoError(t, err)
	assert.Equal(t, len(pktBytes), n)
	assert.Equal(t, pktBytes, dst[:n])
}

func TestConnectEncodeError1(t *testing.T) {
	pkt := NewConnect()

	dst := make([]byte, 4) // < too small buffer
	n, err := pkt.Encode(dst)

	assert.Error(t, err)
	assert.Equal(t, 0, n)
}

func TestConnectEncodeError2(t *testing.T) {
	pkt := NewConnect()
	pkt.Will = &Message{
		Topic: "t",
		QOS:   3, // < wrong qos
	}

	dst := make([]byte, pkt.Len())
	n, err := pkt.Encode(dst)

	assert.Error(t, err)
	assert.Equal(t, 9, n)
}

func TestConnectEncodeError3(t *testing.T) {
	pkt := NewConnect()
	pkt.ClientID = string(make([]byte, 65536)) // < too big

	dst := make([]byte, pkt.Len())
	n, err := pkt.Encode(dst)

	assert.Error(t, err)
	assert.Equal(t, 14, n)
}

func TestConnectEncodeError4(t *testing.T) {
	pkt := NewConnect()
	pkt.Will = &Message{
		Topic: string(make([]byte, 65536)), // < too big
	}

	dst := make([]byte, pkt.Len())
	n, err := pkt.Encode(dst)

	assert.Error(t, err)
	assert.Equal(t, 16, n)
}

func TestConnectEncodeError5(t *testing.T) {
	pkt := NewConnect()
	pkt.Will = &Message{
		Topic:   "t",
		Payload: make([]byte, 65536), // < too big
	}

	dst := make([]byte, pkt.Len())
	n, err := pkt.Encode(dst)

	assert.Error(t, err)
	assert.Equal(t, 19, n)
}

func TestConnectEncodeError6(t *testing.T) {
	pkt := NewConnect()
	pkt.Username = string(make([]byte, 65536)) // < too big

	dst := make([]byte, pkt.Len())
	n, err := pkt.Encode(dst)

	assert.Error(t, err)
	assert.Equal(t, 16, n)
}

func TestConnectEncodeError7(t *testing.T) {
	pkt := NewConnect()
	pkt.Password = "p" // < missing username

	dst := make([]byte, pkt.Len())
	n, err := pkt.Encode(dst)

	assert.Error(t, err)
	assert.Equal(t, 14, n)
}

func TestConnectEncodeError8(t *testing.T) {
	pkt := NewConnect()
	pkt.Username = "u"
	pkt.Password = string(make([]byte, 65536)) // < too big

	dst := make([]byte, pkt.Len())
	n, err := pkt.Encode(dst)

	assert.Error(t, err)
	assert.Equal(t, 19, n)
}

func TestConnectEncodeError9(t *testing.T) {
	pkt := NewConnect()
	pkt.Will = &Message{
		// < missing topic
	}

	dst := make([]byte, pkt.Len())
	n, err := pkt.Encode(dst)

	assert.Error(t, err)
	assert.Equal(t, 9, n)
}

func TestConnectEncodeError10(t *testing.T) {
	pkt := NewConnect()
	pkt.Version = 255

	dst := make([]byte, pkt.Len())
	n, err := pkt.Encode(dst)

	assert.Error(t, err)
	assert.Equal(t, 2, n)
}

func TestConnectEncodeError11(t *testing.T) {
	pkt := NewConnect()
	pkt.CleanSession = false // < client id is empty

	dst := make([]byte, pkt.Len())
	n, err := pkt.Encode(dst)

	assert.Error(t, err)
	assert.Equal(t, 9, n)
}

func TestConnectEqualDecodeEncode(t *testing.T) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		58,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		4,   // Protocol level 4
		238, // Connect Flags
		0,   // Keep Alive MSB
		10,  // Keep Alive LSB
		0,   // Client ID MSB
		6,   // Client ID LSB
		'g', 'o', 'm', 'q', 't', 't',
		0, // Will Topic MSB
		4, // Will Topic LSB
		'w', 'i', 'l', 'l',
		0,  // Will Message MSB
		12, // Will Message LSB
		's', 'e', 'n', 'd', ' ', 'm', 'e', ' ', 'h', 'o', 'm', 'e',
		0, // Username ID MSB
		6, // Username ID LSB
		'g', 'o', 'm', 'q', 't', 't',
		0,  // Password ID MSB
		10, // Password ID LSB
		'v', 'e', 'r', 'y', 's', 'e', 'c', 'r', 'e', 't',
	}

	pkt := NewConnect()
	n, err := pkt.Decode(pktBytes)

	assert.NoError(t, err)
	assert.Equal(t, len(pktBytes), n)

	dst := make([]byte, pkt.Len())
	n2, err := pkt.Encode(dst)

	assert.NoError(t, err)
	assert.Equal(t, len(pktBytes), n2)
	assert.Equal(t, pktBytes, dst[:n2])

	n3, err := pkt.Decode(dst)

	assert.NoError(t, err)
	assert.Equal(t, len(pktBytes), n3)
}

func BenchmarkConnectEncode(b *testing.B) {
	pkt := NewConnect()
	pkt.Will = &Message{
		Topic:   "w",
		Payload: []byte("m"),
		QOS:     QOSAtLeastOnce,
	}
	pkt.CleanSession = true
	pkt.ClientID = "i"
	pkt.KeepAlive = 10
	pkt.Username = "u"
	pkt.Password = "p"

	buf := make([]byte, pkt.Len())

	for i := 0; i < b.N; i++ {
		_, err := pkt.Encode(buf)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkConnectDecode(b *testing.B) {
	pktBytes := []byte{
		byte(CONNECT << 4),
		25,
		0, // Protocol String MSB
		4, // Protocol String LSB
		'M', 'Q', 'T', 'T',
		4,   // Protocol level 4
		206, // Connect Flags
		0,   // Keep Alive MSB
		10,  // Keep Alive LSB
		0,   // Client ID MSB
		1,   // Client ID LSB
		'i',
		0, // Will Topic MSB
		1, // Will Topic LSB
		'w',
		0, // Will Message MSB
		1, // Will Message LSB
		'm',
		0, // Username ID MSB
		1, // Username ID LSB
		'u',
		0, // Password ID MSB
		1, // Password ID LSB
		'p',
	}

	pkt := NewConnect()

	for i := 0; i < b.N; i++ {
		_, err := pkt.Decode(pktBytes)
		if err != nil {
			panic(err)
		}
	}
}
