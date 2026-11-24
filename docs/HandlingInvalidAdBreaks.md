# Handling Invalid Ad Breaks in the Playlist

This document describes how the parser handles invalid ad breaks in HLS playlists, using the methods `TrimInvalidBreaks`, `removeInvalidBreakTags`, `removeDuplicateBreakTags`, and `removeProgramStartDaterange`.

## What are Invalid Ad Breaks?

Invalid ad breaks are `#EXT-X-DATERANGE` tags that have issues such as:

- Malformed or missing `START-DATE`.
- Invalid `StartMediaSequence`.
- Duplicate ad breaks (same start date and duration).
- Missing, zero, negative, or excessive (>600s) `PLANNED-DURATION`.
By [RFC](https://datatracker.ietf.org/doc/html/rfc8216), the PLANNED-DURATION attribute is optional, but in this library we validate this field because it is important for handling ad breaks.

## Example of an Invalid Ad Break

The playlists mocks `/mocks/media/withInvalidBreaks.m3u8`, `/mocks/media/withDuplicatedBreaks.m3u8` and `/mocks/media/withProgramStartDuplicatedBreak.m3u8` have examples of playlist with invalid breaks.

## How Parser Removes Invalid Ad Breaks

### 1. `TrimInvalidBreaks` Method

This method iterates over all ad breaks and checks:

- If `START-DATE` cannot be parsed as RFC3339Nano, the break is removed.
- If `StartMediaSequence` is not a valid integer, the break is removed.
- If `PLANNED-DURATION` is zero, invalid, or greater than 600 seconds (10 minutes), the break is removed.
- If the ad break is preceded by the `## splice_insert(Program Start)` comment and has the same `START-DATE` and `SCTE35-OUT` as the previous ad break, it is removed with `removeProgramStartDaterange`. `PLANNED-DURATION` is not compared, since it may differ between both breaks.
- Otherwise, if there are duplicate ad breaks (same `START-DATE` and `PLANNED-DURATION`), the duplicate is removed with `removeDuplicateBreakTags`. In this case, the `## splice_insert(Program Start)` comment does not appear before the duplicate ad break.

We need to remove all tags that is related to the ad break, so `CUE-OUT`, `CUE-IN`, `PROGRAM-DATE-TIME` tags and related comments are removed together with the `DATERANGE` tag.

### 2. `removeInvalidBreakTags` Method

When an invalid ad break is found, the parser removes:

- `CueOut` tags immediately after the ad break.
- `ProgramDateTime` tags inside the ad break.
- `CueIn` tags inside the ad break, as well as associated comments and program date time tags.
- Specific comments like `## splice_insert(auto_return)` immediately before the ad break.
- The ad break itself.


### 3. `removeDuplicateBreakTags` Method

When duplicate ad breaks are found, the parser removes just the tags that is related to one of the duplicated break:

- `CueOut` tags immediately after the duplicate ad break.
- `CueIn` tags inside the duplicate ad break.
- Specific comments like `## splice_insert(auto_return)` immediately before the duplicate ad break.
- The duplicate ad break itself.

### 4. `removeProgramStartDaterange` Method

Some playlists signal the same ad break twice: one `DATERANGE` preceded by `## splice_insert(auto_return)` and another preceded by `## splice_insert(Program Start)`, both sharing a single `CUE-IN`. In this case, the `CUE-IN` belongs to the remaining break and must be kept, so the parser removes only:

- The `## splice_insert(Program Start)` comment.
- The Program Start `DATERANGE` tag.
- The `CueOut` tag immediately after it.

```m3u8
## splice_insert(auto_return)
#EXT-X-DATERANGE:ID="3221225474-1794316738",START-DATE="2026-11-10T13:18:58.966666Z",PLANNED-DURATION=40,SCTE35-OUT=0xFC30...
#EXT-X-CUE-OUT:40
## splice_insert(Program Start)
#EXT-X-DATERANGE:ID="4035-1794316738",START-DATE="2026-11-10T13:18:58.966666Z",PLANNED-DURATION=40,SCTE35-OUT=0xFC30...
#EXT-X-CUE-OUT:40
#EXT-X-PROGRAM-DATE-TIME:2026-11-10T13:18:58.966666Z
...
## Auto Return Mode
#EXT-X-CUE-IN
```

After running `TrimInvalidBreaks`:

```m3u8
## splice_insert(auto_return)
#EXT-X-DATERANGE:ID="3221225474-1794316738",START-DATE="2026-11-10T13:18:58.966666Z",PLANNED-DURATION=40,SCTE35-OUT=0xFC30...
#EXT-X-CUE-OUT:40
#EXT-X-PROGRAM-DATE-TIME:2026-11-10T13:18:58.966666Z
...
## Auto Return Mode
#EXT-X-CUE-IN
```

## Example

Suppose the playlist contains two ad breaks with the same `START-DATE` and `PLANNED-DURATION`:

```m3u8
## splice_insert(auto_return)
#EXT-X-DATERANGE:ID="4026554821-1747772692",START-DATE="2025-05-20T20:24:52.699999Z",PLANNED-DURATION=20,SCTE35-OUT=0xFC3025000000000BB800FFF01405F00059C67FEFFF253C6500FE001B7740000101010000CB09D539
#EXT-X-CUE-OUT:20
## splice_insert(auto_return)
#EXT-X-DATERANGE:ID="4026554822-1747772692",START-DATE="2025-05-20T20:24:52.699999Z",PLANNED-DURATION=20,SCTE35-OUT=0xFC3025000000000BB800FFF01405F00059C67FEFFF253C6500FE001B7740000101010000CB09D539
#EXT-X-CUE-OUT:20
#EXT-X-PROGRAM-DATE-TIME:2025-05-20T20:24:52.699999Z
```

After running `TrimInvalidBreaks`, only one ad break and one cue out will remain in the playlist.

```m3u8
## splice_insert(auto_return)
#EXT-X-DATERANGE:ID="4026554822-1747772692",START-DATE="2025-05-20T20:24:52.699999Z",PLANNED-DURATION=20,SCTE35-OUT=0xFC3025000000000BB800FFF01405F00059C67FEFFF253C6500FE001B7740000101010000CB09D539
#EXT-X-CUE-OUT:20
#EXT-X-PROGRAM-DATE-TIME:2025-05-20T20:24:52.699999Z
```

## Summary

The implemented logic ensures that the final playlist contains only valid and non-duplicate ad breaks, removing associated tags that could cause inconsistencies in ad playback.

For more details, see the methods in `playlist.go`:

- `TrimInvalidBreaks`
- `removeInvalidBreakTags`
- `removeDuplicateBreakTags`
- `removeProgramStartDaterange`