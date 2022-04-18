import { Guess } from './gameSlice'

const highlight = (word: string, result: number[]) => {
    return (<h2 style={{
        fontFamily: "Andale Mono,AndaleMono,monospace"
    }}>{word.split("").map((letter: string, index: number) => {
        let bgColor = "none"

        if (result[index] === 2) {
            bgColor = "lime"
        } else if (result[index] === 1) {
            bgColor = "orange"
        }
        return <span key={index} style={{ backgroundColor: bgColor, borderRadius: "4px", padding: "2px", margin: "2px" }}>{letter}</span>
    })
        }</h2>)
}

function GuessCard({ word, result, player }: Guess) {
    return (
        <div>
            {highlight(word, result)}
        </div>
    )
}

export default GuessCard
