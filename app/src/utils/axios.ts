import axios from 'axios'

export const AxiosInstance = axios.create({
  baseURL: 'https://virtualfs.herokuapp.com/api',
})


