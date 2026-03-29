import { Link } from "react-router-dom";
import { format } from "date-fns";
import type { Post } from "@/types/post";

function estimateReadingTime(content: string): number {
  const words = content.trim().split(/\s+/).length;
  return Math.max(1, Math.round(words / 250));
}

interface PostCardProps {
  post: Post;
  featured?: boolean;
}

export function PostCard({ post, featured = false }: PostCardProps) {
  const readingTime = estimateReadingTime(post.content);
  const formattedDate = format(new Date(post.created_at), "MMMM d, yyyy");

  return (
    <article className={`group ${featured ? "mb-10" : "mb-8"}`}>
      {post.category && (
        <span className="inline-block text-xs font-medium uppercase tracking-wider text-primary mb-2">
          {post.category}
        </span>
      )}
      <h2
        className={`font-serif font-bold text-text-primary dark:text-text-primary-dark leading-tight mb-2 ${
          featured ? "text-3xl" : "text-xl"
        }`}
      >
        <Link
          to={`/post/${post.slug}`}
          className="hover:opacity-70 transition-opacity"
        >
          {post.title}
        </Link>
      </h2>
      {post.excerpt && (
        <p
          className={`text-text-secondary leading-relaxed mb-3 ${
            featured ? "text-base" : "text-sm"
          }`}
        >
          {post.excerpt}
        </p>
      )}
      <div className="flex items-center gap-3 text-xs text-text-secondary">
        <span>{post.author}</span>
        <span>&middot;</span>
        <time dateTime={post.created_at}>{formattedDate}</time>
        <span>&middot;</span>
        <span>{readingTime} min read</span>
      </div>
    </article>
  );
}
