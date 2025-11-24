# Title
One-line imperative summary, e.g. `feat: add --foo global flag and deprecate --bar`

# What changed
- One short, concrete bullet per change (what was moved/removed/added).
- Use present tense: "Add", "Remove", "Move", "Fix".
- Keep lines short; reviewers should be able to scan quickly.
- Group related changes together.
  - Example of a related change in a sub-bullet
> **Tip:** 90% of the time this section is enough. Fill other sections only when needed.

<details>
<summary>Optional / Expanded sections</summary>

### Migration / Impact
- List breaking changes or deprecated flags.
- How to adapt commands, flags, configs (examples if useful).

### Before / After examples
```bash
# Before
ionosctl my-resource command --bar

# After
ionosctl my-resource command --bar
Error: flag '--bar' has been deprecated, use '--foo' instead
```
</details>
