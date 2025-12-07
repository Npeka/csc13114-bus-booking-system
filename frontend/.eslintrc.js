module.exports = {
  root: true,
  extends: ["next/core-web-vitals", "next/typescript"],
  rules: {
    "@typescript-eslint/no-unused-vars": [
      "warn",
      {
        argsIgnorePattern: "^_",
        varsIgnorePattern: "^_",
      },
    ],
    "@typescript-eslint/no-explicit-any": "error",
    "react/no-unescaped-entities": "error",
    "react-hooks/set-state-in-effect": "error",
    "react-hooks/purity": "error",
    "@next/next/no-img-element": "warn",
  },
};
