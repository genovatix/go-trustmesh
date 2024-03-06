package common

import (
	"bytes"
	"encoding/binary"
	"github.com/davecgh/go-spew/spew"
	"github.com/zeebo/blake3"
	"sync"
	"time"
)

type Nonce struct {
	Address   string
	value     uint32
	hash      []byte
	Timestamp int64
}

type NonceMapEntry struct {
	Value     uint32
	Hash      []byte
	Timestamp int64
}

// Maps an address to its nonce information
var (
	nonces      = make(map[string]NonceMapEntry)
	noncesMutex sync.RWMutex
)

func (n Nonce) Reset() {
	n.value = 0
	n.hash = nil
}

// GenerateOrUpdateNonce creates or updates a nonce for the given address.
func GenerateOrUpdateNonce(address string) *Nonce {
	noncesMutex.Lock()
	defer noncesMutex.Unlock()

	// Check if a nonce already exists for the address
	if entry, exists := nonces[address]; exists {
		// If nonce exists, simply return it (or update it based on your requirements)
		return &Nonce{
			Address:   address,
			value:     entry.Value,
			hash:      entry.Hash,
			Timestamp: entry.Timestamp,
		}
	}

	// Generate a new nonce for the address
	value := uint32(1) // Example value generation, consider using a more secure method
	hash := generateNonceHash(value)
	timestamp := time.Now().Unix()

	nonces[address] = NonceMapEntry{Value: value, Hash: hash, Timestamp: timestamp}

	spew.Dump(nonces)

	return &Nonce{
		Address:   address,
		value:     value,
		hash:      hash,
		Timestamp: timestamp,
	}
}

// ValidateNonce checks if a nonce associated with the address is valid.
func ValidateNonce(address string, nonce Nonce) bool {
	noncesMutex.RLock()
	defer noncesMutex.RUnlock()

	if entry, exists := nonces[address]; exists {
		return nonce.value == entry.Value && bytes.Equal(nonce.hash, entry.Hash) && time.Now().Unix()-entry.Timestamp <= nonceLifetime
	}
	return false
}

const nonceLifetime = 3600 // Define a suitable nonce lifetime

func pruneExpiredNonces() {
	noncesMutex.Lock()
	defer noncesMutex.Unlock()

	currentTimestamp := time.Now().Unix()
	for address, entry := range nonces {
		if currentTimestamp-entry.Timestamp > nonceLifetime {
			delete(nonces, address)
		}
	}
}

// generateNonceHash generates a hash for a given nonce value using Blake3.
func generateNonceHash(value uint32) []byte {
	// Convert the nonce value to a byte slice
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, value)

	// Generate and return the hash
	h := blake3.Sum256(buf)
	return h[:]
}
