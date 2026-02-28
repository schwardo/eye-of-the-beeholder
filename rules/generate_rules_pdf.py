#!/usr/bin/env python3
"""Generate a PDF rules guide for Eye of the Bee-holder."""

from reportlab.lib.pagesizes import letter
from reportlab.lib.styles import getSampleStyleSheet, ParagraphStyle
from reportlab.lib.units import inch
from reportlab.lib.colors import HexColor
from reportlab.lib.enums import TA_CENTER, TA_LEFT
from reportlab.platypus import (
    SimpleDocTemplate, Paragraph, Spacer, Table, TableStyle,
    HRFlowable, KeepTogether, Image
)
import os

# Colors
HONEY_GOLD = HexColor("#D4A017")
DARK_BROWN = HexColor("#3E2723")
WARM_BG = HexColor("#FFF8E7")
TABLE_HEADER_BG = HexColor("#F5E6CC")
TABLE_STRIPE = HexColor("#FFFDF5")
RULE_COLOR = HexColor("#C8A96E")

def build_styles():
    styles = getSampleStyleSheet()

    styles.add(ParagraphStyle(
        'GameTitle',
        parent=styles['Title'],
        fontSize=28,
        leading=34,
        textColor=DARK_BROWN,
        spaceAfter=4,
        alignment=TA_CENTER,
    ))
    styles.add(ParagraphStyle(
        'GameSubtitle',
        parent=styles['Normal'],
        fontSize=12,
        leading=16,
        textColor=HexColor("#6D4C41"),
        alignment=TA_CENTER,
        spaceAfter=12,
    ))
    styles.add(ParagraphStyle(
        'Flavor',
        parent=styles['Normal'],
        fontSize=11,
        leading=15,
        textColor=HexColor("#5D4037"),
        alignment=TA_CENTER,
        fontName='Times-Italic',
        spaceAfter=16,
        spaceBefore=4,
    ))
    styles.add(ParagraphStyle(
        'SectionHead',
        parent=styles['Heading1'],
        fontSize=16,
        leading=20,
        textColor=DARK_BROWN,
        spaceBefore=18,
        spaceAfter=8,
        borderWidth=0,
    ))
    styles.add(ParagraphStyle(
        'SubsectionHead',
        parent=styles['Heading2'],
        fontSize=13,
        leading=16,
        textColor=HexColor("#5D4037"),
        spaceBefore=12,
        spaceAfter=6,
    ))
    styles.add(ParagraphStyle(
        'Body',
        parent=styles['Normal'],
        fontSize=10,
        leading=14,
        textColor=HexColor("#333333"),
        spaceAfter=6,
    ))
    styles.add(ParagraphStyle(
        'BodyBold',
        parent=styles['Normal'],
        fontSize=10,
        leading=14,
        textColor=HexColor("#333333"),
        fontName='Helvetica-Bold',
        spaceAfter=6,
    ))
    styles.add(ParagraphStyle(
        'BulletItem',
        parent=styles['Normal'],
        fontSize=10,
        leading=14,
        textColor=HexColor("#333333"),
        leftIndent=18,
        spaceAfter=3,
        bulletIndent=6,
    ))
    styles.add(ParagraphStyle(
        'SubBullet',
        parent=styles['Normal'],
        fontSize=10,
        leading=14,
        textColor=HexColor("#333333"),
        leftIndent=36,
        spaceAfter=2,
        bulletIndent=24,
    ))
    styles.add(ParagraphStyle(
        'NumberedStep',
        parent=styles['Normal'],
        fontSize=10,
        leading=14,
        textColor=HexColor("#333333"),
        leftIndent=18,
        spaceAfter=4,
    ))
    styles.add(ParagraphStyle(
        'QuickRef',
        parent=styles['Normal'],
        fontSize=10,
        leading=14,
        textColor=HexColor("#333333"),
        spaceAfter=4,
        leftIndent=6,
    ))
    styles.add(ParagraphStyle(
        'StrategyTip',
        parent=styles['Normal'],
        fontSize=9.5,
        leading=13,
        textColor=HexColor("#444444"),
        leftIndent=18,
        spaceAfter=5,
        bulletIndent=6,
    ))
    styles.add(ParagraphStyle(
        'Footer',
        parent=styles['Normal'],
        fontSize=9,
        leading=12,
        textColor=HexColor("#888888"),
        alignment=TA_CENTER,
        spaceBefore=16,
    ))

    return styles


def hr():
    return HRFlowable(
        width="100%", thickness=1, color=RULE_COLOR,
        spaceBefore=8, spaceAfter=8
    )


def make_table(headers, rows, col_widths=None):
    data = [headers] + rows
    style_cmds = [
        ('BACKGROUND', (0, 0), (-1, 0), TABLE_HEADER_BG),
        ('TEXTCOLOR', (0, 0), (-1, 0), DARK_BROWN),
        ('FONTNAME', (0, 0), (-1, 0), 'Helvetica-Bold'),
        ('FONTSIZE', (0, 0), (-1, -1), 10),
        ('LEADING', (0, 0), (-1, -1), 14),
        ('ALIGN', (0, 0), (-1, -1), 'LEFT'),
        ('VALIGN', (0, 0), (-1, -1), 'TOP'),
        ('GRID', (0, 0), (-1, -1), 0.5, RULE_COLOR),
        ('TOPPADDING', (0, 0), (-1, -1), 4),
        ('BOTTOMPADDING', (0, 0), (-1, -1), 4),
        ('LEFTPADDING', (0, 0), (-1, -1), 6),
        ('RIGHTPADDING', (0, 0), (-1, -1), 6),
    ]
    # Stripe odd data rows
    for i in range(1, len(data)):
        if i % 2 == 0:
            style_cmds.append(('BACKGROUND', (0, i), (-1, i), TABLE_STRIPE))

    t = Table(data, colWidths=col_widths, repeatRows=1)
    t.setStyle(TableStyle(style_cmds))
    return t


def build_pdf(output_path):
    styles = build_styles()

    doc = SimpleDocTemplate(
        output_path,
        pagesize=letter,
        leftMargin=0.75 * inch,
        rightMargin=0.75 * inch,
        topMargin=0.65 * inch,
        bottomMargin=0.65 * inch,
    )

    story = []
    usable = doc.width

    # --- Title ---
    story.append(Spacer(1, 20))
    story.append(Paragraph("Eye of the Bee-holder", styles['GameTitle']))
    story.append(Paragraph(
        "A bee beauty pageant card game for 2\u20135 players  |  30 minutes",
        styles['GameSubtitle']
    ))
    story.append(Paragraph(
        "In the bee world, beauty is in the eye of the bee-holder! "
        "Compete to present the most desirable bees according to the Queen\u2019s "
        "ever-changing standards. Draft the criteria, play your best bees, "
        "and manipulate what the hive considers beautiful.",
        styles['Flavor']
    ))
    story.append(hr())

    # --- Components ---
    story.append(Paragraph("Components", styles['SectionHead']))
    for item in [
        "<b>64 Bee Cards</b> \u2014 each with a unique combination of 6 attributes",
        "<b>6 Double-Sided Facet Tiles</b> \u2014 one per attribute, showing the Queen\u2019s current preference",
        "<b>1 Queen\u2019s Favor Tile</b> \u2014 determines first player and breaks ties",
    ]:
        story.append(Paragraph(f"\u2022  {item}", styles['BulletItem']))
    story.append(Spacer(1, 6))

    # Bee Attributes
    story.append(Paragraph("Bee Attributes", styles['SubsectionHead']))
    story.append(Paragraph(
        "Each bee card has 6 binary attributes. A bee has one trait or the other:",
        styles['Body']
    ))

    # Build attribute table with icons
    icon_dir = os.path.join(os.path.dirname(os.path.abspath(__file__)), "squib", "icons")
    icon_size = 22

    def icon_img(name):
        path = os.path.join(icon_dir, f"{name}.png")
        if os.path.exists(path):
            return Image(path, width=icon_size, height=icon_size)
        return ""

    attr_rows = [
        ["Texture",  icon_img("fuzzy"),     "Fuzzy",     icon_img("shiny"),     "Shiny"],
        ["Antennae", icon_img("feathered"), "Feathered", icon_img("whips"),     "Whips"],
        ["Weapon",   icon_img("stinger"),   "Stinger",   icon_img("mandibles"), "Mandibles"],
        ["Pattern",  icon_img("striped"),   "Striped",   icon_img("solid"),     "Solid"],
        ["Wings",    icon_img("sleek"),     "Sleek",     icon_img("flutter"),   "Flutter"],
        ["Payload",  icon_img("honey"),     "Honey",     icon_img("pollen"),    "Pollen"],
    ]

    icon_col = 0.055 * usable
    label_col = 0.20 * usable
    name_col = 0.245 * usable
    attr_col_widths = [label_col, icon_col, name_col, icon_col, name_col]

    attr_data = [["Attribute", "", "Side A", "", "Side B"]] + attr_rows
    attr_style_cmds = [
        ('BACKGROUND', (0, 0), (-1, 0), TABLE_HEADER_BG),
        ('TEXTCOLOR', (0, 0), (-1, 0), DARK_BROWN),
        ('FONTNAME', (0, 0), (-1, 0), 'Helvetica-Bold'),
        ('FONTSIZE', (0, 0), (-1, -1), 10),
        ('LEADING', (0, 0), (-1, -1), 14),
        ('ALIGN', (0, 0), (0, -1), 'LEFT'),
        ('ALIGN', (1, 0), (1, -1), 'CENTER'),
        ('ALIGN', (2, 0), (2, -1), 'LEFT'),
        ('ALIGN', (3, 0), (3, -1), 'CENTER'),
        ('ALIGN', (4, 0), (4, -1), 'LEFT'),
        ('VALIGN', (0, 0), (-1, -1), 'MIDDLE'),
        ('GRID', (0, 0), (-1, -1), 0.5, RULE_COLOR),
        ('TOPPADDING', (0, 0), (-1, -1), 3),
        ('BOTTOMPADDING', (0, 0), (-1, -1), 3),
        ('LEFTPADDING', (0, 0), (-1, -1), 6),
        ('RIGHTPADDING', (0, 0), (-1, -1), 4),
        # Span the header across icon+name columns
        ('SPAN', (1, 0), (2, 0)),
        ('SPAN', (3, 0), (4, 0)),
    ]
    for i in range(1, len(attr_data)):
        if i % 2 == 0:
            attr_style_cmds.append(('BACKGROUND', (0, i), (-1, i), TABLE_STRIPE))

    attr_table = Table(attr_data, colWidths=attr_col_widths, repeatRows=1)
    attr_table.setStyle(TableStyle(attr_style_cmds))
    story.append(attr_table)
    story.append(Spacer(1, 4))

    # Example card + explanation side by side
    card_path = os.path.join(os.path.dirname(os.path.abspath(__file__)), "example-card.png")
    if os.path.exists(card_path):
        card_img = Image(card_path, width=1.1 * inch, height=1.5 * inch)
        card_text = Paragraph(
            "Every possible combination of these 6 attributes appears on exactly "
            "one card, giving 64 unique bees. Each card shows the bee\u2019s name, "
            "illustration, and its six attribute icons along the left edge.",
            styles['Body']
        )
        text_col_w = usable - 1.3 * inch
        card_col_w = 1.3 * inch
        card_layout = Table(
            [[card_text, card_img]],
            colWidths=[text_col_w, card_col_w],
        )
        card_layout.setStyle(TableStyle([
            ('VALIGN', (0, 0), (-1, -1), 'MIDDLE'),
            ('LEFTPADDING', (0, 0), (-1, -1), 0),
            ('RIGHTPADDING', (0, 0), (-1, -1), 0),
            ('TOPPADDING', (0, 0), (-1, -1), 0),
            ('BOTTOMPADDING', (0, 0), (-1, -1), 0),
        ]))
        card_layout.hAlign = 'LEFT'
        story.append(card_layout)
    else:
        story.append(Paragraph(
            "Every possible combination of these 6 attributes appears on exactly "
            "one card, giving 64 unique bees.",
            styles['Body']
        ))

    # Facet Tiles
    story.append(Paragraph("Facet Tiles", styles['SubsectionHead']))
    story.append(Paragraph(
        "Each facet tile is double-sided, corresponding to one attribute \u2014 "
        "one trait per side. The face-up side shows which value the Queen "
        "currently desires for that attribute.",
        styles['Body']
    ))
    story.append(Paragraph(
        "During the game, facet tiles are arranged in a circle around the "
        "Queen\u2019s Favor tile. Position matters: the tile nearest the player "
        "the Queen\u2019s Favor points to is <b>Slot 1</b> (most important), "
        "and slots are numbered clockwise from there through "
        "<b>Slot 6</b> (least important).",
        styles['Body']
    ))

    # Layout diagram
    diagram_path = os.path.join(os.path.dirname(os.path.abspath(__file__)), "rules_layout_diagram.png")
    if os.path.exists(diagram_path):
        diagram_img = Image(diagram_path, width=3.2 * inch, height=3.7 * inch)
        diagram_table = Table([[diagram_img]], colWidths=[usable])
        diagram_table.setStyle(TableStyle([
            ('ALIGN', (0, 0), (-1, -1), 'CENTER'),
            ('TOPPADDING', (0, 0), (-1, -1), 4),
            ('BOTTOMPADDING', (0, 0), (-1, -1), 4),
        ]))
        story.append(diagram_table)

    story.append(hr())

    # --- Setup ---
    story.append(Paragraph("Setup", styles['SectionHead']))
    setup_steps = [
        "<b>Place the 6 facet tiles</b> in the center of the table.",
        "<b>Place the Queen\u2019s Favor</b> in the center of the table, pointing to a random player.",
        "<b>Shuffle all 64 bee cards</b> and deal <b>7 cards</b> to each player. "
        "Set remaining cards aside. Players look at their hands but keep them secret.",
        "<b>Draft the Queen\u2019s Favor</b> (see below).",
    ]
    for i, step in enumerate(setup_steps, 1):
        story.append(Paragraph(f"{i}.  {step}", styles['NumberedStep']))
    story.append(hr())

    # --- Overview ---
    story.append(Paragraph("Overview", styles['SectionHead']))
    story.append(Paragraph(
        "The game is played over a series of <b>hands</b>, each consisting of "
        "7 <b>rounds</b> (one per card in hand). Each round has three phases:",
        styles['Body']
    ))
    for phase in [
        "<b>Present</b> \u2014 All players simultaneously play a bee card face-down, then reveal.",
        "<b>Judge</b> \u2014 The played bees are compared against the facet tiles to determine a winner.",
        "<b>Manipulate</b> \u2014 Each player takes one action to alter the facet tiles.",
    ]:
        story.append(Paragraph(f"\u2022  {phase}", styles['BulletItem']))
    story.append(Spacer(1, 4))
    story.append(Paragraph(
        "After all 7 rounds, check for a winner. If no one has won yet, "
        "deal a new hand and continue.",
        styles['Body']
    ))
    story.append(hr())

    # --- Drafting ---
    story.append(Paragraph("Drafting the Queen\u2019s Favor", styles['SectionHead']))
    story.append(Paragraph(
        "At the start of each hand, players draft the 6 facet tiles into "
        "specific positions around the Queen\u2019s Favor tile. Drafting determines "
        "which attributes occupy each slot and which side is face-up.",
        styles['Body']
    ))
    story.append(Paragraph(
        "The player that the Queen\u2019s Favor tile is pointing towards will draft "
        "last. Proceeding counter-clockwise from them, players take turns "
        "choosing from among the remaining attribute tiles, choosing one side "
        "to be face-up, and adding them around the Queen\u2019s Favor tile in the "
        "slots near them.",
        styles['Body']
    ))

    story.append(Paragraph("Drafting by Player Count", styles['SubsectionHead']))
    story.append(Paragraph(
        "The number and position of tiles drafted depends on player count:",
        styles['Body']
    ))

    cell_style = ParagraphStyle('TableCell', parent=styles['Body'], spaceAfter=0, spaceBefore=0)
    draft_table = make_table(
        ["Players", "Drafting Order"],
        [
            ["2", Paragraph("Players alternate picking from Slot 6 down to Slot 1. First player drafts Slots 6, 4, 2. Other player drafts Slots 5, 3, 1.", cell_style)],
            ["3", Paragraph("First player drafts Slots 6, 5. Next player drafts Slots 4, 3. Last player drafts Slots 2, 1.", cell_style)],
            ["4", Paragraph("First player drafts Slots 6, 5. Next drafts Slots 4, 3. Next drafts Slot 2. Last drafts Slot 1.", cell_style)],
            ["5", Paragraph("First player drafts Slots 6, 5. Remaining players each draft 1 slot (4, 3, 2, 1) in counter-clockwise order.", cell_style)],
        ],
        col_widths=[usable * 0.12, usable * 0.88],
    )
    story.append(draft_table)
    story.append(hr())

    # --- Playing a Round ---
    story.append(Paragraph("Playing a Round", styles['SectionHead']))

    story.append(Paragraph("Phase 1: Present", styles['SubsectionHead']))
    story.append(Paragraph(
        "All players simultaneously choose one bee card from their hand and "
        "place it <b>face-down</b> in front of them. Once everyone has chosen, "
        "flip all cards face-up.",
        styles['Body']
    ))

    story.append(Paragraph("Phase 2: Judge", styles['SubsectionHead']))
    story.append(Paragraph(
        "Compare the played bees against the facet tiles, starting from Slot 1:",
        styles['Body']
    ))
    story.append(Paragraph(
        "1.  <b>Check Slot 1.</b> Does the played bee match the desired "
        "attribute value on this tile?",
        styles['NumberedStep']
    ))
    for sub in [
        "If <b>some bees match and others don\u2019t</b>: eliminate all non-matching bees.",
        "If <b>no bees match</b>: skip this slot (all bees survive).",
        "If <b>one bee remains</b>: that bee wins. Stop judging.",
    ]:
        story.append(Paragraph(f"\u2013  {sub}", styles['SubBullet']))
    story.append(Paragraph(
        "2.  <b>Repeat</b> for Slot 2, then Slot 3, and so on through Slot 6. "
        "This will always result in a single winning bee.",
        styles['NumberedStep']
    ))
    story.append(Spacer(1, 4))
    story.append(Paragraph(
        "The winning player scores <b>1 point</b> for the round. "
        "(You may collect the played bee cards into a score pile to track "
        "points, but each round win is worth exactly 1 point regardless "
        "of the number of players.)",
        styles['Body']
    ))

    story.append(Paragraph("Phase 3: Manipulate", styles['SubsectionHead']))
    story.append(Paragraph(
        "Starting with the player that won this round and proceeding clockwise, "
        "each player takes <b>one action</b>, either:",
        styles['Body']
    ))
    for action in [
        "<b>Flip</b> one facet tile to its opposite side, OR",
        "<b>Swap</b> the positions of any two facet tiles around the Queen\u2019s Favor.",
    ]:
        story.append(Paragraph(f"\u2022  {action}", styles['BulletItem']))
    story.append(Spacer(1, 4))
    story.append(Paragraph(
        "<b>Restriction:</b> You cannot repeat the exact action taken by the "
        "player immediately before you. (You may perform the same type of "
        "action on different tiles.)",
        styles['Body']
    ))
    story.append(Paragraph(
        "After all players have taken an action, the round is over. "
        "Begin the next round with Phase 1.",
        styles['Body']
    ))
    story.append(Paragraph(
        "<b>Note:</b> After judging the final round of a hand (round 7), "
        "skip the Manipulate phase \u2014 the facet tiles are about to be "
        "re-drafted anyway.",
        styles['Body']
    ))
    story.append(hr())

    # --- Winning ---
    story.append(Paragraph("Winning the Game", styles['SectionHead']))
    story.append(Paragraph(
        "After completing a hand (all 7 rounds), count each player\u2019s total "
        "score pile. <b>The first player to reach 10 points wins.</b>",
        styles['Body']
    ))
    for cond in [
        "If <b>one player</b> has 10 or more points and leads outright: that player wins!",
        "If <b>multiple players</b> are tied at 10 or more points: enter <b>Sudden Death</b>.",
        "If <b>no player</b> has 10 points yet: deal a new hand and continue.",
    ]:
        story.append(Paragraph(f"\u2022  {cond}", styles['BulletItem']))

    story.append(Paragraph("Sudden Death", styles['SubsectionHead']))
    story.append(Paragraph(
        "Shuffle all cards, deal a new hand of 7, and draft the Queen\u2019s Favor "
        "as normal. Play rounds one at a time. After each round\u2019s Judge phase, "
        "check: does any single player now lead outright? If so, that player "
        "wins immediately. Otherwise, proceed to Manipulation as normal. "
        "If players remain tied at the top after all rounds, deal another "
        "hand and continue.",
        styles['Body']
    ))

    story.append(Paragraph("Between Hands", styles['SubsectionHead']))
    story.append(Paragraph(
        "When a hand ends without a winner:",
        styles['Body']
    ))
    between_steps = [
        "The player with the most points starts with the Queen\u2019s Favor (rotate the tile to point to them).",
        "Shuffle all 64 cards together, deal 7 to each player, and set the rest aside.",
        "Draft the attribute tiles around the Queen\u2019s Favor again.",
    ]
    for i, step in enumerate(between_steps, 1):
        story.append(Paragraph(f"{i}.  {step}", styles['NumberedStep']))
    story.append(hr())

    # --- Quick Reference ---
    story.append(Paragraph("Quick Reference", styles['SectionHead']))
    qr_items = [
        ("<b>Round structure:</b>  Present (simultaneous) \u2192 "
         "Judge (Slot 1 through 6) \u2192 Manipulate (flip or swap)"),
        "<b>Judging priority:</b>  Slot 1 > Slot 2 > Slot 3 > Slot 4 > Slot 5 > Slot 6",
        "<b>Manipulation actions:</b>  Flip 1 tile OR swap 2 tiles. "
        "Cannot repeat the previous player\u2019s exact action.",
        "<b>Win condition:</b>  First to 10 points with a clear lead.",
    ]
    for item in qr_items:
        story.append(Paragraph(item, styles['QuickRef']))
    story.append(hr())

    # --- Strategy Tips ---
    story.append(Paragraph("Strategy Tips", styles['SectionHead']))
    tips = [
        "<b>During the draft:</b> Place your strongest attributes in the slots "
        "you control. Remember, Slot 1 dominates \u2014 if your best card matches "
        "Slot 1, it beats cards that match Slots 2\u20136 but miss Slot 1.",

        "<b>When presenting:</b> Play the card that survives the most filters. "
        "A card matching Slots 1 and 2 will beat a card matching Slots 3, 4, 5, and 6.",

        "<b>Manipulation is key:</b> A well-timed flip or swap before the next "
        "round can transform a weak hand into a winning one. Think about which "
        "changes help your remaining cards while hurting opponents.",

        "<b>Watch the endgame:</b> As players approach 10 points, sudden death "
        "tension rises. Controlling Slot 1 through manipulation becomes critical.",

        "<b>Drafting trade-offs:</b> By drafting first it may feel like you\u2019re "
        "just choosing tiebreakers, but you\u2019re choosing which tiles <b>will not</b> "
        "be eligible to be in the first few spots.",
    ]
    for tip in tips:
        story.append(Paragraph(f"\u2022  {tip}", styles['StrategyTip']))

    doc.build(story)
    print(f"PDF generated: {output_path}")


if __name__ == "__main__":
    script_dir = os.path.dirname(os.path.abspath(__file__))
    build_pdf(os.path.join(script_dir, "EyeOfTheBeeholderRules.pdf"))
