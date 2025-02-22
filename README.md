<div align="center">

[![Release](https://github.com/Spinozanilast/aseprite-assets-cli/actions/workflows/release.yml/badge.svg)](https://github.com/Spinozanilast/aseprite-assets-cli/actions/workflows/release.yml)
# Aseprite Assets CLI
<p>

![Aseprite Logo](https://github.com/aseprite/aseprite/blob/main/data/icons/ase128.png?raw=true) ![X](https://github.com/Spinozanilast/spinozanilast/blob/master/assets/X.png?raw=true") ![X](https://github.com/Spinozanilast/spinozanilast/blob/master/assets/spinozanilast.gif?raw=true")

</span>
</div>

## Description
CLI interface for aseprite assets interaction. For in-terminal opening of aseprite files and more.

## Commands Tree

``` bash
aseprite-assets
├── list (l)
│   └── List existing aseprite assets
├── config (cfg) [ARG]
│   ├── info: Display configuration information
│   ├── edit: Edit configuration using a TUI
│   └── open: Open configuration file
│       └── --app-path (-a): Specify app to open the config file
├── create (cr) [ARG]
│   └── Create a new aseprite asset with the specified options
│       ├── --name (-n): The name of the asset
│       ├── --ui (-u): Whether to open aseprite after asset creation
│       ├── --width (-w): The width of the asset
│       ├── --height: The height of the asset
│       ├── --mode (-m): The color mode of the asset (indexed, rgb, gray, tilemap)
│       └── --path (-p): The output path for the asset
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
aseprite-assets create --name "new_asset" --width 64 --height 64 --mode "rgb" --path "path/to/save"
```

## TO DO
- Add more features
  - Add support for opening multiple files at once *(currently only one file can be opened at a time)*
  - Add support for assets management (removing, creating, renaming, copying)
  - Add integration with aseprite cli for better and easier (with less steps and tui) interaction
  - Add support for palletes importing (maybe)
- Add more documentation
- Add deployment to winget and brew