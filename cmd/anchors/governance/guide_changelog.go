package governance

// changelogGuide is the ruler for the two changelogs: the TECHNICAL one `anchors changelog`
// builds from the commits, and the PRODUCT one an agent writes from it. The first is
// mechanical and complete; the second is for people who use the product, and only a
// reader of the first can decide what matters to them.
const changelogGuide = `# Changelog guide (the technical one, and the product one made from it)

## Two changelogs, two readers

` + "`anchors changelog`" + ` builds a TECHNICAL changelog: it reads the commits between two tags
and lists them by type. Its reader is whoever works on the code — a contributor, a
reviewer, whoever upgrades a dependency on this project and needs to know what moved.
It is complete and literal: every ` + "`feat`" + `, every ` + "`fix`" + `, every breaking change, in the
words the commit used, with its hash.

It is NOT a changelog for the product's users. A user does not know what a scope is,
does not care that a function was renamed, and reads "fix(budget): a cut batch gets a
TERM and a grace before the KILL" as noise. Publishing the technical changelog as the
product's release notes is handing them the kitchen's order tickets.

So there are two, and the second is made from the first:

    commits ──anchors changelog──▶ technical changelog ──an agent synthesizes──▶ product changelog

## The technical changelog

It is only as good as the commits, so the ruler is the commit (` + "`anchors guide work`" + `):

- the TYPE decides the section: ` + "`feat`" + ` is a feature, ` + "`fix`" + ` is a fix, a ` + "`!`" + ` after the
  type or a ` + "`BREAKING CHANGE:`" + ` footer is a breaking change, whatever the type;
- a ` + "`fix`" + ` of a defect that SHIPPED carries ` + "`Bug: <where it was seen>`" + `, and goes to
  "Bugs fixed"; a fix of work that never reached anyone stays in "Fixes";
- ` + "`refactor`" + `, ` + "`test`" + `, ` + "`chore`" + `, ` + "`docs`" + `, ` + "`ci`" + ` are history, not changelog: they are left out.

Generate it, do not write it: ` + "`anchors changelog --write`" + ` adds the releases the file does
not hold yet and keeps the rest as it is. The subject is the entry — write the subject
as the line you want to read in the changelog.

## The product changelog (recommended)

Write it with an agent, from the technical changelog of the release — not from the
commits, not from memory. One product changelog entry often sums several technical ones.

What goes in:

- **Breaking changes**, always — rewritten as what the user must do differently.
- **Features** the user can see or use, in the user's words: what they can do now, not
  what the code does.
- **Bugs fixed** (the "Bugs fixed" section, ` + "`fix`" + ` with ` + "`Bug:`" + `): the user may have met them,
  so say what was wrong in terms they would recognise.

What stays out:

- **Fixes** without ` + "`Bug:`" + `: they correct work that never reached anyone. The user never
  saw the defect, so there is nothing to tell them.
- **Internal features** — a new gate flag for maintainers, a refactor that happens to be
  a ` + "`feat`" + ` — unless they change what the user experiences.
- **Chores** and the other internal types — they are not even in the technical
  changelog. The one exception: a change the user feels (a dependency upgrade that drops
  support for an old platform, a faster startup, a new minimum version). If it matters to
  the product, say it, in the product's words; read the commits of the release for that.

How it reads:

- grouped by what the user does (new, changed, fixed), not by scope or module;
- no hashes, no scopes, no file names, no internal codes;
- short: a release with forty technical entries is often five product lines;
- nothing invented — every line traces back to an entry of the technical changelog (or
  a commit, for the chore exception). When unsure whether something matters to the
  user, leave it out, or ask who owns the product.

Where it lives is the project's choice (a ` + "`RELEASES.md`" + `, the release notes of the tag, a
page of the site). Keep it apart from the technical ` + "`CHANGELOG.md`" + `, so neither reader has
to skip the other's lines.
`
