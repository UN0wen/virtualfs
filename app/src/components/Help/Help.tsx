import {
  createStyles,
  IconButton,
  makeStyles,
  Theme,
  Tooltip,
  Typography,
  Zoom,
} from '@material-ui/core'
import { Home } from '@material-ui/icons'
import React from 'react'
import { useHistory } from 'react-router-dom'

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    returnButton: {
      position: 'fixed',
      bottom: 0,
      left: 0,
      margin: theme.spacing(1),
    },
    largeIcon: {
      fontSize: '2em',
    },
    title: {
      display: 'flex',
      justifyContent: 'center',
      marginTop: theme.spacing(1),
      fontWeight: 400,
    }
  })
)

export default function MainTerminal() {
  const history = useHistory()
  const classes = useStyles()
  return (
    <div>
      <Tooltip
        title="Return to terminal"
        TransitionComponent={Zoom}
        placement="right"
      >
        <IconButton
          className={classes.returnButton}
          onClick={() => {
            history.push('/')
          }}
        >
          <Home className={classes.largeIcon} />
        </IconButton>
      </Tooltip>
      <Typography variant="h3" className={classes.title} >Documentation</Typography>
    </div>
  )
}
