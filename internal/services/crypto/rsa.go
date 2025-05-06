package crypto

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
)

type RSA interface {
	GenerateKeys(bitlen int) (*PublicKey, *PrivateKey, error)
	Encrypt(pub *PublicKey, m []byte) ([]byte, error)
	Decrypt(priv *PrivateKey, c []byte) ([]byte, error)
}

type RSAService struct {
}

type PublicKey struct {
	N *big.Int
	E *big.Int
}

type PrivateKey struct {
	N *big.Int
	D *big.Int
}

func NewRsaService() *RSAService {
	return &RSAService{}
}

func (r *RSAService) GenerateKeys(bitlen int) (*PublicKey, *PrivateKey, error) {
	numRetries := 0

	for {
		numRetries++
		if numRetries == 10 {
			panic("retrying too many times, something is wrong")
		}

		p, err := rand.Prime(rand.Reader, bitlen/2)
		if err != nil {
			return nil, nil, err
		}
		q, err := rand.Prime(rand.Reader, bitlen/2)
		if err != nil {
			return nil, nil, err
		}

		// n is pq
		n := new(big.Int).Set(p)
		n.Mul(n, q)

		if n.BitLen() != bitlen {
			continue
		}

		// theta(n) = (p-1)(q-1)
		p.Sub(p, big.NewInt(1))
		q.Sub(q, big.NewInt(1))
		totient := new(big.Int).Set(p)
		totient.Mul(totient, q)

		// e as recommended by PKCS#1 (RFC 2313)
		e := big.NewInt(65537)

		d := new(big.Int).ModInverse(e, totient)
		if d == nil {
			continue
		}

		pub := &PublicKey{N: n, E: e}
		priv := &PrivateKey{N: n, D: d}
		return pub, priv, nil
	}
}

func (r *RSAService) Encrypt(pub *PublicKey, m []byte) ([]byte, error) {
	// Compute length of key in bytes, rounding up.
	keyLen := (pub.N.BitLen() + 7) / 8
	if len(m) > keyLen-11 {
		return nil, fmt.Errorf("len(m)=%v, too long", len(m))
	}

	// Following RFC 2313, using block type 02 as recommended for encryption:
	// EB = 00 || 02 || PS || 00 || D
	psLen := keyLen - len(m) - 3
	eb := make([]byte, keyLen)
	eb[0] = 0x00
	eb[1] = 0x02

	// Fill PS with random non-zero bytes.
	for i := 2; i < 2+psLen; {
		_, err := rand.Read(eb[i : i+1])
		if err != nil {
			return nil, err
		}
		if eb[i] != 0x00 {
			i++
		}
	}
	eb[2+psLen] = 0x00

	// Copy the message m into the rest of the encryption block.
	copy(eb[3+psLen:], m)

	// Now the encryption block is complete; we take it as a m-byte big.Int and
	// RSA-encrypt it with the public key.
	mnum := new(big.Int).SetBytes(eb)
	c := r.encrypt(pub, mnum)

	padLen := keyLen - len(c.Bytes())
	for i := 0; i < padLen; i++ {
		eb[i] = 0x00
	}
	copy(eb[padLen:], c.Bytes())
	return eb, nil
}

func (r *RSAService) encrypt(pub *PublicKey, c *big.Int) *big.Int {
	m := new(big.Int)
	m.Exp(c, pub.E, pub.N)
	return m
}

func (r *RSAService) Decrypt(priv *PrivateKey, c []byte) ([]byte, error) {
	keyLen := (priv.N.BitLen() + 7) / 8
	if len(c) != keyLen {
		return nil, fmt.Errorf("len(c)=%v, want keyLen=%v", len(c), keyLen)
	}

	// Convert c into a bit.Int and decrypt it using the private key.
	cnum := new(big.Int).SetBytes(c)
	mnum := r.decrypt(priv, cnum)

	// Write the bytes of mnum into m, left-padding if needed.
	m := make([]byte, keyLen)
	copy(m[keyLen-len(mnum.Bytes()):], mnum.Bytes())

	// Expect proper block 02 beginning.
	if m[0] != 0x00 {
		return nil, fmt.Errorf("m[0]=%v, want 0x00", m[0])
	}
	if m[1] != 0x02 {
		return nil, fmt.Errorf("m[1]=%v, want 0x02", m[1])
	}

	// Skip over random padding until a 0x00 byte is reached. +2 adjusts the index
	// back to the full slice.
	endPad := bytes.IndexByte(m[2:], 0x00) + 2
	if endPad < 2 {
		return nil, fmt.Errorf("end of padding not found")
	}

	return m[endPad+1:], nil
}

func (r *RSAService) decrypt(priv *PrivateKey, m *big.Int) *big.Int {
	c := new(big.Int)
	c.Exp(m, priv.D, priv.N)
	return c
}
