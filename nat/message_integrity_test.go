package nat

import (
	"math/rand"
	"testing"

	"github.com/golang/glog"
)

func TestCheckMsgIntegrity(t *testing.T) {
	tests := []struct {
		desc     string
		bytes    []byte
		username string
		password string
		realm    string
	}{
		{
			"RFC-5769 Valid - short term",
			RFC5769SampleRequest,
			"evtj:h6vY",
			"VOkJxbRl1RmTxUk/WvJxBt",
			"",
		},
		{
			"RFC-5769 Valid - long term",
			RFC5769SampleRequestWithLongTermCreds,
			"\u30de\u30c8\u30ea\u30c3\u30af\u30b9",
			"TheMatrIX",
			"example.org",
		},
	}

	for _, tc := range tests {
		mp := &MessageParser{DefaultRegistry}

		check := func(bytes []byte, isCorrupt bool) bool {
			msg, err := mp.Parse(bytes)
			if err != nil {
				t.Log(isCorrupt, bytes, err)
				return false
			}

			mia, ok := msg.MessageIntegrity()
			if !ok {
				//...this can happen if we accidentially modify the MESSAGE-INTEGRITY bytes.
				t.Log(msg)
				t.Logf("[%s] MESSAGE-INTEGRITY not found:", tc.desc)
				return false
			}

			corrupt := mia.IsCorrupt(tc.username, tc.password, tc.realm)
			if corrupt != isCorrupt {
				t.Log(bytes)
				t.Fatalf("[%s] wrong corrupt result; got:%v want:%v", tc.desc, corrupt, isCorrupt)
			}

			return true
		}

		// The bytes in the testcase should be an OK STUN message. First
		// verify that this is true.
		if !check(tc.bytes, false) {
			t.Fatal("check failed to parse bytes")
		}

		//...now flip some random bits and verify that we detect it.
		attempts := 0
		for attempts < 1000 {
			glog.Info(attempts, "--------------------------------------------")

			tmp := make([]byte, len(tc.bytes))
			copy(tmp, tc.bytes)

			//...we assume there is a fingerprint at the end and since we only
			// check for corruption up and until the fingerprint we should not
			// flip any bytes outside of this range.
			max := uint16(len(tmp)) - TLVHeaderSize - FingerprintAttributeSize

			// Pick a byte to modify but don't modify the message length bytes
			// since changes to these are not possible to detect due to how
			// FINGERPRINT and MESSAGE-INTEGRITY require the length to match
			// certain criteria.
			at := rand.Int31n(int32(max))
			for at == 2 || at == 3 {
				at = rand.Int31n(int32(max))
			}
			curr := tmp[at]

			// Pick a new value for the byte.
			mod := byte(rand.Int31())
			for mod == curr {
				mod = byte(rand.Int31())
			}

			t.Logf("%d flip; at:%d to:%d from:%d", attempts, at, mod, curr)
			tmp[at] = mod

			if check(tmp, true) {
				// Flipping random bytes may cause the parser to fail which
				// is OK. We run this test until we have reached the goal attempts.
				attempts++
			}
		}
	}
}
