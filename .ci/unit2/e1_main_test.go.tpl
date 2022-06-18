package main

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tour/tree"
)

func TestMain(t *testing.T) {
	assert.False(t, Same(tree.New(1), tree.New(2)))
	assert.True(t, Same(tree.New(2), tree.New(2)))
	ch := make(chan int)
	nums := []int{}
	tree2 := []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}
	go Walk(tree.New(2), ch)
	for num := range ch {
		t.Log(num)
		nums = append(nums, num)
	}
	sort.Ints(nums)
	assert.Equal(t, tree2, nums, "Tree walk function produces corrupted data for tree.New(2)")
}
