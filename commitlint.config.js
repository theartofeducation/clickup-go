const conventionalConfig = require("@commitlint/config-conventional")
const { rules: defaultRules } = conventionalConfig

/**
 * Extend "@commitlint/config-conventional" with two adjustments to
 * types allowed: change "feat" to "feature", and add "wip"
 */
module.exports = Object.assign({}, conventionalConfig, {
    rules: Object.assign({}, defaultRules, {
        "type-enum": [2, "always", [
            "build",
            "chore",
            "ci",
            "docs",
            "feature",
            "fix",
            "performance",
            "refactor",
            "revert",
            "style",
            "test",
            "wip"
        ]]
    })
})