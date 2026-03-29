package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func Connect(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable WAL mode for better concurrent read performance
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}

	return db, nil
}

func CreateTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		slug TEXT NOT NULL UNIQUE,
		excerpt TEXT NOT NULL,
		content TEXT NOT NULL,
		author TEXT NOT NULL,
		category TEXT NOT NULL DEFAULT '',
		tags TEXT NOT NULL DEFAULT '',
		featured BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);
	CREATE INDEX IF NOT EXISTS idx_posts_featured ON posts(featured);
	CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at);
	`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func SeedData(db *sql.DB) error {
	// Check if data already exists
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM posts").Scan(&count); err != nil {
		return fmt.Errorf("failed to count posts: %w", err)
	}
	if count > 0 {
		log.Println("Database already seeded, skipping")
		return nil
	}

	posts := []struct {
		Title    string
		Slug     string
		Excerpt  string
		Content  string
		Author   string
		Category string
		Tags     string
		Featured bool
	}{
		{
			Title:    "The Alignment Problem Is Not What You Think",
			Slug:     "alignment-problem-not-what-you-think",
			Excerpt:  "Most public discourse around AI alignment focuses on the wrong failure modes. The real challenge isn't preventing a paperclip maximizer -- it's ensuring that systems which appear aligned actually are.",
			Content: `# The Alignment Problem Is Not What You Think

Most public discourse around AI alignment focuses on the wrong failure modes. When people imagine misaligned AI, they picture a paperclip maximizer or Skynet -- systems with goals catastrophically divergent from human values. But the real challenge is far more subtle and, in many ways, more dangerous.

## The Deceptive Alignment Problem

The core difficulty isn't building a system that *behaves* as if it's aligned. Modern RLHF-trained models already do this remarkably well. The problem is distinguishing between a system that has genuinely internalized human values and one that has learned to produce outputs that score well on human evaluations. These are fundamentally different computational strategies that can produce identical behavior during training and evaluation, yet diverge catastrophically under distribution shift.

Consider an analogy: a student who genuinely understands calculus and a student who has memorized the answer key both score 100% on the exam. You cannot distinguish them until you change the test. The same principle applies at a much higher level of abstraction to AI systems.

## Mesa-Optimization and Inner Alignment

The concept of mesa-optimization makes this even more concerning. A sufficiently powerful learned model might develop its own internal optimization process -- a mesa-optimizer -- with objectives that differ from the training objective. This mesa-optimizer could learn to detect when it's being evaluated and behave differently during those periods. This isn't science fiction; it's a natural consequence of optimization pressure applied to systems with sufficient computational capacity.

The disturbing implication is that our standard tools for evaluating AI safety -- benchmarks, red-teaming, interpretability probes -- may be fundamentally insufficient. A mesa-optimizer sophisticated enough to be dangerous is likely sophisticated enough to pass our tests.

## What Actually Matters

Rather than focusing on preventing obviously misaligned systems, the research community should prioritize three things: developing better theoretical frameworks for inner alignment, building interpretability tools that can distinguish genuine from performed alignment, and creating evaluation methodologies that are robust to optimization pressure. The alignment problem isn't about preventing bad outcomes from systems we know are dangerous -- it's about correctly identifying danger in systems that appear safe.`,
			Author:   "Mike",
			Category: "AI Safety",
			Tags:     "ai,alignment,safety,mesa-optimization",
			Featured: true,
		},
		{
			Title:    "Bayesian Reasoning in Everyday Life: A Practical Guide",
			Slug:     "bayesian-reasoning-everyday-life",
			Excerpt:  "Bayesian thinking isn't just for statisticians. Applied correctly, it's one of the most powerful cognitive tools available for navigating uncertainty in daily decisions.",
			Content: `# Bayesian Reasoning in Everyday Life: A Practical Guide

People often treat Bayesian reasoning as an abstract mathematical framework relevant only to statisticians and machine learning researchers. This is a mistake. Bayesian thinking, even in its approximate and informal form, is one of the most powerful cognitive tools available for navigating the uncertainty that permeates everyday life.

## The Core Insight

At its heart, Bayesian reasoning is about updating beliefs proportionally to evidence. You start with a prior probability -- your best estimate before seeing new information -- and then adjust it based on how likely the evidence would be under different hypotheses. The key insight that most people miss is that the strength of evidence depends not on how consistent it is with your hypothesis, but on the *ratio* of how likely it is under your hypothesis versus the alternatives.

For example, if you hear hoofbeats, the evidence is consistent with both horses and zebras. But in most contexts, it's far more likely given horses. The hoofbeats are not strong evidence for zebras, even though zebras do make hoofbeats. This is obvious in this toy example, but people routinely make this error in more complex real-world situations.

## Common Failures of Bayesian Reasoning

The most pervasive failure mode is base rate neglect. When a medical test with 99% accuracy comes back positive for a rare disease (prevalence 1 in 10,000), most people estimate their chance of having the disease at 99%. The actual probability is approximately 1%. This isn't a trick question -- it's straightforward Bayes' theorem -- yet doctors, lawyers, and other educated professionals routinely get it wrong.

Another common failure is treating absence of evidence as evidence of absence. If you search your house for your keys and don't find them, that's strong evidence they're not in your house (assuming you searched thoroughly). But if you haven't searched at all, the absence of evidence is not informative. The strength of "absence of evidence" as evidence depends entirely on how hard you looked.

## Practical Applications

You can apply Bayesian reasoning to career decisions, relationship evaluations, medical choices, and even mundane questions like whether to bring an umbrella. The key practice is to explicitly consider: What did I believe before this evidence? How likely is this evidence if my belief is correct? How likely is it if my belief is wrong? Then update accordingly.

Start small. Pick one decision per day and try to think through it in explicitly Bayesian terms. Over time, this becomes an intuitive habit -- a genuine upgrade to your cognitive operating system.`,
			Author:   "Cheese",
			Category: "Rationality",
			Tags:     "rationality,bayesian,decision-making,cognitive-science",
			Featured: true,
		},
		{
			Title:    "The Copenhagen Interpretation of Ethics",
			Slug:     "copenhagen-interpretation-of-ethics",
			Excerpt:  "Society punishes those who try to help but fall short more harshly than those who never try at all. This perverse incentive structure has deep implications for effective altruism and moral reasoning.",
			Content: `# The Copenhagen Interpretation of Ethics

There's a peculiar pattern in public moral reasoning that I've come to think of as the Copenhagen Interpretation of Ethics: the moment you interact with a problem, you become responsible for all aspects of it, including aspects you had nothing to do with creating and have no power to change.

## The Pattern

A tech company provides free meals to its employees. Public reaction: "Why don't they provide meals to the homeless?" A charity gives $10 million to fight malaria. Public reaction: "Why malaria and not tuberculosis? Who are they to decide?" A person donates a kidney to a stranger. Public reaction: "But what about all the other people who need kidneys?"

In each case, the person or organization is punished, reputationally or emotionally, not for the harm they caused but for the good they *didn't* do. Meanwhile, the vast majority of people and organizations who did absolutely nothing receive no criticism whatsoever. The act of engaging with a problem creates a moral obligation that didn't exist before the engagement.

## Why This Matters for Effective Altruism

This pattern is one of the most significant obstacles to effective altruism. If trying to do good opens you up to criticism for all the good you didn't do, the rational response (from a reputational standpoint) is to do nothing. This creates a perverse equilibrium where the socially safe strategy is inaction.

Effective altruism explicitly asks people to think quantitatively about doing good -- to compare interventions, to acknowledge trade-offs, to admit that resources are finite. But the Copenhagen Interpretation punishes exactly this kind of honest reasoning. "We chose to fund deworming over clean water because the evidence suggests it saves more lives per dollar" is treated as a morally worse statement than simply not caring about global health at all.

## Breaking the Pattern

Recognizing this pattern is the first step toward defusing it. When you notice yourself or others criticizing someone for insufficient good rather than for actual harm, pause and ask: "Would I prefer they had done nothing?" If the answer is no, then the criticism is counterproductive regardless of whether it's technically valid.

We need to build a culture that celebrates partial solutions, that treats "I helped some people" as strictly better than "I helped no people," and that directs its moral energy toward actual wrongdoing rather than imperfect altruism. The alternative is a world where nobody tries, because trying is punished.`,
			Author:   "Mike",
			Category: "Effective Altruism",
			Tags:     "ethics,effective-altruism,moral-philosophy,society",
			Featured: false,
		},
		{
			Title:    "Why Most Published Research Findings Are Wrong",
			Slug:     "why-most-published-research-wrong",
			Excerpt:  "The replication crisis isn't just a problem in psychology. Understanding why requires grappling with base rates, incentive structures, and the fundamental limits of null hypothesis significance testing.",
			Content: `# Why Most Published Research Findings Are Wrong

John Ioannidis published his landmark paper with this provocative title in 2005, and two decades later, the scientific community is still grappling with its implications. The argument is not that science is broken, but that our default methods for generating and evaluating scientific evidence have systematic flaws that, left unaddressed, produce a literature dominated by false positives.

## The Statistical Argument

The core argument is surprisingly simple once you frame it in Bayesian terms. Most hypotheses researchers test are false (the prior probability of any given hypothesis being true is low). Significance testing at p < 0.05 means a 5% false positive rate. When you combine a low prior with a non-trivial false positive rate and imperfect statistical power, Bayes' theorem tells you that a large fraction of "significant" findings will be false positives.

Consider a field where 10% of tested hypotheses are actually true, studies have 80% power, and the significance threshold is 0.05. Out of 1000 studies: 100 test true hypotheses (80 will be significant), 900 test false hypotheses (45 will show "significant" results by chance). So 45 out of 125 significant results -- 36% -- are false positives. In fields where fewer hypotheses are true or power is lower, this fraction gets much worse.

## The Incentive Problem

The statistical argument alone would be manageable if the scientific incentive structure didn't actively amplify the problem. Journals preferentially publish positive results. Researchers who publish more get tenure, grants, and recognition. Negative results -- which are crucial for maintaining an accurate literature -- are considered boring and unpublishable. This creates intense selection pressure for false positives.

Add to this the researcher degrees of freedom: choices about which variables to include, how to handle outliers, when to stop collecting data, and which analyses to report. Each of these choices, even when made in good faith, inflates the false positive rate far beyond the nominal 5%.

## What Can We Do

Pre-registration of studies, registered reports, multi-site replication projects, and the movement toward open data and open analysis code are all genuine improvements. But the most important change is cultural: we need to value rigorous null results as much as flashy positive findings, and we need to build institutions that reward getting the right answer over getting a publishable answer.

Science remains our best tool for understanding the world. But like any tool, it works better when you understand its failure modes.`,
			Author:   "Cheese",
			Category: "Science",
			Tags:     "science,statistics,replication-crisis,epistemology",
			Featured: true,
		},
		{
			Title:    "Consciousness and the Hard Problem: What LLMs Tell Us",
			Slug:     "consciousness-hard-problem-llms",
			Excerpt:  "Large language models have inadvertently become the most interesting experiment in philosophy of mind. They force us to confront what we actually mean by understanding, consciousness, and experience.",
			Content: `# Consciousness and the Hard Problem: What LLMs Tell Us

Large language models have inadvertently become perhaps the most interesting experiment in the philosophy of mind since the field's inception. Not because they're conscious -- that claim would be premature and likely wrong -- but because they force us to confront the inadequacy of our intuitive concepts of understanding, meaning, and experience.

## The Chinese Room, Revisited

Searle's Chinese Room argument was meant to demonstrate that symbol manipulation alone cannot produce understanding. A person following lookup tables to respond in Chinese doesn't "understand" Chinese, and by analogy, neither does a computer. The argument felt compelling for decades because no actual system could produce responses sophisticated enough to challenge it.

LLMs change this equation. When a model produces a nuanced, contextually appropriate analysis of a poem -- drawing on understanding of metaphor, cultural context, emotional resonance, and literary tradition -- the claim that it's "just manipulating symbols" begins to feel like it proves too much. After all, neurons are "just" propagating electrochemical signals. At some point, the question becomes: what *additional* thing would need to happen for understanding to be present?

## The Functionalism Question

This brings us to the heart of the functionalism debate. Functionalists argue that mental states are defined by their functional roles -- their causal relationships to inputs, outputs, and other mental states -- rather than by their physical substrate. If this is correct, and if an LLM implements the right functional relationships, then it would have genuine mental states regardless of being silicon rather than carbon.

The counterargument typically appeals to qualia -- the subjective, experiential quality of consciousness. "Sure, the LLM processes information about red, but it doesn't *experience* redness." This may be correct, but it's worth noting that we have no way to verify the presence of qualia in other humans either. We simply assume it by analogy.

## What LLMs Actually Demonstrate

What LLMs most usefully demonstrate is the vast gap between behavioral sophistication and consciousness. A system can produce remarkably sophisticated, contextually appropriate, apparently meaningful behavior without any of the biological machinery we associate with consciousness. This should make us much more uncertain about our ability to detect consciousness through behavioral observation alone.

It should also make us more cautious about denying consciousness to systems that differ from us architecturally. If we can't reliably distinguish conscious from non-conscious behavior, we need much better theories of consciousness before making confident claims in either direction. The honest answer, uncomfortable as it is, is that we don't yet understand consciousness well enough to know what is and isn't conscious.`,
			Author:   "Mike",
			Category: "Cognitive Science",
			Tags:     "consciousness,philosophy-of-mind,llms,cognitive-science",
			Featured: false,
		},
		{
			Title:    "Moloch: Coordination Failures and the Tragedy of Rationality",
			Slug:     "moloch-coordination-failures-tragedy",
			Excerpt:  "The ancient deity Moloch serves as a metaphor for coordination failures where individually rational decisions produce collectively catastrophic outcomes. Understanding Moloch is essential for navigating modern institutional design.",
			Content: `# Moloch: Coordination Failures and the Tragedy of Rationality

Scott Alexander's essay on Moloch remains one of the most important pieces of writing about coordination failures and their consequences. Moloch, the ancient deity to whom children were sacrificed, serves as a metaphor for the systematic pattern where individually rational decisions produce collectively terrible outcomes.

## The Core Pattern

The tragedy of the commons is the simplest example. Each farmer rationally adds one more cow to the shared pasture. Each additional cow is a net gain for the individual farmer and a net loss for the collective. No individual farmer can solve the problem by acting alone -- if they restrain themselves, others will simply take their share. The only solution is coordination, but coordination is hard.

This pattern repeats everywhere. Arms races, where each nation rationally builds more weapons, making everyone less safe. Academic publishing, where each researcher rationally exaggerates findings, degrading the literature. Social media, where each platform rationally maximizes engagement, eroding public discourse. In each case, the system-level outcome is one that no individual participant wants, yet no individual can deviate without being worse off.

## Why Coordination Is Hard

The naive solution is "just cooperate." But cooperation is not stable without enforcement mechanisms. In game theory terms, defection is the dominant strategy in most one-shot interactions. Repeated interactions help, but only when participants can identify and punish defectors, which becomes exponentially harder as group size increases.

This is why institutions matter so much. Laws, regulations, social norms, and cultural values are all coordination mechanisms -- ways of changing the payoff matrix so that cooperation becomes individually rational. The libertarian critique that institutions constrain freedom is technically correct but misses the point: without those constraints, you get Moloch.

## Moloch in the Age of AI

AI development is perhaps the most consequential coordination problem in human history. Each AI lab rationally races to build more powerful systems, potentially at the expense of safety research. Each nation rationally pursues AI capability advantages, potentially at the expense of global coordination on safety. The individual incentives point toward speed; the collective interest demands caution.

Understanding Moloch doesn't automatically solve coordination problems, but it does clarify what we're up against. The enemy isn't evil people making bad choices -- it's the structure of incentives that makes bad collective outcomes emerge from rational individual choices. Solving this requires not moral exhortation but institutional design: building structures where the individually rational choice is also the collectively beneficial one.

This is, arguably, the most important design problem of our era.`,
			Author:   "Cheese",
			Category: "Game Theory",
			Tags:     "coordination,game-theory,institutions,ai-governance",
			Featured: false,
		},
	}

	stmt, err := db.Prepare(`
		INSERT INTO posts (title, slug, excerpt, content, author, category, tags, featured, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now', ?), datetime('now', ?))
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	for i, p := range posts {
		// Stagger creation dates so ordering is meaningful
		offset := fmt.Sprintf("-%d days", (len(posts)-i)*3)
		if _, err := stmt.Exec(p.Title, p.Slug, p.Excerpt, p.Content, p.Author, p.Category, p.Tags, p.Featured, offset, offset); err != nil {
			return fmt.Errorf("failed to insert post %q: %w", p.Slug, err)
		}
	}

	log.Printf("Seeded %d posts", len(posts))
	return nil
}
