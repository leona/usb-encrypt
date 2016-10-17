package values

const vaultDir = "vault"

// Patterns
const patternValidateName = `"^[a-zA-Z0-9_-]*$"`
const patternValidateDir = `"^[a-zA-Z0-9_-/.]*$"`

// Error messages
const ErrorEncryptFromUsb = "Do not encrypt files stored on the current USB as residual files can be recovered."
const ErrorVaultName = "Vault name can only contain upper/lower case letters, numbers, underscores or dashes."
const ErrorVaultPath = "Invalid path"

// Messages
const MessageVaultName = "Enter vault name: "
const MessageVaultPath = "Enter directory to encrypt from: "