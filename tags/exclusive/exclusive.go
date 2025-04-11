//	Media or Multivariant Playlist Tags (Section 4.4.2 in RFC)
//
// The tags in this section can appear in either Multivariant Playlists
// or Media Playlists.
//
// If one of these tags appears in a Multivariant
// Playlist, it SHOULD NOT appear in any Media Playlist referenced by
// that Multivariant Playlist.  A tag that appears in both MUST have the
// same value; otherwise, clients SHOULD ignore the value in the Media
// Playlist(s).
//
// Tags in this section MUST NOT appear more than once in a Playlist.
// If one does, clients MUST fail to parse the Playlist.  The only
// exception to this rule is EXT-X-DEFINE, which MAY appear more than
// once.
package exclusive

import (
	"strings"

	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
)

const (
	IndependentSegmentsTag = "#EXT-X-INDEPENDENT-SEGMENTS"
	StartTag               = "#EXT-X-START"  //todo: has attributes
	VariableDefineTag      = "#EXT-X-DEFINE" //todo: has attributes
)

type IndependentSegmentsParser struct{}

type IndependentSegmentsEncoder struct{}

func (p IndependentSegmentsParser) Parse(tag string, playlist *pl.Playlist) error {
	playlist.Insert(&internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "IndependentSegments",
			Attrs: map[string]string{
				IndependentSegmentsTag: "",
			},
		},
	})
	return nil
}

func (e IndependentSegmentsEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	_, err := builder.WriteString(IndependentSegmentsTag + "\n")
	return err
}
