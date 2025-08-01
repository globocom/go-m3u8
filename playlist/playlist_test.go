package playlist_test

import (
	"os"
	"testing"

	m3u8 "github.com/globocom/go-m3u8"
	"github.com/globocom/go-m3u8/internal"
	"github.com/stretchr/testify/assert"
)

func TestVersionValue(t *testing.T) {
	file, _ := os.Open("./../mocks/media/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	version := playlist.VersionValue()
	assert.Equal(t, version, "3")
}

func TestVersionTag(t *testing.T) {
	file, _ := os.Open("./../mocks/media/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	node, found := playlist.VersionTag()
	assert.True(t, found)
	assert.NotNil(t, node)
	assert.Equal(t, node.HLSElement.Attrs["#EXT-X-VERSION"], "3")
}

func TestMediaSequenceValue(t *testing.T) {
	file, _ := os.Open("./../mocks/media/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	mediaSequence := playlist.MediaSequenceValue()
	assert.Equal(t, mediaSequence, "364042169")
}

func TestMediaSequenceTag(t *testing.T) {
	file, _ := os.Open("./../mocks/media/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	node, found := playlist.MediaSequenceTag()
	assert.True(t, found)
	assert.NotNil(t, node)
	assert.Equal(t, node.HLSElement.Attrs["#EXT-X-MEDIA-SEQUENCE"], "364042169")
}

func TestDiscontinuitySequenceValue(t *testing.T) {
	file, _ := os.Open("./../mocks/media/withDiscontinuity.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	discSequence := playlist.DiscontinuitySequenceValue()
	assert.Equal(t, discSequence, "87498")
}

func TestDiscontinuitySequenceTag(t *testing.T) {
	file, _ := os.Open("./../mocks/media/withDiscontinuity.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	node, found := playlist.DiscontinuitySequenceTag()
	assert.True(t, found)
	assert.NotNil(t, node)
	assert.Equal(t, node.HLSElement.Attrs["#EXT-X-DISCONTINUITY-SEQUENCE"], "87498")
}

func TestVariants(t *testing.T) {
	file, _ := os.Open("./../mocks/multivariant/multivariant.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	nodes := playlist.Variants()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 8)
}

func TestSegments(t *testing.T) {
	file, _ := os.Open("./../mocks/media/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	nodes := playlist.Segments()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 27)
}

func TestEncryptionTags(t *testing.T) {
	file, _ := os.Open("./../mocks/media/encryption/withAES128.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	nodes := playlist.EncryptionTags()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 1)
	assert.Equal(t, nodes[0].HLSElement.Attrs["METHOD"], "AES-128")
	assert.Equal(t, nodes[0].HLSElement.Attrs["URI"], "https://example.com/keys/key1.bin")
	assert.Equal(t, nodes[0].HLSElement.Attrs["IV"], "0x0123456789abcdef0123456789abcdef")
}

func TestCueOutEvents(t *testing.T) {
	file, _ := os.Open("./../mocks/media/scte35/withCompleteAdBreak.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	nodes := playlist.CueOutEvents()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 1)
	assert.Equal(t, nodes[0].HLSElement.Attrs["#EXT-X-CUE-OUT"], "60.033333")
}

func TestCueInEvents(t *testing.T) {
	file, _ := os.Open("./../mocks/media/scte35/withCompleteAdBreak.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	nodes := playlist.CueInEvents()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 1)
	assert.Equal(t, nodes[0].HLSElement.Attrs["#EXT-X-CUE-IN"], "")
}

func TestBreaks(t *testing.T) {
	file, _ := os.Open("./../mocks/media/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	nodes := playlist.Breaks()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 1)
	assert.Equal(t, nodes[0].HLSElement.Attrs["SCTE35-OUT"], "0xFC3025000000000BB802FFF01405000000017FEFFFE86CE9387E0052717800010000000097E91FE5")
}

func TestSCTE35InTags(t *testing.T) {
	file, _ := os.Open("./../mocks/media/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	nodes := playlist.SCTE35InTags()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 1)
	assert.Equal(t, nodes[0].HLSElement.Attrs["SCTE35-IN"], "0xFC3025000000000BB802FFF01405000000017F6FFFE8BF5AB07E00000000000100000000E7CF6C5A")
}

func TestFindSegmentInsideAdBreak(t *testing.T) {
	file, _ := os.Open("./../mocks/media/media.m3u8")
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
	file, _ := os.Open("./../mocks/media/media.m3u8")
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
