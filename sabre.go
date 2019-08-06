package main

import "os"
import "fmt"

func main() {
	var mode rune
	args := os.Args[1:]
	if len(args) > 0 {
		if(args[0] == "--help" || args[0] == "-help" || args[0] == "-h") {
			mode = 'h'
		}
		if(args[0] == "e" || args[0] == "encrypt" || args[0] == "-e") {
			mode = 'e'
		}
		if(args[0] == "d" || args[0] == "decrypt" || args[0] == "-d") {
			mode = 'd'
		}
		switch mode {
		case 'e':
			fmt.Println("Encrypt")
		case 'd':
			fmt.Println("Decrypt")
		case 'h':
			fmt.Println(`
sabre                          User Commands                          sabre
NAME
       sabre - encrypt or decrypt message using RC4 stream cypher
SYNOPSIS
       sabre [OPTION] [MESSAGE/FILENAME]
DESCRIPTION
       Sabre uses a custom golang implementation of RC4 to encrypt or or decrypt
       messages supplied to the program via a message string in the command line
       or via the -t flag and a filename to be (en/de)crypted

       -e, --encrypt
              take the supplied [string] or [filename] and encrypt the contents
              printing outputting the results to the console or to [filename].enc

       -d, --decrypt
              take the supplied [string] or [filename].enc and decrypt the contents
              printing outputting the results to the console or to [filename]
       -t, --target
              when this flag is passed sabre will en/decrypt the contents of the
              supplied [FILENAME]
       -h, --help
              prints this man page`)
		}
	}
}
