package gtools

import (
	"math/rand"
	"time"
)

// GenerateRandom 生成一个区间范围的随机数,左闭右开
func GenerateRandom(min, max int) int {
	if min >= max {
		return max
	}

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max - min)
	randNum += min

	return randNum
}

// GenerateRandom64 生成一个区间范围的随机数,左闭右开
func GenerateRandom64(min, max int64) int64 {
	if min >= max {
		return max
	}

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Int63n(max - min)
	randNum += min

	return randNum
}

func AbsInt64(num int64) int64 {
	if num >= 0 {
		return num
	} else {
		return -num
	}
}

// MaxInt64 求最大值.
func MaxInt64(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func MinInt64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}
