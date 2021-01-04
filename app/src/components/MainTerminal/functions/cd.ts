import { apiGetExactPath } from '../../../utils/api'
import { getPath } from '../../../utils/path'

export async function cd(workDir: string, callback, ...args) {
  if (args.length > 1 || args.length == 0) {
    return 'Please provide only 1 arg to cd'
  }

  if (args[0] == '') {
    return 'Please provide a valid path!'
  }

  const intendedPath = getPath(workDir, args[0])
  const data = await apiGetExactPath(intendedPath)

  if (data) {
    if (data.filedirs.filetype == 'directory') {
      callback(intendedPath)
      return
    }
    return "Destination is a file"
  } 

  return "Destination doesn't exist"

}
