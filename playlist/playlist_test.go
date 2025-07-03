package playlist_test

import (
	"fmt"
	"os"
	"testing"

	m3u8 "github.com/globocom/go-m3u8"
	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
	"github.com/stretchr/testify/assert"
)

func TestVersionValue(t *testing.T) {
	file, _ := os.Open("./../testdata/media/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	version := playlist.VersionValue()
	assert.Equal(t, version, "3")
}

func TestVersion(t *testing.T) {
	file, _ := os.Open("./../testdata/media/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	node, found := playlist.Version()
	assert.True(t, found)
	assert.NotNil(t, node)
	assert.Equal(t, node.HLSElement.Attrs["#EXT-X-VERSION"], "3")
}

func TestMediaSequenceValue(t *testing.T) {
	file, _ := os.Open("./../testdata/media/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	mediaSequence := playlist.MediaSequenceValue()
	assert.Equal(t, mediaSequence, "364042169")
}

func TestMediaSequence(t *testing.T) {
	file, _ := os.Open("./../testdata/media/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	node, found := playlist.MediaSequenceTag()
	assert.True(t, found)
	assert.NotNil(t, node)
	assert.Equal(t, node.HLSElement.Attrs["#EXT-X-MEDIA-SEQUENCE"], "364042169")
}

func TestBreaks(t *testing.T) {
	file, _ := os.Open("./../testdata/media/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	nodes := playlist.Breaks()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 1)
	for _, node := range nodes {
		assert.Equal(t, node.HLSElement.Attrs["SCTE35-OUT"], "0xFC3025000000000BB802FFF01405000000017FEFFFE86CE9387E0052717800010000000097E91FE5")
	}
}

func TestVariants(t *testing.T) {
	file, _ := os.Open("./../testdata/multivariant/multivariant.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	nodes := playlist.Variants()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 8)
}

func TestSegments(t *testing.T) {
	file, _ := os.Open("./../testdata/media/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	nodes := playlist.Segments()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 27)
}

func TestReplaceBreaksURI(t *testing.T) {
	playlist := &pl.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{},
	}

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

	playlist.Insert(node1)
	playlist.Insert(node2)
	playlist.Insert(node3)
	playlist.Insert(node4)

	transform := func(uri string) string {
		return fmt.Sprintf("ad-server-bucket-%s", uri)
	}

	err := playlist.ReplaceBreaksURI(transform)
	assert.NoError(t, err)
	assert.Equal(t, "ad-server-bucket-1.ts", node2.HLSElement.URI)
	assert.Equal(t, "ad-server-bucket-2.ts", node3.HLSElement.URI)
	assert.Equal(t, "0xFFFF", node1.HLSElement.Attrs["SCTE35-OUT"])
	assert.Equal(t, "0xFFFD", node4.HLSElement.Attrs["SCTE35-IN"])
}

func TestFindSegmentInsideAdBreak(t *testing.T) {
	file, _ := os.Open("./../testdata/media/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	// #EXTINF:5.3, no desc
	// channel-audio_1=96000-video=3442944-364042175.ts
	adBreakSegment := playlist.Segments()[6]
	assert.Equal(t, adBreakSegment.HLSElement.URI, "channel-audio_1=96000-video=3442944-364042175.ts")

	// #EXT-X-DATERANGE:ID="1-1747402436",START-DATE="2025-05-16T13:33:56.266666Z",PLANNED-DURATION=60.033333,SCTE35-OUT=0xFC3025000000000BB802FFF01405000000017FEFFFE86CE9387E0052717800010000000097E91FE5
	expectedAdBreak := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "DateRange",
			Attrs: map[string]string{
				"ID":               "1-1747402436",
				"START-DATE":       "2025-05-16T13:33:56.266666Z",
				"PLANNED-DURATION": "60.033333",
				"SCTE35-OUT":       "0xFC3025000000000BB802FFF01405000000017FEFFFE86CE9387E0052717800010000000097E91FE5",
			},
		},
	}

	adBreak, _ := playlist.FindNodeInsideAdBreak(adBreakSegment)

	assert.NotNil(t, adBreak)
	assert.Equal(t, adBreak.HLSElement.Name, expectedAdBreak.HLSElement.Name)
	assert.Equal(t, adBreak.HLSElement.URI, expectedAdBreak.HLSElement.URI)
	assert.Equal(t, adBreak.HLSElement.Attrs["SCTE35-OUT"], expectedAdBreak.HLSElement.Attrs["SCTE35-OUT"])
}

func TestFindSegmentOutsideAdBreak(t *testing.T) {
	file, _ := os.Open("./../testdata/media/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	// after break
	// #EXTINF:2.8666, no desc
	// channel-audio_1=96000-video=3442944-364042187.ts
	afterBreakSegment := playlist.Segments()[18]
	assert.Equal(t, afterBreakSegment.HLSElement.URI, "channel-audio_1=96000-video=3442944-364042187.ts")

	// before break
	// #EXTINF:4.3, no desc
	// channel-audio_1=96000-video=3442944-364042174.ts
	beforeBreakSegment := playlist.Segments()[5]
	assert.Equal(t, beforeBreakSegment.HLSElement.URI, "channel-audio_1=96000-video=3442944-364042174.ts")

	afterAdBreak, _ := playlist.FindNodeInsideAdBreak(afterBreakSegment)
	beforeAdBreak, _ := playlist.FindNodeInsideAdBreak(beforeBreakSegment)

	assert.Nil(t, afterAdBreak)
	assert.Nil(t, beforeAdBreak)
}
