# capinfo
Capinfo dumps information for a JavaCard CAP file. It is written in Go 
using the [skythen/cap](https://github.com/skythen/cap) package. At the current stage
the tools is fairly basic and just accepts a single argument interpreted as the cap filename.

I have written the tool mainly to glance over a collection of CAP files we maintain as part of a project. 
It can be easily integrated into JetBrains IntelliJ (see description below) or some other IDE of your choice.

## Build

Simply use go build as follows.

```bash
go build capinfo.go
```

## Run

Usage: capinfo <filename>

## Integration with JetBrains IntelliJ

### Add External Command

* Open Settings
* Select "External Tools"
* Press "+" on the External Tools pane
  * Program: provide absolute path to the capinfo executable
  * Arguments: $FilePath$
  * Enable Checkbox "Make console active on message in stdout"

### Add Shortcut for External Command

* Open Settings
* Select "Keymap"
* Enter "capinfo" into the serach field on the Keymaps pane
* Choose "Add Keyboard Shortcut" on the highlighted "External Tool" item
* Assign a keystroke by pressing the appropriate key(s). For me <crtl><x> worked out fine as it was unassigned before. 
