export interface Post {
  id: number;
  title: string;
  slug: string;
  excerpt: string;
  content: string;
  author: string;
  category: string;
  tags: string;
  featured: boolean;
  created_at: string;
  updated_at: string;
}
