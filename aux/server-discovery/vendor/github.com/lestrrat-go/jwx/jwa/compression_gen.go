// this file was auto-generated by internal/cmd/gentypes/main.go: DO NOT EDIT

package jwa

import (
	"fmt"
	"sort"
	"sync"

	"github.com/pkg/errors"
)

// CompressionAlgorithm represents the compression algorithms as described in https://tools.ietf.org/html/rfc7518#section-7.3
type CompressionAlgorithm string

// Supported values for CompressionAlgorithm
const (
	Deflate    CompressionAlgorithm = "DEF" // DEFLATE (RFC 1951)
	NoCompress CompressionAlgorithm = ""    // No compression
)

var allCompressionAlgorithms = map[CompressionAlgorithm]struct{}{
	Deflate:    {},
	NoCompress: {},
}

var listCompressionAlgorithmOnce sync.Once
var listCompressionAlgorithm []CompressionAlgorithm

// CompressionAlgorithms returns a list of all available values for CompressionAlgorithm
func CompressionAlgorithms() []CompressionAlgorithm {
	listCompressionAlgorithmOnce.Do(func() {
		listCompressionAlgorithm = make([]CompressionAlgorithm, 0, len(allCompressionAlgorithms))
		for v := range allCompressionAlgorithms {
			listCompressionAlgorithm = append(listCompressionAlgorithm, v)
		}
		sort.Slice(listCompressionAlgorithm, func(i, j int) bool {
			return string(listCompressionAlgorithm[i]) < string(listCompressionAlgorithm[j])
		})
	})
	return listCompressionAlgorithm
}

// Accept is used when conversion from values given by
// outside sources (such as JSON payloads) is required
func (v *CompressionAlgorithm) Accept(value interface{}) error {
	var tmp CompressionAlgorithm
	if x, ok := value.(CompressionAlgorithm); ok {
		tmp = x
	} else {
		var s string
		switch x := value.(type) {
		case fmt.Stringer:
			s = x.String()
		case string:
			s = x
		default:
			return errors.Errorf(`invalid type for jwa.CompressionAlgorithm: %T`, value)
		}
		tmp = CompressionAlgorithm(s)
	}
	if _, ok := allCompressionAlgorithms[tmp]; !ok {
		return errors.Errorf(`invalid jwa.CompressionAlgorithm value`)
	}

	*v = tmp
	return nil
}

// String returns the string representation of a CompressionAlgorithm
func (v CompressionAlgorithm) String() string {
	return string(v)
}
