import {
  apiDeletePath,
  apiGetPath,
} from '../../../utils/api'
import { getPath, ltreeToPath } from '../../../utils/path'

export async function rm(workDir, ...args) {
  let retString = ""
  const deletePaths: string[] = []
  if (args.length != 1) {
    return 'Please provide only 1 arg to rm'
  }

  if (args[0] == '') {
    return 'Please provide a valid path!'
  }

  const intendedPath = getPath(workDir, args[0])
  const splitPath = intendedPath.split('/')
  const name = splitPath[splitPath.length - 1]

  // Selecting everything in the directory if the file name is *
  if (name == '*' || name == '.*') {
    const parentPath = splitPath.slice(0, -1).join('/') || '/'
    const dirData = await apiGetPath(parentPath, 1)
    if (dirData.length > 0) {
      dirData.forEach((data) => {
        const fileDir = data.filedirs
        const childPath = ltreeToPath(fileDir.path)
        deletePaths.push(childPath)
      })
    }
  } else {
    deletePaths.push(intendedPath)
  }

  for (let i = 0; i < deletePaths.length; i++) {
    const path = deletePaths[i]
    const data = await apiDeletePath(path)
    if (data.numrows == 0) {
      retString += `Could not delete ${path}\n`
    } else {
      retString += `Deleted ${path}\n`
    }
  }

  return retString
}
