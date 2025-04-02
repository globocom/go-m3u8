package go_m3u8_test

import (
	"fmt"
	"os"
	"testing"

	m3u8 "github.com/globocom/go-m3u8"
	"github.com/globocom/go-m3u8/internal"
	"github.com/stretchr/testify/assert"
)

func TestVersionValue(t *testing.T) {
	file, err := os.Open("./testdata/default/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	version := playlist.VersionValue()
	assert.Equal(t, version, "3")
}

func TestVersion(t *testing.T) {
	file, err := os.Open("./testdata/default/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	node, found := playlist.Version()
	assert.True(t, found)
	assert.NotNil(t, node)
	assert.Equal(t, node.Attrs["#EXT-X-VERSION"], "3")
}

func TestMediaSequenceValue(t *testing.T) {
	file, err := os.Open("./testdata/default/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	mediaSequence := playlist.MediaSequenceValue()
	assert.Equal(t, mediaSequence, "360948012")
}

func TestMediaSequence(t *testing.T) {
	file, err := os.Open("./testdata/default/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	node, found := playlist.MediaSequenceTag()
	assert.True(t, found)
	assert.NotNil(t, node)
	assert.Equal(t, node.Attrs["#EXT-X-MEDIA-SEQUENCE"], "360948012")
}

func TestBreaks(t *testing.T) {
	file, err := os.Open("./testdata/default/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	nodes := playlist.Breaks()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 1)
	for _, node := range nodes {
		assert.Equal(t, node.Attrs["SCTE35-OUT"], "0xFC3025000000000BB802FFF01405000000017FEFFF1A7B3B607E005288E8000100000000F799F45B")
	}
}

func TestVariants(t *testing.T) {
	file, err := os.Open("./testdata/default/master.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	nodes := playlist.Variants()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 8)
}

func TestSegments(t *testing.T) {
	file, err := os.Open("./testdata/default/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	nodes := playlist.Segments()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 16)
}

func TestReplaceBreaksURI(t *testing.T) {
	playlist := &m3u8.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{},
	}

	node1 := &internal.Node{
		Name: "DateRange",
		Attrs: map[string]string{
			"SCTE35-OUT": "0xFFFF",
		},
	}
	node2 := &internal.Node{
		Name: "ExtInf",
		URI:  "1.ts",
	}
	node3 := &internal.Node{
		Name: "ExtInf",
		URI:  "2.ts",
	}
	node4 := &internal.Node{
		Name: "DateRange",
		Attrs: map[string]string{
			"SCTE35-IN": "0xFFFD",
		},
	}

	playlist.Insert(node1)
	playlist.Insert(node2)
	playlist.Insert(node3)
	playlist.Insert(node4)

	transform := func(uri string) string {
		return fmt.Sprintf("ad-server-bucket-%s", uri)
	}

	err := playlist.ReplaceBreaksURI(transform)
	assert.NoError(t, err)
	assert.Equal(t, "ad-server-bucket-1.ts", node2.URI)
	assert.Equal(t, "ad-server-bucket-2.ts", node3.URI)
	assert.Equal(t, "0xFFFF", node1.Attrs["SCTE35-OUT"])
	assert.Equal(t, "0xFFFD", node4.Attrs["SCTE35-IN"])
}
