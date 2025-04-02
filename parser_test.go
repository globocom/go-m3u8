package go_m3u8_test

import (
	"io"
	"strings"
	"testing"
	"time"

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
	assert.Equal(t, "", node.Attrs["#EXTM3U"])
}

func TestVersionParser(t *testing.T) {
	playlist := "#EXT-X-VERSION:3"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, found := p.Find("Version")
	assert.True(t, found)
	assert.Equal(t, "3", node.Attrs["#EXT-X-VERSION"])
}

func TestMediaSequenceParser(t *testing.T) {
	playlist := "#EXT-X-MEDIA-SEQUENCE:360948012"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, found := p.Find("MediaSequence")
	assert.True(t, found)
	assert.Equal(t, "360948012", node.Attrs["#EXT-X-MEDIA-SEQUENCE"])
}

func TestIndependentSegmentsParser(t *testing.T) {
	playlist := "#EXT-X-INDEPENDENT-SEGMENTS"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, found := p.Find("IndependentSegments")
	assert.True(t, found)
	assert.Equal(t, "", node.Attrs["#EXT-X-INDEPENDENT-SEGMENT"])
}

func TestTargetDurationParser(t *testing.T) {
	playlist := "#EXT-X-TARGETDURATION:7"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, ok := p.Find("TargetDuration")
	assert.True(t, ok)
	assert.Equal(t, "7", node.Attrs["#EXT-X-TARGETDURATION"])
}

func TestUspTimestampMapParser(t *testing.T) {
	playlist := "#USP-X-TIMESTAMP-MAP:MPEGTS=900000,LOCAL=2025-01-01T12:34:56Z"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, ok := p.Find("UspTimestampMap")
	assert.True(t, ok)
	assert.Equal(t, "900000", node.Attrs["MPEGTS"])
	assert.Equal(t, "2025-01-01T12:34:56Z", node.Attrs["LOCAL"])
}

func TestProgramDateTimeParser(t *testing.T) {
	playlist := "#EXT-X-PROGRAM-DATE-TIME:2025-01-01T12:34:56Z"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, ok := p.Find("ProgramDateTime")
	assert.True(t, ok)
	assert.Equal(t, "2025-01-01T12:34:56Z", node.Attrs["#EXT-X-PROGRAM-DATE-TIME"])
}

func TestDateRangeParser(t *testing.T) {
	playlist := "#EXT-X-DATERANGE:SCTE35-OUT=0xFF0000,ID=\"break1\",START-DATE=\"2025-01-01T00:00:00Z\""

	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, found := p.Find("DateRange")
	assert.True(t, found)
	assert.Equal(t, p.CurrentDateRange, node)
	assert.Equal(t, "0xFF0000", node.Attrs["SCTE35-OUT"])
	assert.Equal(t, "break1", node.Attrs["ID"])
	assert.Equal(t, "2025-01-01T00:00:00Z", node.Attrs["START-DATE"])
}

func TestDateRangeParserObjects(t *testing.T) {
	playlist := strings.Join([]string{
		"#EXT-X-MEDIA-SEQUENCE:360948012",
		"#EXT-X-PROGRAM-DATE-TIME:2025-01-01T12:34:00Z",
		"#EXTINF:4.8, no desc",
		"0.ts",
		"#EXTINF:4.8, no desc",
		"1.ts",
		"#EXT-X-DATERANGE:SCTE35-OUT=0xFF0000,ID=\"break1\",START-DATE=\"2025-01-01T12:43:06Z\",PLANNED-DURATION=60",
		"#EXT-X-CUE-OUT:48",
		"#EXTINF:4.8, no desc",
		"3.ts",
	}, "\n")

	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, found := p.Find("DateRange")
	dateRange, ok := node.Object.(*m3u8.DateRange)

	assert.True(t, found)
	assert.True(t, ok)
	assert.Equal(t, dateRange.MediaSequence, 360948014)
	assert.Equal(t, dateRange.PlannedDuration, float64(60))
	assert.Equal(t, dateRange.StartDate, time.Date(2025, 01, 01, 12, 43, 6, 0, time.UTC))
}

func TestCueOutParser(t *testing.T) {
	playlist := "#EXT-X-CUE-OUT:30"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, ok := p.Find("CueOut")
	assert.True(t, ok)
	assert.Equal(t, "30", node.Attrs["#EXT-X-CUE-OUT"])
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
	assert.Equal(t, p.CurrentDateRange, &internal.Node{})
	assert.Equal(t, "", node.Attrs["#EXT-X-CUE-IN"])
}

func TestDiscontinuityParser(t *testing.T) {
	playlist := "#EXT-X-DISCONTINUITY"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, ok := p.Find("Discontinuity")
	assert.True(t, ok)
	assert.Equal(t, "", node.Attrs["#EXT-X-DISCONTINUITY"])
}

func TestExtInfParser(t *testing.T) {
	playlist := "#EXTINF:4.8, no desc"
	p, err := setupPlaylist(playlist)

	assert.NoError(t, err)
	assert.Equal(t, 4.8, p.CurrentSegment.Duration)
}

func TestExtInfParserObjects(t *testing.T) {
	playlistLines := []string{
		"#EXT-X-MEDIA-SEQUENCE:360948012",
		"#EXT-X-PROGRAM-DATE-TIME:2025-01-01T12:34:56.0000000Z",
		"#EXTINF:4.8, no desc",
		"0.ts",
		"#EXTINF:4.8, no desc",
		"1.ts",
		"#EXT-X-DATERANGE:ID=\"id\",START-DATE=\"2025-01-01T12:35:05.6000000Z\",PLANNED-DURATION=48,SCTE35-OUT=0xFC3025000000000BB800FFF01405C00000007FEFFFBB373098FE0029560895740002000058CF9EC9",
		"#EXT-X-CUE-OUT:48",
		"#EXT-X-PROGRAM-DATE-TIME:2025-04-01T12:35:05.6000000Z",
		"#EXTINF:4.8, no desc",
		"2.ts",
	}

	playlist := strings.Join(playlistLines, "\n")
	p, err := setupPlaylist(playlist)

	firstSegment := p.Segments()[0].Object.(*m3u8.Segment)
	secondSegment := p.Segments()[1].Object.(*m3u8.Segment)
	thirdSegment := p.Segments()[2].Object.(*m3u8.Segment)

	assert.NoError(t, err)
	assert.Equal(t, 14.4, p.DVR)
	assert.Equal(t, 3, p.SegmentsCounter)
	assert.Equal(t, 360948012, p.MediaSequence)

	assert.Equal(t, 4.8, firstSegment.Duration)
	assert.Equal(t, time.Date(2025, 01, 01, 12, 34, 56, 0, time.UTC), firstSegment.ProgramDateTime)
	assert.Equal(t, 360948012, firstSegment.MediaSequence)
	assert.Nil(t, firstSegment.DateRange)

	assert.Equal(t, 4.8, secondSegment.Duration)
	assert.Equal(t, time.Date(2025, 01, 01, 12, 35, 00, 800000000, time.UTC), secondSegment.ProgramDateTime)
	assert.Equal(t, 360948013, secondSegment.MediaSequence)
	assert.Nil(t, secondSegment.DateRange)

	assert.Equal(t, 4.8, thirdSegment.Duration)
	assert.Equal(t, time.Date(2025, 01, 01, 12, 35, 05, 600000000, time.UTC), thirdSegment.ProgramDateTime)
	assert.Equal(t, 360948014, thirdSegment.MediaSequence)
	assert.NotNil(t, thirdSegment.DateRange)
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
	assert.Equal(t, "1.ts", node.URI)
	assert.Equal(t, "4.8", node.Attrs["Duration"])
}

func TestHandleNonTags_StreamInf(t *testing.T) {
	playlist := `#EXT-X-STREAM-INF:BANDWIDTH=206000,AVERAGE-BANDWIDTH=187000,CODECS="mp4a.40.2,avc1.64001F",RESOLUTION=256x144,FRAME-RATE=30
              channel-audio_1=96000-video=80000.m3u8`
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)
	assert.Nil(t, p.CurrentStreamInf)

	node, found := p.Find("StreamInf")
	assert.True(t, found)
	assert.Equal(t, "channel-audio_1=96000-video=80000.m3u8", node.URI)
	assert.Equal(t, "206000", node.Attrs["BANDWIDTH"])
	assert.Equal(t, "187000", node.Attrs["AVERAGE-BANDWIDTH"])
	assert.Equal(t, "mp4a.40.2,avc1.64001F", node.Attrs["CODECS"])
	assert.Equal(t, "256x144", node.Attrs["RESOLUTION"])
	assert.Equal(t, "30", node.Attrs["FRAME-RATE"])
}

func TestHandleNonTags_Comment(t *testing.T) {
	playlist := `## splice_insert(SCTE35-IN matches Auto Return Mode)`
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, found := p.Find("Comment")
	assert.True(t, found)
	assert.Equal(t, "", node.URI)
	assert.Equal(t, "## splice_insert(SCTE35-IN matches Auto Return Mode)", node.Attrs["Comment"])
}

func setupPlaylist(input string) (*m3u8.Playlist, error) {
	return m3u8.ParsePlaylist(io.NopCloser(strings.NewReader(input)))
}
