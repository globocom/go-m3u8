//	Multivariant Playlist Tags (Section 4.4.6 on RFC)
//
// Multivariant Playlist tags define the Variant Streams, Renditions,
// and other global parameters of the presentation.
//
// Multivariant Playlist tags MUST NOT appear in a Media Playlist;
// clients MUST fail to parse any Playlist that contains both a
// Multivariant Playlist tag and either a Media Playlist tag or a Media
// Segment tag.
package multivariant

import (
	"strings"

	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
)

var (
	StreamInfTag      = "#EXT-X-STREAM-INF"         //todo: has attributes
	MediaTag          = "#EXT-X-MEDIA"              //todo: has attributes
	IFrameStramInfTag = "#EXT-X-I-FRAME-STREAM-INF" //todo: has attributes
	SessionKey        = "#EXT-X-SESSION-KEY"        //todo: has attributes
)

type StreamInfParser struct{}

type StreamInfEncoder struct{}

func (p StreamInfParser) Parse(tag string, playlist *pl.Playlist) error {
	playlist.CurrentStreamInf = pl.GetStreamInfData(pl.TagsToMap(tag))
	return nil
}

func (e StreamInfEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	order := []string{"BANDWIDTH", "AVERAGE-BANDWIDTH", "CODECS", "RESOLUTION", "FRAME-RATE"}
	if err := pl.EncodeTagWithAttributes(builder, StreamInfTag, node.HLSElement.Attrs, order); err != nil {
		return err
	}
	if node.HLSElement.URI != "" {
		_, err := builder.WriteString(node.HLSElement.URI + "\n")
		return err
	}
	return nil
}
