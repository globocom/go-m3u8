package go_m3u8_test

import (
	"github.com/stretchr/testify/assert"
	m3u8 "gitlab.globoi.com/webmedia/media-delivery-advertising/go-m3u8"
	"testing"
)

func TestParseMasterPlaylist(t *testing.T) {
	t.Parallel()
	type testCaseParams struct {
		name           string
		path           string
		masterPlaylist *m3u8.MasterPlaylist
	}

	testCases := []testCaseParams{
		{
			name: "Parse MasterPlaylist playlist with variants",
			path: "./testdata/default/master.m3u8",
			masterPlaylist: &m3u8.MasterPlaylist{
				Version: "3",
				Variants: []m3u8.StreamInf{
					{
						Codecs:           []string{"mp4a.40.2", "avc1.64001F"},
						Bandwidth:        "206000",
						AverageBandwidth: "187000",
						Resolution:       "256x144",
						FrameRate:        "30",
						URI:              "channel-audio_1=96000-video=80000.m3u8",
					},
					{
						Codecs:           []string{"mp4a.40.2", "avc1.64001F"},
						Bandwidth:        "299000",
						AverageBandwidth: "272000",
						Resolution:       "384x216",
						FrameRate:        "30",
						URI:              "channel-audio_1=96000-video=160000.m3u8",
					}, {
						Codecs:           []string{"mp4a.40.2", "avc1.64001F"},
						Bandwidth:        "479000",
						AverageBandwidth: "435000",
						Resolution:       "512x288",
						FrameRate:        "30",
						URI:              "channel-audio_1=96000-video=313984.m3u8",
					}, {
						Codecs:           []string{"mp4a.40.2", "avc1.64001F"},
						Bandwidth:        "764000",
						AverageBandwidth: "695000",
						Resolution:       "640x360",
						FrameRate:        "30",
						URI:              "channel-audio_1=96000-video=558976.m3u8",
					}, {
						Codecs:           []string{"mp4a.40.2", "avc1.64001F"},
						Bandwidth:        "1224000",
						AverageBandwidth: "1112000",
						Resolution:       "768x432",
						FrameRate:        "30",
						URI:              "channel-audio_1=96000-video=952960.m3u8",
					}, {
						Codecs:           []string{"mp4a.40.2", "avc1.64001F"},
						Bandwidth:        "1835000",
						AverageBandwidth: "1668000",
						Resolution:       "1280x720",
						FrameRate:        "30",
						URI:              "channel-audio_1=96000-video=1476992.m3u8",
					}, {
						Codecs:           []string{"mp4a.40.2", "avc1.64001F"},
						Bandwidth:        "2751000",
						AverageBandwidth: "2501000",
						Resolution:       "1280x720",
						FrameRate:        "30",
						URI:              "channel-audio_1=96000-video=2262976.m3u8",
					}, {
						Codecs:           []string{"mp4a.40.2", "avc1.640029"},
						Bandwidth:        "4127000",
						AverageBandwidth: "3752000",
						Resolution:       "1920x1080",
						FrameRate:        "30",
						URI:              "channel-audio_1=96000-video=3442944.m3u8",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			playlist, err := m3u8.ParseMasterPlaylist(tc.path)
			assert.NoError(t, err)
			assert.NotNil(t, playlist)
			assert.Equal(t, tc.masterPlaylist, playlist)
		})
	}
}

func TestParseMediaPlaylist(t *testing.T) {
	t.Parallel()
	type testCaseParams struct {
		name          string
		path          string
		mediaPlaylist *m3u8.MediaPlaylist
	}

	testCases := []testCaseParams{
		{
			name: "Parse media playlist",
			path: "./testdata/default/media.m3u8",
			mediaPlaylist: &m3u8.MediaPlaylist{
				Version:        "3",
				TargetDuration: "7",
				MediaSequence:  "360948012",
				DateRanges: []m3u8.DateRange{
					{
						Scte35Mark: map[string]string{
							"OUT": "0xFC3025000000000BB802FFF01405000000017FEFFF1A7B3B607E005288E8000100000000F799F45B",
						},
						Id:              "1-1732551382",
						StartDate:       "2024-11-25T16:16:22.933333Z",
						EndDate:         "",
						PlannedDuration: "60.1",
						Duration:        "",
					}, {
						Scte35Mark: map[string]string{
							"IN": "0xFC3025000000000BB802FFF01405000000017F6FFF1ACDC4487E00000000000100000000E9B751B0",
						},
						Id:              "1-1732551382",
						StartDate:       "2024-11-25T16:16:22.933333Z",
						EndDate:         "2024-11-25T16:17:23.033333Z",
						PlannedDuration: "",
						Duration:        "60.1",
					},
				},
				ProgramsDateTime: []m3u8.ProgramDateTime{
					{
						DateTime: "2024-11-25T16:00:53.200000Z",
					}, {
						DateTime: "2024-11-25T16:16:22.933333Z",
					}, {
						DateTime: "2024-11-25T16:17:23.033333Z",
					},
				},
				Segments: []m3u8.Segment{
					{
						Duration: "4.8",
						URI:      "channel-audio_1=96000-video=2262976-360948204.ts",
					}, {
						Duration: "3.3333",
						URI:      "channel-audio_1=96000-video=2262976-360948205.ts",
					}, {
						Duration: "6.2666",
						URI:      "channel-audio_1=96000-video=2262976-360948206.ts",
					}, {
						Duration: "4.8",
						URI:      "channel-audio_1=96000-video=2262976-360948207.ts",
					}, {
						Duration: "4.8",
						URI:      "channel-audio_1=96000-video=2262976-360948208.ts",
					}, {
						Duration: "4.8",
						URI:      "channel-audio_1=96000-video=2262976-360948209.ts",
					}, {
						Duration: "4.8",
						URI:      "channel-audio_1=96000-video=2262976-360948210.ts",
					}, {
						Duration: "4.8",
						URI:      "channel-audio_1=96000-video=2262976-360948211.ts",
					}, {
						Duration: "4.8",
						URI:      "channel-audio_1=96000-video=2262976-360948212.ts",
					}, {
						Duration: "4.8",
						URI:      "channel-audio_1=96000-video=2262976-360948213.ts",
					}, {
						Duration: "4.8",
						URI:      "channel-audio_1=96000-video=2262976-360948214.ts",
					}, {
						Duration: "4.8",
						URI:      "channel-audio_1=96000-video=2262976-360948215.ts",
					}, {
						Duration: "4.8",
						URI:      "channel-audio_1=96000-video=2262976-360948216.ts",
					}, {
						Duration: "5.8333",
						URI:      "channel-audio_1=96000-video=2262976-360948217.ts",
					}, {
						Duration: "3.7666",
						URI:      "channel-audio_1=96000-video=2262976-360948218.ts",
					}, {
						Duration: "4.8",
						URI:      "channel-audio_1=96000-video=2262976-360948219.ts",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			playlist, err := m3u8.ParseMediaPlaylist(tc.path)
			assert.NoError(t, err)
			assert.NotNil(t, playlist)
			assert.Equal(t, tc.mediaPlaylist, playlist)
		})
	}
}
