package playlist_test

import (
	"os"
	"testing"

	m3u8 "github.com/globocom/go-m3u8"
	"github.com/globocom/go-m3u8/internal"
	"github.com/globocom/go-m3u8/playlist"
	"github.com/stretchr/testify/assert"
)

func TestValidateSuccessStartDate(t *testing.T) {
	adBreak := &internal.Node{
		HLSElement: &internal.HLSElement{
			Attrs: map[string]string{
				"START-DATE": "2025-05-20T20:24:52.699999Z",
			},
		},
	}

	startDate, err := playlist.ValidateStartDate(adBreak)
	assert.NoError(t, err)
	expectedDate := "2025-05-20 20:24:52.699999 +0000 UTC"
	assert.Equal(t, expectedDate, startDate.String())
}

func TestValidateFailStartDate(t *testing.T) {
	adBreak := &internal.Node{
		HLSElement: &internal.HLSElement{
			Attrs: map[string]string{
				"START-DATE": "invalid-date-format",
			},
		},
	}

	_, err := playlist.ValidateStartDate(adBreak)
	assert.Error(t, err)
}

func TestValidateSuccessMediaSequence(t *testing.T) {
	adBreak := &internal.Node{
		HLSElement: &internal.HLSElement{
			Details: map[string]string{
				"StartMediaSequence": "42",
			},
		},
	}

	mediaSequence, err := playlist.ValidateMediaSequence(adBreak)
	assert.NoError(t, err)
	assert.Equal(t, 42, mediaSequence)
}

func TestValidateFailMediaSequence(t *testing.T) {
	adBreak := &internal.Node{
		HLSElement: &internal.HLSElement{
			Details: map[string]string{
				"StartMediaSequence": "not-an-integer",
			},
		},
	}

	_, err := playlist.ValidateMediaSequence(adBreak)
	assert.Error(t, err)
}

func TestValidatePlannedDuration(t *testing.T) {
	adBreak := &internal.Node{
		HLSElement: &internal.HLSElement{
			Attrs: map[string]string{
				"PLANNED-DURATION": "60",
			},
		},
	}

	err := playlist.ValidatePlannedDuration(adBreak)
	assert.NoError(t, err)
}

func TestValidatePlannedDurationExceedsLimit(t *testing.T) {
	adBreak := &internal.Node{
		HLSElement: &internal.HLSElement{
			Attrs: map[string]string{
				"PLANNED-DURATION": "700",
			},
		},
	}

	err := playlist.ValidatePlannedDuration(adBreak)
	assert.Equal(t, playlist.ErrPlannedDurationExceedsLimit, err)
}

func TestValidatePlannedDurationMissingPlannedDuration(t *testing.T) {
	adBreak := &internal.Node{
		HLSElement: &internal.HLSElement{
			Attrs: map[string]string{
				"PLANNED-DURATION": "0",
			},
		},
	}

	err := playlist.ValidatePlannedDuration(adBreak)
	assert.Equal(t, playlist.ErrMissingPlannedDuration, err)
}

func TestIsDuplicatedBreak(t *testing.T) {
	file, _ := os.Open("./../mocks/media/withDuplicatedBreaks.m3u8")
	pl, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	adBreaks := pl.Breaks()
	assert.GreaterOrEqual(t, len(adBreaks), 2)

	lastBreak := adBreaks[1]
	previousBreak := adBreaks[0]

	isDuplicate := playlist.IsDuplicatedBreak(lastBreak, previousBreak)
	assert.True(t, isDuplicate)
}

func TestIsCueOutFromBreak(t *testing.T) {
	file, _ := os.Open("./../mocks/media/media.m3u8")
	pl, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	cueOutNode, _ := pl.Find("CueOut")

	breakMediaSequence := 364042175
	result := pl.IsCueOutFromBreak(breakMediaSequence, cueOutNode)
	assert.True(t, result)
}

func TestIsCueInFromBreak(t *testing.T) {
	file, _ := os.Open("./../mocks/media/media.m3u8")
	pl, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)
	cueInNode, _ := pl.Find("CueIn")
	breakMediaSequence := 364042175

	result := pl.IsCueInFromBreak(breakMediaSequence, cueInNode)
	assert.True(t, result)

}
