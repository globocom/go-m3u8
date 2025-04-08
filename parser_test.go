package go_m3u8_test

import (
	"io"
	"strings"
	"testing"

	m3u8 "github.com/globocom/go-m3u8"
	"github.com/globocom/go-m3u8/internal"
	"github.com/stretchr/testify/assert"
)

func TestIdentifierParser(t *testing.T) {
	playlist := "#EXTM3U"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, found := p.Find("M3u8Identifier")
	assert.True(t, found)
	assert.Equal(t, "", node.HLSElement.Attrs["#EXTM3U"])
}

func TestVersionParser(t *testing.T) {
	playlist := "#EXT-X-VERSION:3"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, found := p.Find("Version")
	assert.True(t, found)
	assert.Equal(t, "3", node.HLSElement.Attrs["#EXT-X-VERSION"])
}

func TestMediaSequenceParser(t *testing.T) {
	playlist := "#EXT-X-MEDIA-SEQUENCE:360948012"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, found := p.Find("MediaSequence")
	assert.True(t, found)
	assert.Equal(t, "360948012", node.HLSElement.Attrs["#EXT-X-MEDIA-SEQUENCE"])
}

func TestIndependentSegmentsParser(t *testing.T) {
	playlist := "#EXT-X-INDEPENDENT-SEGMENTS"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, found := p.Find("IndependentSegments")
	assert.True(t, found)
	assert.Equal(t, "", node.HLSElement.Attrs["#EXT-X-INDEPENDENT-SEGMENT"])
}

func TestTargetDurationParser(t *testing.T) {
	playlist := "#EXT-X-TARGETDURATION:7"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, ok := p.Find("TargetDuration")
	assert.True(t, ok)
	assert.Equal(t, "7", node.HLSElement.Attrs["#EXT-X-TARGETDURATION"])
}

func TestUspTimestampMapParser(t *testing.T) {
	playlist := "#USP-X-TIMESTAMP-MAP:MPEGTS=900000,LOCAL=2025-01-01T12:34:56Z"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, ok := p.Find("UspTimestampMap")
	assert.True(t, ok)
	assert.Equal(t, "900000", node.HLSElement.Attrs["MPEGTS"])
	assert.Equal(t, "2025-01-01T12:34:56Z", node.HLSElement.Attrs["LOCAL"])
}

func TestProgramDateTimeParser(t *testing.T) {
	playlist := "#EXT-X-PROGRAM-DATE-TIME:2025-01-01T12:34:56Z"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, ok := p.Find("ProgramDateTime")
	assert.True(t, ok)
	assert.Equal(t, "2025-01-01T12:34:56Z", node.HLSElement.Attrs["#EXT-X-PROGRAM-DATE-TIME"])
}

func TestDateRangeParser(t *testing.T) {
	playlist := "#EXT-X-DATERANGE:SCTE35-OUT=0xFF0000,ID=\"break1\",START-DATE=\"2025-01-01T00:00:00Z\""

	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, found := p.Find("DateRange")
	assert.True(t, found)
	assert.Equal(t, p.CurrentDateRange, node.HLSElement.ToDateRangeType(p.MediaSequence, p.SegmentsCounter))
	assert.Equal(t, "0xFF0000", node.HLSElement.Attrs["SCTE35-OUT"])
	assert.Equal(t, "break1", node.HLSElement.Attrs["ID"])
	assert.Equal(t, "2025-01-01T00:00:00Z", node.HLSElement.Attrs["START-DATE"])
}

func TestCueOutParser(t *testing.T) {
	playlist := "#EXT-X-CUE-OUT:30"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, ok := p.Find("CueOut")
	assert.True(t, ok)
	assert.Equal(t, "30", node.HLSElement.Attrs["#EXT-X-CUE-OUT"])
}

func TestCueInParser(t *testing.T) {
	playlist := strings.Join([]string{
		"#EXT-X-DATERANGE:SCTE35-OUT=0xFF0000,ID=\"break1\",START-DATE=\"2025-01-01T00:00:00Z\"",
		"#EXT-X-CUE-IN",
	}, "\n")

	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, ok := p.Find("CueIn")
	assert.True(t, ok)
	assert.Equal(t, p.CurrentDateRange, &internal.DateRange{})
	assert.Equal(t, "", node.HLSElement.Attrs["#EXT-X-CUE-IN"])
}

func TestDiscontinuityParser(t *testing.T) {
	playlist := "#EXT-X-DISCONTINUITY"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, ok := p.Find("Discontinuity")
	assert.True(t, ok)
	assert.Equal(t, "", node.HLSElement.Attrs["#EXT-X-DISCONTINUITY"])
}

func TestExtInfParser(t *testing.T) {
	playlist := "#EXTINF:4.8, no desc"
	p, err := setupPlaylist(playlist)

	assert.NoError(t, err)
	assert.Equal(t, 4.8, p.CurrentSegment.Duration)
}

func TestStreamInfParser(t *testing.T) {
	playlist := "#EXT-X-STREAM-INF:BANDWIDTH=300000,CODECS=\"avc1.42c00a\",RESOLUTION=1280x720"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)
	assert.Equal(t, "300000", p.CurrentStreamInf.Bandwidth)
	assert.Equal(t, []string{"avc1.42c00a"}, p.CurrentStreamInf.Codecs)
	assert.Equal(t, "1280x720", p.CurrentStreamInf.Resolution)
}

func TestHandleNonTags_Segments(t *testing.T) {
	playlist := `#EXTINF:4.8, no desc
              1.ts`
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)
	assert.Nil(t, p.CurrentSegment)

	node, found := p.Find("ExtInf")
	assert.True(t, found)
	assert.Equal(t, "1.ts", node.HLSElement.URI)
	assert.Equal(t, "4.8", node.HLSElement.Attrs["Duration"])
}

func TestHandleNonTags_StreamInf(t *testing.T) {
	playlist := `#EXT-X-STREAM-INF:BANDWIDTH=206000,AVERAGE-BANDWIDTH=187000,CODECS="mp4a.40.2,avc1.64001F",RESOLUTION=256x144,FRAME-RATE=30
              channel-audio_1=96000-video=80000.m3u8`
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)
	assert.Nil(t, p.CurrentStreamInf)

	node, found := p.Find("StreamInf")
	assert.True(t, found)
	assert.Equal(t, "channel-audio_1=96000-video=80000.m3u8", node.HLSElement.URI)
	assert.Equal(t, "206000", node.HLSElement.Attrs["BANDWIDTH"])
	assert.Equal(t, "187000", node.HLSElement.Attrs["AVERAGE-BANDWIDTH"])
	assert.Equal(t, "mp4a.40.2,avc1.64001F", node.HLSElement.Attrs["CODECS"])
	assert.Equal(t, "256x144", node.HLSElement.Attrs["RESOLUTION"])
	assert.Equal(t, "30", node.HLSElement.Attrs["FRAME-RATE"])
}

func TestHandleNonTags_Comment(t *testing.T) {
	playlist := `## splice_insert(SCTE35-IN matches Auto Return Mode)`
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, found := p.Find("Comment")
	assert.True(t, found)
	assert.Equal(t, "", node.HLSElement.URI)
	assert.Equal(t, "## splice_insert(SCTE35-IN matches Auto Return Mode)", node.HLSElement.Attrs["Comment"])
}

func setupPlaylist(input string) (*m3u8.Playlist, error) {
	return m3u8.ParsePlaylist(io.NopCloser(strings.NewReader(input)))
}
