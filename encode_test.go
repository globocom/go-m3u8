package go_m3u8_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	m3u8 "gitlab.globoi.com/webmedia/media-delivery-advertising/go-m3u8"
	"gitlab.globoi.com/webmedia/media-delivery-advertising/go-m3u8/internal"
)

func TestM3u8IdentifierEncoder(t *testing.T) {
	node := &internal.Node{
		Name: "M3u8Identifier",
	}
	playlist := &m3u8.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node,
			Tail: node,
		},
	}

	p, err := m3u8.EncodePlaylist(playlist)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "#EXTM3U\n", p)
}

func TestVersionEncoder(t *testing.T) {
	node := &internal.Node{
		Name: "Version",
		Attrs: map[string]string{
			"#EXT-X-VERSION": "3",
		},
	}
	playlist := &m3u8.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node,
			Tail: node,
		},
	}
	p, err := m3u8.EncodePlaylist(playlist)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "#EXT-X-VERSION:3\n", p)
}

func TestExtInfEncoder(t *testing.T) {
	node := &internal.Node{
		Name: "ExtInf",
		Attrs: map[string]string{
			"Duration": "4.8",
		},
		URI: "1.ts",
	}
	playlist := &m3u8.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node,
			Tail: node,
		},
	}
	p, err := m3u8.EncodePlaylist(playlist)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "#EXTINF:4.8\n1.ts\n", p)
}
func TestStreamInfEncoder(t *testing.T) {
	node := &internal.Node{
		Name: "StreamInf",
		Attrs: map[string]string{
			"BANDWIDTH":         "206000",
			"AVERAGE-BANDWIDTH": "187000",
			"CODECS":            "mp4a.40.2,avc1.64001F",
			"RESOLUTION":        "1280x720",
			"FRAME-RATE":        "30",
			"CLOSED-CAPTIONS":   "cc",
			"SUBTITLES":         "subtitle",
			"AUDIO":             "audio",
		},
		URI: "playlist.m3u8",
	}
	playlist := &m3u8.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node,
			Tail: node,
		},
	}
	p, err := m3u8.EncodePlaylist(playlist)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	expected := `#EXT-X-STREAM-INF:BANDWIDTH=206000,AVERAGE-BANDWIDTH=187000,CODECS="mp4a.40.2,avc1.64001F",RESOLUTION="1280x720",FRAME-RATE=30,AUDIO="audio",CLOSED-CAPTIONS="cc",SUBTITLES="subtitle"` + "\n" + "playlist.m3u8\n"
	assert.Equal(t, expected, p)
}

func TestCommentEncoder(t *testing.T) {
	node := &internal.Node{
		Name: "Comment",
		Attrs: map[string]string{
			"Comment": "## splice_insert(SCTE35-IN matches Auto Return Mode)",
		},
	}
	playlist := &m3u8.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node,
			Tail: node,
		},
	}
	p, err := m3u8.EncodePlaylist(playlist)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "## splice_insert(SCTE35-IN matches Auto Return Mode)\n", p)
}

func TestDateRangeEncoder(t *testing.T) {
	node1 := &internal.Node{
		Name: "DateRange",
		Attrs: map[string]string{
			"ID":                  "ID1",
			"START-DATE":          "2025-01-01T16:16:22.933333Z",
			"PLANNED-DURATION":    "60.1",
			"SCTE35-OUT":          "0x0001",
			"CUSTOM-ATTR":         "custom-value",
			"ANOTHER-CUSTOM-ATTR": "another-custom-value",
		},
	}
	playlist := &m3u8.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node1,
			Tail: node1,
		},
	}
	p, err := m3u8.EncodePlaylist(playlist)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	expected := `#EXT-X-DATERANGE:ID="ID1",START-DATE="2025-01-01T16:16:22.933333Z",PLANNED-DURATION=60.1,SCTE35-OUT=0x0001,ANOTHER-CUSTOM-ATTR="another-custom-value",CUSTOM-ATTR="custom-value"` + "\n"
	assert.Equal(t, expected, p)

	node2 := &internal.Node{
		Name: "DateRange",
		Attrs: map[string]string{
			"ID":         "ID2",
			"START-DATE": "2025-01-01T16:16:22.933333Z",
			"END-DATE":   "2025-01-02T16:17:23.033333Z",
			"DURATION":   "60.1",
			"SCTE35-IN":  "0x0002",
		},
	}
	playlist = &m3u8.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node2,
			Tail: node2,
		},
	}
	p, err = m3u8.EncodePlaylist(playlist)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	expected = `#EXT-X-DATERANGE:ID="ID2",START-DATE="2025-01-01T16:16:22.933333Z",END-DATE="2025-01-02T16:17:23.033333Z",DURATION=60.1,SCTE35-IN=0x0002` + "\n"
	assert.Equal(t, expected, p)
}

func TestIndependentSegmentsEncoder(t *testing.T) {
	node := &internal.Node{
		Name: "IndependentSegments",
	}
	playlist := &m3u8.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node,
			Tail: node,
		},
	}
	p, err := m3u8.EncodePlaylist(playlist)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "#EXT-X-INDEPENDENT-SEGMENTS\n", p)
}

func TestDiscontinuityEncoder(t *testing.T) {
	node := &internal.Node{
		Name: "Discontinuity",
	}
	playlist := &m3u8.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node,
			Tail: node,
		},
	}
	p, err := m3u8.EncodePlaylist(playlist)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "#EXT-X-DISCONTINUITY\n", p)
}

func TestUspTimestampMapEncoder(t *testing.T) {
	node := &internal.Node{
		Name: "UspTimestampMap",
		Attrs: map[string]string{
			"MPEGTS": "90000",
			"LOCAL":  "2025-01-01T00:00:00Z",
		},
	}
	playlist := &m3u8.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node,
			Tail: node,
		},
	}
	p, err := m3u8.EncodePlaylist(playlist)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	expected := `#USP-X-TIMESTAMP-MAP:MPEGTS=90000,LOCAL=2025-01-01T00:00:00Z` + "\n"
	assert.Equal(t, expected, p)
}

func TestCueOutEncoder(t *testing.T) {
	node := &internal.Node{
		Name: "CueOut",
		Attrs: map[string]string{
			"#EXT-X-CUE-OUT": "30",
		},
	}
	playlist := &m3u8.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node,
			Tail: node,
		},
	}
	p, err := m3u8.EncodePlaylist(playlist)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "#EXT-X-CUE-OUT:30\n", p)
}

func TestCueInEncoder(t *testing.T) {
	node := &internal.Node{
		Name: "CueIn",
		Attrs: map[string]string{
			"#EXT-X-CUE-IN": "",
		},
	}
	playlist := &m3u8.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node,
			Tail: node,
		},
	}
	p, err := m3u8.EncodePlaylist(playlist)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "#EXT-X-CUE-IN\n", p)
}

func TestEncodeMasterPlaylist(t *testing.T) {
	node1 := &internal.Node{
		Name: "M3u8Identifier",
	}
	node2 := &internal.Node{
		Name: "Version",
		Attrs: map[string]string{
			"#EXT-X-VERSION": "3",
		},
	}
	node3 := &internal.Node{
		Name: "Comment",
		Attrs: map[string]string{
			"Comment": "## Created with Unified Streaming Platform (version=1.11.23-28141)",
		},
	}

	node4 := &internal.Node{
		Name: "Comment",
		Attrs: map[string]string{
			"Comment": "# variants",
		},
	}
	node5 := &internal.Node{
		Name: "StreamInf",
		Attrs: map[string]string{
			"BANDWIDTH":         "206000",
			"AVERAGE-BANDWIDTH": "187000",
			"CODECS":            "mp4a.40.2,avc1.64001F",
			"RESOLUTION":        "1280x720",
			"FRAME-RATE":        "30",
			"CLOSED-CAPTIONS":   "cc",
			"SUBTITLES":         "subtitle",
			"AUDIO":             "audio",
		},
		URI: "playlist.m3u8",
	}
	playlist := &m3u8.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node1,
			Tail: node1,
		},
	}
	playlist.Insert(node1)
	playlist.Insert(node2)
	playlist.Insert(node3)
	playlist.Insert(node4)
	playlist.Insert(node5)
	p, err := m3u8.EncodePlaylist(playlist)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	expected := `#EXTM3U
#EXT-X-VERSION:3
## Created with Unified Streaming Platform (version=1.11.23-28141)
# variants
#EXT-X-STREAM-INF:BANDWIDTH=206000,AVERAGE-BANDWIDTH=187000,CODECS="mp4a.40.2,avc1.64001F",RESOLUTION="1280x720",FRAME-RATE=30,AUDIO="audio",CLOSED-CAPTIONS="cc",SUBTITLES="subtitle"
playlist.m3u8
`
	assert.Equal(t, expected, p)
}

func TestEncodeMediaPlaylist(t *testing.T) {
	node1 := &internal.Node{
		Name: "M3u8Identifier",
	}
	node2 := &internal.Node{
		Name: "Version",
		Attrs: map[string]string{
			"#EXT-X-VERSION": "3",
		},
	}
	node3 := &internal.Node{
		Name: "Comment",
		Attrs: map[string]string{
			"Comment": "## Created with Unified Streaming Platform (version=1.11.23-28141)",
		},
	}

	node4 := &internal.Node{
		Name: "MediaSequence",
		Attrs: map[string]string{
			"#EXT-X-MEDIA-SEQUENCE": "360948012",
		},
	}
	node5 := &internal.Node{
		Name: "IndependentSegments",
		Attrs: map[string]string{
			"#EXT-X-INDEPENDENT-SEGMENTS": "",
		},
	}
	node6 := &internal.Node{
		Name: "TargetDuration",
		Attrs: map[string]string{
			"#EXT-X-TARGETDURATION": "7",
		},
	}
	node7 := &internal.Node{
		Name: "UspTimestampMap",
		Attrs: map[string]string{
			"MPEGTS": "5048974016",
			"LOCAL":  "2024-11-25T16:00:53.200000Z",
		},
	}
	node8 := &internal.Node{
		Name: "ProgramDateTime",
		Attrs: map[string]string{
			"#EXT-X-PROGRAM-DATE-TIME": "2024-11-25T16:00:53.200000Z",
		},
	}
	node9 := &internal.Node{
		Name: "ExtInf",
		Attrs: map[string]string{
			"Duration": "4.8, no desc",
		},
		URI: "1.ts",
	}

	playlist := &m3u8.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node1,
			Tail: node1,
		},
	}
	playlist.Insert(node1)
	playlist.Insert(node2)
	playlist.Insert(node3)
	playlist.Insert(node4)
	playlist.Insert(node5)
	playlist.Insert(node6)
	playlist.Insert(node7)
	playlist.Insert(node8)
	playlist.Insert(node9)
	p, err := m3u8.EncodePlaylist(playlist)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	expected := `#EXTM3U
#EXT-X-VERSION:3
## Created with Unified Streaming Platform (version=1.11.23-28141)
#EXT-X-MEDIA-SEQUENCE:360948012
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:7
#USP-X-TIMESTAMP-MAP:MPEGTS=5048974016,LOCAL=2024-11-25T16:00:53.200000Z
#EXT-X-PROGRAM-DATE-TIME:2024-11-25T16:00:53.200000Z
#EXTINF:4.8, no desc
1.ts
`
	assert.Equal(t, expected, p)
}
