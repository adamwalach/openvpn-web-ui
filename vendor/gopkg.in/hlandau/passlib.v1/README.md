passlib for go
==============

[![GoDoc](https://godoc.org/gopkg.in/hlandau/passlib.v1?status.svg)](https://godoc.org/gopkg.in/hlandau/passlib.v1) [![Build Status](https://travis-ci.org/hlandau/passlib.svg?branch=master)](https://travis-ci.org/hlandau/passlib)

[Python's passlib](https://pythonhosted.org/passlib/) is quite an amazing
library. I'm not sure there's a password library in existence with more thought
put into it, or with more support for obscure password formats.

This is a skeleton of a port of passlib to Go. It dogmatically adopts the
modular crypt format, which [passlib has excellent documentation for](https://pythonhosted.org/passlib/modular_crypt_format.html#modular-crypt-format).

Currently, it supports sha256-crypt, sha512-crypt, scrypt-sha256, bcrypt and
passlib's bcrypt-sha256 variant. By default, it will hash using scrypt-sha256
and verify existing hashes using any of these schemes.

Example Usage
-------------
There's a default context for ease of use. Most people need only concern
themselves with the functions `Hash` and `Verify`:

```go
// Hash a plaintext, UTF-8 password.
func Hash(password string) (hash string, err error)

// Verifies a plaintext, UTF-8 password using a previously derived hash.
// Returns non-nil err if verification fails.
//
// Also returns an upgraded password hash if the hash provided is
// deprecated.
func Verify(password, hash string) (newHash string, err error)
```

Here's a rough skeleton of typical usage.

```go
import "gopkg.in/hlandau/passlib.v1"

func RegisterUser() {
  (...)

  password := get a (UTF-8, plaintext) password from somewhere

  hash, err := passlib.Hash(password)
  if err != nil {
    // couldn't hash password for some reason
    return
  }

  (store hash in database, etc.)
}

func CheckPassword() bool {
  password := get the password the user entered
  hash := the hash you stored from the call to Hash()

  newHash, err := passlib.Verify(password, hash)
  if err != nil {
    // incorrect password, malformed hash, etc.
    // either way, reject
    return false
  }

  // The context has decided, as per its policy, that
  // the hash which was used to validate the password
  // should be changed. It has upgraded the hash using
  // the verified password.
  if newHash != "" {
    (store newHash in database, replacing old hash)
  }

  return true
}
```

scrypt Modular Crypt Format
---------------------------
Since scrypt does not have a pre-existing modular crypt format standard, I made one. It's as follows:

    $s2$N$r$p$salt$hash

...where `N`, `r` and `p` are the respective difficulty parameters to scrypt as positive decimal integers without leading zeroes, and `salt` and `hash` are base64-encoded binary strings. Note that the RFC 4648 base64 encoding is used (not the one used by sha256-crypt and sha512-crypt).

TODO
----

  - PBKDF2

Licence
-------
passlib is partially derived from Python's passlib and so maintains its BSD license.

    © 2008-2012 Assurance Technologies LLC.  (Python passlib)  BSD License
    © 2014 Hugo Landau <hlandau@devever.net>  BSD License

