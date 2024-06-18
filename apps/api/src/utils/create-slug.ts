export function createSlug(text: string): string {
  // Convert to lowercase
  let slug = text.toLowerCase()

  // Remove accents and diacritics
  slug = slug.normalize('NFD').replace(/[\u0300-\u036f]/g, '')

  // Remove non-alphanumeric characters except for hyphens
  slug = slug.replace(/[^a-z0-9\s-]/g, '')

  // Replace spaces and consecutive hyphens with a single hyphen
  slug = slug.trim().replace(/\s+/g, '-').replace(/-+/g, '-')

  // Trim leading and trailing hyphens
  slug = slug.replace(/^-+|-+$/g, '')

  return slug
}
