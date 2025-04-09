package main

import "fmt"

func main() {
	K := []int{1, 2, 3, 4, 5}
	N := 2
	for i := 0; i < len(K); i++ {
		sum := (N + (K[i]-1)*23) + 902
		fmt.Println(sum)
	}
}
