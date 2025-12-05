export const security = {
  /**
   * Session cookie expiration time in seconds
   * 7 days = 604800 seconds
   * Consider shorter duration (1-24h) with refresh tokens for higher security
   */
  SESSION_EXPIRY_SECONDS: 60 * 60 * 24 * 7,

  /**
   * OAuth state token expiration time in seconds (10 minutes)
   */
  OAUTH_STATE_EXPIRY_SECONDS: 600,
} as const
