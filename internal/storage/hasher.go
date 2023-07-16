package storage

import (
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/cshum/imagor/internal/imagorpath"
	"github.com/cshum/imagor/pkg/params"
)

// StorageHasher define image key for storage
type StorageHasher interface {
	Hash(image string) string
}

// ResultStorageHasher define key for result storage
type ResultStorageHasher interface {
	HashResult(p params.Params) string
}

// StorageHasherFunc StorageHasher handler func
type StorageHasherFunc func(image string) string

// Hash implements StorageHasher interface
func (h StorageHasherFunc) Hash(image string) string {
	return h(image)
}

// ResultStorageHasherFunc ResultStorageHasher handler func
type ResultStorageHasherFunc func(p params.Params) string

// HashResult implements ResultStorageHasher interface
func (h ResultStorageHasherFunc) HashResult(p params.Params) string {
	return h(p)
}

func hexDigestPath(path string) string {
	var digest = sha1.Sum([]byte(path))
	var hash = hex.EncodeToString(digest[:])
	return hash[:2] + "/" + hash[2:4] + "/" + hash[4:]
}

// DigestStorageHasher StorageHasher using SHA digest
var DigestStorageHasher = StorageHasherFunc(hexDigestPath)

// DigestResultStorageHasher  ResultStorageHasher using SHA digest
var DigestResultStorageHasher = ResultStorageHasherFunc(func(p params.Params) string {
	if p.Path == "" {
		p.Path = imagorpath.BuildUnsafe(p)
	}
	return hexDigestPath(p.Path)
})

// SuffixResultStorageHasher  ResultStorageHasher using storage path with digest suffix
var SuffixResultStorageHasher = ResultStorageHasherFunc(func(p params.Params) string {
	if p.Path == "" {
		p.Path = imagorpath.BuildUnsafe(p)
	}
	var digest = sha1.Sum([]byte(p.Path))
	var hash = "." + hex.EncodeToString(digest[:])[:20]
	var dotIdx = strings.LastIndex(p.Image, ".")
	var slashIdx = strings.LastIndex(p.Image, "/")
	if dotIdx > -1 && slashIdx < dotIdx {
		ext := p.Image[dotIdx:]
		if p.Meta {
			ext = ".json"
		} else {
			for _, filter := range p.Filters {
				if filter.Name == "format" {
					ext = "." + filter.Args
				}
			}
		}
		return p.Image[:dotIdx] + hash + ext // /abc/def.{digest}.jpg
	}
	return p.Image + hash // /abc/def.{digest}
})

// SizeSuffixResultStorageHasher  ResultStorageHasher using storage path with digest and size suffix
var SizeSuffixResultStorageHasher = ResultStorageHasherFunc(func(p params.Params) string {
	if p.Path == "" {
		p.Path = imagorpath.BuildUnsafe(p)
	}
	var digest = sha1.Sum([]byte(p.Path))
	var hash = "." + hex.EncodeToString(digest[:])[:20]
	if p.Width != 0 || p.Height != 0 {
		hash += "_" + strconv.Itoa(p.Width) + "x" + strconv.Itoa(p.Height)
	}
	var dotIdx = strings.LastIndex(p.Image, ".")
	var slashIdx = strings.LastIndex(p.Image, "/")
	if dotIdx > -1 && slashIdx < dotIdx {
		ext := p.Image[dotIdx:]
		if p.Meta {
			ext = ".json"
		} else {
			for _, filter := range p.Filters {
				if filter.Name == "format" {
					ext = "." + filter.Args
				}
			}
		}
		return p.Image[:dotIdx] + hash + ext // /abc/def.{digest}_{width}x{height}.jpg
	}
	return p.Image + hash // /abc/def.{digest}_{width}x{height}
})
