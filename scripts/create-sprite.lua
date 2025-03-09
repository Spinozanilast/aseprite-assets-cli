--[[
Aseprite Create sprite script
======================
Creates sprite with specified dimensions and color mode.

Features:
- Supports custom sprite dimensions
- Supports custom color mode

Usage:
Requires Aseprite CLI parameters:
  --output-path    Output sprite path
  --width          Sprite width (default: 32)
  --height         Sprite height (default: 32)
  --color-mode     Sprite color mode (default: rgb)
]]

local DEFAULT_SIZE = 32
local DEFAULT_COLOR_MODE = 'rgb'
local COLOR_MODES = {
    ["indexed"] = ColorMode.INDEXED,
    ["rgb"] = ColorMode.RGB,
    ["gray"] = ColorMode.GRAYSCALE,
}

local output_path = app.params["output-path"]
if not output_path or output_path == "" then
    error("Missing required parameter: output-path")
end

local function parse_positive_number(param, default)
    local num = tonumber(param)
    return (num and num > 0) and math.floor(num) or default
end

local sprite_width = parse_positive_number(tonumber(app.params["width"]), DEFAULT_SIZE)
local sprite_height = parse_positive_number(tonumber(app.params["height"]), DEFAULT_SIZE)

local color_mode = COLOR_MODES[app.params["color-mode"]] or COLOR_MODES[DEFAULT_COLOR_MODE]

-- https://www.aseprite.org/api/command/NewFile#newfile
app.command.NewFile {
    ui = false,
    width = sprite_width,
    height = sprite_height,
    colorMode = color_mode,
    fromClipboard = false
}

-- https://www.aseprite.org/api/command/SaveFile#savefile
app.command.SaveFileAs {
    filename = output_path,
}
