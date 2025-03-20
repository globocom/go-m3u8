package go_m3u8_test

import (
	"fmt"
	"os"
	"testing"

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
		name   string
		path   string
		source m3u8.Source
		error  bool
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
			path:  "./testdata/default/masterMissingVersion.m3u8",
			error: true,
		},
		{
			name: "Parse media playlist",
			path: "./testdata/default/media.m3u8",
		},
		{
			name: "Parse media with discontinuity playlist",
			path: "./testdata/default/mediaWithDiscontinuity.m3u8",
		},
		{
			name: "Parse master playlist with variants",
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
		})
	}
}
