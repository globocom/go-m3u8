package go_m3u8_test

import (
	"testing"

	m3u8 "github.com/globocom/go-m3u8"
	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
	"github.com/stretchr/testify/assert"
)

func TestM3u8IdentifierEncoder(t *testing.T) {
	node := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "M3u8Identifier",
		},
	}

	playlist := &pl.Playlist{
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
		HLSElement: &internal.HLSElement{
			Name: "Version",
			Attrs: map[string]string{
				"#EXT-X-VERSION": "3",
			},
		},
	}
	playlist := &pl.Playlist{
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
		HLSElement: &internal.HLSElement{
			Name: "ExtInf",
			Attrs: map[string]string{
				"Duration": "4.8",
				"Title":    " no desc",
			},
			URI: "1.ts",
		},
	}
	playlist := &pl.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node,
			Tail: node,
		},
	}

	p, err := m3u8.EncodePlaylist(playlist)

	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "#EXTINF:4.8, no desc\n1.ts\n", p)
}
func TestStreamInfEncoder(t *testing.T) {
	node := &internal.Node{
		HLSElement: &internal.HLSElement{
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
		},
	}
	playlist := &pl.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node,
			Tail: node,
		},
	}

	expectedPlaylist := `#EXT-X-STREAM-INF:BANDWIDTH=206000,AVERAGE-BANDWIDTH=187000,CODECS="mp4a.40.2,avc1.64001F",RESOLUTION="1280x720",FRAME-RATE=30,AUDIO="audio",CLOSED-CAPTIONS="cc",SUBTITLES="subtitle"` + "\n" + "playlist.m3u8\n"

	p, err := m3u8.EncodePlaylist(playlist)

	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, expectedPlaylist, p)
}

func TestCommentEncoder(t *testing.T) {
	node := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "Comment",
			Attrs: map[string]string{
				"Comment": "## splice_insert(SCTE35-IN matches Auto Return Mode)",
			},
		},
	}
	playlist := &pl.Playlist{
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
		HLSElement: &internal.HLSElement{
			Name: "DateRange",
			Attrs: map[string]string{
				"ID":                  "ID1",
				"START-DATE":          "2025-01-01T16:16:22.933333Z",
				"PLANNED-DURATION":    "60.1",
				"SCTE35-OUT":          "0x0001",
				"CUSTOM-ATTR":         "custom-value",
				"ANOTHER-CUSTOM-ATTR": "another-custom-value",
			},
		},
	}
	playlist := &pl.Playlist{
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
		HLSElement: &internal.HLSElement{
			Name: "DateRange",
			Attrs: map[string]string{
				"ID":         "ID2",
				"START-DATE": "2025-01-01T16:16:22.933333Z",
				"END-DATE":   "2025-01-02T16:17:23.033333Z",
				"DURATION":   "60.1",
				"SCTE35-IN":  "0x0002",
			},
		},
	}
	playlist = &pl.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node2,
			Tail: node2,
		},
	}
	expectedPlaylist := `#EXT-X-DATERANGE:ID="ID2",START-DATE="2025-01-01T16:16:22.933333Z",END-DATE="2025-01-02T16:17:23.033333Z",DURATION=60.1,SCTE35-IN=0x0002` + "\n"

	p, err = m3u8.EncodePlaylist(playlist)

	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, expectedPlaylist, p)
}

func TestIndependentSegmentsEncoder(t *testing.T) {
	node := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "IndependentSegments",
		},
	}
	playlist := &pl.Playlist{
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
		HLSElement: &internal.HLSElement{
			Name: "Discontinuity",
		},
	}
	playlist := &pl.Playlist{
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
		HLSElement: &internal.HLSElement{
			Name: "UspTimestampMap",
			Attrs: map[string]string{
				"MPEGTS": "90000",
				"LOCAL":  "2025-01-01T00:00:00Z",
			},
		},
	}
	playlist := &pl.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node,
			Tail: node,
		},
	}

	expectedPlaylist := `#USP-X-TIMESTAMP-MAP:MPEGTS=90000,LOCAL=2025-01-01T00:00:00Z` + "\n"

	p, err := m3u8.EncodePlaylist(playlist)

	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, expectedPlaylist, p)
}

func TestCueOutEncoder(t *testing.T) {
	node := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "CueOut",
			Attrs: map[string]string{
				"#EXT-X-CUE-OUT": "30",
			},
		},
	}
	playlist := &pl.Playlist{
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
		HLSElement: &internal.HLSElement{
			Name: "CueIn",
			Attrs: map[string]string{
				"#EXT-X-CUE-IN": "",
			},
		},
	}
	playlist := &pl.Playlist{
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

func TestDiscontinuitySequenceEncoder(t *testing.T) {
	node := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "DiscontinuitySequence",
			Attrs: map[string]string{
				"#EXT-X-DISCONTINUITY-SEQUENCE": "18",
			},
		},
	}
	playlist := &pl.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node,
			Tail: node,
		},
	}

	p, err := m3u8.EncodePlaylist(playlist)

	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "#EXT-X-DISCONTINUITY-SEQUENCE:18\n", p)
}

func TestEncodeMasterPlaylist(t *testing.T) {
	node1 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "M3u8Identifier",
		},
	}
	node2 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "Version",
			Attrs: map[string]string{
				"#EXT-X-VERSION": "3",
			},
		},
	}
	node3 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "Comment",
			Attrs: map[string]string{
				"Comment": "## Created with Unified Streaming Platform (version=1.11.23-28141)",
			},
		},
	}

	node4 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "Comment",
			Attrs: map[string]string{
				"Comment": "# variants",
			},
		},
	}
	node5 := &internal.Node{
		HLSElement: &internal.HLSElement{
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
		},
	}

	playlist := &pl.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node1,
			Tail: node1,
		},
	}

	expectedPlaylist := `#EXTM3U
#EXT-X-VERSION:3
## Created with Unified Streaming Platform (version=1.11.23-28141)
# variants
#EXT-X-STREAM-INF:BANDWIDTH=206000,AVERAGE-BANDWIDTH=187000,CODECS="mp4a.40.2,avc1.64001F",RESOLUTION="1280x720",FRAME-RATE=30,AUDIO="audio",CLOSED-CAPTIONS="cc",SUBTITLES="subtitle"
playlist.m3u8
`

	playlist.Insert(node1)
	playlist.Insert(node2)
	playlist.Insert(node3)
	playlist.Insert(node4)
	playlist.Insert(node5)

	p, err := m3u8.EncodePlaylist(playlist)

	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, expectedPlaylist, p)
}

func TestEncodeMediaPlaylist(t *testing.T) {
	node1 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "M3u8Identifier",
		},
	}
	node2 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "Version",
			Attrs: map[string]string{
				"#EXT-X-VERSION": "3",
			},
		},
	}
	node3 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "Comment",
			Attrs: map[string]string{
				"Comment": "## Created with Unified Streaming Platform (version=1.11.23-28141)",
			},
		},
	}

	node4 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "MediaSequence",
			Attrs: map[string]string{
				"#EXT-X-MEDIA-SEQUENCE": "360948012",
			},
		},
	}
	node5 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "IndependentSegments",
			Attrs: map[string]string{
				"#EXT-X-INDEPENDENT-SEGMENTS": "",
			},
		},
	}
	node6 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "TargetDuration",
			Attrs: map[string]string{
				"#EXT-X-TARGETDURATION": "7",
			},
		},
	}
	node7 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "UspTimestampMap",
			Attrs: map[string]string{
				"MPEGTS": "5048974016",
				"LOCAL":  "2024-11-25T16:00:53.200000Z",
			},
		},
	}
	node8 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "ProgramDateTime",
			Attrs: map[string]string{
				"#EXT-X-PROGRAM-DATE-TIME": "2024-11-25T16:00:53.200000Z",
			},
		},
	}
	node9 := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "ExtInf",
			Attrs: map[string]string{
				"Duration": "4.8",
				"Title":    " no desc",
			},
			URI: "1.ts",
		},
	}

	playlist := &pl.Playlist{
		DoublyLinkedList: &internal.DoublyLinkedList{
			Head: node1,
			Tail: node1,
		},
	}

	expectedPlaylist := `#EXTM3U
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
	assert.Equal(t, expectedPlaylist, p)
}
