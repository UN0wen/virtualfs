import { apiGetPath } from "../../../utils/api"
import { countDepth,  ltreeToPath } from "../../../utils/path"
import { formatCreated, formatUpdated } from "../../../utils/time"

export async function ls(workDir: string, ...args) {
  if (args.length > 2) {
    return "Too many arguments!"
  }
  let full = false
  let level = 1

  const index = args.indexOf('-l')
  if (index >= 0) {
    full = true
    args.splice(index, 1)
  }

  
  // must be limit, try to convert 
  if (args.length == 1) {
    level = Number(args[0])
    if (!level) {
      return "Level must be a number"
    }
  } else if (args.length > 1) { //
      return "Usage: ls [-l] [LIMIT]"
  }
  
  const data = await apiGetPath(workDir, level)

  const parentDepth = countDepth(workDir)
  const strings: string[] = []

  strings.push(`${workDir}`)
  if (Array.isArray(data)) {
    data.forEach((r) => {
      strings.push(generateString(r.filedirs, full, parentDepth))
    })
  }

  return `${strings.join('\n')}`
}

function generateString(fileDir: any, full: boolean, baseDepth: number) {
  const path = ltreeToPath(fileDir.path)
  const relativeDepth = countDepth(path) - baseDepth
  let prepend = ""
  for (let i = 0 ; i < relativeDepth; i++) {
    prepend += '||||'
  }

  prepend += '---'

  if (full) {
    return `${prepend} ${path} | ${fileDir.filetype} | Size: ${fileDir.size} | ${formatCreated(fileDir.created)} | ${formatUpdated(fileDir.updated)} `
  } else {
    return `${prepend} ${path} | ${fileDir.filetype}`
  }
}