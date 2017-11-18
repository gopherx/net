package nat

// PasswordResolverFunc resolves passwords from usernames.
type PasswordResolverFunc func(username string) ([]byte, error)
