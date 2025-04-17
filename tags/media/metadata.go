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

	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
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

	playlist.Insert(dateRangeNode)
	return nil
}

func (e DateRangeEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	order := []string{"ID", "START-DATE", "PLANNED-DURATION", "END-DATE", "DURATION", "SCTE35-OUT", "SCTE35-IN"}
	return pl.EncodeTagWithAttributes(builder, DateRangeTag, node.HLSElement.Attrs, order)
}
