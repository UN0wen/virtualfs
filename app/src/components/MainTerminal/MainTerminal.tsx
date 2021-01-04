import {
  createStyles,
  IconButton,
  makeStyles,
  Theme,
  Tooltip,
  Zoom,
} from '@material-ui/core'
import { Help } from '@material-ui/icons'
import React, { useState } from 'react'
import Terminal from 'react-console-emulator'
import { useHistory } from 'react-router-dom'
import { cd } from './functions/cd'
import { cr } from './functions/cr'
import { ls } from './functions/ls'
import { mv } from './functions/mv'
import { rm } from './functions/rm'

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    helpButton: {
      position: 'fixed',
      bottom: 0,
      right: 0,
    },
    largeIcon: {
      fontSize: '2em',
    },
  })
)

export default function MainTerminal() {
  const classes = useStyles()
  const history = useHistory()
  const [workDir, setWorkDir] = useState('/')

  const commands = {
    echo: {
      description: 'Echo a passed string.',
      usage: 'echo <string>',
      fn: function (...args) {
        return `${Array.from(args).join(' ')}`
      },
    },
    cd: {
      description: 'Change current working directory.',
      usage: 'cd <location>',
      fn: (...args) => cd(workDir, setWorkDir, ...args),
    },
    ls: {
      description:
        "Show what's in the current and below directories with LIMIT",
      usage: 'ls [-l] [LIMIT]',
      fn: (...args) => ls(workDir, ...args),
    },
    cr: {
      description: 'Create new files and directories',
      usage: 'cr [-p] PATH [DATA]',
      fn: (...args) => cr(workDir, ...args),
    },
    mv: {
      description: 'Move all source files from source to dest',
      usage: 'mv SOURCE... DEST',
      fn: (...args) => mv(workDir, ...args),
    },
    rm: {
      description: 'Remove a file/folder and all its children',
      usage: 'rm PATH',
      fn: (...args) => rm(workDir, ...args),
    },
  }

  return (
    <div>
      <Terminal
        commands={commands}
        welcomeMessage={'Type help or click the question mark to get started'}
        promptLabel={`me@React:${workDir}$`}
        style={{ minHeight: '100vh' }}
      />
      <Tooltip title="Help page" TransitionComponent={Zoom} placement="left">
        <IconButton
          className={classes.helpButton}
          onClick={() => {
            history.push('/help')
          }}
        >
          <Help className={classes.largeIcon} />
        </IconButton>
      </Tooltip>
    </div>
  )
}
