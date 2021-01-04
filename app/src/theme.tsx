import { blue, deepPurple, lightBlue, purple, red } from '@material-ui/core/colors'
import { createMuiTheme } from '@material-ui/core/styles'

// A custom theme for this app
export const lightTheme = createMuiTheme({
  palette: {
    type: 'light',
    primary: {
      main: lightBlue[500],
    },
    secondary: {
      main: deepPurple[500],
    },
    error: {
      main: red.A400,
    },
  },
})

export const darkTheme = createMuiTheme({
  palette: {
    type: 'dark',
    primary: {
      main: blue[200],
    },
    secondary: {
      main: purple[200],
    },
    error: {
      main: red[500],
    },
  },
})
