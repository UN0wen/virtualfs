import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  createStyles,
  IconButton,
  makeStyles,
  Theme,
  Tooltip,
  Typography,
  Zoom,
} from '@material-ui/core'
import { ExpandMore, Home } from '@material-ui/icons'
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
      margin: theme.spacing(3),
      fontWeight: 400,
    },
    root: {
      width: '100%',
      padding: theme.spacing(3),
    },
    heading: {
      fontSize: theme.typography.pxToRem(15),
      fontWeight: theme.typography.fontWeightRegular,
    },
    available: {
      marginBottom: theme.spacing(2),
    },
    details: {
      flexDirection: 'column',
    },
    example: {
      marginTop: theme.spacing(1),
    },
  })
)

export default function Help() {
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
      <Typography variant="h3" className={classes.title}>
        Documentation
      </Typography>
      <div className={classes.root}>
        <Typography variant="h5" className={classes.available}>
          Available functions
        </Typography>
        <Accordion>
          <AccordionSummary expandIcon={<ExpandMore />} id="cd">
            <Typography className={classes.heading}>
              cd - Change Directory
            </Typography>
          </AccordionSummary>
          <AccordionDetails className={classes.details}>
            <Typography>
              <b>Usage</b>: <code>cd DIR</code>
            </Typography>
            <Typography>
              Change the current working directroy to DIR.
            </Typography>
            <Typography className={classes.example}>
              <b>Example</b>: <code>cd /var/lib</code>
            </Typography>
          </AccordionDetails>
        </Accordion>
        <Accordion>
          <AccordionSummary expandIcon={<ExpandMore />} id="cat">
            <Typography className={classes.heading}>
              cat - Concatenate
            </Typography>
          </AccordionSummary>
          <AccordionDetails className={classes.details}>
            <Typography>
              <b>Usage</b>: <code>cat FILE</code>
            </Typography>
            <Typography>
              Prints the contents of a file to the terminal
            </Typography>
            <Typography className={classes.example}>
              <b>Example</b>: <code>cat /var/lib/file</code>
            </Typography>
          </AccordionDetails>
        </Accordion>
        <Accordion>
          <AccordionSummary expandIcon={<ExpandMore />} id="ls">
            <Typography className={classes.heading}>ls - List</Typography>
          </AccordionSummary>
          <AccordionDetails className={classes.details}>
            <Typography>
              <b>Usage</b>: <code>ls [-l] [LIMIT]</code>
            </Typography>
            <Typography>
              List all files under the current working directory.
            </Typography>
            <Typography>
              If the -l flag is provided, the file's full details are listed as
              well. If a LIMIT is specified, files up to LIMIT levels deep will
              be shown.
            </Typography>
            <Typography className={classes.example}>
              <b>Example</b>: <code>ls -l 4</code>
            </Typography>
            <Typography>
              This command lists all files and their details 4 levels deep from
              current working directory, showing for example
              /var/lib/path/to/file from /var
            </Typography>
          </AccordionDetails>
        </Accordion>
        <Accordion>
          <AccordionSummary expandIcon={<ExpandMore />} id="cr">
            <Typography className={classes.heading}>cr - Create</Typography>
          </AccordionSummary>
          <AccordionDetails className={classes.details}>
            <Typography>
              <b>Usage</b>: <code>cr [-p] PATH [DATA]</code>
            </Typography>
            <Typography>
              cr creates a new file or directory at the specified PATH.
            </Typography>
            <Typography>
              If DATA is provided, the created entity is a file with text data
              taken from DATA. If it is not, a directory is created instead.
            </Typography>
            <Typography>
              If the -p flag is provided, all of the parent directories are
              created if they don't already exist. If not, the command will
              raise an error if the parent directory does not exist.
            </Typography>
            <Typography className={classes.example}>
              <b>Example</b>:{' '}
              <code>cr -p /var/lib/path/to/test_file this is a file</code>
            </Typography>
            <Typography>
              This command creates the file with name <code>test_file</code>,
              putting in it the text data <code>this is a file</code>. Since the
              -p option is specified, the parent directories (/var, /var/lib,
              /var/lib/path, /var/lib/path/to) will be created if they don't
              already exists.
            </Typography>
          </AccordionDetails>
        </Accordion>
        <Accordion>
          <AccordionSummary expandIcon={<ExpandMore />} id="mv">
            <Typography className={classes.heading}>mv - Move</Typography>
          </AccordionSummary>
          <AccordionDetails className={classes.details}>
            <Typography>
              <b>Usage</b>: <code>mv SOURCE... DEST</code>
            </Typography>
            <Typography>
              mv moves all files and directories (recursively) specified by
              SOURCE... to DEST.
            </Typography>
            <Typography>
              If more than 2 path names are provided, the last path is
              considered the destination DEST, and all previous paths are
              considered source file/folders SOURCE.
            </Typography>
            <Typography>
              If there is a conflict at the DEST (for example: moving /var to
              /lib where /lib/var already exists), the SOURCE will be moved
              under the existing folder if it is a folder, overwrite the file if
              they are both files, and throw an error otherwise.
            </Typography>
            <Typography>
              This command accepts the * operator, selecting all files in a
              folder. By default, if the SOURCE is a directory, itself and all
              of its children will be moved under DEST.
            </Typography>
            <Typography className={classes.example}>
              <b>Example</b>: <code>mv /var/* /bin /lib</code>
            </Typography>
            <Typography>
              This command moves every file and directory under /var, as well as
              the file/directory /bin, to under /lib. Assuming that /var
              contains /var/file0 and /var/file1, and /bin contains /bin/file0
              and /bin/file1, the resulting directory structure in /lib is
              /lib/file0, /lib/file1, /lib/bin/file0, /lib/bin/file1.
            </Typography>
          </AccordionDetails>
        </Accordion>
        <Accordion>
          <AccordionSummary expandIcon={<ExpandMore />} id="rm">
            <Typography className={classes.heading}>rm - Remove</Typography>
          </AccordionSummary>
          <AccordionDetails className={classes.details}>
            <Typography>
              <b>Usage</b>: <code>rm PATH</code>
            </Typography>
            <Typography>
              rm deletes the specified file or directory (recursively) at PATH.
            </Typography>
            <Typography>
              This command accepts the * operator, deleting all files and
              directories in a folder. By default, if the PATH is a directory,
              itself and all of its children will be deleted.
            </Typography>
            <Typography className={classes.example}>
              <b>Example</b>:{' '}
              <code>rm /lib/*</code>
            </Typography>
            <Typography>
              This command deletes everything under the folder /lib, leaving /lib an empty folder.
            </Typography>
          </AccordionDetails>
        </Accordion>
      </div>
    </div>
  )
}
