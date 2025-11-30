# TODO.AI.md - Jokes API Full SPEC Compliance

**Project**: jokes
**Organization**: apimgr
**Current Status**: Partial implementation, needs full SPEC compliance
**Started**: 2025-11-25

---

## Overview

Rebuild the jokes API to match the complete SPEC requirements. The current implementation has basic functionality but needs significant enhancements for full compliance.

---

## Phase 1: Project Structure Reorganization

- [ ] Reorganize all source files into `./src` directory structure
- [ ] Move all scripts to `./scripts` directory
- [ ] Create `./tests` directory for test files
- [ ] Update all imports and paths after reorganization
- [ ] Verify build still works after reorganization

---

## Phase 2: Configuration File Enhancement

- [ ] Design comprehensive YAML config structure with all SPEC requirements
- [ ] Add server configuration (address, port, dual-port support)
- [ ] Add SSL/TLS configuration (Let's Encrypt support)
- [ ] Add web-ui configuration (theme, logo, etc.)
- [ ] Add web-robots configuration (allow/deny rules)
- [ ] Add web-security configuration (security contact)
- [ ] Add CORS configuration (default: '*')
- [ ] Add scheduler configuration
- [ ] Add notification configuration
- [ ] Add comprehensive single-line comments (under 140 chars)
- [ ] Implement config auto-creation with all defaults
- [ ] Implement config live-reload support

---

## Phase 3: Let's Encrypt & SSL/TLS Support

- [ ] Implement Let's Encrypt DNS-01 challenge support (all providers + RFC2136)
- [ ] Implement Let's Encrypt TLS-ALPN-01 challenge support
- [ ] Implement Let's Encrypt HTTP-01 challenge support
- [ ] Check `/etc/letsencrypt/live` for existing certificates
- [ ] Save certificates to `/etc/{projectname}/ssl/certs`
- [ ] Implement dual-port HTTP/HTTPS support
- [ ] Add certificate management to scheduler

---

## Phase 4: Built-in Scheduler

- [ ] Design scheduler system for periodic tasks
- [ ] Implement certificate renewal scheduling
- [ ] Implement notification scheduling
- [ ] Add scheduler configuration options
- [ ] Make scheduler configurable via config file

---

## Phase 5: Notification System

- [ ] Design notification system (bell icon in UI)
- [ ] Implement notification storage (config-based)
- [ ] Implement notification display in web UI
- [ ] Add notification management via config file
- [ ] Support admin announcements (downtime, updates, etc.)
- [ ] Add notification scheduling support

---

## Phase 6: Service Management Implementation

- [ ] Implement `--service` CLI flag
- [ ] Implement service detection (systemd, runit, launchd, Windows, BSD rc.d)
- [ ] Implement `start` command
- [ ] Implement `stop` command
- [ ] Implement `restart` command
- [ ] Implement `reload` command
- [ ] Implement `--install` command
- [ ] Implement `--uninstall` command
- [ ] Implement `--disable` command
- [ ] Implement `--help` for service commands
- [ ] Update all installation scripts to use built-in service management

---

## Phase 7: Web Frontend Enhancements

- [ ] Add CORS support (default: '*', configurable)
- [ ] Implement robots.txt support (file or config-based)
- [ ] Implement security.txt support (file or config-based)
- [ ] Add notification bell to UI
- [ ] Implement logo support (local or remote URL)
- [ ] Implement favicon support (local or remote URL)
- [ ] Add logo/favicon scaling if needed
- [ ] Ensure footer is always centered and at bottom
- [ ] Verify accessibility compliance
- [ ] Create missing web pages (browse, random, categories, api-docs)
- [ ] Add comprehensive tooltips/documentation where needed

---

## Phase 8: API Enhancements

- [ ] Add `.txt` extension support for all API endpoints
- [ ] Ensure web (/) and API (/api) match functionality
- [ ] Implement proper scoped routes
- [ ] Verify all routes are intuitive and simple
- [ ] Add comprehensive input validation
- [ ] Add input sanitization where appropriate
- [ ] Implement "save only valid, clear only invalid" pattern

---

## Phase 9: Docker & Container Updates

- [ ] Update Dockerfile to match SPEC (Alpine with curl, bash)
- [ ] Set internal port to 80
- [ ] Add proper volume mounts (/data, /config, /data/db if needed)
- [ ] Add all required meta labels
- [ ] Update docker-compose.yml (remove version, no build)
- [ ] Create custom network with proper naming
- [ ] Update volume paths (./rootfs/data, ./rootfs/config)
- [ ] Set port mapping for production (172.17.0.1:{randomport}:80)
- [ ] Update Makefile docker target for buildx (arm64 + amd64)
- [ ] Test Docker build and run

---

## Phase 10: Makefile Improvements

- [ ] Simplify Makefile structure
- [ ] Implement VERSION env var support
- [ ] Implement auto-version increment via release.txt
- [ ] Update build target for all platforms
- [ ] Update release target to use `gh` CLI
- [ ] Add tag deletion before release if exists
- [ ] Implement strip for `-musl` binaries
- [ ] Update docker target for multi-arch buildx
- [ ] Ensure test target runs all tests
- [ ] Verify binary naming: {projectname}-{os}-{arch}

---

## Phase 11: Installation Scripts Updates

- [ ] Create scripts/README.md with install instructions at top
- [ ] Update install.sh for full OS/distro agnosticism
- [ ] Update linux.sh for SPEC compliance
- [ ] Update macos.sh for SPEC compliance
- [ ] Update windows.ps1 for SPEC compliance
- [ ] Implement user creation (system user, UID/GID 100-999)
- [ ] Set user home to config or data directory
- [ ] Use built-in service management where possible

---

## Phase 12: Logging System

- [ ] Implement access.log in Apache format
- [ ] Make access log format configurable
- [ ] Implement proper application logging
- [ ] Follow best practices for all logs
- [ ] Configure log rotation

---

## Phase 13: ReadTheDocs Configuration

- [ ] Create .readthedocs.yml configuration
- [ ] Set theme to Dracula
- [ ] Configure naming: apimgr-jokes.readthedocs.io
- [ ] Create documentation source files
- [ ] Test documentation build

---

## Phase 14: Security & Validation

- [ ] Implement comprehensive input validation everywhere
- [ ] Implement input sanitization where appropriate
- [ ] Ensure security doesn't block usability
- [ ] Review all endpoints for security issues
- [ ] Add security headers to responses
- [ ] Implement rate limiting properly

---

## Phase 15: Mobile & Accessibility

- [ ] Verify mobile-first design
- [ ] Test on screens ≥720px (90% width)
- [ ] Test on screens <720px (98% width)
- [ ] Ensure full accessibility compliance
- [ ] Add ARIA labels where needed
- [ ] Test keyboard navigation
- [ ] Test screen reader compatibility

---

## Phase 16: Testing & Quality Assurance

- [ ] Create comprehensive test suite in ./tests
- [ ] Test all CLI commands
- [ ] Test configuration auto-creation
- [ ] Test configuration live-reload
- [ ] Test Let's Encrypt integration
- [ ] Test scheduler functionality
- [ ] Test notification system
- [ ] Test service management
- [ ] Test all API endpoints
- [ ] Test all web pages
- [ ] Test multi-platform builds
- [ ] Test Docker container
- [ ] Test installation scripts (all OSes)

---

## Phase 17: Documentation Updates

- [ ] Update README.md with full SPEC compliance
- [ ] Ensure production instructions before development
- [ ] Update SPEC.md with complete details
- [ ] Update AI.md with current status
- [ ] Add comprehensive comments to code
- [ ] Document all configuration options
- [ ] Create user-friendly help text

---

## Phase 18: Final Cleanup

- [ ] Update .gitignore for SPEC compliance
- [ ] Update .dockerignore for SPEC compliance
- [ ] Ensure base directory is organized and clean
- [ ] Remove any database references
- [ ] Verify all file-based configuration
- [ ] Clean up temporary files
- [ ] Run final build and tests
- [ ] Verify all platforms build successfully

---

## Phase 19: Validation Checklist

- [ ] Single static binary with all assets embedded
- [ ] Builds for all platforms (Linux, BSD, macOS, Windows × AMD64, ARM64)
- [ ] File-based YAML configuration (no database)
- [ ] Configuration auto-creates on first run
- [ ] Configuration live-reload works
- [ ] Let's Encrypt support (all challenge types)
- [ ] Built-in scheduler works
- [ ] Notification system works
- [ ] Service management works (all service managers)
- [ ] Web frontend is fully functional
- [ ] All web pages work
- [ ] PWA support works
- [ ] Both themes work (dark default, light)
- [ ] Mobile responsive (720px breakpoint)
- [ ] Full accessibility
- [ ] REST API works (/api/v1)
- [ ] GraphQL works
- [ ] Swagger works
- [ ] All .txt endpoints work
- [ ] CORS configured properly
- [ ] robots.txt works
- [ ] security.txt works
- [ ] Logging works (Apache format)
- [ ] Docker container works
- [ ] Docker compose works
- [ ] All installation scripts work
- [ ] Version auto-increment works
- [ ] GitHub release works
- [ ] ReadTheDocs configured
- [ ] All CLI commands work
- [ ] --help, --version, --status work without sudo
- [ ] Input validation everywhere
- [ ] Input sanitization where needed
- [ ] Security best practices followed
- [ ] Console output is pretty with emojis
- [ ] Documentation is comprehensive
- [ ] AI.md is in sync

---

## Notes

- This is a comprehensive rebuild to full SPEC compliance
- Each phase should be completed before moving to the next
- Test thoroughly after each major change
- Ask questions if anything is unclear
- Keep AI.md updated throughout the process

---

## Current Phase

**Starting Phase**: 1 - Project Structure Reorganization

---

## Completion Status

**Total Tasks**: ~150+
**Completed**: 0
**In Progress**: 0
**Pending**: ~150+

---

**Last Updated**: 2025-11-25
