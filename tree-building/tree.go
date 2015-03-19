// +build !example

package tree

import (
	"errors"
	"fmt"
)

const TestVersion = 1

type Record struct {
	Id, Parent int
}

type Node struct {
	Id       int
	Children []*Node
}

type Mismatch struct{}

func (m Mismatch) Error() string {
	return "c"
}

func Build(records []Record) (*Node, error) {
	if len(records) == 0 {
		return nil, nil
	}
	root := &Node{}
	todo := []*Node{root}
	n := 1
	for {
		if len(todo) == 0 {
			break
		}
		newTodo := []*Node(nil)
		for _, c := range todo {
			for _, r := range records {
				if r.Parent == c.Id {
					if r.Id < c.Id {
						return nil, errors.New("a")
					} else if r.Id == c.Id {
						if r.Id != 0 {
							return nil, fmt.Errorf("b")
						}
					} else {
						n++
						switch len(c.Children) {
						case 0:
							nn := &Node{Id: r.Id}
							c.Children = []*Node{nn}
							newTodo = append(newTodo, nn)
						case 1:
							nn := &Node{Id: r.Id}
							if c.Children[0].Id < r.Id {
								c.Children = []*Node{c.Children[0], nn}
								newTodo = append(newTodo, nn)
							} else {
								c.Children = []*Node{nn, c.Children[0]}
								newTodo = append(newTodo, nn)
							}
						default:
							nn := &Node{Id: r.Id}
							newTodo = append(newTodo, nn)
						breakpoint:
							for _ = range []bool{false} {
								for i, cc := range c.Children {
									if cc.Id > r.Id {
										a := make([]*Node, len(c.Children)+1)
										copy(a, c.Children[:i])
										copy(a[i+1:], c.Children[i:])
										copy(a[i:i+1], []*Node{nn})
										c.Children = a
										break breakpoint
									}
								}
								c.Children = append(c.Children, nn)
							}
						}
					}
				}
			}
		}
		todo = newTodo
	}
	if n != len(records) {
		return nil, Mismatch{}
	}
	if err := chk(root, len(records)); err != nil {
		return nil, err
	}
	return root, nil
}

func chk(n *Node, m int) (err error) {
	if n.Id > m {
		return fmt.Errorf("z")
	} else if n.Id == m {
		return fmt.Errorf("y")
	} else {
		for i := 0; i < len(n.Children); i++ {
			err = chk(n.Children[i], m)
			if err != nil {
				return
			}
		}
		return
	}
}
