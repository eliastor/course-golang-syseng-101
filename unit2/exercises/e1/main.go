package main

import (
	"fmt"
	"regexp"
	"strconv"

	"golang.org/x/tour/tree"
)

// type Tree struct {
// 	Value int
// }

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	//arr := []int{}

	re := regexp.MustCompile("[0-9]+")
	str := re.FindAllString(t.String(), -1)

	for i := 0; i < len(str); i++ {
		a, _ := strconv.Atoi(str[i])
		ch <- a
	}
	close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)

	arr1 := []int{}
	arr2 := []int{}

	for val := range ch1 {
		arr1 = append(arr1, val)
	}
	for val := range ch2 {
		arr2 = append(arr2, val)
	}
	for i := range arr1 {
		if arr1[i] != arr2[i] {
			return false
		}
	}
	return true
}

func main() {

	fmt.Println(Same(tree.New(1), tree.New(2)))
	fmt.Println(Same(tree.New(2), tree.New(2)))

}
