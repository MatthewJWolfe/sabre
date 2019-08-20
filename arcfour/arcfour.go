package arcfour

import (
    "bufio"
    "encoding/hex"
    "github.com/MatthewJWolfe/sabre/util"
    "fmt"
  )
	//ARC is the stucture that holds the vital components for a valid RC4 implementation
	//Including a pointer to a PRGA (Psuedo Random Generation Algorithm)
	type ARC struct {
    iv [10]byte
		key [256]byte
		keylen uint8
		Rgen PRGA
	}
  //Init will generate 10 high entropy bytes for the IV
  //IV is appended to the user supplied key and the total length is stored
	func (a *ARC) Init(key_ []byte, iv_ ...[]byte) {
    var initial_vector []byte
    _ = initial_vector
    //get 10 high entropy bytes if not provided by caller
    if len(iv_) == 0 {
      initial_vector = util.EntropyBytes(10)
    } else {
      initial_vector = iv_[0][:]
    }
    copy(a.iv[:], initial_vector[0:10])
    key_ = append(key_, initial_vector...)

		//ensure you are passing a slice that is the exact length of the key
		a.keylen = uint8(copy(a.key[:], key_))
    // fmt.Println("Just IV: ", a.iv)
    // fmt.Println("Key + IV: ", a.key)
    // fmt.Println("Key Length:", a.keylen)
		a.Rgen = PRGA{}
		a.Rgen.KSA(a.key, a.keylen)
	}
	func (a *ARC) Encode(p_stream *bufio.Reader) ([]byte){
		var (
			err error
		  plain_byte, e_byte, k_byte byte
			enc_slice []byte
		)
		for {
			plain_byte, err = p_stream.ReadByte()
      if err == nil{
        k_byte = a.Rgen.Kstream()
  			e_byte = plain_byte ^ k_byte
  			enc_slice = append(enc_slice, e_byte)
      } else {
        break
      }
		}
    hex_encoded := make([]byte, hex.EncodedLen(len(enc_slice)))
    hex.Encode(hex_encoded, enc_slice)
    fmt.Printf("%s", hex_encoded)
		return hex_encoded
	}
	func (a *ARC) Decode(p_stream *bufio.Reader) ([]byte){
		var (
			err error
			plain_byte, d_byte, k_byte byte
			temp_slice []byte
		)
		for {
			plain_byte, err = p_stream.ReadByte()
      if err == nil {
        temp_slice = append(temp_slice, plain_byte)
      } else {
        break
      }
		}
    hex_decoded := make([]byte, hex.DecodedLen(len(temp_slice)))
    _, err = hex.Decode(hex_decoded, temp_slice[:len(temp_slice)])
    temp_slice = nil
    for _, c_byte := range hex_decoded {
      k_byte = a.Rgen.Kstream()
			d_byte =  c_byte ^ k_byte
      temp_slice = append(temp_slice, d_byte)
    }
    if err != nil {
		    fmt.Println(err)
	  }
		return temp_slice
	}
  func (a *ARC) WriteEncFile(cyphertext []byte, name string) (error) {
    //creates file with name provided by arg2, file is always in a safe empty state
    //the 10 byte IV is prepended to the cyphertext and the result is written to file
    payload := append(a.iv[:], cyphertext...)
  	return util.Dump2File(payload, name)
  }
  // func (a *ARC) ReadEncFile(cyphertext []byte, name string) (byte[10], error) {
  //   //parses a CypherSabre file (has extension .cs1)
  //   //The IV removed and returned as a byte array
  //   //The remaining cyphertext is copied into the cyphertext buffer
  // }
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
//  the key scheduling algorithm is used to initialize the array S
//  first S must be set to its Identity Permutation (IP)
//  ->in this case the IP is an array of size 256 of all possible 1 Byte values
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
