import { apiGetExactPath } from '../../../utils/api'
import { getPath } from '../../../utils/path'

export async function cat(workDir: string, ...args) {
  if (args.length > 1 || args.length == 0) {
    return 'Please provide only 1 arg to cd'
  }

  if (args[0] == '') {
    return 'Please provide a valid path!'
  }

  const intendedPath = getPath(workDir, args[0])
  const data = await apiGetExactPath(intendedPath)

  if (data) {
    if (data.filedirs.filetype == 'file') {
      return data.filedirs.data
    }
    return "Provided path is a directory"
  } 

  return "Provided path doesn't exist"

}
