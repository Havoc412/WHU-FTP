/// <reference types="vite/client" />

declare module "*.vue" {
  import Vue from "vue";
  export default Vue;
}

// 添加对 "@/axios" 的类型声明
declare module '@/axios' {
  import { AxiosInstance } from 'axios';
  const axiosInstance: AxiosInstance;
  export default axiosInstance;
}

declare module '@/const/electron.js' {
  export const fsFilterType;
}