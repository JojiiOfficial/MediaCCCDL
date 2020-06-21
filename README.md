# MediaCCCDl
A tool to retrieve download URLs for [media.ccc.de](media.ccc.de) videos.

# Usage
`mediacccdl <url> [--format] [--lang]`

### Supported formats:
<b>Video:</b>
- mp4
- webm
<b>Audio:</b>
- mp3
- opus

### Supported languages
- deu/de
- eng/en

# Exit codes
0 - no error
1 - flag error
2 - HTTP error
3 - No URL found
