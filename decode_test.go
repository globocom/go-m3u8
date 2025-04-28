package go_m3u8_test

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	m3u8 "github.com/globocom/go-m3u8"
	pl "github.com/globocom/go-m3u8/playlist"
	"github.com/stretchr/testify/assert"
)

type fakeSource struct{}

func (f fakeSource) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("fake error")
}

func (f fakeSource) Close() error {
	return nil
}

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

func TestCommentParser(t *testing.T) {
	playlist := `#EXTM3U
							#EXT-X-VERSION:4

							## Created with Unified Streaming Platform  (version=1.14.4-30793)
							# AUDIO groups
							#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="audio-aacl-96",LANGUAGE="pt",NAME="Portuguese",DEFAULT=YES,AUTOSELECT=YES,CHANNELS="2"`
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	nodes := p.FindAll("Comment")

	assert.True(t, len(nodes) == 2)
	assert.Equal(t, "## Created with Unified Streaming Platform  (version=1.14.4-30793)", nodes[0].HLSElement.Attrs["Comment"])
	assert.Equal(t, "# AUDIO groups", nodes[1].HLSElement.Attrs["Comment"])
}

func TestMultiLineHLSElements_Segments(t *testing.T) {
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

func TestMultiLineHLSElements_StreamInf(t *testing.T) {
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

func TestCompleteMultivariantPlaylist(t *testing.T) {
	playlist := `#EXTM3U
							#EXT-X-VERSION:3
							## Created with Unified Streaming Platform  (version=1.14.4-30793)

							# variants
							#EXT-X-STREAM-INF:BANDWIDTH=759000,AVERAGE-BANDWIDTH=690000,CODECS=\"mp4a.40.2,avc1.64001F\",RESOLUTION=640x360,FRAME-RATE=30
							coelhodai-audio_1=96000-video=558976.m3u8?dvr_window_length=600
							#EXT-X-STREAM-INF:BANDWIDTH=759000,AVERAGE-BANDWIDTH=690000,CODECS=\"mp4a.40.2,avc1.64001F\",RESOLUTION=720x480,FRAME-RATE=30
							coelhodai-audio_1=96000-video=123456.m3u8?dvr_window_length=600`

	p, err := setupPlaylist(playlist)

	assert.NoError(t, err)
	assert.Equal(t, p.Head.HLSElement.Name, "M3u8Identifier")
	assert.Equal(t, p.Tail.HLSElement.Name, "StreamInf")
	assert.Equal(t, len(p.Variants()), 2)
	assert.Nil(t, p.CurrentSegment)
	assert.Nil(t, p.CurrentStreamInf)
}

func TestCompleteMediaPlaylist(t *testing.T) {
	playlist := `#EXTM3U
							#EXT-X-VERSION:3
							## Created with Unified Streaming Platform  (version=1.14.4-30793)

							# variants
							#EXT-X-STREAM-INF:BANDWIDTH=759000,AVERAGE-BANDWIDTH=690000,CODECS=\"mp4a.40.2,avc1.64001F\",RESOLUTION=640x360,FRAME-RATE=30
							coelhodai-audio_1=96000-video=558976.m3u8?dvr_window_length=600
							#EXT-X-STREAM-INF:BANDWIDTH=759000,AVERAGE-BANDWIDTH=690000,CODECS=\"mp4a.40.2,avc1.64001F\",RESOLUTION=720x480,FRAME-RATE=30
							coelhodai-audio_1=96000-video=123456.m3u8?dvr_window_length=600`

	p, err := setupPlaylist(playlist)

	assert.NoError(t, err)
	assert.Equal(t, p.Head.HLSElement.Name, "M3u8Identifier")
	assert.Equal(t, p.Tail.HLSElement.Name, "StreamInf")
	assert.NotEqual(t, p.Variants(), 0)
	assert.Nil(t, p.CurrentSegment)
	assert.Nil(t, p.CurrentStreamInf)
}

func TestParsePlaylist(t *testing.T) {
	type testCaseParams struct {
		name           string
		kind           string
		path           string
		pdt            time.Time
		source         m3u8.Source
		dvr            float64
		segmentCounter int
		error          bool
	}
	testCases := []testCaseParams{
		{
			name:   "Error parsing playlist",
			path:   "not exist",
			source: fakeSource{},
			error:  true,
		},
		{
			name:  "Parse multivariant playlist without EXT-X-VERSION tag",
			kind:  "multivariant",
			path:  "./testdata/multivariant/missingVersion.m3u8",
			error: true,
		},
		{
			name:           "Parse media playlist",
			kind:           "media",
			path:           "./testdata/media/media.m3u8",
			pdt:            time.Date(2024, 11, 25, 16, 0, 53, 200000000, time.UTC),
			dvr:            76.7998,
			segmentCounter: 16,
		},
		{
			name:           "Parse media playlist with EXT-X-DISCONTINUITY tag",
			kind:           "media",
			path:           "./testdata/media/withDiscontinuity.m3u8",
			pdt:            time.Date(2024, 11, 25, 16, 0, 53, 200000000, time.UTC),
			dvr:            76.7998,
			segmentCounter: 16,
		},
		{
			name: "Parse multivariant playlist",
			kind: "multivariant",
			path: "./testdata/multivariant/multivariant.m3u8",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			file, _ := os.Open(tc.path)
			playlist, err := m3u8.ParsePlaylist(file)
			if tc.error {
				switch tc.name {
				case "Parse multivariant playlist without EXT-X-VERSION tag":
					file, _ = os.Open(tc.path)
					playlist, err = m3u8.ParsePlaylist(file)
					assert.ErrorContains(t, err, "invalid version tag")
					assert.Nil(t, playlist)
					return
				case "Error parsing playlist":
					playlist, err = m3u8.ParsePlaylist(tc.source)
					assert.ErrorContains(t, err, "fake error")
					assert.Nil(t, playlist)
					return
				}
			}

			assert.NoError(t, err)
			assert.NotNil(t, playlist)
			assert.NotNil(t, playlist.Head)
			assert.NotNil(t, playlist.Tail)

			if tc.kind == "media" {
				assert.Equal(t, playlist.ProgramDateTime, tc.pdt)
				assert.Equal(t, playlist.SegmentsCounter, tc.segmentCounter)
				assert.Equal(t, playlist.DVR, tc.dvr)
			}
		})
	}
}

func setupPlaylist(input string) (*pl.Playlist, error) {
	return m3u8.ParsePlaylist(io.NopCloser(strings.NewReader(input)))
}
