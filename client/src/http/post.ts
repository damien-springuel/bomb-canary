import type { AxiosRequestConfig, AxiosResponse } from "axios";

export interface HttpPost {
  post: (url: string, data?: any, config?: AxiosRequestConfig) => Promise<AxiosResponse<{}>>
}