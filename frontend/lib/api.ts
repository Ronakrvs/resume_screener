import axios from "axios";

const BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081/api";

const api = axios.create({ baseURL: BASE_URL });

api.interceptors.request.use((config) => {
  const token = typeof window !== "undefined" ? localStorage.getItem("token") : null;
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

export interface User {
  id: string;
  email: string;
  created_at: string;
}

export interface Resume {
  id: string;
  user_id: string;
  file_name: string;
  file_url: string;
  created_at: string;
}

export interface Job {
  id: string;
  user_id: string;
  title: string;
  description: string;
  created_at: string;
}

export interface Result {
  id: string;
  resume_id: string;
  job_id: string;
  score: number;
  strengths: string[];
  missing_skills: string[];
  recommendation: string;
  created_at: string;
  resume_file_name: string;
  job_title: string;
}

export const auth = {
  register: (email: string, password: string) =>
    api.post<{ token: string; user: User }>("/auth/register", { email, password }),
  login: (email: string, password: string) =>
    api.post<{ token: string; user: User }>("/auth/login", { email, password }),
};

export const resumes = {
  upload: (file: File) => {
    const fd = new FormData();
    fd.append("file", file);
    return api.post<Resume>("/resume/upload", fd, {
      headers: { "Content-Type": "multipart/form-data" },
    });
  },
  list: () => api.get<{ resumes: Resume[] }>("/resume"),
  get: (id: string) => api.get<Resume>(`/resume/${id}`),
};

export const jobs = {
  create: (title: string, description: string) =>
    api.post<Job>("/job", { title, description }),
  list: () => api.get<{ jobs: Job[] }>("/job"),
};

export const analysis = {
  analyze: (resume_id: string, job_id: string) =>
    api.post<Result>("/analyze", { resume_id, job_id }),
  getResult: (id: string) => api.get<Result>(`/results/${id}`),
  getRanking: (jobId: string) =>
    api.get<{ results: Result[]; total: number }>(`/ranking/${jobId}`),
};

export default api;
