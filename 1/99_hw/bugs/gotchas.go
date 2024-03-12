package main

import (
	"fmt"
	"sort"
)


/*
	сюда вам надо писать функции, которых не хватает, чтобы проходили тесты в gotchas_test.go

	IntSliceToString
	MergeSlices
	GetMapValuesSortedByKey
*/

func IntSliceToString(buf []int) (s string) {
	idx := 0
	for idx < len(buf) {
		s += fmt.Sprintf("%d", buf[idx]) 
		idx++
	} 
	return s
}


func MergeSlices(buf1 []float32, buf2 []int32) (result []int) {
	idx := 0
	for idx < len(buf1) {
		result = append(result, int(buf1[idx])) 
		idx++
	} 
	idx = 0
	for idx < len(buf2) {
		result = append(result, int(buf2[idx])) 
		idx++
	}
	 return result 
}


func GetMapValuesSortedByKey(input map[int]string) (result []string){
	result = make([]string,len(input))
	sortedKeys := make([]int, 0, len(input))
	for k := range input {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Ints(sortedKeys)
	i := 0
	for k := range sortedKeys {
		result[i] = input[sortedKeys[k]]
		i++
	}
	return result
}

