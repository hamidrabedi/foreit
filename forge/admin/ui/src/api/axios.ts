import axios from "axios";
import { getAdminConfig } from "../lib/config";

export const axiosInstance = axios.create({
  baseURL: getAdminConfig().apiBase,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true,
});

export default axiosInstance;
