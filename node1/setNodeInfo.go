package main

var node_info = make(map[int][]string)
var total_Servers = []string{
	"http://localhost:4005",
	"http://localhost:4006",
	"http://localhost:4007",
}
func Init() {
	for i := 2; i >=0; i++ {
		firstNode := total_Servers[i%len(total_Servers)]
		secondNode := total_Servers[(i+1)%len(total_Servers)]
		node_info[i] = append(node_info[i], firstNode)
		node_info[i] = append(node_info[i], secondNode)
	}
}
