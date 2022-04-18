import { Stack, Button } from "@mui/material"
import React from "react"
import { Link } from "react-router-dom"

function Welcome() {
    return (
        <Stack justifyContent="center" flexDirection="column" alignItems="center" minHeight="60vh">
            <img src="logo.svg" alt="welcome" width="40%" />
            <br />
            <Stack direction="row" spacing={2}>
                <Button component={Link} variant="contained" to="start">
                    Start Game
                </Button>
                <Button component={Link} variant="contained" to="play">
                    Join Matchmaking
                </Button>
            </Stack>
        </Stack>
    )
}

export default Welcome