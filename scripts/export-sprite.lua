--[[
Aseprite Export Script
======================
Exports sprites with scaling/sizing options while maintaining original file structure

Features:
- Multiple export scales (1x, 2x, etc.)
- Multiple export sizes (64x64, 128x128, etc.)
- Automatic filename generation
- Proper resource cleanup

Usage:
Requires Aseprite CLI parameters:
  --sprite-filename  Input sprite path
  --output-filename  Base output path
  --format           Output file format
  --sizes            Comma-separated size list (WxH)
  --scales           Comma-separated scale factors
]]

--TODO: add support for disabling multiple frames export

local sprite_filename = app.params["sprite-filename"]
local output_filename = app.params["output-filename"]
-- extension if that's more clear for you
local format = app.params["format"]
local sizes = app.params["sizes"]
local scales = app.params["scales"]

local function validate_parameters()
    if not sprite_filename or sprite_filename == "" then
        error("Missing required parameter: sprite-filename")
    end

    if not output_filename and not format then
        error("Either output-filename or format must be specified")
    end

    if sizes and scales then
        error("Cannot specify both sizes and scales - choose one resize method")
    end
end

local function parse_comma_list(input, validator)
    if not input then return nil end
    local items = {}

    for item in string.gmatch(input, "([^,]+)") do
        item = item:match("^%s*(.-)%s*$") -- Trim whitespace
        if not validator(item) then
            error("Invalid list item: " .. item)
        end
        table.insert(items, item)
    end

    return #items > 0 and items or nil
end

local function generate_output_path(base_path, suffix, ext, separator)
    local pattern = "(.-)(%..+)$"
    local name, extension = string.match(base_path, pattern)

    if not separator then separator = "_" end

    if not name then
        name = base_path
        extension = ""
    end

    extension = ext and ("." .. ext) or extension
    return string.format("%s%s%s%s", name, separator, suffix, extension)
end

local function save_scaled_versions(sprite, base_path, scales)
    local original_sprite = app.sprite
    app.sprite = sprite

    for _, scale in ipairs(scales) do
        local numeric_scale = tonumber(scale)
        if not numeric_scale or numeric_scale <= 0 then
            error("Invalid scale factor: " .. scale)
        end

        local output_path = generate_output_path(
            base_path,
            string.format("%dx", numeric_scale),
            format
        )

        print(output_path)
        app.command.SaveFileCopyAs {
            filename = output_path,
            scale = numeric_scale
        }
    end

    app.sprite = original_sprite
end

local function save_sized_versions(sprite, base_path, sizes)
    local original_sprite = app.sprite

    for _, size in ipairs(sizes) do
        local width, height = size:match("^(%d+)x(%d+)$")
        width = tonumber(width)
        height = tonumber(height)

        if not width or not height or width <= 0 or height <= 0 then
            error("Invalid size format: " .. size .. " - use WIDTHxHEIGHT (e.g. 128x64)")
        end

        local resized_sprite = Sprite(sprite)
        app.sprite = resized_sprite
        resized_sprite:resize(width, height)

        local output_path = generate_output_path(
            base_path,
            string.format("%dx%d", width, height),
            format
        )

        print(output_path)
        app.command.SaveFileCopyAs { filename = output_path }
        resized_sprite:close()
    end

    app.sprite = original_sprite
end

local function main()
    validate_parameters()

    local sprite = Sprite { fromFile = sprite_filename }
    if not sprite then
        error("Failed to load sprite: " .. sprite_filename)
    end

    local base_output = output_filename or
        generate_output_path(sprite_filename, "", format, "")

    if not sizes and not scales then
        app.command.SaveFileCopyAs {
            filename = base_output
        }
        sprite:close()
        return
    end

    if scales then
        local scale_list = parse_comma_list(scales, function(s)
            return tonumber(s) and true or false
        end)
        save_scaled_versions(sprite, base_output, scale_list)
    else
        local size_list = parse_comma_list(sizes, function(s)
            return s:match("^%d+x%d+$") and true or false
        end)
        save_sized_versions(sprite, base_output, size_list)
    end

    sprite:close()
end



local function error_handler(error)
    app.alert {
        title = "Export Error",
        text = string.format("Error during export:\n%s", error)
    }
end

xpcall(main, error_handler)
