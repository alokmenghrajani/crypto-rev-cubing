package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

/**
 * This code will help you know if you have correctly solved a given face of the crypto-rev-cube.
 * You can try it out with `go run crypto-rev-cubing.go "testinput"`.
 *
 * crypto-rev-cubing:
 *   a combination of Rubik's Cube, figuring out what some piece of code does and "cryptography".
 */
func main() {
	if !regexp.MustCompile("^[a-zA-Z01]{9}$").MatchString(os.Args[1]) {
		fmt.Println("invalid input")
		return
	}

	quit := make(chan bool)
	go compute1(quit)

	hash := make(chan string)
	result := make(chan string)
	go compute2(hash, result)

	input := []byte(os.Args[1])
	hash <- string(input)

	for ctr := 0; !done(input); ctr++ {
		hash <- string(input)

		str := fmt.Sprintf("%d", ctr)
		r := arcHash([]byte(str))
		i := int(r[0]) % len(input)
		j := int(r[1]) % len(input)
		input[i], input[j] = input[j], input[i]

		hash <- string(input)
	}
	quit <- true
	close(hash)
	r := <-result

	valid := map[string]string{
		"09715a463f5bf3078e3f7c7fea24ac25": "👽", // test string for "testinput"
		"5e285ea1b22dec0ed0b92641845dd7de": "😎",
		"a34e5abc573bb55a2aa4a63673bfda96": "😎",
		"a3abde443c717bb2b8ca539479f129a8": "😎",
		"ab99e8d1bb19818b5e953597872bcf66": "😎",
		"1b61365cf854d0edb6762bb77772d0a6": "😎",
		"09628453f36e1f4dadae81b6d69a2feb": "😎",
	}
	if v, ok := valid[r]; ok {
		fmt.Println(v)
	} else {
		fmt.Println("😢")
	}
}

/**
 * pro-tip: you should always roll your own crypto. This prevents the NSA or other attackers from using
 * off-the-shelf tools to defeat your system.
 */
const Output = 16
const Space = 1024

func arcHash(data []byte) []byte {
	state := make([]int, Space)
	j := 0
	i := 0
	for i = range state {
		state[i] = i
	}

	for t := 0; t < Space; t++ {
		i = (i + 1) % Space
		j = (j + state[i] + int(data[i%len(data)])) % Space
		state[i], state[j] = state[j], state[i]
	}

	r := make([]byte, Output)
	for t := 0; t < Output; t++ {
		i = (i + 1) % Space
		j = state[(state[i]+state[j])%Space]
		r[t] = byte(j & 0xff)
	}
	return r
}

/**
 * important computation, part 1 of 2
 */
func compute1(quit chan bool) {
	fmt.Printf("\n\n\n\n\n\n")
	frames := []string{
		"\033[6A\r               X \n                  \n               O  \n             Y/|\\Z\n               |\n              / \\\n",
		"\033[6A\r                 \n             X    \n             Y_O  \n               |\\Z\n               |\n              / \\\n",
		"\033[6A\r                 \n             XY   \n              (O  \n               |\\Z\n               |\n              / \\\n",
		"\033[6A\r              Y  \n                  \n             X_O  \n               |\\Z\n               |\n              / \\\n",
		"\033[6A\r               Y \n                  \n               O  \n             X/|\\Z\n               |\n              / \\\n",
		"\033[6A\r                 \n                 Y\n               O_Z\n             X/|  \n               |\n              / \\\n",
		"\033[6A\r                 \n                ZY\n               O) \n             X/|  \n               |\n              / \\\n",
		"\033[6A\r                Z\n                  \n               O_Y\n             X/|  \n               |\n              / \\\n",
	}
	ctr := 0
	for {
		select {
		case <-quit:
			return
		case <-time.Tick(time.Duration(250) * time.Millisecond):
			ctr++
			s := frames[ctr%len(frames)]
			x := []byte("\033[32mo\033[39m")
			y := []byte("\033[34mo\033[39m")
			z := []byte("\033[35mo\033[39m")
			for t := 0; t < ctr/len(frames)%3; t++ {
				x = xor_slice(xor_slice(x, y), z)
				y = xor_slice(xor_slice(x, y), z)
				z = xor_slice(xor_slice(x, y), z)
				x = xor_slice(xor_slice(x, y), z)
			}
			s = strings.Replace(s, "X", string(x), 1)
			s = strings.Replace(s, "Y", string(y), 1)
			s = strings.Replace(s, "Z", string(z), 1)
			fmt.Print(s)
		}
	}
}

/**
 * important computation, part 2 of 2
 */
func compute2(hash chan string, result chan string) {
	r := make([]byte, Output)
	for {
		data, ok := <-hash
		if !ok {
			result <- hex.EncodeToString(r)
			return
		}
		r = xor_slice(r, arcHash([]byte(data)))
	}
}

/**
 * A boring helper function
 */
func xor_slice(a []byte, b []byte) []byte {
	r := make([]byte, len(a))
	for i, v := range a {
		r[i] = v ^ b[i]
	}
	return r
}

/**
 * Another boring helper function
 */
func done(arr []byte) (r bool) {
	r = true
	for i, v := range arr {
		for j, w := range arr {
			r = r && (i > j || v <= w)
		}
	}
	return
}
