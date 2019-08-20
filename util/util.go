package util

import (
	"os"
	"io"
	"io/ioutil"
	"bufio"
	"fmt"
	"syscall"
	"golang.org/x/crypto/ssh/terminal"
	)


func CheckPanic(e error) {
    if e != nil {
        panic(e)
    }
}
//If e is not null will print error and return false
func CheckWarn(e error) (bool){
	if e != nil {
		fmt.Print("\033[0;31;40mERR > ")
		fmt.Println(e)
		fmt.Print("\033[0;37;40m")
		return true
	} else {
		return false
	}
}
func IVfile(r *bufio.Reader) ([]byte) {
	b1 := make([]byte, 10)
	_, err := io.ReadFull(r, b1)
	CheckPanic(err)
	return b1
}
func AskForKey(r *bufio.Reader ) ([]byte, error){
	fmt.Print("Please input your key: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return bytePassword, err
}
func Dump2File(d []byte, n string) (error) {
	return ioutil.WriteFile(n, d, 0644)
}

//Function that uses the entropy pool of /dev/urandom
//to generate a slice of bytes for: rv := EntropyBits(n)
//asserts: 0 < n < 256 AND len(rv) == n
func EntropyBytes(bytes uint8) ([]byte) {
	var randbuff [256]byte
	devrand, err := os.Open("/dev/urandom")
	if err == nil {
		reader := bufio.NewReader(devrand)
		io.ReadFull(reader, randbuff[:bytes])
		return randbuff[:bytes]
	} else {
		return randbuff[:bytes]
	}
}
//function that provides a ux for the client to navigate the program and select
//the mode of operation they would like to preform (enc/decryption)
func Menu(r *bufio.Reader) (rune){
	fmt.Print("\033[1;35;40m")
	fmt.Println(`

		     _______.     ___      .______   .______       _______
		    /       |    /   \     |   _  \  |   _  \     |   ____|
		   |   (----'   /  ^  \    |  |_)  | |  |_)  |    |  |__
		    \   \      /  /_\  \   |   _  <  |      /     |   __|
		.----)   |    /  _____  \  |  |_)  | |  |\  \----.|  |____
		|_______/    /__/     \__\ |______/  | _| \______||_______|`)
	fmt.Print("\033[0;37;40m\n")
	fmt.Println("Options:")
	fmt.Println("  (1) \033[2;32;40mE\033[0;37;40mncrypt")
	fmt.Println("  (2) \033[2;32;40mD\033[0;37;40mecrypt")
	fmt.Println("  (3) \033[2;32;40mH\033[0;37;40melp")
	fmt.Println("  (4) \033[2;32;40mQ\033[0;37;40muit")
	fmt.Print("Select an above option (1-4) ")

	selection, _,err := r.ReadLine()
	CheckPanic(err)
	if rune(selection[0]) == 'e' || rune(selection[0]) == 'E' ||  selection[0] == 0x31 {
		return 'e'
	}
	if rune(selection[0]) == 'd' || rune(selection[0]) == 'D' ||  selection[0] == 0x32 {
		return 'd'
	}
	if rune(selection[0]) == 'h' || rune(selection[0]) == 'H' ||  selection[0] == 0x33 {
		return 'h'
	}
	if rune(selection[0]) == 'q' || rune(selection[0]) == 'Q' || selection[0] == 0x34 {
		return 'q'
	}
	return 't'
}
