import { isAbsolute } from 'path'
import base64url from 'base64url'

function removeTrailingSlash(str: string) {
  if (str == '/') return str
  return str.endsWith('/') ? str.slice(0, -1) : str
}

function absolute(base, relative) {
  const stack = base.split("/"),
      parts = relative.split("/");
  for (let i=0; i<parts.length; i++) {
      if (parts[i] == ".")
          continue;
      if (parts[i] == "..")
          stack.pop();
      else
          stack.push(parts[i]);
  }
  return stack.join("/");
}

export function getPath(workDir: string, path: string) {
  if (isAbsolute(path)) return path
  // for compatibility with the absolute func
  let newWorkDir = workDir
  path = removeTrailingSlash(path)
  if (workDir == '/') {
    if (path == '.' || path == '..') {
      return workDir
    } 
    newWorkDir = ''
  }
  const abs = absolute(newWorkDir, path)
  return removeTrailingSlash(abs ? abs : "/")
}

export function ltreeToPath(ltree: string) {
  if (ltree == 'root') return '/'

  const ltreeNoRoot = ltree.replace('root.', '/')
  return ltreeNoRoot.replaceAll('.', '/')
}


export function countDepth(path: string) {
  return (path.match(/\//g) || []).length
}

function padRight(str, n, pad) {
  const length = str.length
  if (n > length) for (let i = 0; i < n - length; i++) str += pad
  return str
}

export function encodePath(path: string) {
  const encoded = base64url.encode(path)
  return padRight(encoded, encoded.length + (4 - encoded.length % 4) % 4, '=')
}

export function getParentPaths(path: string) {
  const paths: string[] = []
  for(let i = 1 ; i < path.length ; i++) {
    if(path[i]==='/') paths.push(path.substr(0,i));
  }

  return paths
}