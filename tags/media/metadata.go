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
	"strings"
	"time"

	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
	"github.com/rs/zerolog/log"
)

var (
	DateRangeTag   = "#EXT-X-DATERANGE"    //todo: has attributes
	SkipTag        = "#EXT-X-SKIP"         //todo: has attributes
	PreLoadHintTag = "#EXT-X-PRELOAD-HINT" //todo: has attributes
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
			Details: map[string]string{
				"MediaSequence": fmt.Sprintf("%d", playlist.MediaSequence+playlist.SegmentsCounter),
			},
		},
	}

	// date range tag with scte-35 attribute (i.e. indicates ad break start)
	if dateRangeNode.HLSElement.Attrs["SCTE35-OUT"] != "" {
		startDate, err := time.Parse(time.RFC3339Nano, dateRangeNode.HLSElement.Attrs["START-DATE"])
		if err != nil {
			return fmt.Errorf("failed to parse start time of date range tag: %s", tag)
		}

		// date range's media sequence should equal the ad break's first segment's media sequence
		// however, when the ad break has begun and segments are outside DVR limit, we have lost the first segment's information
		// as a solution, we will set the date range's media sequence to zero

		// in this case, the program date time tag could come AFTER the date range tag in the playlist, being parsed as a zero value
		// or it comes BEFORE the date range tag, and its timestamp is later than date range's start date
		if playlist.ProgramDateTime.IsZero() || playlist.ProgramDateTime.After(startDate) {
			log.Info().Msg("break about to leave dvr limit, media sequence will be zero")
			dateRangeNode.HLSElement.Details["MediaSequence"] = "0"
		}
	}

	playlist.Insert(dateRangeNode)
	return nil
}

func (e DateRangeEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	order := []string{"ID", "START-DATE", "PLANNED-DURATION", "END-DATE", "DURATION", "SCTE35-OUT", "SCTE35-IN"}
	return pl.EncodeTagWithAttributes(builder, DateRangeTag, node.HLSElement.Attrs, order)
}
