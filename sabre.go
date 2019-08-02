package main

import "os"
import "fmt"

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "soap" {
			fmt.Println("In George Lucas' Star Wars trilogy, Jedi Knights were expected to make their own light sabers. The message was clear: a warrior confronted by a powerful empire bent on totalitarian control must be self-reliant. As we face a real threat of a ban on the distribution of strong cryptography, in the United States and possibly world-wide, we should emulate the Jedi masters by learning how to build strong cryptography programs all by ourselves. If this can be done, strong cryptography will become impossible to suppress.")
		} else {
			fmt.Println("I shall spare you the politics")
		}
	}
}
