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
   * JWT token expiration time
   */
  JWT_EXPIRY: '7d',

  /**
   * OAuth state token expiration time in seconds (10 minutes)
   */
  OAUTH_STATE_EXPIRY_SECONDS: 600,
} as const
