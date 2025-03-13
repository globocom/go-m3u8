package go_m3u8_test

import (
	"os"
	"testing"

	m3u8 "github.com/globocom/go-m3u8"
	"github.com/stretchr/testify/assert"
)

type fakeSource struct{}

func (f fakeSource) Read(p []byte) (n int, err error) {
	return 0, m3u8.ErrOpenPlaylist
}

func (f fakeSource) Close() error {
	return nil
}

func TestParsePlaylist(t *testing.T) {
	type testCaseParams struct {
		name     string
		path     string
		source   m3u8.Source
		elements []m3u8.Node
		error    bool
	}
	testCases := []testCaseParams{
		{
			name:     "Error parsing playlist",
			path:     "not exist",
			elements: nil,
			source:   fakeSource{},
			error:    true,
		},
		{
			name:     "Missing version in master playlist",
			path:     "./testdata/default/masterMissingVersion.m3u8",
			elements: nil,
			error:    true,
		},
		{
			name: "Parse media playlist",
			path: "./testdata/default/media.m3u8",
			elements: []m3u8.Node{
				{Tags: "3"},
				{Tags: "360948012"},
				{Tags: "7"},
				{Tags: &m3u8.ProgramDateTime{DateTime: "2024-11-25T16:00:53.200000Z"}},
				{Tags: &m3u8.Segment{Duration: "4.8", URI: "channel-audio_1=96000-video=2262976-360948204.ts"}},
				{Tags: &m3u8.Segment{Duration: "3.3333", URI: "channel-audio_1=96000-video=2262976-360948205.ts"}},
				{Tags: &m3u8.DateRange{
					Id:              "1-1732551382",
					StartDate:       "2024-11-25T16:16:22.933333Z",
					PlannedDuration: "60.1",
					Scte35Mark:      map[string]string{"OUT": "0xFC3025000000000BB802FFF01405000000017FEFFF1A7B3B607E005288E8000100000000F799F45B"},
				}},
				{Tags: &m3u8.ProgramDateTime{DateTime: "2024-11-25T16:16:22.933333Z"}},
				{Tags: &m3u8.Segment{Duration: "6.2666", URI: "channel-audio_1=96000-video=2262976-360948206.ts"}},
				{Tags: &m3u8.Segment{Duration: "4.8", URI: "channel-audio_1=96000-video=2262976-360948207.ts"}},
				{Tags: &m3u8.Segment{Duration: "4.8", URI: "channel-audio_1=96000-video=2262976-360948208.ts"}},
				{Tags: &m3u8.Segment{Duration: "4.8", URI: "channel-audio_1=96000-video=2262976-360948209.ts"}},
				{Tags: &m3u8.Segment{Duration: "4.8", URI: "channel-audio_1=96000-video=2262976-360948210.ts"}},
				{Tags: &m3u8.Segment{Duration: "4.8", URI: "channel-audio_1=96000-video=2262976-360948211.ts"}},
				{Tags: &m3u8.Segment{Duration: "4.8", URI: "channel-audio_1=96000-video=2262976-360948212.ts"}},
				{Tags: &m3u8.Segment{Duration: "4.8", URI: "channel-audio_1=96000-video=2262976-360948213.ts"}},
				{Tags: &m3u8.Segment{Duration: "4.8", URI: "channel-audio_1=96000-video=2262976-360948214.ts"}},
				{Tags: &m3u8.Segment{Duration: "4.8", URI: "channel-audio_1=96000-video=2262976-360948215.ts"}},
				{Tags: &m3u8.Segment{Duration: "4.8", URI: "channel-audio_1=96000-video=2262976-360948216.ts"}},
				{Tags: &m3u8.Segment{Duration: "5.8333", URI: "channel-audio_1=96000-video=2262976-360948217.ts"}},
				{Tags: &m3u8.DateRange{
					Id:         "1-1732551382",
					StartDate:  "2024-11-25T16:16:22.933333Z",
					EndDate:    "2024-11-25T16:17:23.033333Z",
					Duration:   "60.1",
					Scte35Mark: map[string]string{"IN": "0xFC3025000000000BB802FFF01405000000017F6FFF1ACDC4487E00000000000100000000E9B751B0"},
				}},
				{Tags: &m3u8.ProgramDateTime{DateTime: "2024-11-25T16:17:23.033333Z"}},
				{Tags: &m3u8.Segment{Duration: "3.7666", URI: "channel-audio_1=96000-video=2262976-360948218.ts"}},
				{Tags: &m3u8.Segment{Duration: "4.8", URI: "channel-audio_1=96000-video=2262976-360948219.ts"}},
			},
		},
		{
			name: "Parse master playlist with variants",
			path: "./testdata/default/master.m3u8",
			elements: []m3u8.Node{
				{Tags: "3"},
				{Tags: &m3u8.StreamInf{
					Codecs:           []string{"mp4a.40.2", "avc1.64001F"},
					Bandwidth:        "206000",
					AverageBandwidth: "187000",
					Resolution:       "256x144",
					FrameRate:        "30",
					URI:              "channel-audio_1=96000-video=80000.m3u8",
				}},
				{Tags: &m3u8.StreamInf{
					Codecs:           []string{"mp4a.40.2", "avc1.64001F"},
					Bandwidth:        "299000",
					AverageBandwidth: "272000",
					Resolution:       "384x216",
					FrameRate:        "30",
					URI:              "channel-audio_1=96000-video=160000.m3u8",
				}},
				{Tags: &m3u8.StreamInf{
					Codecs:           []string{"mp4a.40.2", "avc1.64001F"},
					Bandwidth:        "479000",
					AverageBandwidth: "435000",
					Resolution:       "512x288",
					FrameRate:        "30",
					URI:              "channel-audio_1=96000-video=313984.m3u8",
				}},
				{Tags: &m3u8.StreamInf{
					Codecs:           []string{"mp4a.40.2", "avc1.64001F"},
					Bandwidth:        "764000",
					AverageBandwidth: "695000",
					Resolution:       "640x360",
					FrameRate:        "30",
					URI:              "channel-audio_1=96000-video=558976.m3u8",
				}},
				{Tags: &m3u8.StreamInf{
					Codecs:           []string{"mp4a.40.2", "avc1.64001F"},
					Bandwidth:        "1224000",
					AverageBandwidth: "1112000",
					Resolution:       "768x432",
					FrameRate:        "30",
					URI:              "channel-audio_1=96000-video=952960.m3u8",
				}},
				{Tags: &m3u8.StreamInf{
					Codecs:           []string{"mp4a.40.2", "avc1.64001F"},
					Bandwidth:        "1835000",
					AverageBandwidth: "1668000",
					Resolution:       "1280x720",
					FrameRate:        "30",
					URI:              "channel-audio_1=96000-video=1476992.m3u8",
				}},
				{Tags: &m3u8.StreamInf{
					Codecs:           []string{"mp4a.40.2", "avc1.64001F"},
					Bandwidth:        "2751000",
					AverageBandwidth: "2501000",
					Resolution:       "1280x720",
					FrameRate:        "30",
					URI:              "channel-audio_1=96000-video=2262976.m3u8",
				}},
				{Tags: &m3u8.StreamInf{
					Codecs:           []string{"mp4a.40.2", "avc1.640029"},
					Bandwidth:        "4127000",
					AverageBandwidth: "3752000",
					Resolution:       "1920x1080",
					FrameRate:        "30",
					URI:              "channel-audio_1=96000-video=3442944.m3u8",
				}},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			file, err := os.Open(tc.path)
			playlist, err := m3u8.ParsePlaylist(file)

			if tc.error {
				if tc.name == "Missing version in master playlist" {
					file, err = os.Open(tc.path)
					playlist, err = m3u8.ParsePlaylist(file)
					assert.ErrorContains(t, err, m3u8.ErrParseLine.Error())
					assert.Nil(t, playlist)
					return
				} else if tc.name == "Error parsing playlist" {
					playlist, err = m3u8.ParsePlaylist(tc.source)
					assert.Equal(t, err, m3u8.ErrOpenPlaylist)
					assert.Nil(t, playlist)
					return
				}
			}

			assert.NoError(t, err)
			assert.NotNil(t, playlist)
			assert.Len(t, playlist.Elements, len(tc.elements))
			assert.NotNil(t, playlist.Head)
			assert.NotNil(t, playlist.Tail)
			for i := range tc.elements {
				assert.Equal(t, tc.elements[i].Tags, playlist.Elements[i].Tags)
			}
		})
	}
}

func TestMediaPlaylistInsert(t *testing.T) {
	type testCaseParams struct {
		name             string
		initialState     *m3u8.Playlist
		nodeToInsert     *m3u8.Node
		expectedHead     *m3u8.Node
		expectedTail     *m3u8.Node
		expectedElements []m3u8.Node
		expectedPrev     *m3u8.Node
		expectedNext     *m3u8.Node
	}

	testCases := []testCaseParams{
		{
			name:             "Empty list, first node insertion",
			initialState:     &m3u8.Playlist{},
			nodeToInsert:     &m3u8.Node{Tags: "3"},
			expectedHead:     &m3u8.Node{Tags: "3"},
			expectedTail:     &m3u8.Node{Tags: "3"},
			expectedElements: []m3u8.Node{{Tags: "3"}},
			expectedPrev:     nil,
			expectedNext:     nil,
		},
		{
			name: "Existing list, second node insertion",
			initialState: &m3u8.Playlist{
				Head:     &m3u8.Node{Tags: "3"},
				Tail:     &m3u8.Node{Tags: "3"},
				Elements: []m3u8.Node{{Tags: "3"}},
			},
			nodeToInsert: &m3u8.Node{
				Tags: "360948012",
			},
			expectedHead: &m3u8.Node{Tags: "3"},
			expectedTail: &m3u8.Node{Tags: "360948012"},
			expectedElements: []m3u8.Node{
				{Tags: "3"},
				{Tags: "360948012"},
			},
			expectedPrev: &m3u8.Node{Tags: "3"},
			expectedNext: &m3u8.Node{Tags: "360948012"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			playlist := tc.initialState
			nodeToInsert := tc.nodeToInsert
			playlist.Insert(nodeToInsert)

			assert.Equal(t, tc.expectedHead.Tags, playlist.Head.Tags)
			assert.Equal(t, tc.expectedTail.Tags, playlist.Tail.Tags)

			for i := range tc.expectedElements {
				assert.Equal(t, tc.expectedElements[i].Tags, playlist.Elements[i].Tags)
			}

			if tc.name == "Existing list, second node insertion" {
				assert.Equal(t, tc.expectedPrev.Tags, nodeToInsert.Prev.Tags)
				assert.Nil(t, nodeToInsert.Next)
			}
		})
	}
}

func TestMediaPlaylistFind(t *testing.T) {
	type testCaseParams struct {
		name        string
		playlist    *m3u8.Playlist
		expectedTag any
	}
	testCases := []testCaseParams{
		{
			name:        "Empty playlist",
			playlist:    &m3u8.Playlist{},
			expectedTag: nil,
		},
		{
			name: "Existing playlist",
			playlist: &m3u8.Playlist{
				Head: &m3u8.Node{
					Tags: "3",
					Next: &m3u8.Node{
						Tags: "360948012",
					},
				},
				Tail: &m3u8.Node{
					Tags: "360948012",
				},
				Elements: []m3u8.Node{
					{Tags: "3"},
					{Tags: "360948012"},
				},
			},
			expectedTag: "360948012",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			node, found := tc.playlist.Find(tc.expectedTag)
			if tc.expectedTag != nil {
				assert.Equal(t, tc.expectedTag, node.Tags)
				assert.True(t, found)
			} else {
				assert.Nil(t, node)
				assert.False(t, found)
			}
		})
	}

}
