/**
 * @Author: cyj19
 * @Date: 2022/3/13 22:29
 */

package registry

import "testing"

func TestHeartBeat(t *testing.T) {
	HeartBeat("http://localhost:9999/sparrow/registry", "tcp", ":8787", 0)
}
