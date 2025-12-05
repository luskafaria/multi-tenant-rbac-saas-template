export const security = {
  /**
   * bcrypt cost factor (work factor)
   * OWASP recommends minimum of 10, with 12+ preferred for production
   * Each increment doubles the computation time
   */
  BCRYPT_COST_FACTOR: 12,

  /**
   * Password reset token expiration time in minutes
   */
  PASSWORD_RESET_TOKEN_EXPIRY_MINUTES: 30,

  /**
   * JWT/Session cookie expiration time in seconds
   * 7 days = 604800 seconds
   * Consider shorter duration (1-24h) with refresh tokens for higher security
   */
  SESSION_EXPIRY_SECONDS: 60 * 60 * 24 * 7,

  /**
   * JWT token expiration time (for fastify-jwt sign option)
   */
  JWT_EXPIRY: '7d',

  /**
   * OAuth state token expiration time in seconds (10 minutes)
   */
  OAUTH_STATE_EXPIRY_SECONDS: 600,
} as const
