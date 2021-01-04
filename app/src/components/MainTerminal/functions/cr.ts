import { apiCreatePath, CreatePathPayload } from '../../../utils/api'
import { getParentPaths, getPath } from '../../../utils/path'

export async function cr(workDir: string, ...args) {
  if (args.length < 1) {
    return 'Usage: cr [-p] PATH [DATA]'
  }
  const regex = RegExp(/^[a-zA-Z0-9 _-]+$/)
  let data = ''
  let parents = false

  const index = args.indexOf('-p')
  if (index >= 0) {
    parents = true
    args.splice(index, 1)
  }

  // must be path and data respectively, try to convert
  if (args.length < 1) {
    return 'Usage: cr [-p] PATH [DATA]'
  }

  const path = args[0]
  args.splice(0, 1)
  if (path == '') {
    return 'PATH must not be empty'
  }

  if (args.length > 0) {
    data = args.join(' ')
  }
  const intendedPath = getPath(workDir, path)
  const splitPath = intendedPath.split('/')
  const name = splitPath[splitPath.length - 1]
  if (!regex.test(name)) {
    return 'Name must be valid'
  }

  if (parents) {
    const parentPaths = getParentPaths(intendedPath)
    for (let i = 0; i < parentPaths.length; i++) {
      const parentPath = parentPaths[i]
      const parentSplitPath = parentPath.split('/')
      const parentName = parentSplitPath[parentSplitPath.length - 1]
      if (!regex.test(parentName)) {
        return 'Parent name must be valid'
      }
      const parentPayload: CreatePathPayload = {
        filedirs: {
          name: parentName,
          data: '',
          filetype: 'directory',
        },
        path: splitPath.slice(0, -1).join('/') || '/',
      }
      await apiCreatePath(parentPayload)
    }
  }

  const childPayload: CreatePathPayload = {
    filedirs: {
      name: name,
      data: data,
      filetype: data ? 'file' : 'directory',
    },
    path: splitPath.slice(0, -1).join('/') || '/',
  }
  const apiData = await apiCreatePath(childPayload)

  if (apiData) {
    return `${intendedPath} successfully created `
  }

  return `Could not create ${intendedPath}`
}
