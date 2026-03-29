import { useEffect, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { ArrowLeft } from "lucide-react";
import { getPost, createPost, updatePost, deletePost } from "@/lib/api";
import type { Post } from "@/types/post";

function titleToSlug(title: string): string {
  return title
    .toLowerCase()
    .replace(/\s+/g, "-")
    .replace(/[^a-z0-9-]/g, "");
}

const inputClass =
  "w-full px-3 py-2 rounded-md border border-border dark:border-border-dark bg-transparent text-text-primary dark:text-text-primary-dark font-sans text-base focus:outline-none focus:ring-2 focus:ring-primary/50";

const labelClass = "text-sm font-medium text-text-secondary";

export function EditorPage() {
  const { slug } = useParams<{ slug: string }>();
  const navigate = useNavigate();
  const isEdit = Boolean(slug);

  const [loading, setLoading] = useState(isEdit);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [title, setTitle] = useState("");
  const [formSlug, setFormSlug] = useState("");
  const [excerpt, setExcerpt] = useState("");
  const [content, setContent] = useState("");
  const [author, setAuthor] = useState("");
  const [category, setCategory] = useState("");
  const [tags, setTags] = useState("");
  const [featured, setFeatured] = useState(false);
  const [slugManuallyEdited, setSlugManuallyEdited] = useState(false);

  useEffect(() => {
    if (!slug) return;
    async function load() {
      try {
        const post = await getPost(slug!);
        setTitle(post.title);
        setFormSlug(post.slug);
        setExcerpt(post.excerpt);
        setContent(post.content);
        setAuthor(post.author);
        setCategory(post.category);
        setTags(post.tags);
        setFeatured(post.featured);
        setSlugManuallyEdited(true);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load post");
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [slug]);

  function handleTitleChange(value: string) {
    setTitle(value);
    if (!slugManuallyEdited) {
      setFormSlug(titleToSlug(value));
    }
  }

  function handleSlugChange(value: string) {
    setSlugManuallyEdited(true);
    setFormSlug(value);
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setSubmitting(true);
    setError(null);

    const data: Omit<Post, "id" | "created_at" | "updated_at"> = {
      title,
      slug: formSlug,
      excerpt,
      content,
      author,
      category,
      tags,
      featured,
    };

    try {
      if (isEdit) {
        await updatePost(slug!, data);
      } else {
        await createPost(data);
      }
      navigate(`/post/${formSlug}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save post");
      setSubmitting(false);
    }
  }

  async function handleDelete() {
    if (!confirm("Are you sure you want to delete this post? This cannot be undone.")) {
      return;
    }
    try {
      await deletePost(slug!);
      navigate("/");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete post");
    }
  }

  if (loading) {
    return (
      <div className="py-20 text-center text-text-secondary">
        <div className="inline-block h-5 w-5 border-2 border-text-secondary border-t-transparent rounded-full animate-spin" />
      </div>
    );
  }

  return (
    <div>
      <Link
        to="/"
        className="inline-flex items-center gap-1.5 text-sm text-text-secondary hover:text-text-primary dark:hover:text-text-primary-dark transition-colors mb-8"
      >
        <ArrowLeft size={14} />
        Back to home
      </Link>

      <h1 className="font-serif text-3xl font-bold text-text-primary dark:text-text-primary-dark mb-8">
        {isEdit ? "Edit Post" : "New Post"}
      </h1>

      {error && (
        <div className="mb-6 rounded-md border border-red-300 bg-red-50 dark:bg-red-900/20 dark:border-red-800 px-4 py-3 text-sm text-red-700 dark:text-red-400">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label htmlFor="title" className={labelClass}>
            Title
          </label>
          <input
            id="title"
            type="text"
            required
            value={title}
            onChange={(e) => handleTitleChange(e.target.value)}
            className={`${inputClass} mt-1`}
          />
        </div>

        <div>
          <label htmlFor="slug" className={labelClass}>
            Slug
          </label>
          <input
            id="slug"
            type="text"
            required
            value={formSlug}
            onChange={(e) => handleSlugChange(e.target.value)}
            className={`${inputClass} mt-1`}
          />
        </div>

        <div>
          <label htmlFor="excerpt" className={labelClass}>
            Excerpt
          </label>
          <textarea
            id="excerpt"
            rows={3}
            value={excerpt}
            onChange={(e) => setExcerpt(e.target.value)}
            className={`${inputClass} mt-1`}
          />
        </div>

        <div>
          <label htmlFor="content" className={labelClass}>
            Content
          </label>
          <textarea
            id="content"
            rows={20}
            required
            value={content}
            onChange={(e) => setContent(e.target.value)}
            className={`${inputClass} mt-1`}
          />
        </div>

        <div>
          <label htmlFor="author" className={labelClass}>
            Author
          </label>
          <input
            id="author"
            type="text"
            required
            value={author}
            onChange={(e) => setAuthor(e.target.value)}
            className={`${inputClass} mt-1`}
          />
        </div>

        <div>
          <label htmlFor="category" className={labelClass}>
            Category
          </label>
          <input
            id="category"
            type="text"
            value={category}
            onChange={(e) => setCategory(e.target.value)}
            className={`${inputClass} mt-1`}
          />
        </div>

        <div>
          <label htmlFor="tags" className={labelClass}>
            Tags (comma separated)
          </label>
          <input
            id="tags"
            type="text"
            value={tags}
            onChange={(e) => setTags(e.target.value)}
            className={`${inputClass} mt-1`}
          />
        </div>

        <div className="flex items-center gap-2">
          <input
            id="featured"
            type="checkbox"
            checked={featured}
            onChange={(e) => setFeatured(e.target.checked)}
            className="rounded border-border dark:border-border-dark text-primary focus:ring-primary/50"
          />
          <label htmlFor="featured" className={labelClass}>
            Featured
          </label>
        </div>

        <div className="flex items-center justify-between pt-4">
          <button
            type="submit"
            disabled={submitting}
            className="bg-primary text-white px-6 py-2 rounded-md hover:opacity-90 font-medium disabled:opacity-50"
          >
            {submitting ? "Saving..." : isEdit ? "Update" : "Publish"}
          </button>

          {isEdit && (
            <button
              type="button"
              onClick={handleDelete}
              className="text-red-500 hover:text-red-700 text-sm"
            >
              Delete post
            </button>
          )}
        </div>
      </form>
    </div>
  );
}
