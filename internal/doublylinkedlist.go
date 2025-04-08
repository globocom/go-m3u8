package internal

import (
	"fmt"
)

// A HLS Playlist is a doubly-linked list of of Node objects.
// Each Node represents a HLSElement of the Playlist, amounting to one or more lines of the m3u8 file.
// For example, a Media Segment Node will be comprised of two lines: the #EXTINF tag + the segment URI below it.
// Alternatively, a Media Sequence Node is only one line long: the #EXT-X-MEDIA-SEQUENCE tag.
type DoublyLinkedList struct {
	Head, Tail *Node
}

// The Node data type holds the following attributes:
//   - HLSElement: Pointer to HLSElement it represents on the list.
//   - Prev, Next: Pointers to previous or next Node in the list.
type Node struct {
	HLSElement *HLSElement
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
		if current.HLSElement.Name == tagName {
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
		if current.HLSElement.Name == tagName {
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
