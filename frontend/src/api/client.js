import axios from "axios";

const API_BASE_URL = process.env.REACT_APP_API_URL || "/api";

const client = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
});

export function setToken(token) {
  if (token) {
    client.defaults.headers.common.Authorization = `Bearer ${token}`;
  }
}

export function clearToken() {
  delete client.defaults.headers.common.Authorization;
}

export async function register(payload) {
  const response = await client.post("/register", payload);
  return response.data;
}

export async function login(payload) {
  const response = await client.post("/login", payload);
  return response.data;
}

export async function fetchItems() {
  const response = await client.get("/items");
  return response.data.items;
}

export async function createItem(payload) {
  const response = await client.post("/items", payload);
  return response.data;
}

export async function updateItem(id, payload) {
  const response = await client.put(`/items/${id}`, payload);
  return response.data;
}

export async function deleteItem(id) {
  await client.delete(`/items/${id}`);
}

export async function fetchUsers() {
  const response = await client.get("/users");
  return response.data.users;
}

export async function deleteUser(id) {
  await client.delete(`/users/${id}`);
}

const api = {
  setToken,
  clearToken,
  register,
  login,
  fetchItems,
  createItem,
  updateItem,
  deleteItem,
  fetchUsers,
  deleteUser,
};

export default api;
