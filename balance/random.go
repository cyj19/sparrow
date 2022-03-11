/**
 * @Author: cyj19
 * @Date: 2022/3/11 10:35
 */

package balance

import (
	"math/rand"
	"time"
)

// Random 随机
type Random struct {
	r *rand.Rand
}

func NewRandom() *Random {
	return &Random{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r Random) GetModeResult(n int) int {
	return r.r.Intn(n)
}
