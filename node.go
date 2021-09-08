package zql

type WalkFunc func(parent Node, current Node, lvl int) (Node, error)

type Node interface {
	Walk(cb WalkFunc, parent Node, lvl int) (Node, error)
	ChildList() []Node
	Append(n Node)
}
