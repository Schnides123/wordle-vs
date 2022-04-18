import React from "react"

export interface WordCardProps extends React.ComponentProps<'div'> {
	word: string
}

function WordCard({ word, ...rest }: WordCardProps) {
	return (
		<div className="word" {...rest}>
			<span style={{ fontSize: 45 }}>{word}</span>
		</div>
	)
}

export default WordCard
