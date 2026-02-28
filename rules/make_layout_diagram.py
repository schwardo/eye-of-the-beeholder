#!/usr/bin/env python3
"""Create a diagram of the Queen's Favor tile surrounded by 6 facet tiles,
all cropped to flat-top hexagons and packed tightly together."""

import math
from PIL import Image, ImageDraw, ImageFont
import os

BASE = os.path.dirname(os.path.abspath(__file__))
TGC = os.path.join(BASE, "tgc")
OUTPUT = os.path.join(BASE, "rules_layout_diagram.png")


def flat_top_hex_points(cx, cy, R):
    """Return vertices of a flat-top regular hexagon with circumradius R.
    Width = 2R, Height = sqrt(3)*R."""
    points = []
    for i in range(6):
        angle = i * math.pi / 3  # 0, 60, 120, 180, 240, 300
        px = cx + R * math.cos(angle)
        py = cy + R * math.sin(angle)
        points.append((px, py))
    return points


def crop_to_hex(img, R):
    """Crop an image to a flat-top hexagon.
    The hex has width=2R, height=sqrt(3)*R.
    We center-crop the source image to that bounding box, then mask."""
    hex_w = int(2 * R)
    hex_h = int(math.sqrt(3) * R)

    # Center-crop source to hex bounding box
    src_w, src_h = img.size
    left = (src_w - hex_w) // 2
    top = (src_h - hex_h) // 2
    img = img.crop((left, top, left + hex_w, top + hex_h)).convert("RGBA")

    # Create hex mask
    mask = Image.new("L", (hex_w, hex_h), 0)
    draw = ImageDraw.Draw(mask)
    # Shrink radius slightly to crop into the bleed area
    crop_inset = R * 0.03
    points = flat_top_hex_points(hex_w / 2, hex_h / 2, R - crop_inset)
    draw.polygon(points, fill=255)

    img.putalpha(mask)
    return img


# Load images
queen = Image.open(os.path.join(TGC, "queen-tile-illustration.png"))
facet_tiles = [
    Image.open(os.path.join(TGC, "tile-gray-front.png")),   # Texture: Shiny
    Image.open(os.path.join(TGC, "tile-brown-face.png")),    # Antennae: Whips
    Image.open(os.path.join(TGC, "tile-red-front.png")),     # Weapon: Mandibles
    Image.open(os.path.join(TGC, "tile-blue-face.png")),     # Pattern: Solid
    Image.open(os.path.join(TGC, "tile-green-back.png")),    # Wings: Flutter
    Image.open(os.path.join(TGC, "tile-gold-face.png")),     # Payload: Pollen
]

# All source images are 675x600.
# The hex inscribed in those images has circumradius ~ 300 (half of height=600).
# We'll use that as our R for cropping, then scale down for the diagram.
src_R = 290  # circumradius in source pixels (slightly less than 300 to crop bleed)

# Crop all tiles to hex at source resolution
queen_hex = crop_to_hex(queen, src_R)
facet_hexes = [crop_to_hex(t, src_R) for t in facet_tiles]

# Scale down for diagram
display_R = 100  # circumradius in diagram pixels
scale = display_R / src_R

def scale_img(img):
    w, h = img.size
    new_w = int(w * scale)
    new_h = int(h * scale)
    return img.resize((new_w, new_h), Image.LANCZOS)

queen_scaled = scale_img(queen_hex)
facet_scaled = [scale_img(f) for f in facet_hexes]

tile_w, tile_h = facet_scaled[0].size  # should be ~2*display_R x sqrt(3)*display_R

# For flat-top hex packing, neighbor distance (center-to-center) = sqrt(3) * R
neighbor_dist = math.sqrt(3) * display_R
gap = 3
neighbor_dist += gap

# Canvas
margin = 70
canvas_w = int(2 * (neighbor_dist + display_R) + 2 * margin)
canvas_h = int(2 * (neighbor_dist + display_R * math.sqrt(3) / 2) + 2 * margin)
canvas = Image.new("RGBA", (canvas_w, canvas_h), (255, 255, 255, 0))

ccx, ccy = canvas_w // 2, canvas_h // 2


def paste_hex_img(canvas, img, center_x, center_y):
    w, h = img.size
    x = int(center_x - w / 2)
    y = int(center_y - h / 2)
    canvas.paste(img, (x, y), img)


# Paste queen in center
paste_hex_img(canvas, queen_scaled, ccx, ccy)

# 6 neighbors: Slot 1 at top, clockwise
# For flat-top hex, neighbor directions are at 90, 30, -30, -90, -150, 150 degrees
neighbor_angles_deg = [90, 30, -30, -90, -150, 150]

for hex_img, angle_deg in zip(facet_scaled, neighbor_angles_deg):
    angle_rad = math.radians(angle_deg)
    nx = ccx + neighbor_dist * math.cos(angle_rad)
    ny = ccy - neighbor_dist * math.sin(angle_rad)
    paste_hex_img(canvas, hex_img, nx, ny)

draw = ImageDraw.Draw(canvas)

# Crop to content
bbox = canvas.getbbox()
if bbox:
    padding = 12
    crop = (
        max(0, bbox[0] - padding),
        max(0, bbox[1] - padding),
        min(canvas_w, bbox[2] + padding),
        min(canvas_h, bbox[3] + padding),
    )
    canvas = canvas.crop(crop)

# Save with white background
final = Image.new("RGB", canvas.size, (255, 255, 255))
final.paste(canvas, mask=canvas.split()[3])
final.save(OUTPUT, "PNG")
print(f"Diagram saved: {OUTPUT} ({final.size[0]}x{final.size[1]})")
