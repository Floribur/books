import { marked } from 'marked';
import bioRaw from '../content/bio.md?raw';
import './Bio.css';

// bio.md is a local bundled file — no user input — dangerouslySetInnerHTML is safe here
// marked v17 parse() returns string synchronously when no async extensions are used
const bioHtml = marked.parse(bioRaw) as string;

export function Bio() {
  return (
    <section className="bio-section" aria-label="About Florian">
      <div className="bio-layout">
        <div className="bio-photo-wrapper">
          <img
            src="/florian.jpg"
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
