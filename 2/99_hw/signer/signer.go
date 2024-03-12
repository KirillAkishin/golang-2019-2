package main

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var TH = 6 //MultiHash создает TH хешей на каждое входящее значение

var SingleHash = job(func(in, out chan interface{}) {
	fmt.Println("SingleHash:\t\tstarting")
	wg := &sync.WaitGroup{}
	for fibNum := range in {
		wg.Add(1)
		chan1 := make(chan string)
		chan2 := make(chan string)
		var str1, str2 string
		str1, ok := fibNum.(string)
		if !ok {
			dataInt := fibNum.(int)
			str1 = strconv.Itoa(dataInt)
		}

		go func(str string, ch chan string) {
			ch <- DataSignerCrc32(str)
		}(str1, chan2)

		str2 = DataSignerMd5(str1)
		go func(str string, ch chan string) {
			ch <- DataSignerCrc32(str)
		}(str2, chan1)

		go func(ch1, ch2 chan string) {
			var res1, res2 string
			res1 = <-ch1
			res2 = <-ch2
			result := res2 + "~" + res1
			out <- result
			fmt.Println("SingleHash:\t\tsend result =", result)
			wg.Done()
		}(chan1, chan2)
	}
	wg.Wait()
})

var MultiHash = job(func(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	fmt.Println("MultiHash:\t\tstarting")
	for temp := range in {
		dataArr := make([]string, TH)
		wg.Add(1)
		//fmt.Println("MultiHash:\t\tcome temp =", temp)
		go func(tmp interface{}) {
			wg2 := &sync.WaitGroup{}
			data, ok := tmp.(string)
			if !ok {
				dataInt := tmp.(int)
				data = strconv.Itoa(dataInt)
			}
			for i := 0; i < TH; i++ {
				wg2.Add(1)
				go func(arr []string, th int) {
					arr[th] = DataSignerCrc32(strconv.Itoa(th) + data)
					wg2.Done()
				}(dataArr, i)
			}
			wg2.Wait()

			result := strings.Join(dataArr, "")
			out <- result
			wg.Done()
			fmt.Println("MultiHash:\t\tsend result =", result)
		}(temp)
	}
	wg.Wait()
})

var CombineResults = job(func(in, out chan interface{}) {
	fmt.Println("CombineResults:\t\tstarting")
	var dataArr []string
	for temp := range in {
		//fmt.Println("CombineResults:\t\tcome temp =", temp)
		data, ok := temp.(string)
		if !ok {
			dataInt := temp.(int)
			data = strconv.Itoa(dataInt)
		}
		dataArr = append(dataArr, data)
	}
	sort.Strings(dataArr)
	stringOutput := strings.Join(dataArr, "_")
	out <- stringOutput
})

func ExecutePipeline(freeFlowJobs ...job) {
	fmt.Println("ExecutePipeline:\tstarting")
	wg := &sync.WaitGroup{}
	ch2 := make(chan interface{})

	for i := 0; i < len(freeFlowJobs); i++ {
		wg.Add(1)
		ch1 := make(chan interface{})
		go func(i int, input, output chan interface{}) {
			freeFlowJobs[i](input, output)
			close(output)
			wg.Done()
			//fmt.Println("ExecutePipeline:\tLOCK (", i, ")")
		}(i, ch2, ch1)
		ch2 = ch1
	}
	wg.Wait()
	//fmt.Println("ExecutePipeline:\tfinish")
}

func main() {
	runtime.GOMAXPROCS(0)
	//inputData := []int{0, 1, 2, 3, 4, 5, 7}
	//inputData := []int{0,1,2,3,4,5,7,41,23,49,2,0,12,13,7}
	inputData := []int{0, 1, 2, 3}
	//inputData := []int{0}
	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				fmt.Println("fibNum =", fibNum)
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			if !ok {
				fmt.Println("can't convert result data to string")
			}
			fmt.Println("FINAL RESULT =")
			fmt.Println(data)
		}),
	}
	start := time.Now()
	ExecutePipeline(hashSignJobs...)
	end := time.Since(start)
	fmt.Println("ExecutePipeline's time = ", end)
}
