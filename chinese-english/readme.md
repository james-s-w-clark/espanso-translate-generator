# Intro
For now, this package can only do English -> Chinese. 
Chinese -> English not working is a [known issue](https://github.com/federico-terzi/espanso/issues/101), so for now we don't generate that.

# Decision log

For Chinese-English bidirectional translation, I made a few **decisions**:

The following notes reference _disk usage_, not memory usage:
- **When traditional and simplified characters are the same, just save as a single trigger** (minor optimisation: 47MB -> 44MB, no negative effect to user experience)
- Don't do English->Chinese translations when the English definitions have multiple words (major optimisation: 47MB -> 25MB, mostly removing useless triggers, minor effect to user experience)
  - **Allowing English->Chinese translations when the English definition has _at most 2 words_ has 47MB->30MB. This might be a fair trade-off to get better translations.**
- Only doing simplified OR traditional character translations halves file size. All-in-one gives the best user experience, however.
- Filtering to only translate words in HSK ... // TODO? - HSK isn't as accurate as natural word frequency, see below
- **Filtering to only top 2k/5k/10k most frequent words: full-fat 47MB dictionary ->  1MB, 2MB, 4MB respectively**
  - MIT licensed, csv, single characters only - [hanziDB.csv](https://github.com/ruddfawcett/hanziDB.csv/blob/master/data/hanziDB.csv)
  - ? license, tsv? - [namedict](https://github.com/thyrlian/namedict/blob/master/data/Modern%20Chinese%20Character%20Frequency%20List)
  - ? license, tsv, some errors but OK for our use - [pleco forum](https://www.plecoforums.com/threads/word-frequency-list-based-on-a-15-billion-character-corpus-bcc-blcu-chinese-corpus.5859/)
- **If the characters for simplified and traditional are the same, do a single English->Chinese translation (:zh, not :zhs or :zht)**. If people don't like this, can revert
  
Also note that the file `cedict_ts.u8.txt` is licensed under creative commons - https://www.mdbg.net/chinese/dictionary?page=cc-cedict
