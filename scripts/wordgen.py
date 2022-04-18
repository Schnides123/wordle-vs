scrabbleWords = set()
with open("scrabble.txt") as scrabble:
    for line in scrabble:
        scrabbleWords.add(line.strip().lower())

print(len(scrabbleWords))

with open("unigram_freq.csv") as csv:
    wordsraw = [line.strip().split(",") for line in csv]
    wordsSorted = sorted(wordsraw, key=lambda x: int(x[1]), reverse=True)
    words = [word[0] for word in wordsSorted]
    fiveletter = [word for word in words if len(word) == 5 and word in scrabbleWords][:2000]
    sixletter = [word for word in words if len(word) == 6 and word in scrabbleWords][:2000]
    sevenletter = [word for word in words if len(word) == 7 and word in scrabbleWords][:2000]
    eightletter = [word for word in words if len(word) == 8 and word in scrabbleWords][:2000]
    f = open("fivewords.txt", "w")
    for word in fiveletter:
        f.write(word + "\n")
    f.close()
    f = open("sixwords.txt", "w")
    for word in sixletter:
        f.write(word + "\n")
    f.close()
    f = open("sevenwords.txt", "w")
    for word in sevenletter:
        f.write(word + "\n")
    f.close()
    f = open("eightwords.txt", "w")
    for word in eightletter:
        f.write(word + "\n")
    f.close()

csv.close()