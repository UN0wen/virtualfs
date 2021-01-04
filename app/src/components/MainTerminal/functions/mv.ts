import {
  apiGetPath,
  apiUpdatePath,
  UpdatePathPayload,
} from '../../../utils/api'
import { getPath, ltreeToPath } from '../../../utils/path'

export async function mv(workDir, ...args) {
  let retString = ''
  if (args.length < 2) {
    return 'Please provide at least 2 args to mv'
  }
  const sources: string[] = []
  const dest = getPath(workDir, args.pop())

  for (let i = 0; i < args.length; i++) {
    const intendedPath = getPath(workDir, args[i])
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
          sources.push(childPath)
        })
      }
    } else {
      sources.push(intendedPath)
    }
  }

  for (let i = 0; i < sources.length; i++) {
    const source = sources[i]
    const item: UpdatePathPayload = {
      source: source,
      dest: dest,
    }
    const data = await apiUpdatePath(item)
    if (data.length == 0) {
      retString += `Could not move ${source} to ${dest}\n`
    } else {
      retString += `Moved ${source} to ${dest}\n`
    }
  }

  return retString ? retString : "Sources can't be empty"
}
