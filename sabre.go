package main

import (
	"os"
	"bufio"
	"fmt"
	"github.com/MatthewJWolfe/sabre/util"
	RC4 "github.com/MatthewJWolfe/sabre/arcfour"
	)


func main() {
	var (
		mode rune
		quit = false
	)
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
	}
	for !quit {
		if mode == 0x00 {
			mode = util.Menu(cli_reader)
		}
		switch mode {
		case 't':

		case 'e':
			fmt.Println("What would you like to encrypt?\n   1. File\n   2. Message")
			fmt.Print("Select (1 / 2) # ")
			selection, _, err := cli_reader.ReadLine()
			util.CheckPanic(err)
			if len(selection) > 0 {
				switch selection[0]{
				case byte('1'):
					fmt.Println("> FILE ENCRYPTION")
					fmt.Print("Enter filename to be encrypted: ")
					filename, _, _ := cli_reader.ReadLine()
					fmt.Printf("> Filename: \"%s\"\n", filename)
				//Open the specified file and create io reader object
					plainfile, err := os.Open(string(filename))
					if util.CheckWarn(err) {
						break
					}
					plain_reader := bufio.NewReader(plainfile)
				//Request key to encrypt the contents of plainfile
					key, err := util.AskForKey(cli_reader)
					util.CheckPanic(err)
				//Save the key, generate the IV, create PRGA and preform KSA.
					ron.Init(key)
				//Use the Encode function of the RC4 struct to generate cyphertext
					enc_data := ron.Encode(plain_reader)
				//The IV is written as the first 10 Bytes, followed by cyphertext
					ron.WriteEncFile(enc_data, string(append(filename, ".cs1"...)))
				//break the enc loop
					mode = 0

				case byte('2'):
					fmt.Println("> MESSAGE ENCRYPTION")
				case byte('b'):
					fmt.Println("> BACK")
					mode = 0
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
		case 'd':
			fmt.Println("What would you like to decrypt?\n   1. File\n   2. Message")
			fmt.Print("Select (1 / 2) # ")
			selection, _, err := cli_reader.ReadLine()
			util.CheckPanic(err)
			if len(selection) > 0 {
				switch selection[0]{
				case byte('1'):
				//Go through the menu options with user, requesting name of file to be
				//decrypted, if the plaintext should be saved an if so where
					fmt.Println("> FILE DECRYPTION")
					fmt.Print("Enter the name of the file: ")
					filename, _, _ := cli_reader.ReadLine()
					fmt.Printf("> Filename: \"%s\"\n", filename)
					cryptfile, err := os.Open(string(filename))
					if util.CheckWarn(err) {
						break
					}
				//reader object pointing to the cyphertext file
					cypher_reader := bufio.NewReader(cryptfile)
					key, err := util.AskForKey(cli_reader)
					util.CheckPanic(err)
				//iv is removed from file and saved
					iv := util.IVfile(cypher_reader)
				//RC4 object is instantiated with key from user and IV from file
				//KSA is preformed
					fmt.Println(iv)
					ron.Init(key, iv)
				//Cyphertext is decoded and stored and a slice is returned
					plain_data := ron.Decode(cypher_reader)
				//options resulting in plaintext saved to file or printed to console
					fmt.Println("> Save decrypted data to file? (Y/N)")
					fmt.Print("Select (Y / N) # ")
					save, _, err := cli_reader.ReadLine()
					util.CheckPanic(err)
					if save[0] == byte('y') || save[0] == byte('Y') {
						fmt.Print("Name of decoded file: ")
						filename, _, _ := cli_reader.ReadLine()
						util.Dump2File(plain_data, string(filename))
						fmt.Println("Ok saving data to file: ", filename)
						mode = 0
					} else {
						fmt.Println("> Ok printing decoded data in ASCII to terminal...\n")
						fmt.Println("-----BEGIN SABRE MESSAGE-----")
						fmt.Printf("%s\n", plain_data)
						fmt.Println("-----END SABRE MESSAGE-----\n")
						fmt.Println("Press enter to return to menu...")
						cli_reader.ReadLine()
						mode = 0
					}

				case byte('2'):
					fmt.Println("> MESSAGE DECRYPTION")
					mode = 0
				case byte('b'):
					fmt.Println("> BACK")
					mode = 0
				case byte('q'):
					fmt.Println("> RECEIVED QUIT SIG. Exiting...")
					os.Exit(3)
				}
			} else {
				if err == nil {
					fmt.Println("\nERROR no option selected")
				} else {
					fmt.Println("Something went wrong. Try again...")
				}
			}

		case 'h':
			fmt.Println(`
sabre                           User Commands                           sabre

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
				fmt.Println("\nPress enter to return to menu...")
				cli_reader.ReadLine()
				mode = 0
		case 'q':
			fmt.Println("> RECEIVED QUIT SIG. Exiting...")
			os.Exit(3)

		}
	}
}
