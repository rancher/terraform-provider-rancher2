// Basic ESLint configuration file
import eslintPluginJs from '@eslint/js';

export default [
  {
    // Applies to all files by default
    // Recommended rules from ESLint itself
    ...eslintPluginJs.configs.recommended,

    // Define rules specific to your project
    rules: {
      // Example: Require semicolons at the end of statements
      'semi': ['error', 'always'],
      // Example: Enforce 2-space indentation
      'indent': ['error', 2],
      // Example: Prevent unused variables
      'no-unused-vars': 'warn',
    },
  },
];
