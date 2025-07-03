package internal_test

import (
	"testing"

	"github.com/globocom/go-m3u8/internal"
	"github.com/stretchr/testify/assert"
)

func TestDoublyLinkedListInsert(t *testing.T) {
	list := internal.DoublyLinkedList{}

	firstNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "Version",
			Attrs: map[string]string{
				"#EXT-X-VERSION": "3",
			},
		},
	}

	list.Insert(firstNode)

	assert.Equal(t, firstNode, list.Head)
	assert.Equal(t, firstNode, list.Tail)
	assert.Nil(t, firstNode.Prev)
	assert.Nil(t, firstNode.Next)

	secondNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "MediaSequence",
			Attrs: map[string]string{
				"#EXT-X-MEDIA-SEQUENCE": "360948012",
			},
		},
	}

	list.Insert(secondNode)

	assert.Equal(t, firstNode, list.Head)
	assert.Equal(t, secondNode, list.Tail)
	assert.Equal(t, firstNode, secondNode.Prev)
	assert.Equal(t, secondNode, firstNode.Next)

}

func TestDoublyLinkedListInsertAfter(t *testing.T) {
	list := internal.DoublyLinkedList{}

	firstNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "M3u8Identifier",
		},
	}

	secondNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "Comment",
			Attrs: map[string]string{
				"Comment": "## Created with Unified Streaming Platform (version=1.11.23-28141)",
			},
		},
	}

	list.Insert(firstNode)
	list.InsertAfter(firstNode, secondNode)

	assert.Equal(t, firstNode, list.Head)
	assert.Equal(t, secondNode, list.Tail)
	assert.Equal(t, firstNode, secondNode.Prev)
	assert.Equal(t, secondNode, firstNode.Next)

	newNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "Version",
			Attrs: map[string]string{
				"#EXT-X-VERSION": "3",
			},
		},
	}

	list.InsertAfter(firstNode, newNode)

	assert.Equal(t, firstNode, list.Head)
	assert.Equal(t, secondNode, list.Tail)
	assert.Equal(t, newNode, firstNode.Next)
	assert.Equal(t, newNode, secondNode.Prev)
	assert.Equal(t, firstNode, newNode.Prev)
	assert.Equal(t, secondNode, newNode.Next)
}

func TestDoublyLinkedListInsertBefore(t *testing.T) {
	list := internal.DoublyLinkedList{}

	firstNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "M3u8Identifier",
		},
	}

	secondNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "Comment",
			Attrs: map[string]string{
				"Comment": "## Created with Unified Streaming Platform (version=1.11.23-28141)",
			},
		},
	}

	list.Insert(secondNode)
	list.InsertBefore(secondNode, firstNode)

	assert.Equal(t, firstNode, list.Head)
	assert.Equal(t, secondNode, list.Tail)
	assert.Equal(t, firstNode, secondNode.Prev)
	assert.Equal(t, secondNode, firstNode.Next)

	newNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "Version",
			Attrs: map[string]string{
				"#EXT-X-VERSION": "3",
			},
		},
	}

	list.InsertBefore(secondNode, newNode)

	assert.Equal(t, firstNode, list.Head)
	assert.Equal(t, secondNode, list.Tail)
	assert.Equal(t, newNode, firstNode.Next)
	assert.Equal(t, newNode, secondNode.Prev)
	assert.Equal(t, firstNode, newNode.Prev)
	assert.Equal(t, secondNode, newNode.Next)
}

func TestDoublyLinkedListInsertBetween(t *testing.T) {
	list := internal.DoublyLinkedList{}

	firstNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "M3u8Identifier",
		},
	}
	secondNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "Version",
			Attrs: map[string]string{
				"#EXT-X-VERSION": "3",
			},
		},
	}
	thirdNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "Comment",
			Attrs: map[string]string{
				"Comment": "## Created with Unified Streaming Platform (version=1.11.23-28141)",
			},
		},
	}

	list.Insert(firstNode)
	list.Insert(thirdNode)

	assert.Equal(t, firstNode, list.Head)
	assert.Equal(t, thirdNode, list.Tail)

	list.InsertBetween(firstNode, thirdNode, secondNode)

	assert.Equal(t, firstNode, list.Head)
	assert.Equal(t, thirdNode, list.Tail)
	assert.Equal(t, secondNode, firstNode.Next)
	assert.Equal(t, secondNode, thirdNode.Prev)
	assert.Equal(t, firstNode, secondNode.Prev)
	assert.Equal(t, thirdNode, secondNode.Next)
}

func TestDoublyLinkedListFind(t *testing.T) {
	list := internal.DoublyLinkedList{}

	firstNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "Version",
			Attrs: map[string]string{
				"#EXT-X-VERSION": "3",
			},
		},
	}
	secondNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "MediaSequence",
			Attrs: map[string]string{
				"#EXT-X-MEDIA-SEQUENCE": "360948012",
			},
		},
	}

	list.Insert(firstNode)
	list.Insert(secondNode)

	node, found := list.Find("Version")

	assert.True(t, found)
	assert.Equal(t, firstNode, node)
}

func TestDoublyLinkedListFindAll(t *testing.T) {
	list := internal.DoublyLinkedList{}

	firstNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "Version",
			Attrs: map[string]string{
				"#EXT-X-VERSION": "3",
			},
		},
	}
	secondNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "StreamInf",
			Attrs: map[string]string{
				"BANDWIDTH":         "206000",
				"AVERAGE-BANDWIDTH": "187000",
				"CODECS":            "mp4a.40.2,avc1.64001F",
				"RESOLUTION":        "256x144",
				"FRAME-RATE":        "30",
			},
			URI: "channel-audio_1=96000-video=80000.m3u8",
		},
	}
	thirdNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "StreamInf",
			Attrs: map[string]string{
				"BANDWIDTH":         "299000",
				"AVERAGE-BANDWIDTH": "272000",
				"CODECS":            "mp4a.40.2,avc1.64001F",
				"RESOLUTION":        "384x216",
				"FRAME-RATE":        "30",
			},
			URI: "channel-audio_1=96000-video=160000.m3u8",
		},
	}

	list.Insert(firstNode)
	list.Insert(secondNode)
	list.Insert(thirdNode)

	node := list.FindAll("StreamInf")
	assert.Equal(t, len(node), 2)
	assert.Equal(t, node[0], secondNode)
	assert.Equal(t, node[1], thirdNode)
}

func TestModifyNodesBetween(t *testing.T) {
	list := internal.DoublyLinkedList{}

	node1 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "DateRange",
			Attrs: map[string]string{
				"SCTE35-OUT": "0xFFFF",
			},
		},
	}
	node2 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "ExtInf",
			URI:  "1.ts",
		},
	}
	node3 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "ExtInf",
			URI:  "2.ts",
		},
	}
	node4 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "DateRange",
			Attrs: map[string]string{
				"SCTE35-IN": "0xFFFD",
			},
		},
	}

	list.Insert(node1)
	list.Insert(node2)
	list.Insert(node3)
	list.Insert(node4)

	startCondition := func(node *internal.Node) bool {
		return node.HLSElement.Name == "DateRange" && node.HLSElement.Attrs["SCTE35-OUT"] != ""
	}
	endCondition := func(node *internal.Node) bool {
		return node.HLSElement.Name == "DateRange" && node.HLSElement.Attrs["SCTE35-IN"] != ""
	}

	transform := func(node *internal.Node) {
		if node.HLSElement.Name == "ExtInf" && node.HLSElement.URI != "" {
			node.HLSElement.URI = "modified-" + node.HLSElement.URI
		}
	}

	err := list.ModifyNodesBetween(startCondition, endCondition, transform)

	assert.NoError(t, err)
	assert.Equal(t, "0xFFFF", node1.HLSElement.Attrs["SCTE35-OUT"])
	assert.Equal(t, "modified-1.ts", node2.HLSElement.URI)
	assert.Equal(t, "modified-2.ts", node3.HLSElement.URI)
	assert.Equal(t, "0xFFFD", node4.HLSElement.Attrs["SCTE35-IN"])
}

func TestModifyNodesBetween_NoStartNode(t *testing.T) {
	list := internal.DoublyLinkedList{}

	node2 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "ExtInf",
			URI:  "1.ts",
		},
	}
	node3 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "ExtInf",
			URI:  "2.ts",
		},
	}
	node4 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "DateRange",
			Attrs: map[string]string{
				"SCTE35-IN": "0xFFFD",
			},
		},
	}

	list.Insert(node2)
	list.Insert(node3)
	list.Insert(node4)

	startCondition := func(node *internal.Node) bool {
		return node.HLSElement.Name == "DateRange" && node.HLSElement.Attrs["SCTE35-OUT"] != ""
	}
	endCondition := func(node *internal.Node) bool {
		return node.HLSElement.Name == "DateRange" && node.HLSElement.Attrs["SCTE35-IN"] != ""
	}
	transform := func(node *internal.Node) {
		if node.HLSElement.Name == "ExtInf" && node.HLSElement.URI != "" {
			node.HLSElement.URI = "modified-" + node.HLSElement.URI
		}
	}

	err := list.ModifyNodesBetween(startCondition, endCondition, transform)

	assert.Error(t, err)
	assert.Equal(t, "1.ts", node2.HLSElement.URI)
	assert.Equal(t, "2.ts", node3.HLSElement.URI)
}

func TestModifyNodesBetween_NoEndNode(t *testing.T) {
	list := internal.DoublyLinkedList{}

	node1 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "DateRange",
			Attrs: map[string]string{
				"SCTE35-OUT": "true",
			},
		},
	}
	node2 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "ExtInf",
			URI:  "1.ts",
		},
	}
	node3 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "ExtInf",
			URI:  "2.ts",
		},
	}

	list.Insert(node1)
	list.Insert(node2)
	list.Insert(node3)

	startCondition := func(node *internal.Node) bool {
		return node.HLSElement.Name == "DateRange" && node.HLSElement.Attrs["SCTE35-OUT"] != ""
	}
	endCondition := func(node *internal.Node) bool {
		return node.HLSElement.Name == "DateRange" && node.HLSElement.Attrs["SCTE35-IN"] != ""
	}

	transform := func(node *internal.Node) {
		if node.HLSElement.Name == "ExtInf" && node.HLSElement.URI != "" {
			node.HLSElement.URI = "modified-" + node.HLSElement.URI
		}
	}

	err := list.ModifyNodesBetween(startCondition, endCondition, transform)

	assert.Error(t, err)
	assert.Equal(t, "1.ts", node2.HLSElement.URI)
	assert.Equal(t, "2.ts", node3.HLSElement.URI)
}
