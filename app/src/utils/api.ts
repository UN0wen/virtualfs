import { AxiosResponse } from "axios"
import { AxiosInstance } from "./axios"
import { encodePath } from "./path"

export async function apiGetPath(path:string , level:number) {
  const params = {
    level: level
  }
  const encodedPath = encodePath(path)

  let response: AxiosResponse<any>
  try {
      response = await AxiosInstance.get(`/${encodedPath}`, {params})
  }
  catch (err) {
    console.log(err.response.data.error)
    return []
  }

  return response.data ? response.data : []
}

export async function apiGetExactPath(path:string) {
  const encodedPath = encodePath(path)
  let response: AxiosResponse<any>
  try {
    response = await AxiosInstance.get(`/${encodedPath}/exact`)
  }
  catch (err) {
    console.log(err.response.data.error)
    return 0 
  }

  return response.data ? response.data : 0
}

export type CreatePathPayload = {
  filedirs: {
    name: string,
    data: string,
    filetype: string,
  },
  path: string
}

export async function apiCreatePath(item: CreatePathPayload,) {
  const requestOptions = {
    headers: { 'Content-Type': 'application/json' },
    body: item
  }
  let response: AxiosResponse<any>
  try {
    response = await AxiosInstance.post(`/`, item, requestOptions)
  }
  catch (err) {
    console.log(err.response.data.error)
    return 0 
  }

  return response.data ? response.data : 0
}

export type UpdatePathPayload = {
  source:string,
  dest:string
}

export async function apiUpdatePath(item: UpdatePathPayload) {
  const requestOptions = {
    headers: { 'Content-Type': 'application/json' },
    body: item
  }
  let response: AxiosResponse<any>
  try {
    response = await AxiosInstance.put(`/`, item, requestOptions)
  }
  catch (err) {
    console.log(err.response.data.error)
    return []
  }

  return response.data ? response.data : []
}

export async function apiDeletePath(path:string) {

  const encodedPath = encodePath(path)

  let response: AxiosResponse<any>
  try {
      response = await AxiosInstance.delete(`/${encodedPath}`)
  }
  catch (err) {
    console.log(err.response.data.error)
    return []
  }

  return response.data ? response.data : []
}