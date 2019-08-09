package main

import (
	"os"
	"io/ioutil"
	"bufio"
	"fmt"
	"syscall"
	"golang.org/x/crypto/ssh/terminal"
	RC4 "github.com/MatthewJWolfe/sabre/arcfour"
	)

	func check(e error) {
	    if e != nil {
	        panic(e)
	    }
	}
	func askForKey(r *bufio.Reader ) ([]byte, error){
		fmt.Print("Please input your key: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		return bytePassword, err
	}
	func dump2file(d []byte, n string) (error) {
		return ioutil.WriteFile(n, d, 0644)
	}

	func main() {
		var mode rune
		ron := new(RC4.ARC)
		cli_reader := bufio.NewReader(os.Stdin)
		args := os.Args[1:]
		if len(args) > 0 {
			if(args[0] == "--help" || args[0] == "-help" || args[0] == "-h" || args[0] == "h"|| args[0] == "help") {
				mode = 'h'
			}
			if(args[0] == "e" || args[0] == "encrypt" || args[0] == "-e") {
				mode = 'e'
			}
			if(args[0] == "d" || args[0] == "decrypt" || args[0] == "-d") {
				mode = 'd'
			}
			if(args[0] == "t" || args[0] == "-t" || args[0] == "test" || args[0] == "--test") {
				mode = 't'
			}
			switch mode {
			case 't':

			case 'e':
				for {
					fmt.Println("What would you like to encrypt?\n   1. File\n   2. Message")
					fmt.Print("Select (1 / 2) # ")
					selection, _, err := cli_reader.ReadLine()
					check(err)
					if len(selection) > 0 {
						switch selection[0]{
						case byte('1'):
							fmt.Println("> FILE ENCRYPTION")
							fmt.Print("Enter the name of the file: ")
							filename, _, _ := cli_reader.ReadLine()
							fmt.Printf("> Filename: \"%s\"\n", filename)
							plainfile, err := os.Open(string(filename))
							check(err)
							plain_reader := bufio.NewReader(plainfile)
							key, err := askForKey(cli_reader)
							check(err)
							ron.Init(key)
							enc_data := ron.Encode(plain_reader)
							dump2file(enc_data, "cyphertext")

						case byte('2'):
							fmt.Println("> MESSAGE ENCRYPTION")

						case byte('q'):
							fmt.Println("> RECEIVED QUIT SIG. Exiting...")
							os.Exit(3)

						}
					} else {
						if err == nil {
							fmt.Println("\nERROR no option selected")
						} else {
							fmt.Println("Something went wrong. Try again...")
							break
						}
					}
				}
			case 'd':
				for {
					fmt.Println("What would you like to decrypt?\n   1. File\n   2. Message")
					fmt.Print("Select (1 / 2) # ")
					selection, _, err := cli_reader.ReadLine()
					check(err)
					if len(selection) > 0 {
						switch selection[0]{
						case byte('1'):
							fmt.Println("> FILE DECRYPTION")
							fmt.Print("Enter the name of the file: ")
							filename, _, _ := cli_reader.ReadLine()
							fmt.Printf("> Filename: \"%s\"\n", filename)
							cryptfile, err := os.Open(string(filename))
							check(err)
							cypher_reader := bufio.NewReader(cryptfile)
							key, err := askForKey(cli_reader)
							check(err)
							ron.Init(key)
							plain_data := ron.Decode(cypher_reader)
							dump2file(plain_data, "decoded")

						case byte('2'):
							fmt.Println("> MESSAGE DECRYPTION")

						case byte('q'):
							fmt.Println("> RECEIVED QUIT SIG. Exiting...")
							os.Exit(3)

						}
					} else {
						if err == nil {
							fmt.Println("\nERROR no option selected")
						} else {
							fmt.Println("Something went wrong. Try again...")
							break
						}
					}
				}

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
	              printing outputting the results to the console or to [filename]

	       -d, --decrypt
	              take the supplied [string] or [filename] and decrypt the contents
	              printing outputting the results to the console or to [filename]
	       -t, --target
	              when this flag is passed sabre will en/decrypt the contents of the
	              supplied [FILENAME]
	       -h, --help
	              prints this man page`)
			}
		}
}
