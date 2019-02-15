package util

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"hash"
)

var (
	md5_impl    = md5.New()
	sha1_impl   = sha1.New()
	sha256_impl = sha256.New()
	sha512_impl = sha512.New()
)

func Hash(mod, str string) (string, error) {
	bp := []byte(str)
	var hash_impl hash.Hash
	if "md5" == mod {
		hash_impl = md5.New()
	} else if "sha1" == mod {
		hash_impl = sha1.New()
	} else if "sha256" == mod {
		hash_impl = sha256.New()
	} else if "sha512" == mod {
		hash_impl = sha512.New()
	} else {
		return "", errors.New("mod is mot support!!!")
	}
	hash_impl.Write(bp)
	return string(hash_impl.Sum(nil)), nil
}
