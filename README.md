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
├── config (cfg) [ARG]
│   ├── info: Display configuration information
│   ├── edit: Edit configuration using a TUI
│   └── open: Open configuration file
│       └── --app-path (-a): Specify app to open the config file
├── create (cr)
│   └── Create a new aseprite asset with the specified options
├── palette (p)
│   └── Create a new color palette using OpenAI API (surveys used instead of flags)
├── show (sh) [ARGS] [FLAG]
│   └── Preview aseprite sprite or palette in terminal
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

To list all existing aseprite assets, use:

```sh
aseprite-assets list
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

### Create Asset

To create a new aseprite asset:

```sh
aseprite-assets create
```

### Create Palette

To create a new color palette using OpenAI API, follow the interactive prompts:

```sh
aseprite-assets palette
```

### Show Asset or Palette

To preview an aseprite sprite or palette in the terminal:

```sh
aseprite-assets show --filename "path/to/file.aseprite"
```

To show a palette preview with custom colors format and palette size:

```sh
aseprite-assets show --filename "path/to/file.gpl" --color-format "rgb" --output-row-count 10 --palette-preview
```

## Surveys Structure

### Create Asset

The `create` command uses the following survey questions:

1. **Asset name (without extension)**: The name of the asset.
2. **Open aseprite after asset creation?**: Whether to open aseprite after asset creation.
3. **Width**: The width of the asset.
4. **Height**: The height of the asset.
5. **Color mode**: The color mode of the asset (indexed, rgb, gray, tilemap).
6. **Output path**: The output path for the asset.

### Create Palette

The `palette` command uses the following survey questions:

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

## TO DO

- [ ] Features
    - [x] Palette creating with OpenAI API
    - [x] Assets creation
    - [x] Assets listing with open possibility
    - [ ] Exporting assets (multiple sizes in one time (with scale))
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
