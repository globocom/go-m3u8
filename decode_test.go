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
	"github.com/globocom/go-m3u8/tags"
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

func TestDiscontinuitySequenceParser(t *testing.T) {
	playlist := "#EXT-X-DISCONTINUITY-SEQUENCE:18"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, found := p.Find("DiscontinuitySequence")
	assert.True(t, found)
	assert.Equal(t, "18", node.HLSElement.Attrs["#EXT-X-DISCONTINUITY-SEQUENCE"])
}

func TestIndependentSegmentsParser(t *testing.T) {
	playlist := "#EXT-X-INDEPENDENT-SEGMENTS"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, found := p.Find("IndependentSegments")
	assert.True(t, found)
	assert.Equal(t, "", node.HLSElement.Attrs["#EXT-X-INDEPENDENT-SEGMENT"])
}

func TestVariableDefineParser(t *testing.T) {
	// test valid variable define tag with NAME and VALUE
	playlist := "#EXT-X-DEFINE:NAME=\"video_id\",VALUE=\"12345\""
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, found := p.Find("VariableDefine")
	assert.True(t, found)
	assert.Equal(t, "video_id", node.HLSElement.Attrs["NAME"])
	assert.Equal(t, "12345", node.HLSElement.Attrs["VALUE"])

	// test valid variable define tag with QUERYPARAM
	playlist = "#EXT-X-DEFINE:QUERYPARAM=\"video_id\""
	p, err = setupPlaylist(playlist)
	assert.NoError(t, err)

	node, found = p.Find("VariableDefine")
	assert.True(t, found)
	assert.Equal(t, "video_id", node.HLSElement.Attrs["QUERYPARAM"])

	// test invalid variable define tag with NAME but without VALUE
	playlist = "#EXT-X-DEFINE:NAME=\"video_id\""
	_, err = setupPlaylist(playlist)
	assert.Error(t, err)
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

func TestExtKeyParser(t *testing.T) {
	// test valid ext key tag
	playlist := "#EXT-X-KEY:METHOD=SAMPLE-AES,URI=\"drm-uri\",KEYFORMAT=\"com.apple.streamingkeydelivery\",KEYFORMATVERSIONS=\"1\""
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, found := p.Find("ExtKey")
	assert.True(t, found)
	assert.Equal(t, "SAMPLE-AES", node.HLSElement.Attrs["METHOD"])
	assert.Equal(t, "drm-uri", node.HLSElement.Attrs["URI"])
	assert.Equal(t, "com.apple.streamingkeydelivery", node.HLSElement.Attrs["KEYFORMAT"])
	assert.Equal(t, "1", node.HLSElement.Attrs["KEYFORMATVERSIONS"])

	// test invalid ext key without METHOD
	playlist = "#EXT-X-KEY:KEYFORMAT=\"com.apple.streamingkeydelivery\",KEYFORMATVERSIONS=\"1\""
	_, err = setupPlaylist(playlist)
	assert.Error(t, err)

	// test invalid ext key with METHOD not NONE and without URI
	playlist = "#EXT-X-KEY:METHOD=SAMPLE-AES,KEYFORMAT=\"com.apple.streamingkeydelivery\",KEYFORMATVERSIONS=\"1\""
	_, err = setupPlaylist(playlist)
	assert.Error(t, err)

	// test invalid ext key with METHOD AES-128 and without IV
	playlist = "#EXT-X-KEY:METHOD=AES-128,URI=\"https://example.com/keys/key1.bin\""
	_, err = setupPlaylist(playlist)
	assert.Error(t, err)
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

func TestCueOutParserWithoutDuration(t *testing.T) {
	playlist := "#EXT-X-CUE-OUT"
	p, err := setupPlaylist(playlist)
	assert.NoError(t, err)

	node, ok := p.Find("CueOut")
	assert.True(t, ok)
	assert.Equal(t, "0", node.HLSElement.Attrs["#EXT-X-CUE-OUT"])
}

func TestCueInParser(t *testing.T) {
	playlist := strings.Join([]string{
		"#EXT-X-DATERANGE:SCTE35-IN=0xFF0000,ID=\"break1\",START-DATE=\"2025-01-01T00:00:00Z\"",
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
	assert.Equal(t, " no desc", p.CurrentSegment.Title)
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
	assert.Equal(t, " no desc", node.HLSElement.Attrs["Title"])

	playlist = `#EXTINF:3.6
              0.ts`
	p, err = setupPlaylist(playlist)
	assert.NoError(t, err)
	assert.Nil(t, p.CurrentSegment)

	node, found = p.Find("ExtInf")
	assert.True(t, found)
	assert.Equal(t, "0.ts", node.HLSElement.URI)
	assert.Equal(t, "3.6", node.HLSElement.Attrs["Duration"])
	assert.Equal(t, "", node.HLSElement.Attrs["Title"])
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

func TestParseMultivariantPlaylist(t *testing.T) {
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
	assert.Nil(t, p.CurrentSegment)
	assert.Nil(t, p.CurrentStreamInf)
	assert.Equal(t, p.Head.HLSElement.Name, "M3u8Identifier")
	assert.Equal(t, p.Tail.HLSElement.Name, "StreamInf")
	assert.Equal(t, len(p.Variants()), 2)
}

func TestParseMediaPlaylist(t *testing.T) {
	file, _ := os.Open("mocks/media/media.m3u8")
	p, err := m3u8.ParsePlaylist(file)

	assert.NoError(t, err)
	assert.Nil(t, p.CurrentSegment)
	assert.Nil(t, p.CurrentStreamInf)
	assert.Equal(t, p.Head.HLSElement.Name, "M3u8Identifier")
	assert.Equal(t, p.Tail.HLSElement.Name, "ExtInf")
	assert.Equal(t, len(p.Segments()), 27)
	assert.Equal(t, len(p.Breaks()), 1)
	assert.Equal(t, len(p.Variants()), 0)
}

func TestParseMediaPlaylist_WithCompleteAdBreak(t *testing.T) {
	file, _ := os.Open("mocks/media/scte35/withCompleteAdBreak.m3u8")
	p, err := m3u8.ParsePlaylist(file)

	_, foundCueOut := p.Find("CueOut")
	_, foundCueIn := p.Find("CueIn")
	breaks := p.Breaks()

	assert.NoError(t, err)
	assert.Nil(t, p.CurrentSegment)
	assert.Nil(t, p.CurrentStreamInf)
	assert.True(t, foundCueOut)
	assert.True(t, foundCueIn)
	assert.Equal(t, len(breaks), 1)
	assert.Equal(t, breaks[0].HLSElement.Details["StartMediaSequence"], "363969994")
	assert.Equal(t, breaks[0].HLSElement.Details["Status"], tags.BreakStatusComplete)
}

func TestParseMediaPlaylist_WithPartialAdBreak_BeforeDVRLimit(t *testing.T) {
	file, _ := os.Open("mocks/media/scte35/withAdBreakBeforeDVRLimit.m3u8")
	p, err := m3u8.ParsePlaylist(file)

	_, foundCueOut := p.Find("CueOut")
	allBreaks := p.Breaks()
	allPDTs := p.FindAll("ProgramDateTime")

	assert.NoError(t, err)
	assert.True(t, foundCueOut)
	assert.Equal(t, len(allBreaks), 1)

	assert.Equal(t, fmt.Sprintf("%d", p.MediaSequence), "363991004")
	assert.Equal(t, allBreaks[0].HLSElement.Details["StartMediaSequence"], "363991006")
	assert.Equal(t, allBreaks[0].HLSElement.Details["Status"], tags.BreakStatusComplete)

	assert.Equal(t, len(allPDTs), 3)
	assert.NotEqual(t, allPDTs[0].HLSElement.Attrs["#EXT-X-PROGRAM-DATE-TIME"], allBreaks[0].HLSElement.Attrs["START-DATE"])
	assert.Equal(t, allPDTs[1].HLSElement.Attrs["#EXT-X-PROGRAM-DATE-TIME"], allBreaks[0].HLSElement.Attrs["START-DATE"])
}

func TestParseMediaPlaylist_WithPartialAdBreak_OnDVRLimit(t *testing.T) {
	file, _ := os.Open("mocks/media/scte35/withAdBreakOnDVRLimit.m3u8")
	p, err := m3u8.ParsePlaylist(file)

	_, foundCueOut := p.Find("CueOut")
	allBreaks := p.Breaks()
	allPDTs := p.FindAll("ProgramDateTime")

	assert.NoError(t, err)
	assert.True(t, foundCueOut)
	assert.Equal(t, len(allBreaks), 1)

	assert.Equal(t, fmt.Sprintf("%d", p.MediaSequence), "363991006")
	assert.Equal(t, allBreaks[0].HLSElement.Details["StartMediaSequence"], "0")
	assert.Equal(t, allBreaks[0].HLSElement.Details["Status"], tags.BreakStatusLeavingDVR)

	assert.Equal(t, len(allPDTs), 2)
	assert.Equal(t, allPDTs[0].HLSElement.Attrs["#EXT-X-PROGRAM-DATE-TIME"], allBreaks[0].HLSElement.Attrs["START-DATE"])
}

func TestParseMediaPlaylist_WithPartialAdBreak_OutsideDVRLimit(t *testing.T) {
	file, _ := os.Open("mocks/media/scte35/withAdBreakOutsideDVRLimit.m3u8")
	p, err := m3u8.ParsePlaylist(file)

	_, foundCueOut := p.Find("CueOut")
	allBreaks := p.Breaks()
	allPDTs := p.FindAll("ProgramDateTime")

	assert.NoError(t, err)
	assert.False(t, foundCueOut)
	assert.Equal(t, len(allBreaks), 1)

	assert.Equal(t, fmt.Sprintf("%d", p.MediaSequence), "363991008")
	assert.Equal(t, allBreaks[0].HLSElement.Details["StartMediaSequence"], "0")
	assert.Equal(t, allBreaks[0].HLSElement.Details["Status"], tags.BreakStatusLeavingDVR)

	assert.Equal(t, len(allPDTs), 2)
	assert.NotEqual(t, allPDTs[0].HLSElement.Attrs["#EXT-X-PROGRAM-DATE-TIME"], allBreaks[0].HLSElement.Attrs["START-DATE"])
}

func TestParseMediaPlaylist_WithPartialAdBreak_NewNotReady(t *testing.T) {
	file, _ := os.Open("mocks/media/scte35/withAdBreakNewNotReady.m3u8")
	p, err := m3u8.ParsePlaylist(file)

	_, foundCueOut := p.Find("CueOut")
	allBreaks := p.Breaks()
	allPDTs := p.FindAll("ProgramDateTime")

	assert.NoError(t, err)
	assert.False(t, foundCueOut)
	assert.Nil(t, allBreaks[0].Next)
	assert.Equal(t, len(allPDTs), 1)
	assert.Equal(t, allBreaks[0].HLSElement.Details["StartMediaSequence"], "0")
	assert.Equal(t, allBreaks[0].HLSElement.Details["Status"], tags.BreakStatusNotReady)
}

func TestParseMediaPlaylist_WithPartialAdBreak_NewReadyButNoSegment(t *testing.T) {
	file, _ := os.Open("mocks/media/scte35/withAdBreakNewReady.m3u8")
	p, err := m3u8.ParsePlaylist(file)

	_, foundCueOut := p.Find("CueOut")
	allBreaks := p.Breaks()
	allPDTs := p.FindAll("ProgramDateTime")

	assert.NoError(t, err)
	assert.False(t, foundCueOut)
	assert.Nil(t, allBreaks[0].Next)
	assert.Equal(t, len(allPDTs), 1)
	assert.Equal(t, allBreaks[0].HLSElement.Details["StartMediaSequence"], "363969994")
	assert.Equal(t, allBreaks[0].HLSElement.Details["Status"], tags.BreakStatusComplete)
}

func TestParseMediaPlaylist_WithPartialAdBreak_NewReadyButWithSegment(t *testing.T) {
	file, _ := os.Open("mocks/media/scte35/withAdBreakNewReadyWithSegment.m3u8")
	p, err := m3u8.ParsePlaylist(file)

	_, foundCueOut := p.Find("CueOut")
	allBreaks := p.Breaks()
	allPDTs := p.FindAll("ProgramDateTime")
	newestSegment := p.Tail

	assert.NoError(t, err)
	assert.True(t, foundCueOut)
	assert.NotNil(t, allBreaks[0].Next)
	assert.Equal(t, len(allPDTs), 2)
	assert.Equal(t, allBreaks[0].HLSElement.Details["StartMediaSequence"], newestSegment.HLSElement.Details["MediaSequence"])
	assert.Equal(t, allBreaks[0].HLSElement.Details["Status"], tags.BreakStatusComplete)
}

func TestParseMediaPlaylist_WithCompleteAdBreak_BreakStartTimePrecision(t *testing.T) {
	file, _ := os.Open("mocks/media/scte35/withBreakStartTimePrecisionEx1.m3u8")
	p, err := m3u8.ParsePlaylist(file)

	_, foundCueOut := p.Find("CueOut")
	allBreaks := p.Breaks()
	allPDTs := p.FindAll("ProgramDateTime")

	assert.NoError(t, err)
	assert.True(t, foundCueOut)
	assert.Equal(t, len(allBreaks), 1)
	assert.Equal(t, len(allPDTs), 3)
	assert.Equal(t, allBreaks[0].HLSElement.Details["StartMediaSequence"], "547307194")
	assert.Equal(t, allBreaks[0].HLSElement.Details["Status"], tags.BreakStatusComplete)

	file, _ = os.Open("mocks/media/scte35/withBreakStartTimePrecisionEx2.m3u8")
	p, err = m3u8.ParsePlaylist(file)

	_, foundCueOut = p.Find("CueOut")
	allBreaks = p.Breaks()
	allPDTs = p.FindAll("ProgramDateTime")

	assert.NoError(t, err)
	assert.True(t, foundCueOut)
	assert.Equal(t, len(allBreaks), 1)
	assert.Equal(t, len(allPDTs), 3)
	assert.Equal(t, allBreaks[0].HLSElement.Details["StartMediaSequence"], "548062663")
	assert.Equal(t, allBreaks[0].HLSElement.Details["Status"], tags.BreakStatusComplete)
}

func TestParseMediaPlaylistWithDiscontinuity(t *testing.T) {
	file, _ := os.Open("mocks/media/withDiscontinuity.m3u8")
	p, err := m3u8.ParsePlaylist(file)

	assert.NoError(t, err)
	assert.Nil(t, p.CurrentSegment)
	assert.Nil(t, p.CurrentStreamInf)
	assert.Equal(t, p.Head.HLSElement.Name, "M3u8Identifier")
	assert.Equal(t, p.Tail.HLSElement.Name, "ExtInf")

	node, found := p.Find("DiscontinuitySequence")

	assert.True(t, found)
	assert.Equal(t, "87498", node.HLSElement.Attrs["#EXT-X-DISCONTINUITY-SEQUENCE"])
	assert.Equal(t, p.DiscontinuitySequence, 87498)

	nodes := p.FindAll("Discontinuity")

	assert.Len(t, nodes, 2)
	assert.Equal(t, "", nodes[0].HLSElement.Attrs["#EXT-X-DISCONTINUITY"])
}

func TestParseMediaPlaylistWithEncryption_AES128(t *testing.T) {
	file, _ := os.Open("mocks/media/encryption/withAES128.m3u8")
	p, err := m3u8.ParsePlaylist(file)

	assert.NoError(t, err)
	assert.Nil(t, p.CurrentSegment)
	assert.Nil(t, p.CurrentStreamInf)
	assert.Equal(t, p.Head.HLSElement.Name, "M3u8Identifier")
	assert.Equal(t, p.Tail.HLSElement.Name, "ExtInf")

	node, found := p.Find("ExtKey")

	assert.True(t, found)
	assert.Equal(t, "AES-128", node.HLSElement.Attrs["METHOD"])
	assert.Equal(t, "https://example.com/keys/key1.bin", node.HLSElement.Attrs["URI"])
	assert.Equal(t, "0x0123456789abcdef0123456789abcdef", node.HLSElement.Attrs["IV"])
}

func TestParseMediaPlaylistWithEncryption_SampleAES(t *testing.T) {
	file, _ := os.Open("mocks/media/encryption/withSampleAES.m3u8")
	p, err := m3u8.ParsePlaylist(file)

	assert.NoError(t, err)
	assert.Nil(t, p.CurrentSegment)
	assert.Nil(t, p.CurrentStreamInf)
	assert.Equal(t, p.Head.HLSElement.Name, "M3u8Identifier")
	assert.Equal(t, p.Tail.HLSElement.Name, "ExtInf")

	node, found := p.Find("ExtKey")

	assert.True(t, found)
	assert.Equal(t, "SAMPLE-AES", node.HLSElement.Attrs["METHOD"])
	assert.Equal(t, "sample-aes-uri", node.HLSElement.Attrs["URI"])
	assert.Equal(t, "com.apple.streamingkeydelivery", node.HLSElement.Attrs["KEYFORMAT"])
	assert.Equal(t, "1", node.HLSElement.Attrs["KEYFORMATVERSIONS"])
}

func TestParseMediaPlaylistWithEncryptionAndCompleteAdBreak(t *testing.T) {
	file, _ := os.Open("mocks/media/withEncryptionAndSCTE35.m3u8")
	p, err := m3u8.ParsePlaylist(file)

	assert.NoError(t, err)
	assert.Nil(t, p.CurrentSegment)
	assert.Nil(t, p.CurrentStreamInf)
	assert.Equal(t, p.Head.HLSElement.Name, "M3u8Identifier")
	assert.Equal(t, p.Tail.HLSElement.Name, "ExtInf")

	extKeyNodes := p.FindAll("ExtKey")
	assert.Len(t, extKeyNodes, 3)

	_, found1 := p.FindNodeInsideAdBreak(extKeyNodes[0])
	assert.False(t, found1)

	_, found2 := p.FindNodeInsideAdBreak(extKeyNodes[1])
	assert.True(t, found2)

	_, found3 := p.FindNodeInsideAdBreak(extKeyNodes[2])
	assert.False(t, found3)

	dateRangeNodes := p.Breaks()

	assert.Len(t, dateRangeNodes, 1)
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
			path:  "./mocks/multivariant/missingVersion.m3u8",
			error: true,
		},
		{
			name:           "Parse media playlist",
			kind:           "media",
			path:           "./mocks/media/media.m3u8",
			pdt:            time.Date(2025, 05, 16, 13, 33, 27, 966666000, time.UTC),
			dvr:            129.5999,
			segmentCounter: 27,
		},
		{
			name:           "Parse media playlist with EXT-X-DISCONTINUITY tag",
			kind:           "media",
			path:           "./mocks/media/withDiscontinuity.m3u8",
			pdt:            time.Date(2025, time.July, 1, 19, 1, 40, 466666000, time.UTC),
			dvr:            51.1998,
			segmentCounter: 16,
		},
		{
			name: "Parse multivariant playlist",
			kind: "multivariant",
			path: "./mocks/multivariant/multivariant.m3u8",
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
