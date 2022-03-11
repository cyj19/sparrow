/**
 * @Author: cyj19
 * @Date: 2022/3/11 9:20
 */

package balance

type RoundRobin struct {
	index int
}

var _ LoadBalancing = (*RoundRobin)(nil)

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		index: 0,
	}
}

func (r *RoundRobin) GetModeResult(n int) int {
	idx := r.index % n
	r.index = (r.index + 1) % n
	return idx
}
