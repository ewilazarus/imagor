package signer

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"hash"
)

type signer struct {
	secret   []byte
	alg      func() hash.Hash
	truncate int
}

func (s *signer) Sign(path string) string {
	h := hmac.New(s.alg, s.secret)
	h.Write([]byte(path))

	sig := base64.URLEncoding.EncodeToString(h.Sum(nil))
	if s.truncate > 0 && len(sig) > s.truncate {
		return sig[:s.truncate]
	}

	return sig
}

func New(secret string, alg func() hash.Hash, truncate int) *signer {
	return &signer{
		secret:   []byte(secret),
		alg:      alg,
		truncate: 0,
	}
}

func NewDefault(secret string) *signer {
	return New(secret, sha1.New, 0)
}
