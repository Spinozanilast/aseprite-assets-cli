
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
```


## 

### From Releases 
1. Download the latest release from [here](https://github.com/Spinozanilast/aseprite-assets-cli/releases/latest)
2. Extract the archive to a folder of your choice
3. Add the folder to your PATH environment variable
4. And you are ready to go!
   
```pwsh
aseprite-assets --help
```

## TO DO
- Add more features
  - Add support for opening multiple files at once *(currently only one file can be opened at a time)*
  - Add support for assets management (removing, creating, renaming, copying)
  - Add integration with aseprite cli for better and easier (with less steps and tui) interaction
  - Add support for palletes importing (maybe)
- Add more documentation
- Add deployment to winget and brew


