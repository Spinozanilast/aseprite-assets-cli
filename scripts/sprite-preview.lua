--[[
Aseprite Sprite Preview script
======================
Creates formatted ansi preview of sprite for terminal output.

Features:
- Supports custom output size for preview

Usage:
Requires Aseprite CLI parameters:
  --sprite-filename Sprite filename to output
  --size            Output size for preview (default: 0, that mean if size not specified, output will be original size)
]]
local COLOR = {
    WARNING = { 200, 100, 100 },
    INFO = { 50, 200, 255 },
    TIP_HEADER = { 255, 200, 0 }, -- Yellow
    TIP_BODY = { 150, 150, 150 }  -- Grey
}

local TRANSPARENT = "  " -- Two spaces for transparent pixels
local TO_ORIGINAL_SIZE_RESET = 0

local sprite_filename = app.params["sprite-filename"]
if not sprite_filename or sprite_filename == "" then
    error("Missing required parameter: sprite-filename")
end

local target_size = math.max(tonumber(app.params["size"]) or TO_ORIGINAL_SIZE_RESET, TO_ORIGINAL_SIZE_RESET)

local function colorize(text, r, g, b)
    return string.format("\27[38;2;%d;%d;%dm%s\27[0m", r, g, b, text)
end

local function make_ansi_block(r, g, b, a)
    if a == 0 then
        return TRANSPARENT
    end
    return string.format(
            "\27[48;2;%d;%d;%dm  \27[0m",
            math.floor(r * (a / 255)),
            math.floor(g * (a / 255)),
            math.floor(b * (a / 255))
    )
end

local function load_sprite(filename)
    local ok, sprite = pcall(function()
        return Sprite { fromFile = filename }
    end)
    if not ok or not sprite then
        error(styled_text("Failed to load sprite: " .. filename, table.unpack(COLOR.WARNING)))
    end
    return sprite
end

local function process_sprite_resize(sprite, size)
    if size == 0 or (sprite.width == size and sprite.height == size) then
        return sprite
    end

    local resized = Sprite(sprite)
    resized:resize(size, size)
    return resized
end

local original_sprite = load_sprite(sprite_filename)

local processed_sprite
if size ~= TO_ORIGINAL_SIZE_RESET then
    processed_sprite = process_sprite_resize(original_sprite, target_size)
end

local image = Image(processed_sprite or original_sprite)

if image:isEmpty() then
    print(colorize("Empty sprite image", table.unpack(COLOR.WARNING)))
    return
end

local width, height = image.width, image.height
local size_info = target_size ~= 0
        and string.format("compressed to %dx%d", width, height)
        or string.format("original size %dx%d", width, height)

print(colorize(
        string.format("Sprite Preview (%s):\n", size_info),
        table.unpack(target_size ~= 0 and COLOR.WARNING or COLOR.INFO)
))

local output = {}
for y = 0, height - 1 do
    local line = {}
    for x = 0, width - 1 do
        local rgba = image:getPixel(x, y)
        local r = app.pixelColor.rgbaR(rgba)
        local g = app.pixelColor.rgbaG(rgba)
        local b = app.pixelColor.rgbaB(rgba)
        local a = app.pixelColor.rgbaA(rgba)

        table.insert(line, make_ansi_block(r, g, b, a))
    end
    table.insert(output, table.concat(line))
end

print(table.concat(output, "\n"))

if target_size == 0 then
    print(colorize("\n[TIP]: ", table.unpack(COLOR.TIP_HEADER)))
    print(colorize("Zoom out terminal for large sprites\n", table.unpack(COLOR.TIP_BODY)))
end

if processed_sprite ~= original_sprite then
    processed_sprite:close()
end
original_sprite:close()