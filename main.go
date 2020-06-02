package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

const freeFork = -1

var forks = make(map[int]int)

var l = sync.Mutex{}

const numOfPhil = 5

type stat struct {
	thinkTime  time.Duration
	thinkCount int
	eatTime    time.Duration
	eatCount   int
}

var stats [numOfPhil]stat

func main() {
	log.SetFlags(log.Lmicroseconds)

	for i := 0; i < numOfPhil; i++ {
		right := (i + 1) % numOfPhil
		go philosopher(i, i, right)
	}

	go printForksAndStats()

	for {
		time.Sleep(time.Second)
	}
}

func printForksAndStats() {
	for {
		fmt.Println("================")
		total := 0
		for i := 0; i < numOfPhil; i++ {
			fmt.Printf("Fork:%d, userID:%d\n", i, forks[i])
			if forks[i] != freeFork {
				total++
			}
		}
		if total == 5 {
			fmt.Printf(">>>>>>>All forks are being used<<<<<<<\n")
		}
		for i := 0; i < numOfPhil; i++ {
			averThink, averEat := time.Duration(0), time.Duration(0)
			if stats[i].thinkCount > 1 {
				averThink = stats[i].thinkTime / time.Duration(stats[i].thinkCount)
			}
			if stats[i].eatCount > 1 {
				averEat = stats[i].eatTime / time.Duration(stats[i].eatCount)
			}

			fmt.Printf("userID:%d thinkTime(ave):%v, thinkCount:%d, eatTime(ave):%v, eatCount:%d, TotalTime:%v\n",
				i, averThink, stats[i].thinkCount, averEat, stats[i].eatCount, stats[i].thinkTime+stats[i].eatTime)
		}

		time.Sleep(time.Second * 2)
	}
}

func registerFork(id int) {
	_, ok := forks[id]
	if !ok {
		l.Lock()
		forks[id] = freeFork
		l.Unlock()
	}
}

// Returns true if able to pick up
func pickFork(userID, forkID int, timeout time.Duration) bool {
	for {
		l.Lock()
		val, ok := forks[forkID]
		l.Unlock()
		if !ok {
			log.Panic("Unknown fork ", forkID)
		}
		if val == freeFork {
			l.Lock()
			forks[forkID] = userID
			l.Unlock()
			break
		}
		time.Sleep(time.Millisecond * 100)
		timeout -= 100 * time.Millisecond
		if timeout <= 0 {
			return false
		}
	}
	log.Println("Philosopher", userID, "picked up fork", forkID)
	return true
}

func releaseFork(userID, forkID int) {
	l.Lock()
	val, ok := forks[forkID]
	l.Unlock()
	if !ok {
		log.Panic("Unknown fork ", forkID)
	}
	if val == freeFork {
		log.Panic("Unused fork ", forkID)
	}
	l.Lock()
	forks[forkID] = freeFork
	l.Unlock()
	log.Println("Philosopher", userID, "released fork", forkID)
}

func philosopher(userID, leftForkID, rightForkID int) {
	log.Println("Philosopher", userID, "started. Left Fork:", leftForkID, ", Right Fork:", rightForkID)
	registerFork(leftForkID)
	registerFork(rightForkID)

	for {
		timeStart := time.Now()
		picked := false
		for picked == false {
			log.Println("Philosopher", userID, "is waiting to pick up left fork:", leftForkID)
			picked = pickFork(userID, leftForkID, time.Second*100)
			if !picked {
				continue
			}

			log.Println("Philosopher", userID, "is waiting to pick up right fork:", rightForkID)
			picked = pickFork(userID, rightForkID, time.Second*3)
			if !picked {
				log.Println(">>>>>>> Phil:", userID, "couldn't get right fork. Releasing left fork:", leftForkID)
				releaseFork(userID, leftForkID)
				time.Sleep(time.Second)
			}
		}

		diffTime := time.Now().Sub(timeStart)
		stats[userID].thinkTime += diffTime
		stats[userID].thinkCount++
		log.Print("Philosopher ", userID, " thinking time was ", diffTime)
		timeStart = time.Now()

		sleepUpTo(6)
		diffTime = time.Now().Sub(timeStart)
		stats[userID].eatTime += diffTime
		stats[userID].eatCount++
		log.Print("Philosopher ", userID, " eating time was ", diffTime)
		timeStart = time.Now()
		releaseFork(userID, leftForkID)
		releaseFork(userID, rightForkID)
		time.Sleep(time.Second)
	}
}

// Sleep random milliseconds, upto x seconds
func sleepUpTo(second int) {
	time.Sleep(time.Duration(rand.Intn(second*1000)) * time.Millisecond)
}
