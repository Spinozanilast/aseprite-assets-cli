--[[
Aseprite Save Palette as Preset script
======================
Save existing palette to preset for locality of palettes inside aseprite.

Usage:
Requires Aseprite CLI parameters:
  --preset-name Preset name of palette
  --palette-filename Input palette path
]]

local function error_if_empty_string(param, message)
    if not param or param == "" then
        error(message)
    end
end

local preset_name = app.params["preset-name"]
error_if_empty_string(preset_name)

local palette_filename = app.params["palette-filename"]
error_if_empty_string(palette_filename)

-- https://github.com/aseprite/aseprite/blob/f2b870a17f2a3a6c1b00e341b2b1d0216db0dc3b/src/app/commands/cmd_load_palette.cpp#L28
app.command.LoadPalette {
    filename = palette_filename,
}

-- https://github.com/aseprite/aseprite/blob/f2b870a17f2a3a6c1b00e341b2b1d0216db0dc3b/src/app/commands/cmd_save_palette.cpp
app.command.SavePalette {
    preset = preset_name,
    saveAsPreset = true
}
