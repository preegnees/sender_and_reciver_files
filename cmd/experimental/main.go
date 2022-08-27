package main

import (
	"log"
	"sync"
)

func main() {
	memStorage := sync.Pool{New: func() interface{} {
		mem := make([]byte, 1*1024)
		return mem
	}}
	arr1 := memStorage.Get().([]byte)
	log.Println("arr len:", len(arr1))
	arr2 := memStorage.Get().([]byte)
	log.Println("arr len:", len(arr2))
	arr3 := memStorage.Get().([]byte)
	log.Println("arr len:", len(arr3))
}

//ничего не нужно ограничивать, нужно при отправке данных, не ждать пока они дайдут, а сразу читать в память
// а если не долшло, то отменить отправку



