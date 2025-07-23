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
