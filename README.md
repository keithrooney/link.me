# Anchorly

## Tests

### Setup

1. Create a `settings.json` file in the `.vscode` directory.
2. Add the settings below to the file.

```json
{
    "go.inferGopath": false,
    "go.testEnvVars": {
        "ANCHORLY_TOKEN_KEY": "this.is.a.very.secret.value!"
    }
}
```
