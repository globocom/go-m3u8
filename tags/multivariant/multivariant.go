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

const (
	StreamInfName = "StreamInf"
)

var (
	StreamInfTag      = "#EXT-X-STREAM-INF"
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
	orderAttr := []string{"BANDWIDTH", "AVERAGE-BANDWIDTH", "CODECS", "RESOLUTION", "FRAME-RATE", "VIDEO-RANGE", "AUDIO", "VIDEO", "SUBTITLES", "CLOSED-CAPTIONS"}
	shouldQuoteAttr := e.shouldQuoteStreamInf(node)

	if err := pl.EncodeTagWithAttributes(builder, StreamInfTag, node.HLSElement.Attrs, orderAttr, shouldQuoteAttr); err != nil {
		return err
	}
	if node.HLSElement.URI != "" {
		_, err := builder.WriteString(node.HLSElement.URI + "\n")
		return err
	}
	return nil
}

func (e StreamInfEncoder) shouldQuoteStreamInf(node *internal.Node) map[string]bool {
	shouldQuoteAttr := map[string]bool{
		"BANDWIDTH":         false,
		"AVERAGE-BANDWIDTH": false,
		"CODECS":            true,
		"RESOLUTION":        false,
		"FRAME-RATE":        false,
		"VIDEO-RANGE":       false,
		"AUDIO":             true,
		"VIDEO":             true,
		"SUBTITLES":         true,
		"CLOSED-CAPTIONS":   true,
	}

	// the value can be either a quoted-string or an enumerated-string with the value NONE
	if node.HLSElement.Attrs["CLOSED-CAPTIONS"] == "NONE" {
		shouldQuoteAttr["CLOSED-CAPTIONS"] = false
	}

	return shouldQuoteAttr
}
