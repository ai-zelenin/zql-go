package zql

type WalkFunc func(parent Node, current Node, lvl int) error

type Node interface {
	Walk(cb WalkFunc, parent Node, lvl int) error
	Children() []Node
	Append(n Node)
}
