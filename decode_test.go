package go_m3u8_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	m3u8 "github.com/globocom/go-m3u8"
	"github.com/stretchr/testify/assert"
)

type fakeSource struct{}

func (f fakeSource) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("fake error")
}

func (f fakeSource) Close() error {
	return nil
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
			name:  "Missing version in master playlist",
			kind:  "master",
			path:  "./testdata/default/masterMissingVersion.m3u8",
			error: true,
		},
		{
			name:           "Parse media playlist",
			kind:           "media",
			path:           "./testdata/default/media.m3u8",
			pdt:            time.Date(2024, 11, 25, 16, 0, 53, 200000000, time.UTC),
			dvr:            76.7998,
			segmentCounter: 16,
		},
		{
			name:           "Parse media with discontinuity playlist",
			kind:           "media",
			path:           "./testdata/default/mediaWithDiscontinuity.m3u8",
			pdt:            time.Date(2024, 11, 25, 16, 0, 53, 200000000, time.UTC),
			dvr:            76.7998,
			segmentCounter: 16,
		},
		{
			name: "Parse master playlist with variants",
			kind: "master",
			path: "./testdata/default/master.m3u8",
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
