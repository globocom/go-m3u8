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

func TestVariableDefineTag(t *testing.T) {
	file, _ := os.Open("./../mocks/multivariant/withQueryParam.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	node, found := playlist.VariableDefineTag()
	assert.True(t, found)
	assert.NotNil(t, node)
	assert.Equal(t, node.HLSElement.Attrs["QUERYPARAM"], "stream_id")
}

func TestVariants(t *testing.T) {
	file, _ := os.Open("./../mocks/multivariant/multivariant.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	nodes := playlist.Variants()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 8)
}

func TestMediaGroups(t *testing.T) {
	file, _ := os.Open("./../mocks/multivariant/withClosedCaptionGroups.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	mediaGroups := playlist.MediaGroups()
	assert.Len(t, mediaGroups, 4)
	assert.Equal(t, mediaGroups[0].HLSElement.Attrs["TYPE"], "AUDIO")
	assert.Equal(t, mediaGroups[3].HLSElement.Attrs["TYPE"], "CLOSED-CAPTIONS")
}

func TestKeyframes(t *testing.T) {
	file, _ := os.Open("./../mocks/multivariant/withAudioGroups.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	keyframes := playlist.Keyframes()
	assert.Len(t, keyframes, 4)
	assert.Equal(t, keyframes[0].HLSElement.Attrs["URI"], "keyframes/channel-video=558976.m3u8?dvr_window_length=120")
	assert.Equal(t, keyframes[3].HLSElement.Attrs["URI"], "keyframes/channel-video=3442944.m3u8?dvr_window_length=120")
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

func TestComment(t *testing.T) {
	file, _ := os.Open("./../mocks/media/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	node := playlist.Comment("## Created with Unified Streaming Platform")
	assert.NotNil(t, node)
	assert.Equal(t, node.HLSElement.Attrs["Comment"], "## Created with Unified Streaming Platform  (version=1.14.4-30793)")
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

func TestFindLastAdBreak(t *testing.T) {
	file, _ := os.Open("./../mocks/media/withMultipleBreaks.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	// Expected: last Ad Break present in the manifest
	expectedAdBreak := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "DateRange",
			Attrs: map[string]string{
				"START-DATE":       "2025-05-23T19:28:34.299999Z",
				"PLANNED-DURATION": "20",
				"SCTE35-OUT":       "0xFC3025000000000BB800FFF01405F000A7E57FEFFED43025D0FE001B7740000101010000CE90B6B7",
			},
		},
	}

	lastAdBreak, found := playlist.FindLastAdBreak()

	assert.True(t, found)
	assert.NotNil(t, lastAdBreak)
	assert.Equal(t, expectedAdBreak.HLSElement.Name, lastAdBreak.HLSElement.Name)
	assert.Equal(t, expectedAdBreak.HLSElement.Attrs["START-DATE"], lastAdBreak.HLSElement.Attrs["START-DATE"])
	assert.Equal(t, expectedAdBreak.HLSElement.Attrs["PLANNED-DURATION"], lastAdBreak.HLSElement.Attrs["PLANNED-DURATION"])
	assert.Equal(t, expectedAdBreak.HLSElement.Attrs["SCTE35-OUT"], lastAdBreak.HLSElement.Attrs["SCTE35-OUT"])
}

func TestHasDuplicateAdBreak(t *testing.T) {
	file, _ := os.Open("./../mocks/media/withDuplicateBreaks.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	adBreaks := playlist.Breaks()
	assert.GreaterOrEqual(t, len(adBreaks), 2)

	// Expected: Last two Ad Breaks are duplicated
	last := adBreaks[len(adBreaks)-1]
	previous := adBreaks[len(adBreaks)-2]

	expectedAdBreak := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "DateRange",
			Attrs: map[string]string{
				"START-DATE":       last.HLSElement.Attrs["START-DATE"],
				"PLANNED-DURATION": last.HLSElement.Attrs["PLANNED-DURATION"],
			},
		},
	}

	isDuplicate := playlist.HasDuplicateAdBreak()
	assert.True(t, isDuplicate)

	assert.Equal(t, expectedAdBreak.HLSElement.Attrs["START-DATE"], previous.HLSElement.Attrs["START-DATE"])
	assert.Equal(t, expectedAdBreak.HLSElement.Attrs["PLANNED-DURATION"], previous.HLSElement.Attrs["PLANNED-DURATION"])
}
