local preset_name = app.params["preset-name"]
local palette_filename = app.params["palette-filename"]


-- https://github.com/aseprite/aseprite/blob/f2b870a17f2a3a6c1b00e341b2b1d0216db0dc3b/src/app/commands/cmd_load_palette.cpp#L28
app.command.LoadPalette {
    filename = palette_filename,
}

-- https://github.com/aseprite/aseprite/blob/f2b870a17f2a3a6c1b00e341b2b1d0216db0dc3b/src/app/commands/cmd_save_palette.cpp
app.command.SavePalette {
    preset = preset_name,
    saveAsPreset = true
}
