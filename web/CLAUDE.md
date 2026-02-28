# Eye of the Bee-holder - Web Game

## Project Overview

A fully playable web-based implementation of the "Eye of the Bee-holder" board game, built as a single HTML file (`eye-of-the-beeholder.html`). Includes interactive play, AI opponents, batch simulation, and an in-app rules reference.

## Source Assets (inputs)

- `../rules/EyeOfTheBeeholderRules.pdf` - The official game rules. The HTML was built entirely from this PDF.
- `../generate-cards/data.csv` - 64 rows mapping each bee card to its name, six attributes, and original image filename. **Card file numbering follows CSV row order** (row 0 = `card_00.png`), NOT the `monster_XX` number in the `image_filename` column.
- `cards/card_00.png` through `cards/card_63.png` - Card artwork (825x1125 PNG). Numbered by CSV row index.
- `rules_layout_diagram.png` - Hex tile layout diagram used in the Rules tab.

## Generated Output

- `eye-of-the-beeholder.html` - Single-file web app. Open in a browser from this directory (needs `cards/` and `rules_layout_diagram.png` as relative paths).

## How to Regenerate

### 1. Parse card data from CSV

Card IDs are the 0-based row index in `../generate-cards/data.csv` (skipping the header). For each row:
- Extract name from column 0
- Map attribute columns 1-6 to binary: `fuzzy=0/shiny=1, feathered=0/whips=1, stinger=0/mandibles=1, striped=0/solid=1, sleek=0/flutter=1, honey=0/pollen=1`
- Image path: `cards/card_XX.png` where XX is the zero-padded row index

**Important:** Do NOT use the `monster_XX` number from the `image_filename` column as the card ID. The card files are numbered by their row position in the CSV.

### 2. Game engine

Implements the rules from `../rules/EyeOfTheBeeholderRules.pdf` with these clarifications (deviations from the PDF text):

- **Scoring:** Winning a round is worth **1 point** (the PDF's wording about collecting "all played bee cards" each worth 1 point is misleading - it's 1 point per round win, not N points where N is player count).
- **Last round of hand:** Skip the Manipulate phase after round 7, since tiles are re-drafted immediately.
- **Sudden death timing:** Check for an outright leader after the **Judge** phase, before Manipulation (not after the full round).

### 3. AI strategies

Five AI strategies are implemented:
- **Random** - Baseline, picks randomly
- **Greedy** - Plays card with highest weighted facet match score
- **Strategic** - Evaluates cards considering post-play manipulation potential
- **Defensive** - Heavily prioritizes matching slots 1-2; disrupts opponents via manipulation
- **Adaptive** - Plays greedy when ahead, defensive when behind

### 4. Attribute display

Attributes are displayed in this order with these colors (matching the card artwork):
- Texture: gray (#757575)
- Antennae: brown (#6D4C41)
- Weapon: red (#D32F2F)
- Pattern: blue (#1565C0)
- Wings: green (#2E7D32)
- Payload: gold (#BF8C00)

### 5. Three tabs

- **Play Game** - Interactive game with configurable human/AI players (2-5), drafting UI, card selection, judging visualization, manipulation controls, pass-and-play for multiple humans
- **Simulate** - Batch AI-only games (up to 5000). Round-robin across all strategy combos or custom matchups. Shows win rates, average scores, game length, sudden death rate, seat position stats.
- **Rules** - HTML version of `../rules/EyeOfTheBeeholderRules.pdf` including the sample card image (`cards/card_00.png`) and tile layout diagram (`rules_layout_diagram.png`). Reflects the three rule clarifications listed above.
