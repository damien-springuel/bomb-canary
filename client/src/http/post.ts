import type { AxiosRequestConfig, AxiosResponse } from "axios";

export interface HttpPost<T> {
  post: (url: string, data?: any, config?: AxiosRequestConfig) => Promise<AxiosResponse<T>>
}