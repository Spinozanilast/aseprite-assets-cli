local palette_filename = app.params["palette-filename"]
local color_output_type = app.params["color-format"]
local line_elements_count = tonumber(app.params["output-row-count"]) or 5

local function colorize(text, r, g, b)
    return string.format("\27[38;2;%d;%d;%dm%s\27[0m", r, g, b, text)
end

local ColorFormat = {
    HEX = "hex",
    RGB = "rgb",
}

local color_format = {
    ["rgb"] = ColorFormat.RGB,
    ["hex"] = ColorFormat.HEX,
}

local color_format = color_format[color_output_type] or ColorFormat.HEX

local function make_hex_color_block(r, g, b, a)
    local color_s
    if color_format == ColorFormat.HEX then
        color_s = string.format("#%02x%02x%02x%02x ", r, g, b, a)
    else
        color_s = string.format("(%3d,%3d,%3d,%3d) ", r, g, b, a)
    end

    color_s = colorize(color_s, 30, 35, 32)

    return string.format(
        "\27[48;2;%d;%d;%dm %s \27",
        math.floor(r * (a / 255)),
        math.floor(g * (a / 255)),
        math.floor(b * (a / 255)),
        color_s
    )
end

local palette = Palette { fromFile = palette_filename }
local n_colors = #palette

if n_colors == 0 then
    print(colorize("There are no colors in the palette\n", 255, 0, 0))
    return
end

print(colorize(string.format("Palette Preview (%d colors):\n\n", n_colors), 50, 200, 255))

local output = {}
for i = 0, n_colors - 1, line_elements_count do
    local line = {}
    for j = 0, line_elements_count - 1 do
        local index = i + j
        if index >= n_colors then break end

        local color = palette:getColor(index)
        local r = color.red
        local g = color.green
        local b = color.blue
        local a = color.alpha

        table.insert(line, make_hex_color_block(r, g, b, a))
    end
    if #line > 0 then
        table.insert(output, table.concat(line))
    end
end

print(table.concat(output, "\n") .. "\n\n")
