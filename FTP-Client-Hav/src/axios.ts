import axios from 'axios';

const axiosInstance = axios.create({
  baseURL: "http://localhost:8080/v1",
  headers: {
    "Content-Type": "application/json",
  },
  // 自定义响应成功的状态码范围
  validateStatus: function (status) {
    return status >= 200 && status < 400; // 只有 200-299 状态码认为是成功
  },
});

// 添加响应拦截器
axiosInstance.interceptors.response.use(
  response => {
    return response;
  },
  error => {
    // 对响应错误做点什么
    if (error.response && error.response.status >= 300 && error.response.status < 400) {
      // 你可以在这里处理 300-399 状态码的情况
      console.error('重定向错误:', error.response.status);
    }
    // 返回任何非 2xx 的响应码都将触发错误
    return Promise.reject(error);
  }
);

export default axiosInstance;
