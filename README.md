<img src="https://cdn.hackclub.com/01a1bab1-17b7-78de-a4e5-45f16d47bf0d/ng-dl2.svg" alt="ng-dl">

# ng-dl
Like yt-dlp but for Newgrounds: download Newgrounds assets from the command line!
## Support
| Feature | Support | Notes |
|---------|---------|-------|
| Audio | Supported | Not tested that much, but works. |
| Lists/Collections/Series | Supported | Ditto |
| User Assets | Supported | Ditto |
| Movies | Partially supported | Can download movies with limited metadata. Untested. |
| Art | Partially supported | Does not write any EXIF metadata. |
| Games | Unsupported | |
| NG Guard | Unsupported | **See note below.** |
| Other | Unsupported | |

ng-dl is in alpha. Once all features are supported, this table will be removed and ng-dl will leave alpha.
## Usage (with examples)
> [!NOTE]
> If you encounter NG Guard or "Error Forbidden occured while trying to fetch", try opening www.newgrounds.com in a web browser, then try running ng-dl again.

To use ng-dl, open a terminal window and navigate to the directory with the executable file, then follow the usage below:
> [!TIP]
> On Windows, if you place the ng-dl executable in a folder in your PATH, you will be able to use it from anywhere. To edit your PATH, search for "Edit the system environment variables" in the Start menu and open it, click Environment Variables..., click Path, click Edit..., click New, enter the name of the directory, then click OK in all the windows.   
> On Linux, just place the executable in /usr/local/bin/.
```
ng-dl [options] <url>
```
### Command line flags
> [!NOTE]
> The URL to download must be after all command line flags.

- `--dont-rename`: Do not rename downloaded files, keep them with the same filename as the one on the server.
- `--dont-fix-metadata`: Do not apply fetched metadata to downloaded files. Also does not apply album art.
- `--resolution <text>`: Requests a specific resolution for a video (e.g. 1080p). If the resolution cannot be found, it will default to the highest resolution.
- `--no-progress-bar`: Hides the top progress bar.
- `--no-parallel`: Does not download files concurrently; waits for one download to finish before starting another.
- `--parallel-wait <number>`: When downloading files concurrently, how long to wait between starting each download.
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
ng-dl also has special strings that you can put instead of a Newgrounds link:
- `featured-audio` will download the audio featured on the Newgrounds front page.
- `featured-movies` will download the movies featured on the Newgrounds front page.
### Example 2: Downloading featured audio
```
$ ng-dl featured-audio
Fetching https://www.newgrounds.com/audio/listen/XXXXXX
...
Renamed file to Artist - Song.mp3
```
