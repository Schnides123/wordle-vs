import { Card, CardContent, Typography } from '@mui/material'

export interface PlayerCardProps {
    player: string
    isTurn: boolean
}
function PlayerCard(props: PlayerCardProps) {
    return (
        <Card style={{ backgroundColor: props.isTurn ? "rgb(25,205,38)" : "rgb(205,25,38)" }}>
            <CardContent>
                <Typography variant="h5" component="h2">
                    {props.player}
                </Typography>
                <Typography variant="body2" component="p">
                    Player 1's card
                </Typography>
            </CardContent>
        </Card>
    )
}

export default PlayerCard
