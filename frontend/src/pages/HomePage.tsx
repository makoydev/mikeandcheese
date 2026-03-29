import { useEffect, useState } from "react";
import { PostCard } from "@/components/PostCard";
import { getPosts, getFeaturedPosts } from "@/lib/api";
import type { Post } from "@/types/post";

export function HomePage() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [featured, setFeatured] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function load() {
      try {
        const [allPosts, featuredPosts] = await Promise.all([
          getPosts(),
          getFeaturedPosts().catch(() => [] as Post[]),
        ]);
        setPosts(allPosts);
        setFeatured(featuredPosts);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load posts");
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  if (loading) {
    return (
      <div className="py-20 text-center text-text-secondary">
        <div className="inline-block h-5 w-5 border-2 border-text-secondary border-t-transparent rounded-full animate-spin" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="py-20 text-center text-text-secondary">
        <p>{error}</p>
      </div>
    );
  }

  const featuredPost = featured[0];
  const recentPosts = posts.filter((p) => p.id !== featuredPost?.id);

  return (
    <div>
      {featuredPost && (
        <section className="mb-12 pb-10 border-b border-border dark:border-border-dark">
          <PostCard post={featuredPost} featured />
        </section>
      )}

      <section>
        <h2 className="text-sm font-medium uppercase tracking-wider text-text-secondary mb-8">
          Recent Posts
        </h2>
        {recentPosts.length > 0 ? (
          <div className="divide-y divide-border dark:divide-border-dark">
            {recentPosts.map((post) => (
              <div key={post.id} className="py-6 first:pt-0">
                <PostCard post={post} />
              </div>
            ))}
          </div>
        ) : (
          <p className="text-text-secondary">No posts yet.</p>
        )}
      </section>
    </div>
  );
}
