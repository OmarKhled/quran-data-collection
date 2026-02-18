### QPC Glypgs

The QPC Script utilizes a unique glyph-based font system to ensure that every word in the Quran is rendered exactly as it appears in the physical Mushaf. Unlike standard fonts, these scripts use specific character codes to call upon precise calligraphic shapes that maintain the integrity of the Uthmani script.

This repository processes the QPC script to generate a structured dataset. Each Ayah is mapped to its specific glyph text and its corresponding page number. Due to the extreme detail and number of unique characters, the QPC font is too large to load as a single file. To optimize performance and loading times, the font is sharded (split) by page. The script contained in this repo only maps each ayah to its page, so we know whic font shard to load.
