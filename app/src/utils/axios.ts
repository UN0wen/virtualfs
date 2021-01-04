import axios from 'axios'

export const AxiosInstance = axios.create({
  baseURL: 'http://virtualfs.herokuapp.com:8080/api',
})


