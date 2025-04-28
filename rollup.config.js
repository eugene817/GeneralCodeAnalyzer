import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';

export default {
  input: 'lsp-entry.js',
  output: {
    file: 'api/static/vendor/lsp-bundle.js',
    format: 'esm'
  },
  plugins: [ resolve(), commonjs() ]
};

