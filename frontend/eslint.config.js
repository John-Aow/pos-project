import js from '@eslint/js';
import vue from 'eslint-plugin-vue';
import prettier from 'eslint-config-prettier';

export default [
  {
    ignores: ['dist/', 'build/', 'coverage/', 'node_modules/', '*.min.js'],
  },
  js.configs.recommended,
  ...vue.configs['flat/recommended'],
  prettier,
  {
    files: ['**/*.vue', '**/*.ts'],
    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
    },
    rules: {
      'vue/multi-word-component-names': 'off',
    },
  },
  {
    files: ['tests/**/*.ts'],
    languageOptions: {
      globals: {
        Event: 'readonly',
        MediaQueryList: 'readonly',
        window: 'readonly',
      },
    },
  },
];
