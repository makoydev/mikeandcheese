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
