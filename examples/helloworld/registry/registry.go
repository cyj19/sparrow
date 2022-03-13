/**
 * @Author: cyj19
 * @Date: 2022/3/13 22:13
 */

package main

import "github.com/cyj19/sparrow/registry"

func main() {
	registry.Run("tcp", ":9999")
}
