# Marathon event maker

## Memo maker
This is the tools used to make the memo that appear on the board.

### Parameters
- Memo background:
  - The file containing the background of the memo (no text).
  - Default to the common white notes that are used in all except for one or 2 events.
- Memo Pin:
  - The pin used to attach the memo to the board.
  - Default to the pink pin.
  - Commonly used pins like blue pin are also available.
- Pin X, Pin Y:
  - The location to render the pin assets at.
- Memo text:
  - The text that appears on the memo.
  - Use `\` to break between lines.
- Text X, Text Y, Rotation:
  - The position and rotation to draw the text at.
  - Default is -3.8, which is the same as the background white note.
- Font:
  - The font used for the text.
  - Default to FOT-SkipStd-B.otf, which is the font used for JP and EN.
  - For ZH, use DFT_YY6.ttf
  - For KR, use FOTK-YOONGOTHIC760.otf
  - For other fonts, put the font in ``./gui/fonts`` to use them.
- Character Size, Letter Spacing, Line Spacing:
  - These configure the text rendered by the fonts.
- Color:
  - The color of the text, RGBA written in hex.
  - Defaut t0 494949ff, which is the same color as the white background notes.
- Text scaling factor, background scaling on/off:
  - How much to scale the text DOWN, and whether to scale the background UP to match the text.
  - This is done because the original assets would have been made at a higher resolution in the studio.
  - If the assets were made at 800x800 then downscaled, then we can't really recreate the character's shape at 200x200.
  - So we will render text at a higher resolution then downscale them.
  - The background scale is important because scaling background along with the text will blend the background and the text together.
  - The scale factor default to 4, and the background scaling default to on.
- Line offsets:
  - Each of the text lines are "centered", howerver there might be some small differences.
  - So these offsets are used to correct whatever differences there are.
  - The lower the more left the text is (can go to negative).
  - They all default to 0.
- Output scaling factor:
  - This is only for showing the notes.
  - Default to 1.2, which is roughly the same resolution as the note take up in game, assuming a resolution of 1800x900.
- Memo directory:
  - The folder/directory to export / import the result.
  - Exporting will export configs and generated memos (see the exporting part).
  - Importing will import the exported config from the same directory.
### Exporting

After making the note, use the Export buttons to export the result. This will generate a few files in the memo directory:

- `memo_maker_config.json`:
  - The current config, can be used to restore the current state of the maker.
- `output.png`:
  - The output image, using internal scaling algorithms.
- `output_gimp.png`:
  - The output image, using gimp for scaling.
  - This is only generated if the correct command line was setup.
  - Generally, this is "more correct" with the default settings.
- Exporting can be quiet heavy, especially with gimp export, so expect a few seconds delay until it is done.

#### Setting up gimp export
- First install [GIMP](https://www.gimp.org/).
- Then copy `./gui/app/marathon_event_maker/memo-scale.scm` into gimp's scripts folder. This can be in `C:\Program Files\GIMP 2\share\gimp\2.0\scripts` or some other places, depending on the gimp version.
- Then edit the `GIMP_COMMAND_FORMAT` in `./gui/app/marathon_event_maker/memo_maker/memo_maker.go`:
  - By default this is `"C:\Program Files\GIMP 2\bin\gimp-console-2.10" -i -b "(memo-scale \"%s\")" -b "(gimp-quit 0)\`
  - Change `gimp-console-2.10` to the appropriate version. 
- Rebuild
 