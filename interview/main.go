package main

import "fmt"

type Node struct {
	Val   int
	Left  *Node
	Right *Node
}

var sumNum int

func sum(root *Node) {
	if root == nil {
		return
	}
	if root.Right != nil {
		sumNum += root.Right.Val
		sum(root.Right)
	}
	if root.Left != nil {
		sum(root.Left)
	}

}
func NewNode(val int) *Node {
	return &Node{
		Val:   val,
		Left:  nil,
		Right: nil,
	}
}

// type Agent struct {
// 	name string
// 	Dentdpents []Agent
// }
// func NewAgent(name string) Agent {
// 	return Agent{
// 		name: name,
// 		Dentdpents: make([]Agent,0),
// 	}
// }
func main() {
	/*	1
			/ \
			2	3
			/
		   4 5
		    6
	*/
	tree := &Node{1, nil, nil}

	node1 := NewNode(2)
	node2 := NewNode(3)
	tree.Left = node1
	tree.Right = node2
	node3 := NewNode(4)
	node4 := NewNode(5)
	node2.Left = node3
	node2.Right = node4
	node5 := NewNode(6)
	node4.Left = node5
	sum(tree)
	fmt.Println(sumNum)
	// //a->b->
	// Agenta := NewAgent("agenta")
	// Agentb := NewAgent("agentb")
	// Agentc := NewAgent("agentc")
	// //a添加依赖b，c
	// Agenta.Dentdpents = append(Agenta.Dentdpents,Agentb)
	// Agenta.Dentdpents = append(Agenta.Dentdpents,Agentc)
	// //b添加依赖a
	// Agentb.Dentdpents = append(Agentb.Dentdpents,Agenta)
	// m := make(map[string]bool,0)
	// AgentaArray := make([]Agent,0)
	// AgentaArray = append(AgentaArray,Agenta)
	// AgentaArray = append(AgentaArray, Agentb)
	// AgentaArray = append(AgentaArray, Agentc)
	// //
	// for _,agent:= range AgentaArray {
	// 	agentDepdents := agent.Dentdpents
	// 	for _,dependent := agentDepdents {
	// 		//出现依赖重复
	// 		if m[dependent.name] {

	// 		}
	// 		m[dependent.name] = true
	// 	}
	// }
}
