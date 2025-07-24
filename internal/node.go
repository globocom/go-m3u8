package internal

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

// The HLSElement data type holds the following attributes:
//   - Name: The name of the Element (e.g. tag name).
//   - URI: The Uniform Resource Identifier of the Element (if applicable).
//   - Attrs: In-manifest Element attributes, in key-value format.
//   - Details: Not-in-manifest Element attributes, in key-value format.
type HLSElement struct {
	Name    string
	URI     string
	Attrs   map[string]string
	Details map[string]string
}

// Creates a new Node with the given HLSElement attributes.
func (l *DoublyLinkedList) NewNode(name, uri string, attrs, details map[string]string) *Node {
	element := &HLSElement{
		Name:    name,
		URI:     uri,
		Attrs:   attrs,
		Details: details,
	}
	return &Node{HLSElement: element}
}

// Adds a new node to the end of the doubly linked list
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

// Inserts newNode after node in the doubly linked list
//
//	node ---> newNode ---> node.Next
func (l *DoublyLinkedList) InsertAfter(node, newNode *Node) {
	if node == nil {
		return
	}

	newNode.Prev = node
	newNode.Next = node.Next
	node.Next = newNode

	if newNode.Next != nil {
		newNode.Next.Prev = newNode
	} else {
		l.Tail = newNode
	}
}

// Inserts newNode before node in the doubly linked list
//
//	node.Prev ---> newNode ---> node
func (l *DoublyLinkedList) InsertBefore(node, newNode *Node) {
	if node == nil {
		return
	}

	newNode.Next = node
	newNode.Prev = node.Prev
	node.Prev = newNode

	if newNode.Prev != nil {
		newNode.Prev.Next = newNode
	} else {
		l.Head = newNode
	}
}

// Inserts newNode between node1 and node2 in the doubly linked list
//
//	node1 ---> newNode ---> node2
func (l *DoublyLinkedList) InsertBetween(node1, node2, newNode *Node) {
	if node1 == nil || node2 == nil {
		return
	}

	if node1.Next != node2 {
		return
	}

	newNode.Prev = node1
	newNode.Next = node2
	node1.Next = newNode
	node2.Prev = newNode
}

// Searches for a node with the specified element name in the doubly linked list
func (l *DoublyLinkedList) Find(elementName string) (*Node, bool) {
	current := l.Head
	for current != nil {
		if current.HLSElement.Name == elementName {
			return current, true
		}
		current = current.Next
	}

	return nil, false
}

// Searches for all nodes with the specified element name in the doubly linked list
func (l *DoublyLinkedList) FindAll(elementName string) []*Node {
	current := l.Head
	result := make([]*Node, 0)
	for current != nil {
		if current.HLSElement.Name == elementName {
			result = append(result, current)
		}
		current = current.Next
	}
	return result
}
