package internal

import (
	"fmt"
)

type DoublyLinkedList struct {
	Head, Tail *Node
}

// Node refers to each element in a playlist, being represented by a doubly linked list.
//   - Name: The name of the node.
//   - URI: The Uniform Resource Identifier of the node (if applicable)
//   - Attrs: In-manifest node attributes, in key-value format.
//   - Details: Not-in-manifest node attributes, in key-value format.
//   - Prev, Next: Pointers to previous or next node in the linked list.
type Node struct {
	Name       string
	URI        string
	Attrs      map[string]string
	Details    map[string]string
	Prev, Next *Node
}

func (l *DoublyLinkedList) Insert(node *Node) {
	if l.Head == nil {
		l.Head = node
		l.Tail = node
	} else {
		node.Prev = l.Tail
		l.Tail.Next = node
		l.Tail = node
	}
}

func (l *DoublyLinkedList) Find(tagName string) (*Node, bool) {
	current := l.Head
	for current != nil {
		if current.Name == tagName {
			return current, true
		}
		current = current.Next
	}

	return nil, false
}

func (l *DoublyLinkedList) FindAll(tagName string) []*Node {
	current := l.Head
	result := make([]*Node, 0)
	for current != nil {
		if current.Name == tagName {
			result = append(result, current)
		}
		current = current.Next
	}
	return result
}

func (l *DoublyLinkedList) ModifyNodesBetween(
	startCondition func(*Node) bool,
	endCondition func(*Node) bool,
	transform func(*Node),
) error {
	var startNode, endNode *Node
	current := l.Head

	for current != nil {
		if startCondition(current) {
			startNode = current
			break
		}
		current = current.Next
	}

	if startNode == nil {
		return fmt.Errorf("start node not found")
	}

	current = startNode.Next
	for current != nil {
		if endCondition(current) {
			endNode = current
			break
		}
		current = current.Next
	}

	if endNode == nil {
		return fmt.Errorf("end node not found")
	}

	current = startNode.Next
	for current != nil && current != endNode {
		transform(current)
		current = current.Next
	}
	return nil
}
