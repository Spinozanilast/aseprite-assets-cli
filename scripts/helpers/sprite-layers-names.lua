--[[
Aseprite Sprite helper script to get layers names every on separate row
======================
Prints sprite layers names.

Usage:
Requires Aseprite CLI parameters:
  --sprite-filename Sprite filename to output
]]

local sprite_filename = app.params["sprite-filename"]
if not sprite_filename or sprite_filename == "" then
    error("Missing required parameter: sprite-filename")
end

local sprite = Sprite { fromFile = sprite_filename }
if not sprite then
    error("Such sprite doesnt exists or have inner issues")
end

local layers_number = #sprite.layers
if (layers_number < 1) then
    error("This sprite doesnt has some layers or smth went wrong!")
end

for i = 1, layers_number do
    print(sprite.layers[i].name)
end