<div align="center">

[![Release](https://github.com/Spinozanilast/aseprite-assets-cli/actions/workflows/release.yml/badge.svg)](https://github.com/Spinozanilast/aseprite-assets-cli/actions/workflows/release.yml)

# Aseprite Assets CLI

A command-line interface for processing Aseprite files and assets, built with Go.

![logo](https://github.com/Spinozanilast/aseprite-assets-cli/blob/master/www/static/logo128.png?raw=true")

</div>

## Description

CLI interface for aseprite assets interaction. For in-terminal opening of aseprite files and more.

## Commands Tree

> [!NOTE]
> If commands missing arg - it means that they work with surveys (interactive questions) instead of flags (
> watch [surveys-structure](#surveys-structure))

``` bash
aseprite-assets
├── list (l) [FLAG]
│   └── List existing aseprite assets
├── config (cfg) [command]
│   ├── info: Display configuration information
│   ├── edit: Edit configuration using a TUI
│   └── open: Open configuration file
│       └── --app-path (-a): Specify app to open the config file
├── sprite (command)
│   └── create (c, cr): Create a new aseprite sprite with the specified options
├── palette (p)
│   └── create (c, cr): Create a new color palette using OpenAI API (surveys used instead of flags)
├── show (sh) [ARGS] [FLAG]
│   └── Preview aseprite sprite or palette in terminal
├── export (e, exp) [FLAGS]
│   └── Export existing aseprite (ase) files by format or output template path and with optional scales or sizes specified
```

## Installation

### From Releases

1. Download the latest release from [here](https://github.com/Spinozanilast/aseprite-assets-cli/releases/latest)
2. Extract the archive to a folder of your choice
3. Add the folder to your PATH environment variable
4. And you are ready to go!

```pwsh
aseprite-assets --help
```

## Usage Examples

### List Assets

To list all existing aseprite sprites, use:

```sh
aseprite-assets list -s
```
To list all existing aseprite palettes, use:

```sh
aseprite-assets list -p
```

To list all existing aseprite sprites recursively, use:

```sh
aseprite-assets list -pr
```

### Configure CLI

To display the current configuration:

```sh
aseprite-assets config info
```

To edit the configuration using a TUI:

```sh
aseprite-assets config edit
```

To open the configuration file with a specific application:

```sh
aseprite-assets config open --app-path "path/to/application"
```

### Create Sprite

To create a new aseprite sprite:

```sh
aseprite-assets sprite create
```

### Create Palette

To create a new color palette using OpenAI API, follow the interactive prompts:

```sh
aseprite-assets palette create
```

### Show Sprite or Palette

To preview an aseprite sprite or palette in the terminal:

```sh
aseprite-assets show --filename "path/to/file.aseprite"
```

To show a palette preview with custom colors format and palette size:

```sh
aseprite-assets show --filename "path/to/file.gpl" --color-format "rgb" --output-row-count 10 --palette-preview
```

### Show Sprite or Palette
---
To export some existing aseprite sprite to png use:

```sh
aseprite-assets export --sprite-filename "path/to/file.aseprite" --output-filename "path/to/output-file.png"
```
Or (in this case you will export sprite to the same location with changed extension only):

```sh
aseprite-assets export --sprite-filename "path/to/file.aseprite" --format png
```
-----

To export sprite for multiple sizes to the same as sprite location use:
```sh
# args separated by comma and format is numberxnumber
aseprite-assets export --sprite-filename "path/to/file.aseprite" --format png --sizes 32x32,64x64,128x128
```
Or (to export to different location):
```sh
# args separated by comma and format is numberxnumber
aseprite-assets export --sprite-filename "path/to/file.aseprite" --output-filename "path/to/file.png" --sizes 32x32,64x64,128x128
# command will create path/to/file_32x32.png, path/to/file_64x64.png, path/to/file_128x128.png files
```
---

To export sprite for multiple scales to the same as sprite location use:
```sh
# args separated by comma and format is numberxnumber
aseprite-assets export --sprite-filename "path/to/file.aseprite" --format png --scales 2,4,6
```
Or (to export to different location):
```sh
# args separated by comma and format is numberxnumber
aseprite-assets export --sprite-filename "path/to/file.aseprite" --output-filename "path/to/file.png" --scales 2,4,6
# command will create path/to/file_2x.png, path/to/file_4x.png, path/to/file_6x.png files
```
---


## Surveys Structure

### Create Sprite

The `sprite create` command uses the following survey questions:

1. **Asset name (without extension)**: The name of the sprite.
2. **Open aseprite after asset creation?**: Whether to open aseprite after sprite creation.
3. **Width**: The width of the sprite.
4. **Height**: The height of the sprite.
5. **Color mode**: The color mode of the sprite (indexed, rgb, gray, tilemap).
6. **Output path**: The output path for the sprite.

### Create Palette

The `palette create` command uses the following survey questions:

1. **Color palette description**: Description of the palette (e.g. 'love, robots, batman').
2. **Number of colors to generate**: Number of colors to generate (if 0 - generate all colors).
3. **AI model to use**: AI model to use.
   > Now available: `GPT3Dot5Turbo`, `GPT4oMini`, `GPT4o`, `O1Mini`, `GPT4Turbo`, `GPT4`,
4. **Enable advanced mode?**: Enable advanced mode.
5. **Include transparency?**: Include transparency in the colors (only asked if advanced mode is enabled).
6. **Directory to save palettes to**: Directory to save the palette.
7. **Palette name**: Name of the palette.
8. **Select file type**: File type of the palette (gpl, png).
9. **Select save variant**: Save variant (Save as preset, Save as file, Both).

### Export Sprite
1. **Sprite filename**: Path to the Aseprite file to export
    • Auto-complete suggests files from configured sprite folders
    • Validates file existence and extension (.ase, .aseprite)
2. **Output format**: Select export format 
    • Options: Available Aseprite export formats (png, json, gif, etc.)
3. **Output path**: Destination path for exported file 
    • Auto-complete suggests filename based on input file + format
    • Validates non-overwrite and valid extension
4. Configure **scaling/sizing** options?: Yes/No to configure resizing 
    • Default: Yes
5. **Resize mode** (only if Yes to previous): 
    • Options: scaled (scale multipliers), sized (exact dimensions)
6. **Enter scales** (if scaled mode chosen):
    • Comma-separated numbers (e.g., "1,2,3")
    • Auto-suggests common presets like "0.5,1,2"
    • Validates numeric format

OR

7. **Enter sizes** (if sized mode chosen):
    • Comma-separated WxH pairs (e.g., "64x64,128x128")
    • Auto-suggests common presets like "64x64"
    • Validates WxH format

## TO DO

- [ ] Features
    - [x] Palette creating with OpenAI API
    - [x] Assets creation
    - [x] Assets listing with open possibility
    - [x] Exporting assets (multiple sizes in one time (with scale))
    - [ ] Support for opening multiple files at once *(currently only one file can be opened at a time)*
    - [ ] Support for assets management (removing, creating, renaming, copying)
    - [x] ~ Add integration with aseprite cli for better and easier (with less steps and tui) interaction
    - [ ] Add support for palletes importing (maybe)
- [ ] More documentation
- [ ] Add deployment to winget and brew

------
<div align="center">

## With love to [Aseprite](https://www.aseprite.org/)

![Aseprite Logo](https://github.com/aseprite/aseprite/blob/main/data/icons/ase128.png?raw=true) ![X](https://github.com/Spinozanilast/spinozanilast/blob/master/assets/X.png?raw=true") ![spinozanilast](https://github.com/Spinozanilast/spinozanilast/blob/master/assets/spinozanilast.gif?raw=true")

</div>
