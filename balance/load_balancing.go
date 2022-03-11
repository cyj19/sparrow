/**
 * @Author: cyj19
 * @Date: 2022/3/11 9:16
 */

package balance

type SelectMode int

type LoadBalancing interface {
	GetModeResult(n int) int
}
