import type { Post } from "@/types/post";

const API_BASE = "/api";

async function fetchJSON<T>(url: string): Promise<T> {
  const response = await fetch(`${API_BASE}${url}`);
  if (!response.ok) {
    throw new Error(`API error: ${response.status} ${response.statusText}`);
  }
  return response.json() as Promise<T>;
}

export function getPosts(): Promise<Post[]> {
  return fetchJSON<Post[]>("/posts");
}

export function getPost(slug: string): Promise<Post> {
  return fetchJSON<Post>(`/posts/${slug}`);
}

export function getFeaturedPosts(): Promise<Post[]> {
  return fetchJSON<Post[]>("/posts/featured");
}

export async function createPost(data: Omit<Post, "id" | "created_at" | "updated_at">): Promise<Post> {
  const response = await fetch(`${API_BASE}/posts`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!response.ok) throw new Error(`API error: ${response.status}`);
  return response.json();
}

export async function updatePost(slug: string, data: Omit<Post, "id" | "created_at" | "updated_at">): Promise<Post> {
  const response = await fetch(`${API_BASE}/posts/${slug}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  if (!response.ok) throw new Error(`API error: ${response.status}`);
  return response.json();
}

export async function deletePost(slug: string): Promise<void> {
  const response = await fetch(`${API_BASE}/posts/${slug}`, { method: "DELETE" });
  if (!response.ok) throw new Error(`API error: ${response.status}`);
}
