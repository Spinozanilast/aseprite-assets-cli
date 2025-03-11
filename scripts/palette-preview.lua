--[[
Aseprite Palette Preview script
======================
Creates formatted ansi palette preview for terminal output.

Features:
- Supports custom color output format (hex, rgb)
- Supports custom colors num in row for output

Usage:
Requires Aseprite CLI parameters:
  --palette-filename  Input palette path
  --color-format      Output color format (hex, rgb)
  --output-row-count Number of colors in row for output (default: 5)
]]

local DEFAULT_ROW_LENGTH = 5
local MAX_ANSI_COLOR = 255
local ERROR_COLOR = { 255, 0, 0 }
local HEADER_COLOR = { 50, 200, 255 }
local SWATCH_TEXT_COLOR = { 30, 35, 32 }

local COLOR_FORMATS = {
    hex = function(r, g, b, a)
        return string.format("#%02x%02x%02x%02x ", r, g, b, a)
    end,
    rgb = function(r, g, b, a)
        return string.format("(%3d,%3d,%3d,%3d) ", r, g, b, a)
    end
}

--Palette filename param
local palette_filename = app.params["palette-filename"]
if not palette_filename or palette_filename == "" then
    error("Missing required parameter: palette-filename")
end

--Color output type param
local color_output_type = app.params["color-format"] or "hex"
local color_formatter = COLOR_FORMATS[color_output_type:lower()] or COLOR_FORMATS.hex

--Line elements count param
local line_elements_count = math.max(tonumber(app.params["output-row-count"]) or DEFAULT_ROW_LENGTH, 1)

local function ansi_color_escape(r, g, b, background)
    local prefix = background and 48 or 38
    return string.format("\27[%d;2;%d;%d;%dm", prefix, r, g, b)
end

local function styled_text(text, fg_r, fg_g, fg_b, bg_r, bg_g, bg_b)
    local parts = {}
    if bg_r then
        table.insert(parts, ansi_color_escape(bg_r, bg_g, bg_b, true))
    end
    table.insert(parts, ansi_color_escape(fg_r, fg_g, fg_b, false))
    table.insert(parts, text)
    table.insert(parts, "\27[0m")
    return table.concat(parts)
end

local function create_ansi_block(r, g, b, a)
    local base_text = color_formatter(r, g, b, a)
    local display_text = styled_text(base_text, table.unpack(SWATCH_TEXT_COLOR))

    local blended_r = math.floor(r * (a / MAX_ANSI_COLOR))
    local blended_g = math.floor(g * (a / MAX_ANSI_COLOR))
    local blended_b = math.floor(b * (a / MAX_ANSI_COLOR))

    return string.format("%s %s \27[0m",
            ansi_color_escape(blended_r, blended_g, blended_b, true),
            display_text
    )
end

local function load_palette(filename)
    local ok, palette = pcall(function()
        return Palette { fromFile = filename }
    end)
    if not ok or not palette then
        error(styled_text("Failed to load palette file: " .. filename, table.unpack(ERROR_COLOR)))
    end
    return palette
end

local palette = load_palette(palette_filename)
local color_count = #palette

if color_count == 0 then
    print(styled_text("Empty palette file contains no colors", table.unpack(ERROR_COLOR)))
    return
end

local output_chunks = {
    styled_text(string.format("Palette Preview (%d colors):\n", color_count), table.unpack(HEADER_COLOR))
}

for i = 0, color_count - 1, line_elements_count do
    local row = {}
    for j = 0, line_elements_count - 1 do
        local idx = i + j
        if idx >= color_count then
            break
        end

        local color = palette:getColor(idx)
        table.insert(row, create_ansi_block(
                color.red,
                color.green,
                color.blue,
                color.alpha
        ))
    end
    table.insert(output_chunks, table.concat(row, ""))
end

print(table.concat(output_chunks, "\n") .. "\n")