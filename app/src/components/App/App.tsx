
import { CssBaseline, ThemeProvider, useMediaQuery } from '@material-ui/core';
import React from 'react';
import { Route, Switch} from 'react-router-dom';
import { darkTheme, lightTheme } from '../../theme';
import Help from '../Help';
import MainTerminal from '../MainTerminal';
import './App.css';

function App() {
  const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)')

  const theme = prefersDarkMode ? darkTheme : lightTheme

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Switch>
        <Route path='/help' component={Help} />
        <Route path='/' component={MainTerminal} />
      </Switch>
    </ThemeProvider>
  );
}

export default App;
