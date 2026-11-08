<img src="https://user-cdn.hackclub-assets.com/01a1b87b-0fcf-7ee2-a882-1c4a2fdf55fe/ng-dl.svg" alt="ng-dl">

# ng-dl
Like yt-dlp but for Newgrounds: download Newgrounds assets from the command line!
## Support
| Feature | Support | Notes |
|---------|---------|-------|
| Audio | Supported | Not tested that much, but works |
| Lists/Collections/Series | Partially supported | Only audio lists work |
| Movies | Partially supported | Can download movies but doesn't add metadata. Untested. |
| Games | Unsupported | |
| Art | Unsupported | |
| NG Guard | Unsupported | |
| Other | Unsupported | |

ng-dl is in alpha. Once all features are supported, this table will be removed and ng-dl will leave alpha.
## Usage (with examples)
```
ng-dl [options] <url>
```
### Example 1: Download audio
```
$ ng-dl https://www.newgrounds.com/audio/listen/639309
Fetching https://www.newgrounds.com/audio/listen/639309
Song name is Sidestep
Fetched metadata for Sidestep
Applied metadata to file.
Attached album art.
Renamed file to AeronMusic - Sidestep.mp3

```