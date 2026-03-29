export function AboutPage() {
  return (
    <div>
      <h1 className="font-serif text-3xl sm:text-4xl font-bold text-text-primary dark:text-text-primary-dark leading-tight mb-6">
        About
      </h1>
      <div className="prose prose-lg max-w-none">
        <p>
          <strong>mikeandcheese</strong> is a place for exploring ideas at the
          intersection of technology, artificial intelligence, rationality, and
          the many interesting things that emerge when you look closely at the
          world.
        </p>
        <p>
          The name is deliberately playful. The writing aims to be clear,
          honest, and occasionally surprising. Some posts are long explorations
          of complex topics; others are short observations that seemed worth
          sharing.
        </p>
        <p>
          This blog draws inspiration from communities and writers who value
          careful thinking, intellectual honesty, and the willingness to change
          one's mind when presented with better evidence or arguments.
        </p>
        <p>
          If you find something here that makes you think differently about a
          topic, or if it sparks a question you hadn't considered before, then
          it has done its job.
        </p>
      </div>
    </div>
  );
}
