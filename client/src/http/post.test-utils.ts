import type { AxiosRequestConfig, AxiosResponse } from "axios";

export class HttpPostMock {
  public givenUrl: string = "";
  public givenData: any;
  public givenConfig: AxiosRequestConfig;
  constructor(private readonly responseToReturn?: Promise<AxiosResponse<{}>>){}

  post(url: string, data?: any, config?: AxiosRequestConfig): Promise<AxiosResponse<{}>> {
    this.givenUrl = url;
    this.givenData = data;
    this.givenConfig = config;
    return this.responseToReturn
  }
}