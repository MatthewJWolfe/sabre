package arcfour

import (
    "bufio"
    //"fmt"
  )

	//ARC is the stucture that holds the vital components for a valid RC4 implementation
	//Including a pointer to a PRGA (Psuedo Random Generation Algorithm)
	type ARC struct {
		key [256]byte
		keylen uint8
		Rgen PRGA
	}

	func (a *ARC) Init(key_ []byte) {
		//ensure you are passing a slice that is the exact length of the key
		a.keylen = uint8(copy(a.key[:], key_))
		a.Rgen = PRGA{}
		a.Rgen.KSA(a.key, a.keylen)
	}

	func (a *ARC) Encode(p_stream *bufio.Reader) ([]byte){
		var (
			err error
		  plain_byte, e_byte, k_byte byte
			enc_slice []byte
		)
		for err == nil {
			plain_byte, err = p_stream.ReadByte()
			k_byte = a.Rgen.Kstream()
			e_byte = plain_byte ^ k_byte
			//fmt.Printf("plain: %x, key: %x, cypher: %x\n", plain_byte, k_byte, e_byte)
			//fmt.Printf("%x", e_byte)
			enc_slice = append(enc_slice, e_byte)
		}
		return enc_slice[:len(enc_slice) - 1]
	}
	func (a *ARC) Decode(p_stream *bufio.Reader) ([]byte){
		var (
			err error
			plain_byte, d_byte, k_byte byte
			dec_slice []byte
		)
		for err == nil {
			plain_byte, err = p_stream.ReadByte()
			k_byte = a.Rgen.Kstream()
			d_byte = plain_byte ^ k_byte
			//fmt.Printf("plain: %x, key: %x, cypher: %x\n", plain_byte, k_byte, e_byte)
			//fmt.Printf("%x", e_byte)
			dec_slice = append(dec_slice, d_byte)
		}
		return dec_slice[:len(dec_slice) - 1]
	}
	//PRGA is the Psuedo Rando Generation Algorithm it contains:
	// two single byte indexes and a 256 byte state array
	type PRGA struct {
	  i, j uint8
		S [256]byte
	}

	func (gen *PRGA) Kstream() byte{
		gen.i++
		gen.j = gen.j + gen.S[gen.i]
		gen.S[gen.i], gen.S[gen.j] = gen.S[gen.j], gen.S[gen.i]
		var sum uint8 = (gen.S[gen.i] + gen.S[gen.j])
		K := gen.S[sum]
		return K
	}

	func (gen *PRGA) KSA(key [256]byte, keylen uint8) {
		//the key scheduling algorithm is used to initialize the array S
		//first S must be set to its Identity Permutation
		var i, j uint8
		for i = 0; true; i++ {
			gen.S[i] = i
			if i == 255 {
				break
			}
		}
		j = 0
		for i = 0; true; i++ {
			j = j + gen.S[i] + key[i % keylen]
			gen.S[i], gen.S[j] = gen.S[j], gen.S[i]
			if i == 255 {
				break
			}
		}
	}
