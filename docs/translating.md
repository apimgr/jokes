# Translating

The application embeds translation catalogs and resolves a language from the `lang` query parameter, saved language cookie, then `Accept-Language`, falling back to English. The registry currently includes English (`en`), Spanish (`es`), French (`fr`), Arabic (`ar`, right-to-left), German (`de`), and Italian (`it`). Catalog parity is validated by `i18n.ValidateLocales`; user-facing strings must use catalog keys rather than hardcoded text.

## Required catalog format

The English catalog, `src/common/i18n/locales/en.json`, is the source of truth for message keys. Translators change message values only. Each supported locale must have the same message-key set and placeholder names as English, with the CLDR plural categories required for that language. Locale JSON is UTF-8, two-space indented, has no comments or trailing commas, and includes the locale metadata (`language`, `name`, `native_name`, `direction`, and `version`). A translator may include their preferred name or handle in `meta.translators`.

Security-report acknowledgment subject and body use the `email` catalog keys. The other built-in notification templates have not yet been migrated to locale keys.

## Translation and review

- Preserve the meaning and intent of the English source. Write complete sentences rather than fragments that must be joined in code.
- Keep every `{placeholder}` exactly spelled and braced as in English; its position may change to fit the translated sentence.
- Use one consistent register and tone within a locale. Keep product names, API/protocol names, code, URLs, filenames, configuration keys, and literal option values unchanged.
- Use the glossary below consistently. Add a term when it becomes part of recurring UI or documentation language.
- Check text at 320, 768, and 1280 pixel widths; verify keyboard use, accessible names, and right-to-left layout for RTL locales.
- A maintainer should check key parity, placeholders, plural categories, glossary use, layout, and the README language count before merging. Native-speaker review is welcome.

| English term | Usage |
|---|---|
| joke | A single joke item, whether one-line or two-part |
| setup / punchline | The two parts of a two-part joke |
| category | A public grouping used to filter the joke collection |
| owner token | Anonymous credential that grants access to its owner's resources |
| linked token | Separately revocable token created by an owner for a client |
| operator token | Private server credential; never translate or expose its value |
| rotate / revoke | Replace a credential in place / permanently disable a credential |
| report | A security, abuse, DMCA, or privacy submission |

## Currently enabled languages and credits

The six required locale catalogs ship in `src/common/i18n/locales/`. No translator credits are recorded yet. The current catalogs and runtime registry must be updated together whenever a locale or message key changes; do not claim a locale is supported unless its embedded catalog passes parity validation.
