const StealPrompt = (time: number) => {
    if(time > 0) {
    return (
        <h3>Steal in {(time)/1000} seconds...</h3>
    )
    }
    return (<h3>Steal!</h3>)
}

export default StealPrompt
