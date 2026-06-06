# Security Policy

## Supported Versions

Security fixes are applied to the latest released version only. Please make sure
you are running the most recent release before reporting an issue.

| Version | Supported          |
| ------- | ------------------ |
| Latest release | :white_check_mark: |
| Older releases  | :x:                |

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues,
discussions, or pull requests.** Public disclosure can put every user running
this panel at risk before a fix is available.

Instead, report privately through one of these channels:

1. **GitHub Private Vulnerability Reporting** (preferred)
   Go to the **Security** tab of this repository and click
   **"Report a vulnerability"**. This keeps the report private until a fix is
   ready.
2. If that is unavailable, open a minimal issue asking the maintainer to get in
   touch — **without** including any vulnerability details — and we will arrange
   a private channel.

When reporting, please include:

- A description of the vulnerability and its potential impact.
- Steps to reproduce (proof of concept if possible).
- The affected version / commit.
- Any suggested fix or mitigation, if you have one.

**Please redact sensitive data** from your report — do not paste real
credentials, tokens, private keys, certificates, server IPs, or domains.

## What to Expect

- We will acknowledge your report as soon as we are able.
- We will investigate and keep you updated on progress.
- Once a fix is released, we are happy to credit you in the release notes
  (unless you prefer to stay anonymous).

We ask that you give us a reasonable amount of time to release a fix before any
public disclosure.

## Security Best Practices for Users

This project is a proxy management panel and is exposed to the network. To stay
safe:

- **Change the default credentials immediately** after installation.
- Keep the panel **updated to the latest release**.
- Do **not** expose the panel directly to the public internet without
  protection — use a firewall, restrict access by IP, and/or place it behind a
  reverse proxy.
- Always enable **HTTPS/TLS** for the panel.
- Never commit or share your configuration files, database, certificates, or
  credentials.

> **Disclaimer:** This project is provided for personal learning and
> communication only, without warranty of any kind. You are responsible for how
> you deploy and use it.
