## Using the Scripts to Generate a GPX or JSON FILE
<br>

**NOTE: You must have a binary (bin) file in order to generate a GPX or JSON file.**

The scripts must be compiled to an executable before they can be run (only has to be done once). Inside each of these folders run ```go build```.


After the script has been compiled you may use it by running ```[NAME OF UTILITY] -i [PATH/TO/FILE].bin -o [PATH/WHERE/TO/STORE/FILE].gpx```.

In our case we only need to use gopro2gpx utility because the JSON utility doesn't format the JSON the way Fulcrum expects. This is only an issue because we are currently storing the video and GPX file in Fulcrum. This may change in the future.  

### Errors
You might encounter an error due to an incorrect path in an import statement. If this happens just change the import statement to reflect the current path to the file.
