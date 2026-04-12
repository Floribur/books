import { marked, Renderer } from 'marked';
import bioRaw from '../content/bio.md?raw';
import './Bio.css';

// Open all markdown links in a new tab
const renderer = new Renderer();
renderer.link = ({ href, title, text }: { href: string; title?: string | null; text: string }) =>
  `<a href="${href}" target="_blank" rel="noopener noreferrer"${title ? ` title="${title}"` : ''}>${text}</a>`;

// bio.md is a local bundled file — no user input — dangerouslySetInnerHTML is safe here
// marked v17 parse() returns string synchronously when no async extensions are used
const bioHtml = marked.parse(bioRaw, { renderer }) as string;

export function Bio() {
  return (
    <section className="bio-section" aria-label="About Florian">
      <div className="bio-layout">
        <div className="bio-photo-wrapper">
          <img
            src={`${import.meta.env.BASE_URL}florian.jpg`}
            alt="Florian"
            className="bio-photo"
          />
        </div>
        <div
          className="bio-text"
          dangerouslySetInnerHTML={{ __html: bioHtml }}
        />
      </div>
    </section>
  );
}
