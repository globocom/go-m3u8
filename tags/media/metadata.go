//	Media Metadata Tags (Section 4.4.5 on RFC)
//
// Media Metadata tags provide information about the playlist that is
// not associated with specific Media Segments.  There MAY be more than
// one Media Metadata tag of each type in any Media Playlist.  The only
// exception to this rule is EXT-X-SKIP, which MUST NOT appear more than
// once.
package media

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
	"github.com/rs/zerolog/log"
)

var (
	DateRangeTag   = "#EXT-X-DATERANGE"
	SkipTag        = "#EXT-X-SKIP"         //todo: has attributes
	PreLoadHintTag = "#EXT-X-PRELOAD-HINT" //todo: has attributes
)

const (
	BreakStatusLeavingDVR = "leavingDVRLimit"
	BreakStatusNotReady   = "segmentsNotReady"
	BreakStatusComplete   = "complete"
)

type DateRangeParser struct{}

type DateRangeEncoder struct{}

func (p DateRangeParser) Parse(tag string, playlist *pl.Playlist) error {
	params := pl.TagsToMap(tag)
	if len(params) < 1 {
		return fmt.Errorf("invalid date range tag: %s", tag)
	}

	dateRangeNode := &internal.Node{
		HLSElement: &internal.HLSElement{
			Name:  "DateRange",
			Attrs: params,
		},
	}

	// An EXT-X-DATERANGE SCTE35-OUT tag signals the start of an Ad Break
	if dateRangeNode.HLSElement.Attrs["SCTE35-OUT"] != "" {
		mediaSequence, status := getAdBreakDetails(playlist, dateRangeNode)
		dateRangeNode.HLSElement.Details = map[string]string{
			"StartMediaSequence": mediaSequence,
			"Status":             status,
		}
	}

	playlist.Insert(dateRangeNode)
	return nil
}

func (e DateRangeEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	order := []string{"ID", "START-DATE", "PLANNED-DURATION", "END-DATE", "DURATION", "SCTE35-OUT", "SCTE35-IN"}
	return pl.EncodeTagWithAttributes(builder, DateRangeTag, node.HLSElement.Attrs, order)
}

// Returns the Ad Break's media sequence (string) and status (string).
//   - The Break's media sequence will be the media sequence of the first segment inside the break (or zero if Break is incomplete).
//   - The Break's status will be: "complete" or incomplete ("leavingDVRLimit" or "segmentsNotReady").
func getAdBreakDetails(playlist *pl.Playlist, dateRangeNode *internal.Node) (string, string) {
	currentMediaSequence := fmt.Sprintf("%d", playlist.MediaSequence+playlist.SegmentsCounter)
	breakStartDate, _ := time.Parse(time.RFC3339Nano, dateRangeNode.HLSElement.Attrs["START-DATE"])

	// when ad break segments are leaving DVR, we lose the break's first segment's media sequence
	if playlist.ProgramDateTime.IsZero() {
		// if the playlist's PDT tag was not parsed yet, we check if there are any media segments before the date range tag
		if len(playlist.Segments()) == 0 {
			log.Debug().Msg("break is leaving dvr limit, media sequence will be zero")
			return "0", BreakStatusLeavingDVR
		}
	} else {
		// if the playlist's PDT tag was already parsed, we check if the playlist PDT is equal or higher than the break's start date
		if playlist.ProgramDateTime.Equal(breakStartDate) || playlist.ProgramDateTime.After(breakStartDate) {
			log.Debug().Msg("break is leaving dvr limit, media sequence will be zero")
			return "0", BreakStatusLeavingDVR
		}
	}

	// when date range tag exists, but we don't know if we have the break's first media segment yet
	// we check if the break's start date comes later than the estimated next segment's PDT
	nextSegmentEstimatedPDT := playlist.ProgramDateTime.Add(time.Duration(playlist.DVR * float64(time.Second)))
	if (roundUpToSecond(breakStartDate)).After(roundUpToSecond(nextSegmentEstimatedPDT)) {
		log.Debug().Msg("segments for ad break are not ready yet, media sequence will be zero")
		return "0", BreakStatusNotReady
	}

	return currentMediaSequence, BreakStatusComplete
}

func roundUpToSecond(t time.Time) time.Time {
	seconds := float64(t.UnixNano()) / float64(time.Second)
	return time.Unix(int64(math.Ceil(seconds)), 0).UTC()
}
