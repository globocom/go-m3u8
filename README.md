<div align=center><img src="gopher.png" height=100px></div>
<h1 align=center>
go-m3u8
<div align=center><img src="https://github.com/globocom/go-m3u8/actions/workflows/go.yml/badge.svg"></div>
</h1>


## Golang m3u8 parser


This library has the goal of parse a HLS [m3u8](https://tools.ietf.org/html/rfc8216) into a doubly linked list and, for some tags, we can access the tag struct from the node.

### 1. Doubly linked list

This data structure allows us to access the manifest sorted and apply operations (modify, add and remove) to them. So we can access the HLS data (decode), manipulate it, and get it back in string format (encode).

#### Exemples of usefull node operation:

- **Add** discontinuity tag for `SSAI segments manipulation`
- **Change** discontinuity sequence tag count
- Remove DRM for SSAI segments manipulation by **adding** the tag `#EXT-X-KEY:METHOD=NONE`
- **Add** `SGAI` at DateRange tags
- **Remove** packager comment lines


### 2. Objects

To simplify the way we access the tags attributes, some nodes can be accessed via [custom structs](https://github.com/globocom/go-m3u8/blob/main/types.go):

- Media Manifest (StreamInf)
- Segment
- DateRange


## Getting started

```
go get github.com/globocom/go-m3u8
```