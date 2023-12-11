package main

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mr-rabot/GoAssess-MapUp/Models"
)



func processSingle(c *gin.Context) {
	var req models.SortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	startTime := time.Now()
	sortedArrays := make([][]int, len(req.ToSort))

	for i, arr := range req.ToSort {
		sortedArray := make([]int, len(arr))
		copy(sortedArray, arr)
		sort.Ints(sortedArray)
		sortedArrays[i] = sortedArray
	}

	timeTaken := time.Since(startTime).Nanoseconds()

	resp := models.SortResponse{
		SortedArrays: sortedArrays,
		TimeNs:       timeTaken,
	}

	c.JSON(http.StatusOK, resp)
}

func processConcurrent(c *gin.Context) {
	var req models.SortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	startTime := time.Now()
	sortedArrays := make([][]int, len(req.ToSort))
	ch := make(chan int, len(req.ToSort))

	var wg sync.WaitGroup

	for i, arr := range req.ToSort {
		wg.Add(1)
		go func(i int, arr []int) {
			defer wg.Done()
			sortedArray := make([]int, len(arr))
			copy(sortedArray, arr)
			sort.Ints(sortedArray)
			sortedArrays[i] = sortedArray
			ch <- i
		}(i, arr)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for range req.ToSort {
		<-ch
	}

	timeTaken := time.Since(startTime).Nanoseconds()

	resp := models.SortResponse{
		SortedArrays: sortedArrays,
		TimeNs:       timeTaken,
	}

	c.JSON(http.StatusOK, resp)
}

func main() {
	r := gin.Default()

	r.POST("/process-single", processSingle)
	r.POST("/process-concurrent", processConcurrent)

	if err := r.Run(":8000"); err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
