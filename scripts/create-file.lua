local ui = app.params["ui"] == "true" and true or false
local sprite_width = tonumber(app.params["width"]) or 32
local sprite_height = tonumber(app.params["height"]) or 32
local color_mode = app.params["color-mode"]
local output_path = app.params["output-path"]

local color_modes = {
    ["indexed"] = ColorMode.INDEXED,
    ["rgb"] = ColorMode.RGB,
    ["gray"] = ColorMode.GRAYSCALE,
    ["tilemap"] = ColorMode.BITMAP
}

if not color_modes[color_mode] then
    color_mode = "rgb"
end

local get_color_mode = function(color_mode)
    return color_modes[color_mode]
end

color_mode = get_color_mode(color_mode)

-- https://www.aseprite.org/api/command/NewFile#newfile
app.command.NewFile {
    ui = ui,
    width = sprite_width,
    height = sprite_height,
    colorMode = color_mode,
    fromClipboard = false
}

-- https://www.aseprite.org/api/command/SaveFile#savefile
app.command.SaveFileAs {
    filename = output_path,
}
