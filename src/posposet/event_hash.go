package posposet

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"strings"

	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/crypto"
	"github.com/Fantom-foundation/go-lachesis/src/rlp"
)

type (
	// EventHash is a unique identificator of Event.
	// It is a hash of Event.
	EventHash common.Hash

	// EventHashSlice is a sortable slice of EventHash.
	EventHashSlice []EventHash

	// EventHashes provides additional methods of EventHash index.
	EventHashes map[EventHash]struct{}
)

var (
	ZeroEventHash = EventHash{}
)

/*
 * EventHash methods:
 */

// EventHashOf calcs hash of event.
func EventHashOf(e *Event) EventHash {
	buf, err := rlp.EncodeToBytes(e)
	if err != nil {
		panic(err)
	}
	return EventHash(crypto.Keccak256Hash(buf))
}

// Bytes returns value as byte slice.
func (hash EventHash) Bytes() []byte {
	return (common.Hash)(hash).Bytes()
}

// String returns human readable string representation.
func (hash EventHash) String() string {
	if name, ok := EventNameDict[hash]; ok {
		return name
	}
	return (common.Hash)(hash).ShortString()
}

// IsZero returns true if hash is empty.
func (hash *EventHash) IsZero() bool {
	return *hash == EventHash{}
}

/*
 * EventHashes methods:
 */

func newEventHashes(hash ...EventHash) EventHashes {
	hh := EventHashes{}
	hh.Add(hash...)
	return hh
}

// String returns human readable string representation.
func (hh EventHashes) String() string {
	ss := make([]string, 0, len(hh))
	for hash, _ := range hh {
		ss = append(ss, hash.String())
	}
	return "[" + strings.Join(ss, ", ") + "]"
}

// All returns whole index.
func (hh EventHashes) Slice() EventHashSlice {
	arr := make(EventHashSlice, len(hh))
	i := 0
	for h := range hh {
		arr[i] = h
		i++
	}
	return arr
}

// Add appends hash to the index.
func (hh EventHashes) Add(hash ...EventHash) (changed bool) {
	for _, h := range hash {
		if _, ok := hh[h]; !ok {
			hh[h] = struct{}{}
			changed = true
		}
	}
	return
}

// Contains returns true if hash is in.
func (hh EventHashes) Contains(hash EventHash) bool {
	_, ok := hh[hash]
	return ok
}

// EncodeRLP is a specialized encoder to encode index into array.
func (hh EventHashes) EncodeRLP(w io.Writer) error {
	var arr EventHashSlice
	for h := range hh {
		arr = append(arr, h)
	}
	sort.Sort(arr)
	return rlp.Encode(w, arr)
}

// DecodeRLP is a specialized decoder to decode index from array.
func (hh *EventHashes) DecodeRLP(s *rlp.Stream) error {
	var arr EventHashSlice
	err := s.Decode(&arr)
	if err != nil {
		return err
	}

	res := EventHashes{}
	for _, h := range arr {
		if !res.Add(h) {
			return fmt.Errorf("Double value is detected")
		}
	}

	*hh = res
	return nil
}

/*
 * EventHashSlice's methods:
 */

func (hh EventHashSlice) Len() int      { return len(hh) }
func (hh EventHashSlice) Swap(i, j int) { hh[i], hh[j] = hh[j], hh[i] }
func (hh EventHashSlice) Less(i, j int) bool {
	return bytes.Compare(hh[i].Bytes(), hh[j].Bytes()) < 0
}

/*
 * Utils:
 */

func FakeEventHash() (h EventHash) {
	_, err := rand.Read(h[:])
	if err != nil {
		panic(err)
	}
	return
}

func FakeEventHashes(n int) EventHashes {
	res := EventHashes{}
	for i := 0; i < n; i++ {
		res.Add(FakeEventHash())
	}
	return res
}
