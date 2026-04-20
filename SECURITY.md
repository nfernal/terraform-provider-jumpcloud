# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in this project, please report it responsibly.

**Please do NOT open a public GitHub issue for security vulnerabilities.**

Instead, please send an email to security@nfernal.com with:

- A description of the vulnerability
- Steps to reproduce the issue
- Any potential impact
- Suggested fix (if any)

We will acknowledge receipt within 48 hours and aim to provide a fix within a reasonable timeframe depending on severity.

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| latest  | Yes                |

## Security Best Practices

When using this provider:

- Store API keys in environment variables or a secrets manager, never in `.tf` files.
- Use Terraform state encryption for any state that may contain sensitive values.
- Restrict API key permissions to the minimum required scope.
