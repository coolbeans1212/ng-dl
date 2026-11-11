<img src="https://cdn.hackclub.com/01a1bab1-17b7-78de-a4e5-45f16d47bf0d/ng-dl2.svg" alt="ng-dl">

# ng-dl
Like yt-dlp but for Newgrounds: download Newgrounds assets from the command line!
## Support
| Feature | Support | Notes |
|---------|---------|-------|
| Audio | Supported | Not tested that much, but works |
| Lists/Collections/Series | Partially supported | Only audio & movie lists work |
| Movies | Partially supported | Can download movies with limited metadata. Untested. |
| Games | Unsupported | |
| Art | Unsupported | |
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
