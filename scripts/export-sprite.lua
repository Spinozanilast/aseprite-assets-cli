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
  --frames_included range of included frames to export
]]

local FRAMES_INCLUDED_SEPARATOR = ":"

local sprite_filename = app.params["sprite-filename"]
local output_filename = app.params["output-filename"]
-- extension if that's more clear for you
local format = app.params["format"]

local frames_included = app.params["frames-included"]
local sizes = app.params["sizes"]
local scales = app.params["scales"]

local function validate_parameters()
    if not sprite_filename and not sprite_filename then
        error("Either output-filename or format must be specified")
    end

    if not output_filename and not format then
        error("Either output-filename or format must be specified")
    end

    if sizes and scales then
        error("Cannot specify both sizes and scales - choose one resize method")
    end

    if frames_included then
        if not string.match(frames_included, "^%*$")
                and not string.match(frames_included, "^%d+$")
                and not string.match(frames_included, "^%d+:%d+$") then
            error(string.format("Invalid frames format ('%s' used). Use: '*' | 'N' | 'N:M'", frames_included))
        end
    end
end

local function parse_comma_list(input, validator)
    if not input then
        return nil
    end
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

local function parse_frames(frames_str, total_frames)
    if not frames_str or frames_str == "*" then
        return 0, total_frames - 1
    end

    if string.find(frames_included, FRAMES_INCLUDED_SEPARATOR) == nil then
        included_one_frame = tonumber(frames_included)
        return included_one_frame, included_one_frame
    end

    if string.match(frames_str, "^%d+" .. FRAMES_INCLUDED_SEPARATOR .. "%d+$") then
        local start_frame, end_frame = string.match(frames_str, "^(%d+)" .. FRAMES_INCLUDED_SEPARATOR .. "(%d+)$")
        start_frame = tonumber(start_frame)
        end_frame = tonumber(end_frame)
        return start_frame, end_frame
    end

    error("Invalid frames format: " .. frames_str)
end

local function generate_output_path(base_path, suffix, ext, separator)
    local pattern = "(.-)(%..+)$"
    local name, extension = string.match(base_path, pattern)

    if not separator then
        separator = "_"
    end

    if not name then
        name = base_path
        extension = ""
    end

    extension = ext and ("." .. ext) or extension
    return string.format("%s%s%s%s", name, separator, suffix, extension)
end

local function save_scaled_versions(sprite, base_path, sizes_set, from_frame, to_frame)
    local original_sprite = app.sprite
    app.sprite = sprite

    for _, scale in ipairs(sizes_set) do
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
            scale = numeric_scale,
            fromFrame = from_frame,
            toFrame = to_frame
        }
    end

    app.sprite = original_sprite
end

local function save_sized_versions(sprite, base_path, sizes_set, from_frame, to_frame)
    local original_sprite = app.sprite

    for _, size in ipairs(sizes_set) do
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
        app.command.SaveFileCopyAs {
            filename = output_path,
            fromFrame = from_frame,
            toFrame = to_frame
        }

        resized_sprite:close()
    end

    app.sprite = original_sprite
end

local function main()
    validate_parameters()

    local sprite
    local emptyFilename = false
    if sprite_filename == nil or sprite_filename == "" then
        sprite = app.sprite
        emptyFilename = true
    else
        sprite = Sprite { fromFile = sprite_filename }
    end

    if not sprite then
        if emptyFilename then
            error("Active sprite is not exists (open sprite in Aseprite)")
        end
        error("Failed to load sprite: " .. sprite_filename)
    end

    local base_output = output_filename or
            generate_output_path(sprite_filename, "", format, "")

    local total_frames = #sprite.frames
    local from_frame, to_frame = parse_frames(frames_included, total_frames)

    if not sizes and not scales then
        app.command.SaveFileCopyAs {
            filename = base_output,
            fromFrame = from_frame,
            toFrame = to_frame
        }
        sprite:close()
        return
    end

    if scales then
        local scale_list = parse_comma_list(scales, function(s)
            return tonumber(s) and true or false
        end)
        save_scaled_versions(sprite, base_output, scale_list, from_frame, to_frame)
    else
        local size_list = parse_comma_list(sizes, function(s)
            return s:match("^%d+x%d+$") and true or false
        end)
        save_sized_versions(sprite, base_output, size_list, from_frame, to_frame)
    end

    sprite:close()
end

local function error_handler(error)
    print("Error during export:\n%s", error)
end

xpcall(main, error_handler)

