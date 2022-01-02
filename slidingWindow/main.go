package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"sync"
	"time"
)

func main() {
	ep := newErrPercent()
	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 10)
		go func() {
			for {
				time.Sleep(time.Millisecond * 10)
				//模拟添加请求数量和错误数量
				ep.add(time.Now().Unix()%int64(3) == 0)
			}
		}()
	}

	cp := newCpuPercent()
	go func() {
		for {
			time.Sleep(time.Millisecond * 200)
			//添加当前cpu使用量
			cp.add()
		}
	}()

	for {
		time.Sleep(time.Millisecond * 400)
		fmt.Printf("窗口时间内cpu使用量：%f,错误率：%f \n", cp.avgCpuUsage(), ep.avgErrCount())
	}
}

var reqCount float32
var errCount float32

type errCountPercent interface {
	add(isErr bool)
	avgErrCount() float32
}

type errPercent struct {
	list         [5][2]float32
	writeIn      int
	totalNumErrs float32
	totalNumReqs float32
	lastTime     time.Time
}

func (e *errPercent) add(isErr bool) {
	wLock := sync.Mutex{}
	wLock.Lock()

	reqCount += 1
	if isErr {
		errCount += 1
	}
	if time.Now().UnixMilli()-e.lastTime.UnixMilli() >= 200 {
		delIndex := delIn(len(e.list), e.writeIn)
		e.totalNumErrs -= e.list[delIndex][0]
		e.totalNumReqs -= e.list[delIndex][1]
		e.list[e.writeIn][0] = errCount
		e.list[e.writeIn][1] = reqCount
		e.totalNumErrs += errCount
		e.totalNumReqs += reqCount
		e.writeIn = nextWriIn(len(e.list), e.writeIn)
		e.lastTime = time.Now()
	}
	wLock.Unlock()
}

func (e *errPercent) avgErrCount() float32 {
	if e.totalNumReqs == 0 {
		return 0
	}
	return e.totalNumErrs / e.totalNumReqs * float32(100)
}

func newErrPercent() errCountPercent {
	return &errPercent{lastTime: time.Now()}
}

type cpuUsagePervent interface {
	add()
	avgCpuUsage() float64
}

type cpuPercent struct {
	list       [5]float64
	writeIn    int
	totalUsage float64
	size       float64
}

func (c *cpuPercent) add() {
	c.totalUsage -= c.list[delIn(len(c.list), c.writeIn)]
	cpuInfo := getCpuInfo()
	c.list[c.writeIn] = cpuInfo
	c.totalUsage += cpuInfo
	c.writeIn = nextWriIn(len(c.list), c.writeIn)
}

func (c *cpuPercent) avgCpuUsage() float64 {
	return c.totalUsage / c.size
}

func newCpuPercent() cpuUsagePervent {
	return &cpuPercent{size: 5}
}

func getCpuInfo() float64 {
	cp, _ := cpu.Percent(time.Second, true)
	var c float64
	for i := 0; i < len(cp); i++ {
		c += cp[i]
	}
	return c / float64(len(cp))
}

func nextWriIn(size, wriIn int) int {
	if wriIn+1 <= size-1 {
		wriIn++
	} else {
		wriIn = 0
	}
	return wriIn
}

func delIn(size, wriIn int) int {
	var index int
	if wriIn == size-1 {
		index = 0
	} else {
		index = wriIn + 1
	}
	return index
}
