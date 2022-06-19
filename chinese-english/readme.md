For Chinese-English bidirectional translation, I made a few decisions:

The following notes reference _disk usage_, not memory usage:
- When traditional and simplified characters are the same, just save as a trigger (minor optimisation: 47MB -> 44MB, no negative effect to user experience)
- Don't do English->Chinese translations when the English definitions have multiple words (major optimisation: 47MB -> 25MB, mostly removing useless triggers, minor effect to user experience)
  - Allowing English->Chinese translations when the English definition has _at most 2 words_ has 47MB->30MB. This might be a fair trade-off to get better translations.
- Only doing simplified OR traditional character translations halves file size. All-in-one gives the best user experience, but 
- Filtering to only translate words in HSK ... // TODO
- Filtering to only top 2k/5k/10k dictionary lines: full-fat 47MB dictionary ->  1MB, 2MB, 4MB. I tested this just by limiting line count. **This looks like a great optimisation, if we can find Chinese word frequency data**
  - MIT licensed, csv, single characters only - [hanziDB.csv](https://github.com/ruddfawcett/hanziDB.csv/blob/master/data/hanziDB.csv)
  - ? license, tsv? - [namedict](https://github.com/thyrlian/namedict/blob/master/data/Modern%20Chinese%20Character%20Frequency%20List)
  - ? license, tsv, some errors but OK for our use - [pleco forum](https://www.plecoforums.com/threads/word-frequency-list-based-on-a-15-billion-character-corpus-bcc-blcu-chinese-corpus.5859/)
  
Also note that the file `cedict_ts.u8.txt` is licensed under creative commons - https://www.mdbg.net/chinese/dictionary?page=cc-cedict