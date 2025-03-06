local aseprite_filename = app.params["sprite-filename"]
local size = tonumber(app.params["size"] or 0)

local function colorize(text, r, g, b)
    return string.format("\27[38;2;%d;%d;%dm%s\27[0m", r, g, b, text)
end

local function make_color_block(r, g, b, a)
    if a == 0 then
        return "  " -- Two spaces for transparent pixels
    end
    return string.format(
            "\27[48;2;%d;%d;%dm  \27[0m",
            math.floor(r * (a / 255)),
            math.floor(g * (a / 255)),
            math.floor(b * (a / 255))
    )
end

local sprite = Sprite { fromFile = aseprite_filename }

if (size ~= 0) then
    local resizedSprite = Sprite(sprite)
    resizedSprite:resize(size, size)
    sprite = resizedSprite
end

local image = Image(sprite)

if image:isEmpty() then
    return
end

local w = image.width
local h = image.height

if size ~= 0 then
    print(colorize(string.format("Sprite Preview (compressed a lot: %dx%d) :\n\n", size, size), 200, 100, 100))
else
    print(colorize(string.format("Sprite Preview (%dx%d) :\n\n", w, h), 50, 200, 255))
end

local output = {}
for y = 0, h - 1 do
    local line = {}
    for x = 0, w - 1 do
        local rgba = image:getPixel(x, y)
        local r = app.pixelColor.rgbaR(rgba)
        local g = app.pixelColor.rgbaG(rgba)
        local b = app.pixelColor.rgbaB(rgba)
        local a = app.pixelColor.rgbaA(rgba)

        table.insert(line, make_color_block(r, g, b, a))
    end
    table.insert(output, table.concat(line))
end

print(table.concat(output, "\n") .. "\n\n")

if size == 0 then
    print(colorize("[TIP]: ", 255, 200, 0))                                 -- Yellow
    print(colorize("Zoom out terminal for large sprites\n", 150, 150, 150)) -- Grey
end
