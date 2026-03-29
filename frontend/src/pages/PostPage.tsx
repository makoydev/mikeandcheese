import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { ArrowLeft } from "lucide-react";
import { format } from "date-fns";
import { getPost } from "@/lib/api";
import type { Post } from "@/types/post";

function estimateReadingTime(content: string): number {
  const words = content.trim().split(/\s+/).length;
  return Math.max(1, Math.round(words / 250));
}

export function PostPage() {
  const { slug } = useParams<{ slug: string }>();
  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!slug) return;
    async function load() {
      try {
        const data = await getPost(slug!);
        setPost(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load post");
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [slug]);

  if (loading) {
    return (
      <div className="py-20 text-center text-text-secondary">
        <div className="inline-block h-5 w-5 border-2 border-text-secondary border-t-transparent rounded-full animate-spin" />
      </div>
    );
  }

  if (error || !post) {
    return (
      <div className="py-20 text-center">
        <p className="text-text-secondary mb-4">{error ?? "Post not found"}</p>
        <Link to="/" className="text-primary hover:underline text-sm">
          Back to home
        </Link>
      </div>
    );
  }

  const readingTime = estimateReadingTime(post.content);
  const formattedDate = format(new Date(post.created_at), "MMMM d, yyyy");

  return (
    <article>
      <Link
        to="/"
        className="inline-flex items-center gap-1.5 text-sm text-text-secondary hover:text-text-primary dark:hover:text-text-primary-dark transition-colors mb-8"
      >
        <ArrowLeft size={14} />
        Back to home
      </Link>

      <header className="mb-10">
        {post.category && (
          <span className="inline-block text-xs font-medium uppercase tracking-wider text-primary mb-3">
            {post.category}
          </span>
        )}
        <h1 className="font-serif text-3xl sm:text-4xl font-bold text-text-primary dark:text-text-primary-dark leading-tight mb-4">
          {post.title}
        </h1>
        <div className="flex flex-wrap items-center gap-3 text-sm text-text-secondary">
          <span>{post.author}</span>
          <span>&middot;</span>
          <time dateTime={post.created_at}>{formattedDate}</time>
          <span>&middot;</span>
          <span>{readingTime} min read</span>
        </div>
        <Link
          to={`/edit/${post.slug}`}
          className="inline-flex items-center gap-1 text-sm text-text-secondary hover:text-primary transition-colors"
        >
          Edit post
        </Link>
      </header>

      <div
        className="prose prose-lg max-w-none"
        dangerouslySetInnerHTML={{ __html: post.content }}
      />
    </article>
  );
}
