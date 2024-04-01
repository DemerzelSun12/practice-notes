package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	nodesPola []int
	graph     [][]int
	posCount  []int
	negCount  []int
	totalDam  int64
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscanf(reader, "%d\n", &n)

	nodesPolaStr, _ := reader.ReadString('\n')
	nodesPolaStr = strings.TrimSpace(nodesPolaStr)

	nodesPola = make([]int, n+1)
	posCount = make([]int, n+1)
	negCount = make([]int, n+1)
	graph = make([][]int, n+1)

	for i, ch := range nodesPolaStr {
		if ch == '+' {
			nodesPola[i+1] = 1
		} else {
			nodesPola[i+1] = -1
		}
	}

	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscanf(reader, "%d %d\n", &u, &v)
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
	}
	dfs(1, 0)
	fmt.Println(totalDam)

}

func dfs(node, parent int) (int, int) {
	postCount, negeCount := 0, 0
	if nodesPola[node] == 1 {
		postCount++
	} else {
		negeCount++
	}
	for _, child := range graph[node] {
		if child != parent {
			pos, neg := dfs(child, node)
			postCount += pos
			negeCount += neg
		}
	}
	posCount[node], negCount[node] = postCount, negeCount
	totalDam += int64(postCount) * int64(negeCount)
	return postCount, negeCount
}
