/// <reference types="vite/client" />

// Allow Vite's ?raw suffix for raw text imports (used for bio.md)
declare module '*.md?raw' {
  const content: string;
  export default content;
}
