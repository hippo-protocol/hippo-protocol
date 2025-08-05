# Hippo Protocol SDK

---

```bash
# Build for nodejs
wasm-pack build --target=nodejs --out-dir=../js/pkg/node --out-name=core && rm ../js/pkg/node/.gitignore

# Build for web
wasm-pack build --target=bundler --out-dir=../js/pkg/bundle --out-name=core && rm ../js/pkg/bundle/.gitignore

# Package for npm
wasm-pack pack ../js/pkg --help
```

If you don't have `wasm-pack` yet, follow the [installation guide](https://rustwasm.github.io/wasm-pack/installer/)

### Tests

```bash
cargo test
wasm-pack test --node
```
