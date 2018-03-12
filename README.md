## GoPro Utilities

This repo combines the work of:

[github.com/stilldavid/gopro-utils]() <br>
[github.com/paulmach/go.geo]() <br>
[github.com/tkrajina/gpxgo]() <br>

### Setup Notes
Do **NOT** change the repo name (except to remove `-master` if present) or else you will have to go through every file and change the import statements!

Save this repo in your profile's Go directory. `C:\Users\USERNAME\go\src`

### Usage
The following steps will enable you to extract the GPS data from a GoPro video by: first, extracting the metadata from the mp4 file as a binary file; second, converting that file to a GPX file. In order to do this you will need to install a utility called **ffmpeg** and the **Go** language. The download and installation process can be found in the [Wiki](https://github.com/CenturyGIS/cei-gopro-utils/wiki/Installation-Guide).

#### Steps
1. `ffprobe [FILENAME].mp4`.
2. `ffmpeg -y -i [INFILENAME].mp4 -codec copy -map 0:[STREAMID] -f rawvideo [OUTFILENAME].bin`
3. `gopro2gpx -i [INFILENAME].bin -o [OUTFILENAME].gpx`

Steps 1 and 2 must be done with the mp4 file in the `ffmpeg\bin` folder.
